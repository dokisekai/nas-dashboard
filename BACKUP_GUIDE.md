# Restic 备份管理 — 部署与使用指南

本系统提供完整的 Restic 备份管理（仓库、任务、快照、恢复、跨仓库同步），
**全部走 nas-backend 容器自身的 restic 二进制**，不再依赖外部 restic-backup 容器。

## 1. 架构

```
┌───────────────────────────────────────────────────────────┐
│ 宿主机                                                    │
│                                                           │
│   /data  ────────┐                                       │
│   /home/hserver ─┼──┐                                    │
│                  │  │  nas-backend 容器内：              │
│                  │  │   /host/data ←─┐                   │
│                  │  └─ /host/home/hserver                 │
│                  │                  │                    │
│                  │            ┌─────┴──────┐    rclone.conf
│                  │            │  restic    │     (挂载)
│                  │            │  backup    │        │
│                  │            └─────┬──────┘        │
│                  │                  │               │
│                  │                  ▼               ▼
│                  │            ┌─────────────────────────┐
│                  │            │  rclone:123pan:data     │
│                  │            │  rclone:123pan:hserver  │
│                  │            │   （WebDAV → 123 网盘）  │
│                  │            └─────────────────────────┘
│                  │                                         
│   /var/restic-restore ← 恢复目标（容器内 /restore）      
└───────────────────────────────────────────────────────────┘
```

### 关键设计

| 设计点 | 实现 |
|---|---|
| 主机文件访问 | docker-compose 把主机 `/` 挂载到容器 `/host`（只读） |
| 云端访问 | rclone.conf 挂载到容器，restic 通过 `rclone:123pan:data` 直连 |
| 任务执行 | nas-backend 容器内调 `restic backup`，**不依赖任何外部容器** |
| 状态持久化 | SQLite DB `./data/nas-dashboard.db` |
| 配置种子 | `backend/config/backup-seed.json`（项目内可 git 跟踪）|

## 2. 部署

```bash
cd /data/nas-dashboard
docker compose up -d
# 访问 http://localhost:3001
# 默认账号 admin / admin
```

启动后：
1. 后端读取 `config/backup-seed.json`，自动创建其中声明的仓库与任务（已存在的不覆盖）
2. 如需修改默认仓库/任务，编辑 `backend/config/backup-seed.json` 后重启 backend
3. 真实密码用 `${RESTIC_PASSWORD}` 占位符，由 docker-compose 的 env 注入（见下）

### docker-compose.yml 关键挂载

```yaml
backend:
  volumes:
    - /:/host:ro                                            # 主机文件
    - /var/restic-restore:/restore                          # 恢复目标
    - /data/docker/restic/config/rclone.conf:/root/.config/rclone/rclone.conf:ro
    - /var/run/docker.sock:/var/run/docker.sock             # 应用管理用
  environment:
    - RESTIC_PASSWORD=veTjav-8jabty-gupzec                  # 占位符替换源
```

## 3. 五个 Tab 的功能

### 仓库
- 列出所有 restic 仓库（含本地 / S3 / SFTP / Rclone 等）
- 操作：新建（自动 init）、编辑、删除（可选清理磁盘）、**测试连接**、**校验**、**解锁**、刷新统计

### 快照
- 选择仓库 → 看到所有真实快照（含时间、主机、路径、文件数、新增量）
- 单条操作：**详情**、**浏览文件**、**两两对比 diff**、**单文件恢复**、删除
- 多选：批量 forget

### 备份任务
- 任务列表（含状态、上次运行、用时）
- 编辑保留策略（keep_daily / keep_weekly / keep_monthly / keep_yearly / keep_last / keep_within）
- 启用 autoPrune 时，forget 后自动 `restic prune`
- 一键 **立即备份** → 异步执行 + 实时日志

### 仓库同步（restic copy）
- 创建同步任务：源仓库 → 目标仓库
- `restic copy` 增量、加密（只传输目标仓库缺失的快照）
- 典型用途：本地仓库 → 云端镜像，或 123pan → 另一个云盘做异地备份

### 设置
- 全局默认 hostname / 标签 / 排除规则（自动并入所有备份任务）
- Restic 版本、缓存目录、容器主机名等环境信息

## 4. 单文件 / 整目录恢复

1. 仓库 Tab → 选快照 → 点 📁 浏览文件
2. 文件列表每一行都有「恢复」按钮
   - **文件**：精确匹配，只恢复这一个
   - **目录**：恢复该目录下所有内容（自动追加 `/**` glob）
3. 弹出恢复对话框，已自动填好 `include` pattern 和 `target` 路径
4. 点「开始恢复」→ 异步执行 → 日志 Tab 可看进度
5. 完成后文件落在容器内 `/restore/...`，即主机 `/var/restic-restore/...`

### 路径映射
- 备份时：主机 `/data/foo` → 任务里写 `/host/data/foo`
- 恢复时：容器 `/restore/foo` → 主机 `/var/restic-restore/foo`
- restic 会保持快照内的原始路径结构（`/data/...`），所以恢复后实际路径是
  `/var/restic-restore/<你指定的 target>/data/...`

## 5. 种子配置（持久化）

`backend/config/backup-seed.json` 是项目内的版本控制文件，记录了
"这台 NAS 应该有哪些备份仓库与任务"。**首次启动时**会自动应用，
已存在的记录不会被覆盖（用户在 UI 里做的修改不会被还原）。

```json
{
  "repos": [
    {
      "name": "123pan-data",
      "type": "rclone",
      "url": "rclone:123pan:data",
      "password": "${RESTIC_PASSWORD}",   ← docker-compose env 注入
      "env": { "RCLONE_CONFIG": "/root/.config/rclone/rclone.conf" },
      "init": false
    }
  ],
  "tasks": [
    {
      "name": "data-daily",
      "repoName": "123pan-data",
      "sourcePath": "/host/data",
      "retention": { "keep_daily": 7, "keep_weekly": 4 },
      "autoPrune": true
    }
  ]
}
```

## 6. 常见问题

### "repository is already locked exclusively"
中断的 backup / restore 会留下锁。仓库 Tab → 点 **解锁** → 选「强制清除」。
或调 API：
```bash
curl -X POST http://localhost:8888/api/storage/backup/repos/3/unlock?force=true \
  -H "Authorization: Bearer $TOKEN"
```

### 备份很慢
- 首次：会扫描全部数据并上传（27GB 大约 30 分钟）
- 后续：增量，秒级完成（除非大量文件变更）
- 缓存目录 `/data/restic-cache` 持久化后速度更快

### 想看 123pan 上有什么文件
快照 Tab → 选仓库 → 选快照 → 浏览文件，**直接列出云端快照内的所有文件**。

### 跟旧的 `/data/docker/restic/` 容器什么关系？
- **可以共存**：旧容器（restic-backup / restic-home-backup）继续按原方式工作
- **本系统不依赖它们**：所有备份 / 恢复都在 nas-backend 容器内完成
- 应用管理中心的「外部容器」Tab 仍能看到 / 触发旧容器（作为备份方案）

## 7. API 速查

| 操作 | 方法 | 路径 |
|---|---|---|
| 健康检查 | GET | `/api/storage/backup/ping` |
| 全局设置 | GET/PUT | `/api/storage/backup/settings` |
| 仓库 CRUD | GET/POST/PUT/DELETE | `/api/storage/backup/repos` |
| 仓库校验 | POST | `/api/storage/backup/repos/:id/check` |
| 仓库解锁 | POST | `/api/storage/backup/repos/:id/unlock` |
| 快照列表 | GET | `/api/storage/backup/repos/:id/snapshots` |
| 文件浏览 | GET | `/api/storage/backup/repos/:id/snapshots/:sid/ls` |
| 快照对比 | GET | `/api/storage/backup/repos/:id/diff?a=&b=` |
| 恢复 | POST | `/api/storage/backup/repos/:id/restore` |
| 任务 CRUD | GET/POST/PUT/DELETE | `/api/storage/backup/tasks` |
| 立即备份 | POST | `/api/storage/backup/tasks/:id/run` |
| 同步任务 | GET/POST/DELETE | `/api/storage/backup/sync-jobs` |
| 立即同步 | POST | `/api/storage/backup/sync-jobs/:id/run` |
