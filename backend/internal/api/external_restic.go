package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ExternalContainerInfo 描述一个被发现的"外部 restic 容器"。
// 这些容器通常由独立的 docker-compose 项目部署（例如 /data/docker/restic/），
// 不是 nas-dashboard 创建的，但我们仍然可以观察、触发、浏览其快照。
type ExternalContainerInfo struct {
	ID            string            `json:"id"`            // 容器短 ID
	Name          string            `json:"name"`          // 容器名（去前导斜杠）
	Image         string            `json:"image"`         // 镜像名
	State         string            `json:"state"`         // running / exited / ...
	Status        string            `json:"status"`        // docker 给的人类可读状态
	ExitCode      int               `json:"exitCode"`      // 最近一次退出码
	Created       time.Time         `json:"created"`       // 创建时间
	StartedAt     time.Time         `json:"startedAt"`     // 最近一次启动时间
	FinishedAt    time.Time         `json:"finishedAt"`    // 最近一次结束时间
	Command       string            `json:"command"`       // 容器执行的命令（已拼接）
	ComposeProject string           `json:"composeProject"` // docker-compose 项目名（来自 label）
	Repo          string            `json:"repo"`          // 从 env 中解析的 RESTIC_REPOSITORY
	SourcePath    string            `json:"sourcePath"`    // 从 command 中解析的备份源路径
	Retention     string            `json:"retention"`     // 从 command 中解析的 forget 保留参数
	Env           map[string]string `json:"env"`           // env 变量（敏感字段已脱敏）
	Mounts        []MountInfo       `json:"mounts"`        // 卷挂载（用于判断可访问哪些主机路径）
}

type MountInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
}

// ListExternalResticContainers GET /api/storage/backup/external/containers
// 扫描所有 Docker 容器，返回看起来是 restic 相关的容器。
func ListExternalResticContainers(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusOK, gin.H{
			"containers": []any{},
			"available":  false,
			"message":    "Docker socket not accessible; please mount /var/run/docker.sock into backend container",
		})
		return
	}

	all, err := globalDocker.ListContainers(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := []ExternalContainerInfo{}
	for _, cnt := range all {
		if !IsResticContainer(cnt) {
			continue
		}
		info := buildExternalContainerInfo(cnt)
		out = append(out, info)
	}
	c.JSON(http.StatusOK, gin.H{
		"containers": out,
		"total":      len(out),
		"available":  true,
	})
}

func buildExternalContainerInfo(cnt ContainerSummary) ExternalContainerInfo {
	// docker /containers/json 返回 Names 数组（带前导 /），取第一个作为容器名
	name := cnt.Name
	if name == "" && len(cnt.Names) > 0 {
		name = cnt.Names[0]
	}
	info := ExternalContainerInfo{
		ID:             cnt.ID,
		Name:           strings.TrimPrefix(name, "/"),
		Image:          cnt.Image,
		State:          cnt.State,
		Status:         cnt.Status,
		Command:        cnt.Command,
		ComposeProject: cnt.Labels["com.docker.compose.project"],
		Env:            map[string]string{},
	}
	info.Created = dockerTime(cnt.Created)
	for _, m := range cnt.Mounts {
		info.Mounts = append(info.Mounts, MountInfo{
			Source: m.Source, Destination: m.Destination, Mode: m.Mode,
		})
	}

	// 进一步 inspect 拿 env / state 详情
	if detail, err := globalDocker.InspectContainer(cnt.ID); err == nil {
		if detail.State != nil {
			info.ExitCode = detail.State.ExitCode
			info.StartedAt = parseTime(detail.State.StartedAt)
			info.FinishedAt = parseTime(detail.State.FinishedAt)
		}
		if detail.Config != nil {
			info.Command = strings.Join(detail.Config.Cmd, " ")
			for _, e := range detail.Config.Env {
				kv := strings.SplitN(e, "=", 2)
				if len(kv) != 2 {
					continue
				}
				k, v := kv[0], kv[1]
				switch {
				case k == "RESTIC_REPOSITORY":
					info.Repo = v
					info.Env[k] = v
				case k == "RESTIC_PASSWORD", strings.Contains(k, "PASSWORD"),
					strings.Contains(k, "TOKEN"), strings.Contains(k, "SECRET"),
					strings.Contains(k, "KEY"), strings.Contains(k, "PASS"):
					info.Env[k] = maskValue(v)
				default:
					info.Env[k] = v
				}
			}
		}
	}

	// 从 command 中提取源路径与保留策略
	info.SourcePath = extractBackupPath(info.Command)
	info.Retention = extractRetention(info.Command)
	return info
}

// extractBackupPath 从 "restic backup /data ..." 这种命令中提取路径。
func extractBackupPath(cmd string) string {
	parts := strings.Fields(cmd)
	for i, p := range parts {
		if p == "backup" && i+1 < len(parts) {
			// 收集连续的路径（非 flag）参数
			paths := []string{}
			for j := i + 1; j < len(parts); j++ {
				if strings.HasPrefix(parts[j], "-") {
					break
				}
				paths = append(paths, parts[j])
			}
			if len(paths) > 0 {
				return strings.Join(paths, " ")
			}
		}
	}
	return ""
}

// extractRetention 从 "restic forget --keep-daily 7 ..." 提取保留参数。
func extractRetention(cmd string) string {
	idx := strings.Index(cmd, "forget")
	if idx < 0 {
		return ""
	}
	tail := cmd[idx:]
	out := []string{}
	parts := strings.Fields(tail)
	for _, p := range parts {
		if strings.HasPrefix(p, "--keep") {
			out = append(out, p)
		} else if len(out) > 0 && !strings.HasPrefix(p, "-") && !strings.Contains(p, "prune") {
			// 保留策略的参数值（数字 / duration）
			out = append(out, p)
		}
	}
	return strings.Join(out, " ")
}

// maskValue 脱敏：保留首尾 2 字符，中间 ****
func maskValue(v string) string {
	if len(v) <= 4 {
		return "****"
	}
	return v[:2] + "****" + v[len(v)-2:]
}

// dockerTime: docker /containers/json 返回的 Created 是 unix 时间戳（秒）。
func dockerTime(unix int64) time.Time {
	if unix == 0 {
		return time.Time{}
	}
	return time.Unix(unix, 0)
}

// parseTime: docker /containers/{id}/json 里的时间字符串是 RFC3339。
func parseTime(s string) time.Time {
	if s == "" || s == "0001-01-01T00:00:00Z" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// StartExternalContainer POST /api/storage/backup/external/containers/:name/start
// 重新启动一个 restic one-shot 容器（适用于 restart=no 的备份容器）。
// 注意：对 restart=no 的容器，docker 会在执行完命令后退出，但每次 start 都会重新执行。
func StartExternalContainer(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Docker socket not accessible"})
		return
	}
	name := c.Param("name")
	// 先 inspect 检查状态
	info, err := globalDocker.InspectContainer(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found: " + err.Error()})
		return
	}
	if info.State != nil && info.State.Running {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "container is already running",
			"status":  "running",
			"started": info.State.StartedAt,
		})
		return
	}
	if err := globalDocker.StartContainer(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"message": "container started",
		"name":    name,
	})
}

// GetExternalContainerLogs GET /api/storage/backup/external/containers/:name/logs?tail=500
func GetExternalContainerLogs(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Docker socket not accessible"})
		return
	}
	name := c.Param("name")
	tail := 500
	if t := c.Query("tail"); t != "" {
		fmt.Sscanf(t, "%d", &tail)
	}
	logs, err := globalDocker.Logs(name, tail, time.Time{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"logs": logs,
		"name": name,
	})
}

// GetExternalContainerStatus GET /api/storage/backup/external/containers/:name/status
func GetExternalContainerStatus(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Docker socket not accessible"})
		return
	}
	name := c.Param("name")
	info, err := globalDocker.InspectContainer(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	resp := gin.H{
		"name":      strings.TrimPrefix(info.Name, "/"),
		"state":     "unknown",
		"exitCode":  0,
		"startedAt": "",
	}
	if info.State != nil {
		resp["state"] = info.State.Status
		resp["running"] = info.State.Running
		resp["exitCode"] = info.State.ExitCode
		resp["startedAt"] = info.State.StartedAt
		resp["finishedAt"] = info.State.FinishedAt
	}
	c.JSON(http.StatusOK, resp)
}

// GetExternalContainerRawEnv 返回外部容器的完整 env（含未脱敏密码）。
// 仅供后端内部使用（导入仓库），不通过 HTTP 暴露。
func GetExternalContainerRawEnv(name string) (map[string]string, error) {
	info, err := globalDocker.InspectContainer(name)
	if err != nil {
		return nil, err
	}
	if info.Config == nil {
		return map[string]string{}, nil
	}
	out := map[string]string{}
	for _, e := range info.Config.Env {
		kv := strings.SplitN(e, "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out, nil
}
