#!/bin/bash
# 导入环境变量 (Restic密码等)
source /etc/restic/backup.env

# 执行备份 (使用统一的忽略列表)
restic -r rclone:123pan:data backup /data \
    --exclude-file="/data/.resticignore" \
    --verbose
    
# 清理旧快照 (保留最近7天，和过去4周的每周快照)
restic -r rclone:123pan:data forget --keep-daily 7 --keep-weekly 4 --prune
