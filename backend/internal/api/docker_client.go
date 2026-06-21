package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DockerClient 通过 unix socket 直接调用 Docker Engine API。
// 之所以不用 docker SDK / docker CLI：依赖更轻、容器内可用、便于嵌入。
type DockerClient struct {
	socketPath string
	client     *http.Client
}

// NewDockerClient 创建一个连接本地 docker.sock 的客户端。
func NewDockerClient() *DockerClient {
	socket := "/var/run/docker.sock"
	return &DockerClient{
		socketPath: socket,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					d := net.Dialer{Timeout: 5 * time.Second}
					return d.DialContext(ctx, "unix", socket)
				},
			},
		},
	}
}

// Available 检查 docker socket 是否可连通。
func (d *DockerClient) Available() bool {
	conn, err := net.DialTimeout("unix", d.socketPath, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (d *DockerClient) do(method, path string, query url.Values) (*http.Response, error) {
	u := "http://docker" + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}
	return d.client.Do(req)
}

// ContainerSummary docker /containers/json 返回的字段子集。
type ContainerSummary struct {
	ID      string            `json:"Id"`
	Name    string            `json:"Name"`
	Image   string            `json:"Image"`
	Command string            `json:"Command"`
	State   string            `json:"State"`
	Status  string            `json:"Status"`
	Created int64             `json:"Created"`
	Labels  map[string]string `json:"Labels"`
	Names   []string          `json:"Names"`
	Mounts  []struct {
		Type        string `json:"Type"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Mode        string `json:"Mode"`
		RW          bool   `json:"RW"`
	} `json:"Mounts"`
}

// ContainerInspect docker /containers/{id}/json 返回的字段子集。
type ContainerInspect struct {
	ID      string                 `json:"Id"`
	Name    string                 `json:"Name"`
	State   *ContainerState        `json:"State,omitempty"`
	Config  *ContainerConfig       `json:"Config,omitempty"`
	Created string                 `json:"Created"`
}

type ContainerState struct {
	Status     string    `json:"Status"`
	Running    bool      `json:"Running"`
	Paused     bool      `json:"Paused"`
	ExitCode   int       `json:"ExitCode"`
	StartedAt  string    `json:"StartedAt"`
	FinishedAt string    `json:"FinishedAt"`
	Error      string    `json:"Error"`
}

type ContainerConfig struct {
	Image string            `json:"Image"`
	Cmd   []string          `json:"Cmd"`
	Env   []string          `json:"Env"`
	Labels map[string]string `json:"Labels"`
}

// ListContainers 列出所有容器（包括已停止的）。
func (d *DockerClient) ListContainers(all bool) ([]ContainerSummary, error) {
	q := url.Values{}
	if all {
		q.Set("all", "true")
	}
	q.Set("size", "false")
	resp, err := d.do("GET", "/v1.43/containers/json", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker list containers: %s: %s", resp.Status, string(body))
	}
	var list []ContainerSummary
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

// InspectContainer 获取容器详情。
func (d *DockerClient) InspectContainer(idOrName string) (*ContainerInspect, error) {
	raw, err := d.RawInspect(idOrName)
	if err != nil {
		return nil, err
	}
	var info ContainerInspect
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// RawInspect 获取原始 JSON（用于按需提取非结构化字段，如 ports/network）。
func (d *DockerClient) RawInspect(idOrName string) ([]byte, error) {
	resp, err := d.do("GET", "/v1.43/containers/"+idOrName+"/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker inspect: %s: %s", resp.Status, string(body))
	}
	return io.ReadAll(resp.Body)
}

// StartContainer 启动一个已存在的容器（适用于 restart=no 的 one-shot 容器）。
func (d *DockerClient) StartContainer(idOrName string) error {
	resp, err := d.do("POST", "/v1.43/containers/"+idOrName+"/start", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker start: %s: %s", resp.Status, string(body))
	}
	return nil
}

// StopContainer 停止容器。timeout 秒后发 SIGKILL。
func (d *DockerClient) StopContainer(idOrName string, timeoutSec int) error {
	q := url.Values{}
	if timeoutSec <= 0 {
		timeoutSec = 10
	}
	q.Set("t", fmt.Sprintf("%d", timeoutSec))
	resp, err := d.do("POST", "/v1.43/containers/"+idOrName+"/stop", q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker stop: %s: %s", resp.Status, string(body))
	}
	return nil
}

// RestartContainer 重启容器。
func (d *DockerClient) RestartContainer(idOrName string, timeoutSec int) error {
	q := url.Values{}
	if timeoutSec <= 0 {
		timeoutSec = 10
	}
	q.Set("t", fmt.Sprintf("%d", timeoutSec))
	resp, err := d.do("POST", "/v1.43/containers/"+idOrName+"/restart", q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker restart: %s: %s", resp.Status, string(body))
	}
	return nil
}

// RemoveContainer 删除容器。force=true 即使在运行也删除；volumes=true 同时删除匿名卷。
func (d *DockerClient) RemoveContainer(idOrName string, force, removeVolumes bool) error {
	q := url.Values{}
	q.Set("force", boolStr(force))
	q.Set("v", boolStr(removeVolumes))
	resp, err := d.do("DELETE", "/v1.43/containers/"+idOrName, q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("docker remove: %s: %s", resp.Status, string(body))
	}
	return nil
}

// ContainerStats 获取容器资源占用（单次快照，非 stream）。
type ContainerStats struct {
	Name      string  `json:"name"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryUsage int64  `json:"memoryUsage"`
	MemoryLimit int64  `json:"memoryLimit"`
	MemoryPercent float64 `json:"memoryPercent"`
	NetworkRx   uint64 `json:"networkRx"`
	NetworkTx   uint64 `json:"networkTx"`
	ReadTime    string `json:"readTime"`
}

func (d *DockerClient) Stats(idOrName string) (*ContainerStats, error) {
	q := url.Values{}
	q.Set("stream", "false")
	resp, err := d.do("GET", "/v1.43/containers/"+idOrName+"/stats", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("docker stats: %s: %s", resp.Status, string(body))
	}
	// Docker stats 返回结构很复杂，我们只取需要的字段
	var raw struct {
		Name     string `json:"name"`
		Read     string `json:"read"`
		PreCPU   struct {
			CPUUsage struct {
				TotalUsage uint64 `json:"total_usage"`
				PercpuUsage []uint64 `json:"percpu_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage uint64 `json:"system_cpu_usage"`
			OnlineCPUs     int    `json:"online_cpus"`
		} `json:"precpu_stats"`
		CPU struct {
			CPUUsage struct {
				TotalUsage  uint64   `json:"total_usage"`
				PercpuUsage []uint64 `json:"percpu_usage"`
			} `json:"cpu_usage"`
			SystemCPUUsage uint64 `json:"system_cpu_usage"`
			OnlineCPUs     int    `json:"online_cpus"`
		} `json:"cpu_stats"`
		Memory struct {
			Usage uint64 `json:"usage"`
			Limit uint64 `json:"limit"`
		} `json:"memory_stats"`
		Networks map[string]struct {
			RxBytes uint64 `json:"rx_bytes"`
			TxBytes uint64 `json:"tx_bytes"`
		} `json:"networks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	out := &ContainerStats{
		Name:      strings.TrimPrefix(raw.Name, "/"),
		MemoryUsage: int64(raw.Memory.Usage),
		MemoryLimit: int64(raw.Memory.Limit),
		ReadTime:    raw.Read,
	}
	if raw.Memory.Limit > 0 {
		out.MemoryPercent = float64(raw.Memory.Usage) / float64(raw.Memory.Limit) * 100
	}
	// CPU 百分比计算
	cpuDelta := float64(raw.CPU.CPUUsage.TotalUsage) - float64(raw.PreCPU.CPUUsage.TotalUsage)
	sysDelta := float64(raw.CPU.SystemCPUUsage) - float64(raw.PreCPU.SystemCPUUsage)
	onlineCPUs := raw.CPU.OnlineCPUs
	if onlineCPUs == 0 {
		onlineCPUs = len(raw.CPU.CPUUsage.PercpuUsage)
	}
	if sysDelta > 0 && cpuDelta > 0 && onlineCPUs > 0 {
		out.CPUPercent = (cpuDelta / sysDelta) * float64(onlineCPUs) * 100
	}
	for _, nw := range raw.Networks {
		out.NetworkRx += nw.RxBytes
		out.NetworkTx += nw.TxBytes
	}
	return out, nil
}

// Logs 获取容器日志。
// 注意：docker daemon 的 /logs 端点使用 multiplexed stream（每个 frame 8 字节 header），
// 当 container 是 json-file driver 时需要按 header 解 frame。这里使用 follow=false 简化。
func (d *DockerClient) Logs(idOrName string, tail int, since time.Time) (string, error) {
	q := url.Values{}
	q.Set("stdout", "true")
	q.Set("stderr", "true")
	if tail > 0 {
		q.Set("tail", fmt.Sprintf("%d", tail))
	} else {
		q.Set("tail", "all")
	}
	if !since.IsZero() {
		q.Set("since", fmt.Sprintf("%d", since.Unix()))
	}
	resp, err := d.do("GET", "/v1.43/containers/"+idOrName+"/logs", q)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("docker logs: %s: %s", resp.Status, string(body))
	}
	return demultiplexDockerStream(resp.Body), nil
}

// demultiplexDockerStream 解析 docker daemon 的 multiplexed stream 格式。
// 每个 frame 头部 8 字节：[streamType(1), 0,0,0, length(4)]，随后是 payload。
func demultiplexDockerStream(r io.Reader) string {
	var out strings.Builder
	br := bufio.NewReader(r)
	for {
		header := make([]byte, 8)
		if _, err := io.ReadFull(br, header); err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		// header[0] = stream type (1=stdout, 2=stderr); length is big-endian uint32 in header[4:8]
		length := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if length == 0 {
			continue
		}
		buf := make([]byte, length)
		if _, err := io.ReadFull(br, buf); err != nil {
			break
		}
		out.Write(buf)
	}
	// 如果上面解析失败（例如 daemon 已经退化成 plain text），就回退到原始读取。
	if out.Len() == 0 {
		raw, _ := io.ReadAll(r)
		return string(raw)
	}
	return out.String()
}

// IsResticContainer 判断一个容器是否是 restic 相关容器。
// 命中任一条件即视为 restic 容器：
//   - 镜像名包含 "restic"
//   - Cmd/Entrypoint 中包含 "restic"
//   - Env 中存在 RESTIC_REPOSITORY
//   - 标签里有 nas-dashboard 自动注册的标记
func IsResticContainer(s ContainerSummary) bool {
	if strings.Contains(strings.ToLower(s.Image), "restic") {
		return true
	}
	if strings.Contains(strings.ToLower(s.Command), "restic") {
		return true
	}
	if strings.Contains(strings.ToLower(s.Name), "restic") {
		return true
	}
	for _, n := range s.Names {
		if strings.Contains(strings.ToLower(n), "restic") {
			return true
		}
	}
	// 通过 inspect 进一步看 env
	info, err := globalDocker.InspectContainer(s.ID)
	if err == nil && info.Config != nil {
		for _, e := range info.Config.Env {
			if strings.HasPrefix(e, "RESTIC_REPOSITORY=") ||
				strings.HasPrefix(e, "RESTIC_PASSWORD=") {
				return true
			}
		}
		for _, c := range info.Config.Cmd {
			if strings.Contains(strings.ToLower(c), "restic") {
				return true
			}
		}
	}
	return false
}

// 全局单例，避免每次请求都重新创建连接池。
var globalDocker = NewDockerClient()
