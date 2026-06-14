<template>
  <div class="disk-manager">
    <div class="manager-header">
      <h2>磁盘管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="showCreateRAID">
          <el-icon><Plus /></el-icon>
          创建RAID
        </el-button>
        <el-button @click="refreshDisks">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- Disk Overview -->
    <div class="disk-overview">
      <div class="overview-cards">
        <div class="card total-disks">
          <div class="card-icon"><el-icon><Files /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ disks.length }}</div>
            <div class="card-label">总磁盘数</div>
          </div>
        </div>
        <div class="card healthy-disks">
          <div class="card-icon"><el-icon><CircleCheck /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ healthyDisks }}</div>
            <div class="card-label">健康磁盘</div>
          </div>
        </div>
        <div class="card raid-arrays">
          <div class="card-icon"><el-icon><Grid /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ raidArrays.length }}</div>
            <div class="card-label">RAID阵列</div>
          </div>
        </div>
        <div class="card lvm-groups">
          <div class="card-icon"><el-icon><Box /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ volumeGroups.length }}</div>
            <div class="card-label">卷组</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabs for different disk management features -->
    <el-tabs v-model="activeTab" class="disk-tabs">
      <!-- Physical Disks -->
      <el-tab-pane label="物理磁盘" name="disks">
        <div class="disks-grid">
          <div
            v-for="disk in disks"
            :key="disk.device"
            class="disk-card"
            :class="getDiskClass(disk)"
            @click="showDiskDetails(disk)"
          >
            <div class="disk-header">
              <div class="disk-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="disk-name">{{ disk.device }}</div>
              <el-tag :type="getHealthType(disk.health)" size="small">
                {{ getHealthLabel(disk.health) }}
              </el-tag>
            </div>

            <div class="disk-body">
              <div class="disk-info">
                <div class="info-row">
                  <span class="label">型号:</span>
                  <span class="value">{{ disk.model }}</span>
                </div>
                <div class="info-row">
                  <span class="label">容量:</span>
                  <span class="value">{{ formatSize(disk.size) }}</span>
                </div>
                <div class="info-row">
                  <span class="label">温度:</span>
                  <span class="value">{{ disk.temperature }}°C</span>
                </div>
                <div class="info-row">
                  <span class="label">分区:</span>
                  <span class="value">{{ disk.partitions?.length || 0 }}</span>
                </div>
              </div>

              <div class="disk-actions">
                <el-button-group>
                  <el-button size="small" @click.stop="partitionDisk(disk)">
                    <el-icon><Edit /></el-icon>
                    分区
                  </el-button>
                  <el-button size="small" @click.stop="benchmarkDisk(disk)">
                    <el-icon><Timer /></el-icon>
                    测试
                  </el-button>
                  <el-button size="small" @click.stop="viewSMART(disk)">
                    <el-icon><Monitor /></el-icon>
                    SMART
                  </el-button>
                </el-button-group>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>

      <!-- RAID Arrays -->
      <el-tab-pane label="RAID阵列" name="raid">
        <div class="raid-section">
          <el-table :data="raidArrays" style="width: 100%">
            <el-table-column prop="name" label="阵列名称" min-width="150" />
            <el-table-column prop="level" label="RAID级别" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ (RAID_LEVEL_LABELS as any)[row.level] || row.level }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getRAIDStatusType(row.status)" size="small">
                  {{ getRAIDStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="容量" width="120">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="devices" label="设备数量" width="100" align="center">
              <template #default="{ row }">
                {{ row.devices?.length || 0 }}
              </template>
            </el-table-column>
            <el-table-column label="使用情况" width="200">
              <template #default="{ row }">
                <div class="usage-bar">
                  <el-progress
                    :percentage="getUsagePercent(row)"
                    :color="getProgressColor(getUsagePercent(row))"
                    :show-text="false"
                  />
                </div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="manageRAID(row)">管理</el-button>
                  <el-button size="small" type="danger" @click="deleteRAID(row)">删除</el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- LVM Management -->
      <el-tab-pane label="LVM管理" name="lvm">
        <div class="lvm-section">
          <!-- Volume Groups -->
          <div class="lvm-section-header">
            <h3>卷组 (Volume Groups)</h3>
            <el-button size="small" type="primary" @click="createVolumeGroup">
              <el-icon><Plus /></el-icon>
              创建卷组
            </el-button>
          </div>

          <el-table :data="volumeGroups" style="width: 100%; margin-bottom: 30px">
            <el-table-column prop="name" label="卷组名称" min-width="150" />
            <el-table-column prop="size" label="总容量" width="120">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="free" label="可用空间" width="120">
              <template #default="{ row }">
                {{ formatSize(row.free) }}
              </template>
            </el-table-column>
            <el-table-column prop="pvCount" label="物理卷数" width="100" align="center" />
            <el-table-column prop="lvCount" label="逻辑卷数" width="100" align="center" />
            <el-table-column label="使用率" width="200">
              <template #default="{ row }">
                <el-progress
                  :percentage="getLVMUsagePercent(row)"
                  :color="getProgressColor(getLVMUsagePercent(row))"
                />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="createLogicalVolume(row)">创建逻辑卷</el-button>
                  <el-button size="small" type="danger" @click="deleteVolumeGroup(row)">删除</el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>

          <!-- Logical Volumes -->
          <div class="lvm-section-header">
            <h3>逻辑卷 (Logical Volumes)</h3>
          </div>

          <el-table :data="allLogicalVolumes" style="width: 100%">
            <el-table-column prop="name" label="逻辑卷名称" min-width="150" />
            <el-table-column prop="vgName" label="所属卷组" width="120" />
            <el-table-column prop="size" label="容量" width="120">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="path" label="设备路径" min-width="200" />
            <el-table-column prop="mountPoint" label="挂载点" width="150">
              <template #default="{ row }">
                {{ row.mountPoint || '未挂载' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                  {{ row.status === 'active' ? '活动' : '非活动' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="mountLogicalVolume(row)">挂载</el-button>
                  <el-button size="small" type="danger" @click="deleteLogicalVolume(row)">删除</el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- SMART Monitoring -->
      <el-tab-pane label="SMART监控" name="smart">
        <SMARTMonitor />
      </el-tab-pane>

      <!-- Performance Benchmark -->
      <el-tab-pane label="性能测试" name="benchmark">
        <BenchmarkTool />
      </el-tab-pane>
    </el-tabs>

    <!-- RAID Creation Wizard Dialog -->
    <RAIDWizard
      v-model:visible="showRAIDWizard"
      @created="onRAIDCreated"
    />

    <!-- Disk Details Dialog -->
    <DiskDetailsDialog
      v-model:visible="showDiskDetailsDialog"
      :disk="selectedDisk"
    />

    <!-- Partition Editor Dialog -->
    <PartitionEditor
      v-model:visible="showPartitionEditor"
      :disk="selectedDisk"
      @updated="onPartitionUpdated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { diskAPI } from '@/api/disk'
import { ElMessage } from 'element-plus'
import type { DiskInfo, RAIDConfig, VolumeGroup, LogicalVolume } from '@/types/disk'
import { RAID_LEVEL_LABELS, RAID_STATUS_COLORS, DISK_HEALTH_COLORS } from '@/types/disk'
import {
  Plus,
  Refresh,
  Files,
  CircleCheck,
  Grid,
  Box,
  Document,
  Edit,
  Timer,
  Monitor
} from '@element-plus/icons-vue'
import RAIDWizard from '@/components/Disk/RAIDWizard.vue'
import DiskDetailsDialog from '@/components/Disk/DiskDetailsDialog.vue'
import PartitionEditor from '@/components/Disk/PartitionEditor.vue'
import SMARTMonitor from '@/components/Disk/SMARTMonitor.vue'
import BenchmarkTool from '@/components/Disk/BenchmarkTool.vue'

const activeTab = ref('disks')
const disks = ref<DiskInfo[]>([])
const raidArrays = ref<RAIDConfig[]>([])
const volumeGroups = ref<VolumeGroup[]>([])
const logicalVolumes = ref<LogicalVolume[]>([])

const selectedDisk = ref<DiskInfo | null>(null)
const showRAIDWizard = ref(false)
const showDiskDetailsDialog = ref(false)
const showPartitionEditor = ref(false)

// Computed
const healthyDisks = computed(() =>
  disks.value.filter(d => d.health === 'good').length
)

const allLogicalVolumes = computed(() =>
  logicalVolumes.value.filter(lv => lv.vgName)
)

// Methods
const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getDiskClass = (disk: DiskInfo): string => {
  return `disk-${disk.health}`
}

const getHealthType = (health: string): string => {
  return (DISK_HEALTH_COLORS as any)[health] || 'info'
}

const getHealthLabel = (health: string): string => {
  const labels: any = { good: '健康', warning: '警告', failed: '故障' }
  return labels[health] || health
}

const getRAIDStatusType = (status: string): string => {
  return (RAID_STATUS_COLORS as any)[status] || 'info'
}

const getRAIDStatusLabel = (status: string): string => {
  const labels: any = { active: '活动', degraded: '降级', failed: '故障', rebuilding: '重建中' }
  return labels[status] || status
}

const getUsagePercent = (raid: RAIDConfig): number => {
  if (!raid.size || raid.size === 0) return 0
  return (raid.used / raid.size) * 100
}

const getLVMUsagePercent = (vg: VolumeGroup): number => {
  if (!vg.size || vg.size === 0) return 0
  return ((vg.size - vg.free) / vg.size) * 100
}

const getProgressColor = (percent: number): string => {
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const refreshDisks = async () => {
  try {
    disks.value = await diskAPI.getDisks()
    raidArrays.value = await diskAPI.getRAIDArrays()
    volumeGroups.value = await diskAPI.getVolumeGroups()
    const lvs = await diskAPI.getLogicalVolumes()
    logicalVolumes.value = lvs
    ElMessage.success('磁盘信息已刷新')
  } catch (error) {
    ElMessage.error('获取磁盘信息失败')
  }
}

const showDiskDetails = (disk: DiskInfo) => {
  selectedDisk.value = disk
  showDiskDetailsDialog.value = true
}

const partitionDisk = (disk: DiskInfo) => {
  selectedDisk.value = disk
  showPartitionEditor.value = true
}

const benchmarkDisk = (disk: DiskInfo) => {
  activeTab.value = 'benchmark'
  // Pass disk to benchmark tool
  ElMessage.info(`即将对 ${disk.device} 进行性能测试`)
}

const viewSMART = (disk: DiskInfo) => {
  activeTab.value = 'smart'
  // Pass disk to SMART monitor
  ElMessage.info(`查看 ${disk.device} 的SMART信息`)
}

const showCreateRAID = () => {
  showRAIDWizard.value = true
}

const manageRAID = (raid: RAIDConfig) => {
  ElMessage.info(`管理RAID阵列 ${raid.name} 的功能即将推出`)
}

const deleteRAID = async (raid: RAIDConfig) => {
  try {
    await diskAPI.deleteRAID(raid.name)
    ElMessage.success(`RAID阵列 ${raid.name} 已删除`)
    await refreshDisks()
  } catch (error) {
    ElMessage.error('删除RAID阵列失败')
  }
}

const createVolumeGroup = () => {
  ElMessage.info('创建卷组功能即将推出')
}

const deleteVolumeGroup = async (vg: VolumeGroup) => {
  try {
    await diskAPI.deleteVolumeGroup(vg.name)
    ElMessage.success(`卷组 ${vg.name} 已删除`)
    await refreshDisks()
  } catch (error) {
    ElMessage.error('删除卷组失败')
  }
}

const createLogicalVolume = (vg: VolumeGroup) => {
  ElMessage.info(`在卷组 ${vg.name} 中创建逻辑卷的功能即将推出`)
}

const mountLogicalVolume = (lv: LogicalVolume) => {
  ElMessage.info(`挂载逻辑卷 ${lv.name} 的功能即将推出`)
}

const deleteLogicalVolume = async (lv: LogicalVolume) => {
  try {
    await diskAPI.deleteLogicalVolume(lv.name, lv.vgName)
    ElMessage.success(`逻辑卷 ${lv.name} 已删除`)
    await refreshDisks()
  } catch (error) {
    ElMessage.error('删除逻辑卷失败')
  }
}

const onRAIDCreated = () => {
  refreshDisks()
}

const onPartitionUpdated = () => {
  refreshDisks()
}

// Lifecycle
onMounted(() => {
  refreshDisks()
})
</script>

<style scoped lang="scss">
.disk-manager {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    color: #303133;
    font-size: 24px;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.disk-overview {
  margin-bottom: 20px;

  .overview-cards {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 15px;

    .card {
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      display: flex;
      align-items: center;
      gap: 15px;

      .card-icon {
        font-size: 32px;
        color: #409eff;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        height: 60px;
        border-radius: 50%;
        background: #ecf5ff;
      }

      .card-content {
        .card-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          line-height: 1;
          margin-bottom: 5px;
        }

        .card-label {
          font-size: 14px;
          color: #909399;
        }
      }

      &.total-disks .card-icon { color: #409eff; background: #ecf5ff; }
      &.healthy-disks .card-icon { color: #67c23a; background: #f0f9ff; }
      &.raid-arrays .card-icon { color: #e6a23c; background: #fdf6ec; }
      &.lvm-groups .card-icon { color: #909399; background: #f4f4f5; }
    }
  }
}

.disk-tabs {
  flex: 1;
  min-height: 0;

  :deep(.el-tabs__content) {
    height: 100%;
  }

  :deep(.el-tab-pane) {
    height: 100%;
    overflow: auto;
  }
}

.disks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 15px;

  .disk-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s;
    border-left: 4px solid #67c23a;

    &:hover {
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
    }

    &.disk-warning { border-left-color: #e6a23c; }
    &.disk-failed { border-left-color: #f56c6c; }

    .disk-header {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 15px;
      background: #f5f7fa;
      border-bottom: 1px solid #ebeef5;

      .disk-icon {
        font-size: 20px;
        color: #409eff;
      }

      .disk-name {
        flex: 1;
        font-weight: 500;
        color: #303133;
      }
    }

    .disk-body {
      padding: 15px;

      .disk-info {
        margin-bottom: 15px;

        .info-row {
          display: flex;
          justify-content: space-between;
          margin-bottom: 8px;
          font-size: 12px;

          .label { color: #909399; }
          .value { color: #303133; font-weight: 500; }
        }
      }

      .disk-actions {
        display: flex;
        gap: 5px;
      }
    }
  }
}

.raid-section, .lvm-section {
  padding: 15px;

  .lvm-section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;

    h3 {
      margin: 0;
      color: #303133;
      font-size: 16px;
    }
  }

  .usage-bar {
    width: 100%;
  }
}
</style>