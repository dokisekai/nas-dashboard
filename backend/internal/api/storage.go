package api

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

// Disk 磁盘信息
type Disk struct {
	Name       string `json:"name"`
	Size       uint64 `json:"size"`
	Used       uint64 `json:"used"`
	Available  uint64 `json:"available"`
	Usage      string `json:"usage"`
	Type       string `json:"type"`
	Mounted    bool   `json:"mounted"`
	MountPoint string `json:"mountPoint"`
	Label      string `json:"label"`
	UUID       string `json:"uuid"`
}

// SMBShare SMB 共享
type SMBShare struct {
	Name          string `json:"name"`
	Path          string `json:"path"`
	Comment       string `json:"description"`
	ReadOnly      bool   `json:"readOnly"`
	GuestOK       bool   `json:"guest"`
	Browseable    bool   `json:"browseable"`
	IsTimeMachine bool   `json:"isTimeMachine"`
}


// MountRequest 挂载请求
type MountRequest struct {
	Device     string `json:"device" binding:"required"`
	MountPoint string `json:"mountPoint" binding:"required"`
	Type       string `json:"type"`       // 文件系统类型 (ext4, ntfs, etc.)
	Options    string `json:"options"`     // 挂载选项
}

// UmountRequest 卸载请求
type UmountRequest struct {
	MountPoint string `json:"mountPoint" binding:"required"`
}

// SMBShareRequest SMB 共享请求
type SMBShareRequest struct {
	Name          string `json:"name"`
	Path          string `json:"path" binding:"required"`
	Comment       string `json:"description"`
	ReadOnly      bool   `json:"readOnly"`
	GuestOK       bool   `json:"guest"`
	Browseable    bool   `json:"browseable"`
	IsTimeMachine bool   `json:"isTimeMachine"`
}

// FormatRequest 格式化请求
type FormatRequest struct {
	Device string `json:"device" binding:"required"`
	FSType string `json:"fsType" binding:"required"` // ext4, xfs, etc.
	Label  string `json:"label"`
}

// GetDisks 获取磁盘列表
// 使用 lsblk 和 df 命令获取真实的磁盘信息
func GetDisks(c *gin.Context) {
	disks, err := getDisksFromSystem()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get disk info: %v", err)})
		return
	}
	c.JSON(200, gin.H{"disks": disks})
}

// getDisksFromSystem 从系统获取磁盘信息
func getDisksFromSystem() ([]Disk, error) {
	// 使用 lsblk 获取块设备信息
	// -b: 以字节为单位
	// -o: 指定输出列
	lsblkCmd := exec.Command("lsblk", "-b", "-o", "NAME,SIZE,TYPE,FSTYPE,MOUNTPOINT,UUID,LABEL")
	lsblkOutput, err := lsblkCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk command failed: %w", err)
	}

	// 使用 df 获取挂载点和使用情况
	dfCmd := exec.Command("df", "-B1", "--output=source,size,used,avail,target,pcent")
	dfOutput, err := dfCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("df command failed: %w", err)
	}

	dfMap := parseDfOutput(dfOutput)
	disks := parseLsblkOutput(lsblkOutput, dfMap)

	return disks, nil
}

// parseLsblkOutput 解析 lsblk 输出，只返回整块硬盘信息
func parseLsblkOutput(output []byte, dfMap map[string]map[string]string) []Disk {
	var disks []Disk
	lines := strings.Split(string(output), "\n")

	// 获取根分区所在的设备
	rootDevice := GetRootDevice()

	// 第一遍：收集所有设备信息
	type localDeviceRaw struct {
		name       string
		size       uint64
		diskType   string
		fstype     string
		mountPoint string
		uuid       string
		label      string
	}

	rawDevices := make(map[string]*localDeviceRaw)
	var diskNames []string

	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		name := fields[0]
		// 处理树状结构输出 (├─ 或 └─)
		name = strings.TrimPrefix(name, "├─")
		name = strings.TrimPrefix(name, "└─")
		name = strings.TrimPrefix(name, "│")
		name = strings.TrimSpace(name)

		size := parseSize(fields[1])
		diskType := fields[2]
		fstype := ""
		mountPoint := ""
		uuid := ""
		label := ""

		if len(fields) > 3 {
			fstype = fields[3]
		}
		if len(fields) > 4 {
			mountPoint = fields[4]
		}
		if len(fields) > 5 {
			uuid = fields[5]
		}
		if len(fields) > 6 {
			label = fields[6]
		}

		raw := &localDeviceRaw{
			name:       name,
			size:       size,
			diskType:   diskType,
			fstype:     fstype,
			mountPoint: mountPoint,
			uuid:       uuid,
			label:      label,
		}
		rawDevices[name] = raw
		if diskType == "disk" {
			diskNames = append(diskNames, name)
		}
	}

	// 辅助函数：寻找物理磁盘名称
	findParentDiskLocal := func(name string) string {
		reDisk := regexp.MustCompile(`^(sd[a-z]|nvme[0-9]n[0-9]|vd[a-z])`)
		match := reDisk.FindString(name)
		if match != "" {
			return match
		}
		return name
	}

	// 为每个分区寻找父设备并标记系统盘
	isSystemDisk := make(map[string]bool)
	
	// 查找挂载情况
	type usageInfo struct {
		used  uint64
		avail uint64
		mounted bool
		mountPoint string
	}
	diskUsageInfo := make(map[string]*usageInfo)

	for name, dev := range rawDevices {
		deviceName := "/dev/" + name
		parent := findParentDiskLocal(name)
		
		// 检查系统挂载点
		if dev.mountPoint == "/" || dev.mountPoint == "/boot" || strings.HasPrefix(rootDevice, deviceName) {
			isSystemDisk[parent] = true
		}

		// 记录挂载和使用信息
		if dev.mountPoint != "" && dev.mountPoint != "[SWAP]" {
			info, exists := diskUsageInfo[parent]
			if !exists {
				info = &usageInfo{}
				diskUsageInfo[parent] = info
			}
			
			if df, exists := dfMap[dev.mountPoint]; exists {
				used := parseSize(df["used"])
				avail := parseSize(df["avail"])
				info.used += used
				info.avail += avail
				info.mounted = true
				if info.mountPoint == "" {
					info.mountPoint = dev.mountPoint
				} else if !strings.Contains(info.mountPoint, dev.mountPoint) {
					info.mountPoint += ", " + dev.mountPoint
				}
			}
		}
	}

	// 最终生成磁盘列表
	for _, name := range diskNames {
		dev := rawDevices[name]
		info, hasInfo := diskUsageInfo[name]
		
		disk := Disk{
			Name:       "/dev/" + name,
			Size:       dev.size,
			Type:       dev.diskType,
			Label:      dev.label,
			UUID:       dev.uuid,
		}

		if hasInfo {
			disk.Mounted = info.mounted
			disk.MountPoint = info.mountPoint
			disk.Used = info.used
			disk.Available = info.avail
			if info.used + info.avail > 0 {
				disk.Usage = fmt.Sprintf("%.1f", float64(info.used)/float64(info.used+info.avail)*100)
			}
		}
		
		if isSystemDisk[name] {
			if disk.Label == "" {
				disk.Label = "System"
			} else if !strings.Contains(disk.Label, "System") {
				disk.Label = disk.Label + " (System)"
			}
		}
		
		disks = append(disks, disk)
	}

	return disks
}

// GetRootDevice 获取根分区所在的物理设备
func GetRootDevice() string {
	cmd := exec.Command("sh", "-c", "findmnt -n -o SOURCE /")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// FormatDisk 格式化磁盘
func FormatDisk(c *gin.Context) {
	var req FormatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 安全检查：不允许格式化根分区
	rootDevice := GetRootDevice()
	if req.Device == rootDevice || strings.HasPrefix(rootDevice, req.Device) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot format system disk"})
		return
	}

	// 检查磁盘是否已挂载
	if IsDeviceMounted(req.Device) {
		c.JSON(http.StatusConflict, gin.H{"error": "Disk is currently mounted. Please unmount it first."})
		return
	}

	// 验证文件系统类型
	validFS := map[string]bool{"ext4": true, "xfs": true, "btrfs": true, "ntfs": true, "vfat": true}
	if !validFS[req.FSType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file system type"})
		return
	}

	// 构建格式化命令
	mkfsCmd := "mkfs." + req.FSType
	args := []string{req.Device}
	if req.Label != "" {
		if req.FSType == "xfs" {
			args = append([]string{"-L", req.Label}, args...)
		} else {
			args = append([]string{"-L", req.Label}, args...)
		}
	}

	// 如果是 ext4，添加 -F 强制格式化
	if req.FSType == "ext4" {
		args = append([]string{"-F"}, args...)
	}

	cmd := exec.Command(mkfsCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Format failed: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disk formatted successfully", "output": string(output)})
}

// IsDeviceMounted 检查设备是否已挂载
func IsDeviceMounted(device string) bool {
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 1 && fields[0] == device {
			return true
		}
	}
	return false
}

// parseDfOutput 解析 df 输出
func parseDfOutput(output []byte) map[string]map[string]string {
	result := make(map[string]map[string]string)
	lines := strings.Split(string(output), "\n")

	// 跳过标题行
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		source := fields[0]
		info := map[string]string{
			"size":  fields[1],
			"used":  fields[2],
			"avail": fields[3],
			"target": fields[4],
			"pcent": fields[5],
		}

		result[source] = info
		if fields[4] != "" {
			result[fields[4]] = info
		}
	}

	return result
}

// parseSize 解析大小字符串（字节数）
func parseSize(sizeStr string) uint64 {
	// lsblk -b 输出的是原始字节数
	var size uint64
	_, err := fmt.Sscanf(sizeStr, "%d", &size)
	if err != nil {
		return 0
	}
	return size
}

// MountDisk 挂载磁盘
func MountDisk(c *gin.Context) {
	var req MountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查挂载点目录是否存在，不存在则创建
	if _, err := os.Stat(req.MountPoint); os.IsNotExist(err) {
		if err := os.MkdirAll(req.MountPoint, 0755); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create mount point: %v", err)})
			return
		}
	}

	// 构建挂载命令
	args := []string{req.Device, req.MountPoint}
	if req.Type != "" {
		args = []string{"-t", req.Type, req.Device, req.MountPoint}
	}
	if req.Options != "" {
		args = []string{"-o", req.Options, req.Device, req.MountPoint}
	}

	cmd := exec.Command("mount", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to mount disk: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Disk mounted successfully"})
}

// UmountDisk 卸载磁盘
func UmountDisk(c *gin.Context) {
	var req UmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查挂载点是否存在
	if _, err := os.Stat(req.MountPoint); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Mount point does not exist"})
		return
	}

	// 检查是否真的挂载了
	if !isMounted(req.MountPoint) {
		c.JSON(400, gin.H{"error": "Device not mounted"})
		return
	}

	cmd := exec.Command("umount", req.MountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to unmount disk: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Disk unmounted successfully"})
}

// isMounted 检查设备是否已挂载
func isMounted(mountPoint string) bool {
	// 读取 /proc/mounts
	data, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return false
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == mountPoint {
			return true
		}
	}

	return false
}

// GetSMBShares 获取 SMB 共享列表
func GetSMBShares(c *gin.Context) {
	shares, err := getSMBShares()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get SMB shares: %v", err)})
		return
	}
	c.JSON(200, gin.H{"shares": shares})
}

// getSMBShares 读取 Samba 配置
func getSMBShares() ([]SMBShare, error) {
	// Samba 配置文件路径
	smbConfPath := "/etc/samba/smb.conf"

	// 检查文件是否存在
	if _, err := os.Stat(smbConfPath); os.IsNotExist(err) {
		// 返回空列表而不是错误
		return []SMBShare{}, nil
	}

	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return nil, err
	}

	return parseSmbConf(data)
}

// parseSmbConf 解析 Samba 配置文件
func parseSmbConf(data []byte) ([]SMBShare, error) {
	var shares []SMBShare
	var currentShare *SMBShare

	scanner := bufio.NewScanner(bytes.NewReader(data))
	sectionRegex := regexp.MustCompile(`^\[([^\]]+)\]\s*$`)
	commentRegex := regexp.MustCompile(`^\s*#|^\s*;`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || commentRegex.MatchString(line) {
			continue
		}

		// 检查是否是共享定义的开始
		if matches := sectionRegex.FindStringSubmatch(line); len(matches) > 1 {
			// 保存上一个共享
			if currentShare != nil {
				shares = append(shares, *currentShare)
			}

			// 开始新的共享定义
			shareName := matches[1]
			// 跳过全局配置
			if shareName != "global" {
				currentShare = &SMBShare{Name: shareName}
			} else {
				currentShare = nil
			}
			continue
		}

		// 解析共享参数
		if currentShare != nil {
			if strings.HasPrefix(line, "path") {
				currentShare.Path = strings.TrimPrefix(line, "path")
				currentShare.Path = strings.TrimSpace(strings.TrimPrefix(currentShare.Path, "="))
			} else if strings.HasPrefix(line, "comment") {
				currentShare.Comment = strings.TrimPrefix(line, "comment")
				currentShare.Comment = strings.TrimSpace(strings.TrimPrefix(currentShare.Comment, "="))
			} else if strings.HasPrefix(line, "read only") {
				val := strings.TrimPrefix(line, "read only")
				val = strings.TrimSpace(strings.TrimPrefix(val, "="))
				currentShare.ReadOnly = strings.ToLower(val) == "yes"
			} else if strings.HasPrefix(line, "guest ok") {
				val := strings.TrimPrefix(line, "guest ok")
				val = strings.TrimSpace(strings.TrimPrefix(val, "="))
				currentShare.GuestOK = strings.ToLower(val) == "yes"
			} else if strings.HasPrefix(line, "browseable") {
				val := strings.TrimPrefix(line, "browseable")
				val = strings.TrimSpace(strings.TrimPrefix(val, "="))
				currentShare.Browseable = strings.ToLower(val) == "yes"
			} else if strings.Contains(line, "fruit:time machine") {
				val := strings.Split(line, "=")
				if len(val) > 1 {
					currentShare.IsTimeMachine = strings.TrimSpace(strings.ToLower(val[1])) == "yes"
				}
			}
		}
	}

	// 保存最后一个共享
	if currentShare != nil {
		shares = append(shares, *currentShare)
	}

	return shares, nil
}

// CreateSMBShare 创建 SMB 共享
func CreateSMBShare(c *gin.Context) {
	var req SMBShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径是否存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		c.JSON(400, gin.H{"error": "Path does not exist"})
		return
	}

	// 检查 Samba 是否安装
	if _, err := exec.LookPath("smbd"); err != nil {
		c.JSON(500, gin.H{"error": "Samba is not installed on this system"})
		return
	}

	// 检查共享名是否已存在
	shares, _ := getSMBShares()
	for _, s := range shares {
		if s.Name == req.Name {
			c.JSON(400, gin.H{"error": "Share name already exists"})
			return
		}
	}

	// 添加共享到配置文件
	if err := addSMBShareToConfig(req); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create SMB share: %v", err)})
		return
	}

	// 重启 Samba 服务
	if err := restartSambaService(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to restart Samba service: %v", err)})
		return
	}

	// 如果是 Time Machine，注册 Avahi 服务
	if req.IsTimeMachine {
		if err := registerTimeMachineAvahi(req.Name); err != nil {
			fmt.Printf("Warning: Failed to register Avahi service for Time Machine: %v\n", err)
		}
	}

	c.JSON(201, gin.H{"message": "SMB share created successfully"})
}

// registerTimeMachineAvahi 注册 Avahi mDNS 服务以支持 macOS 发现
func registerTimeMachineAvahi(shareName string) error {
	avahiPath := fmt.Sprintf("/etc/avahi/services/nas-tm-%s.service", shareName)
	
	content := fmt.Sprintf(`<?xml version="1.0" standalone='no'?>
<!DOCTYPE service-group SYSTEM "avahi-service.dtd">
<service-group>
  <name replace-wildcards="yes">%%h - %s</name>
  <service>
    <type>_smb._tcp</type>
    <port>445</port>
  </service>
  <service>
    <type>_device-info._tcp</type>
    <port>0</port>
    <txt-record>model=Macmini</txt-record>
  </service>
  <service>
    <type>_adisk._tcp</type>
    <txt-record>dk0=adVN=%s,adVF=0x82</txt-record>
    <txt-record>sys=waMa=0,adVF=0x100</txt-record>
  </service>
</service-group>
`, shareName, shareName)

	err := os.WriteFile(avahiPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	// 重启 Avahi
	exec.Command("systemctl", "restart", "avahi-daemon").Run()
	return nil
}

func unregisterTimeMachineAvahi(shareName string) {
	avahiPath := fmt.Sprintf("/etc/avahi/services/nas-tm-%s.service", shareName)
	os.Remove(avahiPath)
	exec.Command("systemctl", "restart", "avahi-daemon").Run()
}

// UpdateSMBShare 更新 SMB 共享
func UpdateSMBShare(c *gin.Context) {
	name := c.Param("name")
	var req SMBShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径是否存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		c.JSON(400, gin.H{"error": "Path does not exist"})
		return
	}

	if err := updateSMBShareInConfig(name, req); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update SMB share: %v", err)})
		return
	}

	if err := restartSambaService(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to restart Samba service: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "SMB share updated successfully"})
}

// DeleteSMBShare 删除 SMB 共享
func DeleteSMBShare(c *gin.Context) {
	name := c.Param("name")

	if err := deleteSMBShareFromConfig(name); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete SMB share: %v", err)})
		return
	}

	// 尝试取消注册 Avahi (不论是否是 TM，清理一次无害)
	unregisterTimeMachineAvahi(name)

	if err := restartSambaService(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to restart Samba service: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "SMB share deleted successfully"})
}

func updateSMBShareInConfig(name string, req SMBShareRequest) error {
	smbConfPath := "/etc/samba/smb.conf"
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	inSection := false
	sectionFound := false

	sectionHeader := "[" + name + "]"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if trimmed == sectionHeader {
				inSection = true
				sectionFound = true
				// 写入新的配置
				newLines = append(newLines, sectionHeader)
				newLines = append(newLines, fmt.Sprintf("    path = %s", req.Path))
				if req.Comment != "" {
					newLines = append(newLines, fmt.Sprintf("    comment = %s", req.Comment))
				}
				if req.ReadOnly {
					newLines = append(newLines, "    read only = yes")
				} else {
					newLines = append(newLines, "    read only = no")
				}
				if req.GuestOK {
					newLines = append(newLines, "    guest ok = yes")
				} else {
					newLines = append(newLines, "    guest ok = no")
				}
				if req.Browseable {
					newLines = append(newLines, "    browseable = yes")
				} else {
					newLines = append(newLines, "    browseable = no")
				}
				if req.IsTimeMachine {
					newLines = append(newLines, "    vfs objects = catia fruit streams_xattr")
					newLines = append(newLines, "    fruit:time machine = yes")
					newLines = append(newLines, "    fruit:aapl = yes")
					newLines = append(newLines, "    fruit:metadata = netatalk")
					newLines = append(newLines, "    fruit:resource = file")
					newLines = append(newLines, "    fruit:nfs_aces = no")
				}
				continue
			} else {
				inSection = false
			}
		}

		if !inSection {
			newLines = append(newLines, line)
		}
	}

	if !sectionFound {
		return fmt.Errorf("share not found")
	}

	return os.WriteFile(smbConfPath, []byte(strings.Join(newLines, "\n")), 0644)
}

func deleteSMBShareFromConfig(name string) error {
	smbConfPath := "/etc/samba/smb.conf"
	data, err := os.ReadFile(smbConfPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string
	inSection := false
	sectionHeader := "[" + name + "]"

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if trimmed == sectionHeader {
				inSection = true
				continue
			} else {
				inSection = false
			}
		}

		if !inSection {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(smbConfPath, []byte(strings.Join(newLines, "\n")), 0644)
}

// addSMBShareToConfig 添加共享到 Samba 配置
func addSMBShareToConfig(req SMBShareRequest) error {
	smbConfPath := "/etc/samba/smb.conf"

	// 构建共享配置
	var configBuilder strings.Builder
	configBuilder.WriteString(fmt.Sprintf("\n[%s]\n", req.Name))
	configBuilder.WriteString(fmt.Sprintf("    path = %s\n", req.Path))
	if req.Comment != "" {
		configBuilder.WriteString(fmt.Sprintf("    comment = %s\n", req.Comment))
	}
	if req.ReadOnly {
		configBuilder.WriteString("    read only = yes\n")
	} else {
		configBuilder.WriteString("    read only = no\n")
	}
	if req.GuestOK {
		configBuilder.WriteString("    guest ok = yes\n")
	} else {
		configBuilder.WriteString("    guest ok = no\n")
	}
	if req.Browseable {
		configBuilder.WriteString("    browseable = yes\n")
	} else {
		configBuilder.WriteString("    browseable = no\n")
	}
	if req.IsTimeMachine {
		configBuilder.WriteString("    vfs objects = catia fruit streams_xattr\n")
		configBuilder.WriteString("    fruit:time machine = yes\n")
		configBuilder.WriteString("    fruit:aapl = yes\n")
		configBuilder.WriteString("    fruit:metadata = netatalk\n")
		configBuilder.WriteString("    fruit:resource = file\n")
		configBuilder.WriteString("    fruit:nfs_aces = no\n")
	}

	// 追加到配置文件
	f, err := os.OpenFile(smbConfPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(configBuilder.String())
	return err
}

// restartSambaService 重启 Samba 服务
func restartSambaService() error {
	// 尝试使用 systemctl
	cmd := exec.Command("systemctl", "restart", "smb")
	if err := cmd.Run(); err != nil {
		// 如果 systemctl 失败，尝试 service 命令
		cmd = exec.Command("service", "smb", "restart")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to restart Samba service")
		}
	}
	return nil
}

// GetDiskUsage 获取指定路径的磁盘使用情况
func GetDiskUsage(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/"
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get disk usage: %v", err)})
		return
	}

	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - (stat.Bfree * uint64(stat.Bsize))
	usagePercent := 0.0
	if total > 0 {
		usagePercent = float64(used) / float64(total) * 100
	}

	c.JSON(200, gin.H{
		"path":   path,
		"total":  total,
		"used":   used,
		"avail":  available,
		"usage":  fmt.Sprintf("%.1f%%", usagePercent),
	})
}
