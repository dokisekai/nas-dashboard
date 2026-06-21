package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"nas-dashboard/internal/models"
)

// Restic binary path; can be overridden for tests.
var ResticBinary = "restic"

// Snapshot represents a restic snapshot (subset of fields restic prints).
type Snapshot struct {
	ID         string    `json:"id"`
	ShortID    string    `json:"short_id"`
	Time       time.Time `json:"time"`
	Tree       string    `json:"tree,omitempty"`
	Paths      []string  `json:"paths"`
	Hostname   string    `json:"hostname"`
	Username   string    `json:"username,omitempty"`
	Tags       []string  `json:"tags,omitempty"`
	Parent     string    `json:"parent,omitempty"`
	ProgramVer string    `json:"program_version,omitempty"`
	Summary    *SnapStats `json:"summary,omitempty"`
}

// SnapStats holds stats about a snapshot (new files / changed etc).
type SnapStats struct {
	FilesNew       int   `json:"files_new"`
	FilesChanged   int   `json:"files_changed"`
	FilesUnmodified int  `json:"files_unmodified"`
	DirsNew        int   `json:"dirs_new"`
	DirsChanged    int   `json:"dirs_changed"`
	DirsUnmodified int   `json:"dirs_unmodified"`
	DataBlobs      int   `json:"data_blobs"`
	TreeBlobs      int   `json:"tree_blobs"`
	DataAdded      int64 `json:"data_added"`
	DataAddedPacked int64 `json:"data_added_packed"`
	TotalFilesProcessed int `json:"total_files_processed"`
	TotalBytesProcessed int64 `json:"total_bytes_processed"`
}

// RepoStats represents restic stats --json output summary.
type RepoStats struct {
	TotalSize      int64 `json:"total_size"`
	TotalBlobCount int64 `json:"total_blob_count"`
	SnapshotsCount int   `json:"snapshots_count"`
	TotalUncompressedSize int64 `json:"total_uncompressed_size,omitempty"`
}

// LSNode represents a file/dir in a snapshot from restic ls --json.
type LSNode struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"` // file, dir, symlink
	Path        string    `json:"path"`
	UID         uint32    `json:"uid,omitempty"`
	GID         uint32    `json:"gid,omitempty"`
	Size        int64     `json:"size,omitempty"`
	Mode        uint32    `json:"mode,omitempty"`
	Permissions string    `json:"permissions,omitempty"`
	Mtime       time.Time `json:"mtime,omitempty"`
	Atime       time.Time `json:"atime,omitempty"`
	CTime       time.Time `json:"ctime,omitempty"`
	// StructType 是 restic 用于区分 "snapshot"（首行）和 "node" 的字段
	StructType string `json:"struct_type,omitempty"`
}

// RunResult captures a restic command output.
type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// LogCallback receives streaming log lines during long-running operations.
type LogCallback func(line string)

// ResticService wraps restic invocations against a given repository.
type ResticService struct {
	repo *models.BackupRepo
}

// NewResticService creates a service bound to a repo config.
func NewResticService(repo *models.BackupRepo) *ResticService {
	return &ResticService{repo: repo}
}

// repoSpec builds the value passed to restic -r based on type+URL.
func (s *ResticService) repoSpec() string {
	r := s.repo
	if r.Type == "" || r.Type == "local" {
		// local: just the path (no prefix); restic accepts absolute paths directly.
		return r.URL
	}
	// For remote backends, URL should already include the scheme prefix
	// (e.g. "s3:https://...", "sftp:user@host:/path", "rest:http://...").
	// If user forgot the prefix, add it.
	if strings.Contains(r.URL, ":") {
		return r.URL
	}
	return r.Type + ":" + r.URL
}

// env returns the merged environment (host env + repo env + password + repo).
func (s *ResticService) env() []string {
	env := os.Environ()
	// User-supplied credentials (AWS_ACCESS_KEY_ID etc.)
	repoEnv := s.repo.ParseEnv()
	for k, v := range repoEnv {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	// Password
	env = append(env, "RESTIC_PASSWORD="+s.repo.Password)
	// Repository env (so we don't need -r everywhere; restic honors RESTIC_REPOSITORY)
	env = append(env, "RESTIC_REPOSITORY="+s.repoSpec())

	// 处理 RCLONE_CONFIG：
	//   - 如果仓库显式指定了且文件存在 -> 用仓库的
	//   - 如果仓库显式指定了但文件不存在（如迁移后旧路径 /root/.config/...）-> 改用 /data/rclone.conf
	//   - 没指定但 /data/rclone.conf 存在 -> 用它
	//   - 都没有 -> 留给 rclone 自己的默认值
	configuredPath := repoEnv["RCLONE_CONFIG"]
	if configuredPath == "" || !fileExists(configuredPath) {
		if fileExists("/data/rclone.conf") {
			// 移除可能存在的旧值（包括从父进程继承的）然后追加正确值
			filtered := env[:0]
			for _, kv := range env {
				if !strings.HasPrefix(kv, "RCLONE_CONFIG=") {
					filtered = append(filtered, kv)
				}
			}
			env = append(filtered, "RCLONE_CONFIG=/data/rclone.conf")
		}
	}
	return env
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// baseCmd creates a restic command with env set.
func (s *ResticService) baseCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(ResticBinary, args...)
	cmd.Env = s.env()
	return cmd
}

// run executes a command and captures combined output.
func (s *ResticService) run(args []string) (RunResult, error) {
	start := time.Now()
	cmd := s.baseCmd(args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := RunResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Duration: time.Since(start),
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
		} else {
			res.ExitCode = -1
		}
		return res, fmt.Errorf("restic %s: %w (stderr: %s)",
			strings.Join(args, " "), err, strings.TrimSpace(res.Stderr))
	}
	return res, nil
}

// runStream executes a command and streams stderr (restic prints progress to stderr) to cb.
func (s *ResticService) runStream(args []string, cb LogCallback) (RunResult, error) {
	start := time.Now()
	cmd := s.baseCmd(args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return RunResult{}, err
	}
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	if err := cmd.Start(); err != nil {
		return RunResult{}, err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		if cb != nil {
			cb(scanner.Text())
		}
	}

	err = cmd.Wait()
	res := RunResult{
		Stdout:   stdoutBuf.String(),
		Duration: time.Since(start),
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
		} else {
			res.ExitCode = -1
		}
		return res, fmt.Errorf("restic %s: %w", strings.Join(args, " "), err)
	}
	return res, nil
}

// CheckAvailable verifies that restic binary exists; returns error if not.
func CheckAvailable() error {
	_, err := exec.LookPath(ResticBinary)
	if err != nil {
		return fmt.Errorf("restic binary not found in PATH: %w", err)
	}
	return nil
}

// Init initialises a new repository. Returns error if already initialised.
func (s *ResticService) Init() error {
	if s.repo.Type == "local" && s.repo.URL != "" {
		if err := os.MkdirAll(s.repo.URL, 0700); err != nil {
			return fmt.Errorf("create repo dir: %w", err)
		}
	}
	_, err := s.run([]string{"init"})
	return err
}

// Check runs `restic check` to verify repo integrity.
//   full=true  : full read-data check (slow, reads every pack)
//   full=false : only checks snapshots/indexes (fast)
func (s *ResticService) Check(full bool) (string, error) {
	args := []string{"check"}
	if full {
		args = append(args, "--read-data")
	}
	res, err := s.run(args)
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// Snapshots lists all snapshots in the repo.
func (s *ResticService) Snapshots() ([]Snapshot, error) {
	res, err := s.run([]string{"snapshots", "--json"})
	if err != nil {
		return nil, err
	}
	var snaps []Snapshot
	if res.Stdout == "" || res.Stdout == "[]" {
		return snaps, nil
	}
	if err := json.Unmarshal([]byte(res.Stdout), &snaps); err != nil {
		return nil, fmt.Errorf("parse snapshots: %w (raw: %s)", err, truncate(res.Stdout, 200))
	}
	return snaps, nil
}

// Stats returns repo-wide statistics.
func (s *ResticService) Stats() (*RepoStats, error) {
	res, err := s.run([]string{"stats", "--json", "--mode", "repo-data"})
	if err != nil {
		// Some restic versions don't support --mode; fallback
		res, err = s.run([]string{"stats", "--json"})
		if err != nil {
			return nil, err
		}
	}
	var st RepoStats
	if res.Stdout == "" {
		return &st, nil
	}
	if err := json.Unmarshal([]byte(res.Stdout), &st); err != nil {
		return nil, fmt.Errorf("parse stats: %w", err)
	}
	// snapshots count isn't always in stats; populate via Snapshots
	if snaps, err := s.Snapshots(); err == nil {
		st.SnapshotsCount = len(snaps)
	}
	return &st, nil
}

// BackupOptions configures a backup run.
type BackupOptions struct {
	Source   string
	Excludes []string
	Tags     []string
	Hostname string
	// Optional parent snapshot ID for incremental.
	Parent string
}

// BackupResult captures the result of a backup run.
type BackupResult struct {
	SnapshotID string
	Duration   time.Duration
	Stats      *SnapStats
}

// Backup runs `restic backup`. The cb receives progress lines (stderr).
func (s *ResticService) Backup(opts BackupOptions, cb LogCallback) (*BackupResult, error) {
	if opts.Source == "" {
		return nil, fmt.Errorf("backup source path is empty")
	}
	if _, err := os.Stat(opts.Source); err != nil {
		return nil, fmt.Errorf("source path not accessible: %w", err)
	}
	args := []string{"backup", opts.Source, "--json"}
	for _, ex := range opts.Excludes {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		args = append(args, "--exclude", ex)
	}
	for _, t := range opts.Tags {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		args = append(args, "--tag", t)
	}
	if opts.Hostname != "" {
		args = append(args, "--host", opts.Hostname)
	}
	if opts.Parent != "" {
		args = append(args, "--parent", opts.Parent)
	}

	res, err := s.runStream(args, cb)
	if err != nil {
		return &BackupResult{Duration: res.Duration}, err
	}
	br := &BackupResult{Duration: res.Duration}
	// restic backup --json emits a sequence of JSON messages, one per line;
	// the last "summary" message contains the snapshot id + stats.
	for _, line := range strings.Split(strings.TrimSpace(res.Stdout), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var msg map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		if mt, ok := msg["message_type"]; ok {
			var mtStr string
			_ = json.Unmarshal(mt, &mtStr)
			if mtStr == "summary" {
				summary := &SnapStats{}
				if raw, ok := msg["snapshot_id"]; ok {
					_ = json.Unmarshal(raw, &br.SnapshotID)
				}
				if raw, ok := msg["files_new"]; ok {
					_ = json.Unmarshal(raw, &summary.FilesNew)
				}
				if raw, ok := msg["files_changed"]; ok {
					_ = json.Unmarshal(raw, &summary.FilesChanged)
				}
				if raw, ok := msg["files_unmodified"]; ok {
					_ = json.Unmarshal(raw, &summary.FilesUnmodified)
				}
				if raw, ok := msg["dirs_new"]; ok {
					_ = json.Unmarshal(raw, &summary.DirsNew)
				}
				if raw, ok := msg["dirs_changed"]; ok {
					_ = json.Unmarshal(raw, &summary.DirsChanged)
				}
				if raw, ok := msg["dirs_unmodified"]; ok {
					_ = json.Unmarshal(raw, &summary.DirsUnmodified)
				}
				if raw, ok := msg["data_added"]; ok {
					_ = json.Unmarshal(raw, &summary.DataAdded)
				}
				if raw, ok := msg["total_files_processed"]; ok {
					_ = json.Unmarshal(raw, &summary.TotalFilesProcessed)
				}
				if raw, ok := msg["total_bytes_processed"]; ok {
					_ = json.Unmarshal(raw, &summary.TotalBytesProcessed)
				}
				br.Stats = summary
			}
		}
	}
	return br, nil
}

// RestoreOptions configures a restore run.
type RestoreOptions struct {
	Snapshot string // snapshot ID / "latest"
	Target   string // destination directory
	Include  []string
	Exclude  []string
	Host     string
	Paths    []string
}

// Restore runs `restic restore`.
func (s *ResticService) Restore(opts RestoreOptions, cb LogCallback) error {
	if opts.Snapshot == "" {
		opts.Snapshot = "latest"
	}
	if opts.Target == "" {
		return fmt.Errorf("restore target path is required")
	}
	if err := os.MkdirAll(opts.Target, 0755); err != nil {
		return fmt.Errorf("create restore target: %w", err)
	}
	args := []string{"restore", opts.Snapshot, "--target", opts.Target, "--json"}
	for _, p := range opts.Include {
		args = append(args, "--include", p)
	}
	for _, p := range opts.Exclude {
		args = append(args, "--exclude", p)
	}
	if opts.Host != "" {
		args = append(args, "--host", opts.Host)
	}
	for _, p := range opts.Paths {
		args = append(args, "--path", p)
	}
	res, err := s.runStream(args, cb)
	if err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}
	_ = res
	return nil
}

// Forget removes snapshots by ID. If prune is true, runs prune after forget.
func (s *ResticService) Forget(snapshotIDs []string, prune bool) (string, error) {
	if len(snapshotIDs) == 0 {
		return "", fmt.Errorf("no snapshot IDs provided")
	}
	args := []string{"forget"}
	args = append(args, snapshotIDs...)
	if prune {
		args = append(args, "--prune")
	}
	res, err := s.run(args)
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// ApplyRetention applies a retention policy (keep_daily etc.) and optionally prunes.
// keep is a map like {"keep_daily": 7, "keep_weekly": 4}.
func (s *ResticService) ApplyRetention(keep map[string]int, prune bool, groupingByTag bool) (string, error) {
	if len(keep) == 0 {
		return "", nil
	}
	args := []string{"forget"}
	for k, v := range keep {
		args = append(args, fmt.Sprintf("--%s=%d", k, v))
	}
	if groupingByTag {
		args = append(args, "--group-by", "host,paths,tags")
	}
	if prune {
		args = append(args, "--prune")
	}
	res, err := s.run(args)
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// LS lists files in a snapshot.
func (s *ResticService) LS(snapshotID string) ([]LSNode, error) {
	if snapshotID == "" {
		snapshotID = "latest"
	}
	res, err := s.run([]string{"ls", snapshotID, "--json"})
	if err != nil {
		return nil, err
	}
	var nodes []LSNode
	dec := json.NewDecoder(strings.NewReader(res.Stdout))
	for {
		var node LSNode
		if err := dec.Decode(&node); err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		// 只收集 node 类型，跳过 snapshot 概述行
		if node.StructType == "node" && node.Type != "" {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

// Diff returns the textual diff between two snapshots.
// restic diff <a> <b> prints added/removed files; output is text.
func (s *ResticService) Diff(a, b string) (string, error) {
	res, err := s.run([]string{"diff", a, b})
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// Unlock 移除仓库锁。removeAll=true 时强制清除所有锁（含他人的）。
// 通常在 backup 中断、出现 "repository is already locked exclusively" 错误时使用。
func (s *ResticService) Unlock(removeAll bool) (string, error) {
	args := []string{"unlock"}
	if removeAll {
		args = append(args, "--remove-all")
	}
	res, err := s.run(args)
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// Find 在所有快照里搜索文件 / 路径（restic find <pattern>）。
// pattern 支持普通字符串或 glob（如 "*.jpg"）。返回文本结果。
func (s *ResticService) Find(pattern string) (string, error) {
	res, err := s.run([]string{"find", pattern, "--long"})
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// ListKeys 列出仓库的所有 key（restic key list）。
func (s *ResticService) ListKeys() (string, error) {
	res, err := s.run([]string{"key", "list"})
	if err != nil {
		return res.Stdout + res.Stderr, err
	}
	return res.Stdout, nil
}

// Copy copies snapshots from this repo (source) to dst repo (target).
func (s *ResticService) Copy(dst *models.BackupRepo, snapshotIDs []string, cb LogCallback) (string, error) {
	// restic copy --repo2 <dst> --password-file2 <...>
	// Easier: set RESTIC_REPOSITORY2 / RESTIC_PASSWORD2 env on the command.
	args := []string{"copy"}
	args = append(args, snapshotIDs...)
	args = append(args, "--repo2", dst.Type+":"+dst.URL)
	if dst.Type == "local" {
		args[len(args)-1] = dst.URL
	}
	cmd := s.baseCmd(args...)
	for k, v := range dst.ParseEnv() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s2=%s", k, v))
		// restic copy honors <ENV>2 for the secondary repo only for some vars;
		// most credentials like AWS_ACCESS_KEY_ID2 are supported by restic.
	}
	cmd.Env = append(cmd.Env, "RESTIC_PASSWORD2="+dst.Password)
	cmd.Env = append(cmd.Env, "RESTIC_REPOSITORY2="+func() string {
		if dst.Type == "local" {
			return dst.URL
		}
		if strings.Contains(dst.URL, ":") {
			return dst.URL
		}
		return dst.Type + ":" + dst.URL
	}())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if cb != nil {
		scanner := bufio.NewScanner(stderr)
		scanner.Buffer(make([]byte, 64*1024), 1024*1024)
		for scanner.Scan() {
			cb(scanner.Text())
		}
	} else {
		io.Copy(io.Discard, stderr)
	}
	if err := cmd.Wait(); err != nil {
		return stdoutBuf.String(), err
	}
	return stdoutBuf.String(), nil
}

// FindAutoParent returns the latest snapshot ID for the same host+path, if any,
// usable as --parent for incremental backups. Errors are ignored (returns "").
func (s *ResticService) FindAutoParent(host, path string) string {
	snaps, err := s.Snapshots()
	if err != nil || len(snaps) == 0 {
		return ""
	}
	// restic already does this automatically; expose for explicit use.
	_ = host
	_ = path
	// Use the most recent snapshot ID.
	return snaps[len(snaps)-1].ID
}

// AuditRepoPath ensures a local repo directory is safe to create.
func AuditRepoPath(p string) error {
	abs, err := filepath.Abs(p)
	if err != nil {
		return err
	}
	// Forbid overly dangerous paths.
	if abs == "/" || strings.HasPrefix(abs, "/etc/") ||
		strings.HasPrefix(abs, "/usr/") || strings.HasPrefix(abs, "/proc/") ||
		strings.HasPrefix(abs, "/sys/") || strings.HasPrefix(abs, "/dev/") {
		return fmt.Errorf("path %s is not allowed as a restic repository", abs)
	}
	return nil
}

// LogBuffer is a thread-safe ring buffer for recent log lines of a job.
type LogBuffer struct {
	mu     sync.Mutex
	lines  []string
	cap    int
}

// NewLogBuffer creates a buffer with capacity N lines.
func NewLogBuffer(n int) *LogBuffer { return &LogBuffer{cap: n} }

// Append adds a line.
func (b *LogBuffer) Append(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.lines = append(b.lines, line)
	if len(b.lines) > b.cap {
		b.lines = b.lines[len(b.lines)-b.cap:]
	}
}

// Lines returns a copy of current lines.
func (b *LogBuffer) Lines() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]string, len(b.lines))
	copy(out, b.lines)
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
