#!/bin/bash
# Quick fix for remaining TypeScript errors

echo "Fixing TypeScript compilation errors..."

# Fix chart legend position type
sed -i "s/position: 'top'/position: 'top' as const/g" /home/hserver/nas-dashboard/frontend/src/apps/SystemMonitor.vue

# Fix unused imports
sed -i "/ChartBarIcon/d" /home/hserver/nas-dashboard/frontend/src/apps/SystemMonitor.vue
sed -i "/monitorApi/d" /home/hserver/nas-dashboard/frontend/src/apps/SystemMonitor.vue

# Fix UserManager type issue
sed -i "s/:selectedUsers\.username/:selectedUsers === true ? '' : selectedUsers.username/g" /home/hserver/nas-dashboard/frontend/src/apps/UserManager.vue

# Fix Chart.vue type issue
sed -i "s/type?: string;/type?: 'area' | 'line' | 'bar' | 'donut';/g" /home/hserver/nas-dashboard/frontend/src/components/Chart.vue

echo "TypeScript errors fixed"
