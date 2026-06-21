package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// jsonUnmarshal 包内便捷函数
func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// =====================================================================
// 数据结构
// =====================================================================

// AppContainer 一个容器实例的完整描述（用于应用管理 UI）。
type AppContainer struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	State         string            `json:"state"`
	Status        string            `json:"status"`
	Health        string            `json:"health"` // healthy / unhealthy / starting / ""
	ExitCode      int               `json:"exitCode"`
	Created       time.Time         `json:"created"`
	StartedAt     time.Time         `json:"startedAt"`
	FinishedAt    time.Time         `json:"finishedAt"`
	RestartCount  int               `json:"restartCount"`
	Ports         []PortBinding     `json:"ports"`
	IP            string            `json:"ip"`
	Command       string            `json:"command"`
	ComposeProject string           `json:"composeProject"`
	ComposeService string           `json:"composeService"`
	ComposeDir    string            `json:"composeDir"`
	ComposeFile   string            `json:"composeFile"`
	Labels        map[string]string `json:"labels"`
	Mounts        []MountInfo       `json:"mounts"`
	EnvCount      int               `json:"envCount"`
	Stats         *ContainerStats   `json:"stats,omitempty"`
}

// PortBinding 主机端口映射。
type PortBinding struct {
	HostIP        string `json:"hostIp"`
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
	URL           string `json:"url,omitempty"` // 推断的可点击 URL
}

// ComposeProject 一个 compose 项目（一组关联的容器）。
type ComposeProject struct {
	Name          string   `json:"name"`
	ConfigFile    string   `json:"configFile"`
	WorkingDir    string   `json:"workingDir"`
	Containers    []string `json:"containers"` // 容器名列表
	RunningCount  int      `json:"runningCount"`
	TotalCount    int      `json:"totalCount"`
	Category      string   `json:"category"` // 由 catalog 标注：file-server / photo / dev / audio / backup / system
	Description   string   `json:"description"`
	IconHint      string   `json:"iconHint"`
}

// =====================================================================
// 容器列表 / 项目分组
// =====================================================================

// ListAppContainers GET /api/apps/containers
// 列出所有容器（默认包含已停止），按 compose 项目分组排序。
// query: stats=true 时附带 CPU/内存快照（更慢）。
func ListAppContainers(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusOK, gin.H{
			"containers": []any{},
			"available":  false,
			"message":    "Docker socket not accessible",
		})
		return
	}
	withStats := c.Query("stats") == "true"

	all, err := globalDocker.ListContainers(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]AppContainer, 0, len(all))
	for _, cnt := range all {
		app := buildAppContainer(cnt, withStats)
		out = append(out, app)
	}
	// 先按 compose project，再按 service name 排序
	sort.Slice(out, func(i, j int) bool {
		pi, pj := out[i].ComposeProject, out[j].ComposeProject
		if pi != pj {
			if pi == "" {
				return false
			}
			if pj == "" {
				return true
			}
			return pi < pj
		}
		return out[i].Name < out[j].Name
	})

	c.JSON(http.StatusOK, gin.H{
		"containers": out,
		"total":      len(out),
		"available":  true,
	})
}

func buildAppContainer(cnt ContainerSummary, withStats bool) AppContainer {
	name := cnt.Name
	if name == "" && len(cnt.Names) > 0 {
		name = cnt.Names[0]
	}
	name = strings.TrimPrefix(name, "/")

	app := AppContainer{
		ID:             cnt.ID,
		Name:           name,
		Image:          cnt.Image,
		State:          cnt.State,
		Status:         cnt.Status,
		Command:        cnt.Command,
		ComposeProject: cnt.Labels["com.docker.compose.project"],
		ComposeService: cnt.Labels["com.docker.compose.service"],
		ComposeDir:     cnt.Labels["com.docker.compose.project.working_dir"],
		ComposeFile:    cnt.Labels["com.docker.compose.project.config_files"],
		Labels:         cnt.Labels,
		Created:        dockerTime(cnt.Created),
	}
	for _, m := range cnt.Mounts {
		app.Mounts = append(app.Mounts, MountInfo{
			Source: m.Source, Destination: m.Destination, Mode: m.Mode,
		})
	}

	// inspect 拿更详细信息
	if detail, err := globalDocker.InspectContainer(cnt.ID); err == nil {
		if detail.State != nil {
			app.ExitCode = detail.State.ExitCode
			app.StartedAt = parseTime(detail.State.StartedAt)
			app.FinishedAt = parseTime(detail.State.FinishedAt)
		}
		if detail.Config != nil {
			app.EnvCount = len(detail.Config.Env)
			if len(detail.Config.Cmd) > 0 && app.Command == "" {
				app.Command = strings.Join(detail.Config.Cmd, " ")
			}
		}
	}

	// 端口 / URL：docker /containers/json 里端口在 separate field
	// 我们已经丢掉了，从 inspect 再拿一次
	if detail, err := globalDocker.InspectContainer(cnt.ID); err == nil {
		app.Ports = extractPorts(detail)
		// IP
		if h, ok := detail.NetworkSettingsBase(); ok {
			app.IP = h
		}
	}

	if withStats {
		if st, err := globalDocker.Stats(cnt.ID); err == nil {
			app.Stats = st
		}
	}

	return app
}

// helper: 提取端口映射 + 自动生成可访问 URL
func extractPorts(info *ContainerInspect) []PortBinding {
	out := []PortBinding{}
	if info == nil {
		return out
	}
	// ContainerInspect 目前只解析了我们关心的字段；端口映射需要扩展结构
	// 这里从原始 HTTP body 再解一次
	raw, err := globalDocker.RawInspect(info.ID)
	if err != nil {
		return out
	}
	type portMapping struct {
		PrivatePort int    `json:"PrivatePort"`
		PublicPort  int    `json:"PublicPort"`
		Type        string `json:"Type"`
		IP          string `json:"IP"`
	}
	type rawInfo struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostIP   string `json:"HostIp"`
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	var r rawInfo
	if err := jsonUnmarshal(raw, &r); err != nil {
		return out
	}
	for contPort, hosts := range r.NetworkSettings.Ports {
		cp, proto := parsePortSpec(contPort)
		if len(hosts) == 0 {
			// 仅暴露给容器网络、未映射到主机
			out = append(out, PortBinding{
				ContainerPort: cp, Protocol: proto,
			})
			continue
		}
		for _, h := range hosts {
			hp := 0
			fmt.Sscanf(h.HostPort, "%d", &hp)
			pb := PortBinding{
				HostIP:        h.HostIP,
				HostPort:      hp,
				ContainerPort: cp,
				Protocol:      proto,
			}
			if proto == "tcp" && hp > 0 {
				pb.URL = guessURL(hp, cp)
			}
			out = append(out, pb)
		}
	}
	return out
}

func parsePortSpec(spec string) (port int, proto string) {
	proto = "tcp"
	if idx := strings.Index(spec, "/"); idx >= 0 {
		proto = spec[idx+1:]
		spec = spec[:idx]
	}
	fmt.Sscanf(spec, "%d", &port)
	return
}

// guessURL 推测端口对应的可访问 URL（用于 UI 点击跳转）。
func guessURL(hostPort, containerPort int) string {
	// 已知容器端口的协议
	switch containerPort {
	case 80, 8080, 8000, 3000, 5173:
		return fmt.Sprintf("http://localhost:%d", hostPort)
	case 443, 8443:
		return fmt.Sprintf("https://localhost:%d", hostPort)
	case 2283: // immich
		return fmt.Sprintf("http://localhost:%d", hostPort)
	case 5244: // alist
		return fmt.Sprintf("http://localhost:%d", hostPort)
	case 22:
		return fmt.Sprintf("ssh://localhost:%d", hostPort)
	}
	// 默认：如果是 80 类，http；其他未知时也尝试 http（用户可点击看是否通）
	if hostPort >= 80 && hostPort <= 65535 {
		return fmt.Sprintf("http://localhost:%d", hostPort)
	}
	return ""
}

// =====================================================================
// 容器操作
// =====================================================================

// AppContainerAction POST /api/apps/containers/:name/:action
// action: start | stop | restart | remove
func AppContainerAction(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Docker socket not accessible"})
		return
	}
	name := c.Param("name")
	action := c.Param("action")
	timeoutSec := 10
	var err error
	switch action {
	case "start":
		err = globalDocker.StartContainer(name)
	case "stop":
		err = globalDocker.StopContainer(name, timeoutSec)
	case "restart":
		err = globalDocker.RestartContainer(name, timeoutSec)
	case "remove":
		force := c.Query("force") == "true"
		volumes := c.Query("volumes") == "true"
		err = globalDocker.RemoveContainer(name, force, volumes)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action: " + action})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "action": action})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("container %s %s ed", name, action), "action": action})
}

// AppContainerLogs GET /api/apps/containers/:name/logs?tail=500
func AppContainerLogs(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"logs": logs, "name": name})
}

// AppContainerStats GET /api/apps/containers/:name/stats
func AppContainerStats(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Docker socket not accessible"})
		return
	}
	name := c.Param("name")
	stats, err := globalDocker.Stats(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// =====================================================================
// Compose 项目视图
// =====================================================================

// ListComposeProjects GET /api/apps/projects
// 把容器按 compose project 分组，并对已知项目（immich / samba / alist 等）附加元信息。
func ListComposeProjects(c *gin.Context) {
	if !globalDocker.Available() {
		c.JSON(http.StatusOK, gin.H{"projects": []any{}, "available": false})
		return
	}
	all, err := globalDocker.ListContainers(true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 按 project 分组
	type group struct {
		names    []string
		running  int
		total    int
		workDir  string
		cfgFile  string
	}
	groups := map[string]*group{}
	for _, cnt := range all {
		name := cnt.Name
		if name == "" && len(cnt.Names) > 0 {
			name = cnt.Names[0]
		}
		name = strings.TrimPrefix(name, "/")
		proj := cnt.Labels["com.docker.compose.project"]
		if proj == "" {
			proj = "(standalone)"
		}
		g, ok := groups[proj]
		if !ok {
			g = &group{}
			groups[proj] = g
		}
		g.names = append(g.names, name)
		g.total++
		if cnt.State == "running" {
			g.running++
		}
		if g.workDir == "" {
			g.workDir = cnt.Labels["com.docker.compose.project.working_dir"]
		}
		if g.cfgFile == "" {
			g.cfgFile = cnt.Labels["com.docker.compose.project.config_files"]
		}
	}

	projects := []ComposeProject{}
	for name, g := range groups {
		cat, desc, icon := knownProject(name)
		projects = append(projects, ComposeProject{
			Name:         name,
			ConfigFile:   g.cfgFile,
			WorkingDir:   g.workDir,
			Containers:   g.names,
			RunningCount: g.running,
			TotalCount:   g.total,
			Category:     cat,
			Description:  desc,
			IconHint:     icon,
		})
	}
	sort.Slice(projects, func(i, j int) bool {
		if projects[i].Category != projects[j].Category {
			return projects[i].Category < projects[j].Category
		}
		return projects[i].Name < projects[j].Name
	})
	c.JSON(http.StatusOK, gin.H{
		"projects":  projects,
		"total":     len(projects),
		"available": true,
	})
}

// knownProject 返回已知 compose 项目的友好描述（用于前端图标和分类）。
// 未识别的项目返回 category="other"。
func knownProject(name string) (category, desc, icon string) {
	switch strings.ToLower(name) {
	case "immich":
		return "photo", "自托管照片/视频管理（Google Photos 替代）", "photo"
	case "samba":
		return "file-share", "SMB/CIFS 文件共享（Windows / macOS 访问）", "share"
	case "alist":
		return "file-share", "多网盘聚合文件管理", "files"
	case "forgejo":
		return "dev", "自托管 Git 服务（GitHub 替代）", "git"
	case "restic":
		return "backup", "加密增量备份（云端）", "backup"
	case "shairport-sync", "shairport":
		return "audio", "AirPlay 音频接收器", "audio"
	case "nas-dashboard":
		return "system", "本管理系统", "system"
	}
	return "other", "", "container"
}

// =====================================================================
// 服务目录（用户已知服务的固定视图，比 compose 项目更具体）
// =====================================================================

// ServiceCatalogEntry 一个已知服务（应用商店视图）。
type ServiceCatalogEntry struct {
	ID           string            `json:"id"`           // immich / samba / alist ...
	Name         string            `json:"name"`         // 显示名
	Category     string            `json:"category"`
	Description  string            `json:"description"`
	IconHint     string            `json:"iconHint"`
	URL          string            `json:"url"`          // 推断的访问 URL
	Containers   []AppContainer    `json:"containers"`   // 实际关联的容器
	Status       string            `json:"status"`       // running / partial / stopped / missing
	RunningCount int               `json:"runningCount"`
	TotalCount   int               `json:"totalCount"`
	ComposeDir   string            `json:"composeDir"`
	ComposeFile  string            `json:"composeFile"`
	Notes        string            `json:"notes"`        // 安装 / 集成提示
}

// ListServiceCatalog GET /api/apps/catalog
// 返回用户的"应用目录"，自动检测每个已知服务的实际状态。
func ListServiceCatalog(c *gin.Context) {
	catalog := buildServiceCatalog()
	c.JSON(http.StatusOK, gin.H{"catalog": catalog, "total": len(catalog)})
}

func buildServiceCatalog() []ServiceCatalogEntry {
	// 已知服务清单（按用户当前部署）
	known := []ServiceCatalogEntry{
		{
			ID: "immich", Name: "Immich 照片管理", Category: "photo",
			Description: "Google Photos 替代，自动备份手机照片", IconHint: "photo",
			URL: "http://localhost:2283",
			Notes: "必须 Docker 部署。包含 server / machine-learning / redis / postgres 四个容器。",
		},
		{
			ID: "samba", Name: "Samba 文件共享", Category: "file-share",
			Description: "SMB/CIFS 协议，Windows / macOS / Linux 通用访问", IconHint: "share",
			URL: "",
			Notes: "Docker 部署（dperson/samba）。Go 没有生产级 SMB 服务端库，无法原生集成。",
		},
		{
			ID: "alist", Name: "Alist 网盘聚合", Category: "file-share",
			Description: "多网盘统一挂载与 Web 文件管理", IconHint: "files",
			URL: "http://localhost:5244",
			Notes: "Docker 部署。",
		},
		{
			ID: "forgejo", Name: "Forgejo 私有 Git", Category: "dev",
			Description: "GitHub 替代，自托管代码仓库", IconHint: "git",
			URL: "http://localhost:3000",
			Notes: "Docker 部署。",
		},
		{
			ID: "shairport-sync", Name: "AirPlay 接收器", Category: "audio",
			Description: "把 NAS 当 AirPlay 音箱，iPhone/iPad/Mac 可推送音频", IconHint: "audio",
			Notes: "Docker 部署（需要 /dev/snd）。",
		},
		{
			ID: "restic", Name: "云端加密备份", Category: "backup",
			Description: "增量加密备份到 123pan，使用 rclone 后端", IconHint: "backup",
			Notes: "已集成到本系统的「备份管理」（外部容器 Tab）。",
		},
		{
			ID: "nas-dashboard", Name: "NAS Dashboard（本系统）", Category: "system",
			Description: "你正在使用的管理面板", IconHint: "system",
			URL: "http://localhost:3001",
			Notes: "后端 Go + 前端 Vue；本项目。",
		},
	}

	if !globalDocker.Available() {
		for i := range known {
			known[i].Status = "unknown"
		}
		return known
	}

	all, err := globalDocker.ListContainers(true)
	if err != nil {
		return known
	}
	// 按 ID（也是 compose project 名）匹配
	byProject := map[string][]ContainerSummary{}
	for _, cnt := range all {
		proj := cnt.Labels["com.docker.compose.project"]
		// 兼容：samba 容器名/项目可能就一个；shairport-sync 容器名是 shairport-sync 但 compose project 也是 shairport-sync
		if proj == "" {
			// 通过 name / image 推测
			name := strings.TrimPrefix(cnt.Name, "/")
			if name == "" && len(cnt.Names) > 0 {
				name = strings.TrimPrefix(cnt.Names[0], "/")
			}
			proj = guessProjectByName(name, cnt.Image)
		}
		if proj != "" {
			byProject[proj] = append(byProject[proj], cnt)
		}
	}

	for i, svc := range known {
		// ID 匹配多个可能的 compose 项目名
		candidates := []string{svc.ID}
		if svc.ID == "forgejo" {
			candidates = append(candidates, "forgejo-git")
		}
		if svc.ID == "shairport-sync" {
			candidates = append(candidates, "shairport")
		}
		var found []ContainerSummary
		for _, cand := range candidates {
			if list, ok := byProject[cand]; ok {
				found = list
				break
			}
		}
		// fallback: 通过容器名匹配（兼容不同 compose project 命名）
		if found == nil {
			for _, cnt := range all {
				name := strings.TrimPrefix(cnt.Name, "/")
				if name == "" && len(cnt.Names) > 0 {
					name = strings.TrimPrefix(cnt.Names[0], "/")
				}
				if strings.Contains(strings.ToLower(name), svc.ID) ||
					strings.Contains(strings.ToLower(cnt.Image), svc.ID) {
					found = append(found, cnt)
				}
			}
		}
		if len(found) == 0 {
			known[i].Status = "missing"
			continue
		}
		running := 0
		for _, cnt := range found {
			if cnt.State == "running" {
				running++
			}
			known[i].Containers = append(known[i].Containers, buildAppContainer(cnt, false))
		}
		known[i].TotalCount = len(found)
		known[i].RunningCount = running
		switch {
		case running == len(found):
			known[i].Status = "running"
		case running == 0:
			known[i].Status = "stopped"
		default:
			known[i].Status = "partial"
		}
		// compose dir
		if len(found) > 0 {
			known[i].ComposeDir = found[0].Labels["com.docker.compose.project.working_dir"]
			known[i].ComposeFile = found[0].Labels["com.docker.compose.project.config_files"]
		}
	}
	return known
}

func guessProjectByName(name, image string) string {
	name = strings.ToLower(name)
	image = strings.ToLower(image)
	for _, p := range []string{"immich", "samba", "alist", "forgejo", "restic", "shairport", "nas-backend", "nas-frontend"} {
		if strings.Contains(name, p) || strings.Contains(image, p) {
			if p == "nas-backend" || p == "nas-frontend" {
				return "nas-dashboard"
			}
			return p
		}
	}
	return ""
}

// =====================================================================
// helper
// =====================================================================

// 给 inspect struct 加的方法
func (ci *ContainerInspect) NetworkSettingsBase() (string, bool) {
	if ci == nil {
		return "", false
	}
	// 简化：从 RawInspect 取 IPAddress
	raw, err := globalDocker.RawInspect(ci.ID)
	if err != nil {
		return "", false
	}
	var tmp struct {
		NetworkSettings struct {
			IPAddress string `json:"IPAddress"`
		} `json:"NetworkSettings"`
	}
	if err := jsonUnmarshal(raw, &tmp); err != nil {
		return "", false
	}
	if tmp.NetworkSettings.IPAddress == "" {
		return "", false
	}
	return tmp.NetworkSettings.IPAddress, true
}
