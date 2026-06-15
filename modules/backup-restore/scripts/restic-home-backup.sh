#!/bin/bash
# 导入环境变量 (Restic密码等)
source /etc/restic/backup.env

# 执行用户目录备份 (使用统一的忽略列表)
# 注意：已移除对 Downloads 的排除，Downloads 将被同步
restic -r rclone:123pan:hserver backup /home/hserver \
    --exclude-file="/data/.resticignore" \
    --verbose
    
# 清理旧快照
restic -r rclone:123pan:hserver forget --keep-daily 7 --keep-weekly 4 --prune
