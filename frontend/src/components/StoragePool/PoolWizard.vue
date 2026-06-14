<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建存储池"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-steps :active="currentStep" align-center class="wizard-steps">
      <el-step title="基本信息" />
      <el-step title="选择磁盘" />
      <el-step title="配置选项" />
      <el-step title="确认创建" />
    </el-steps>

    <div class="wizard-content">
      <!-- Step 1: Basic Information -->
      <div v-show="currentStep === 0" class="wizard-step">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
          <el-form-item label="存储池名称" prop="name">
            <el-input
              v-model="form.name"
              placeholder="请输入存储池名称"
              maxlength="50"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="存储池类型" prop="type">
            <el-select v-model="form.type" placeholder="选择存储池类型" @change="onTypeChange">
              <el-option label="MergerFS" value="mergerfs">
                <div class="option-content">
                  <span class="option-label">MergerFS</span>
                  <span class="option-desc">合并多个文件系统，提供统一视图</span>
                </div>
              </el-option>
              <el-option label="Btrfs" value="btrfs">
                <div class="option-content">
                  <span class="option-label">Btrfs</span>
                  <span class="option-desc">下一代文件系统，支持快照和压缩</span>
                </div>
              </el-option>
              <el-option label="ZFS" value="zfs">
                <div class="option-content">
                  <span class="option-label">ZFS</span>
                  <span class="option-desc">企业级文件系统，数据完整性保护</span>
                </div>
              </el-option>
              <el-option label="LVM" value="lvm">
                <div class="option-content">
                  <span class="option-label">LVM</span>
                  <span class="option-desc">逻辑卷管理器，灵活的存储管理</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

        <el-form-item label="挂载点" prop="mountPoint">
          <el-input
            v-model="form.mountPoint"
            placeholder="/data (推荐) 或 /mnt/pool-name"
            prefix-icon="Folder"
          />
          <div class="form-tip">
            推荐使用 /data 作为统一存储入口
          </div>
        </el-form-item>

          <el-form-item label="描述">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="3"
              placeholder="可选的存储池描述"
            />
          </el-form-item>
        </el-form>

        <!-- Type Description -->
        <div class="type-description">
          <el-alert
            :title="getTypeDescription(form.type).title"
            :type="getTypeDescription(form.type).type"
            :description="getTypeDescription(form.type).description"
            :closable="false"
            show-icon
          />
        </div>
      </div>

      <!-- Step 2: Select Disks -->
      <div v-show="currentStep === 1" class="wizard-step">
        <div class="disk-selection">
          <div class="available-disks">
            <h4>可用磁盘</h4>
            <el-table
              :data="availableDisks"
              @selection-change="handleDiskSelection"
              style="width: 100%"
            >
              <el-table-column type="selection" width="55" />
              <el-table-column prop="device" label="设备" width="150" />
              <el-table-column prop="size" label="容量" width="120">
                <template #default="{ row }">
                  {{ formatSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="model" label="型号" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'available' ? 'success' : 'danger'" size="small">
                    {{ row.status === 'available' ? '可用' : '占用' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="selected-disks">
            <h4>已选择磁盘 ({{ selectedDisks.length }})</h4>
            <el-table :data="selectedDisks" style="width: 100%">
              <el-table-column prop="device" label="设备" />
              <el-table-column prop="size" label="容量">
                <template #default="{ row }">
                  {{ formatSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column label="模式" width="120">
                <template #default="{ row }">
                  <el-select v-model="row.mode" size="small">
                    <el-option label="读写" value="rw" />
                    <el-option label="只读" value="ro" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="角色" width="120">
                <template #default="{ row }">
                  <el-select 
                    v-model="row.role" 
                    size="small"
                    style="width: 100%"
                  >
                    <el-option label="热数据 (SSD)" value="hot" />
                    <el-option label="温数据 (SATA SSD)" value="warm" />
                    <el-option label="冷数据 (HDD)" value="cold" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="优先级" width="150">
                <template #default="{ row }">
                  <el-input-number 
                    v-model="row.priority" 
                    :min="0" 
                    :max="100" 
                    size="small" 
                    controls-position="right"
                    :step="10"
                  />
                </template>
              </el-table-column>
              <el-table-column label="排序" width="100" align="center">
                <template #default="{ $index }">
                  <div class="sort-buttons">
                    <el-button
                      type="primary"
                      size="small"
                      link
                      @click="moveUp($index)"
                      :disabled="$index === 0"
                    >
                      <el-icon><ArrowUp /></el-icon>
                    </el-button>
                    <el-button
                      type="primary"
                      size="small"
                      link
                      @click="moveDown($index)"
                      :disabled="$index === selectedDisks.length - 1"
                    >
                      <el-icon><ArrowDown /></el-icon>
                    </el-button>
                  </div>
                </template>
              </el-table-column>
              <el-table-column width="60">
                <template #default="{ $index }">
                  <el-button
                    type="danger"
                    size="small"
                    @click="removeDisk($index)"
                    link
                  >
                    移除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <div class="total-capacity">
              总容量: <strong>{{ formatSize(totalCapacity) }}</strong>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 3: Configuration -->
      <div v-show="currentStep === 2" class="wizard-step">
        <div v-if="form.type === 'mergerfs'" class="mergerfs-config">
          <h4>MergerFS 配置</h4>
          <el-alert
            title="存储策略说明"
            type="info"
            :closable="false"
            show-icon
            style="margin-bottom: 20px"
          >
            MergerFS 提供多种文件分配策略来控制新文件如何在磁盘间分布。选择合适的策略可以优化性能和空间利用率。
          </el-alert>

          <el-form :model="form.config" label-width="150px">
            <el-form-item label="文件创建策略" required>
              <el-select v-model="form.config.category" placeholder="选择策略">
                <el-option
                  v-for="cat in mergerFSCategories"
                  :key="cat.value"
                  :label="cat.label"
                  :value="cat.value"
                >
                  <div class="option-content">
                    <span class="option-label">{{ cat.label }}</span>
                    <span class="option-desc">{{ cat.description }}</span>
                  </div>
                </el-option>
              </el-select>
              <div class="form-tip">
                {{ getCategoryDescription(form.config.category) }}
              </div>
            </el-form-item>

            <el-form-item label="策略推荐">
              <div class="strategy-recommendations">
                <el-card shadow="hover" class="strategy-card" :class="{ recommended: form.config.category === 'epmfs' }">
                  <template #header>
                    <div class="card-header">
                      <span>推荐策略</span>
                      <el-tag type="success" size="small">推荐</el-tag>
                    </div>
                  </template>
                  <div class="strategy-content">
                    <div class="strategy-item">
                      <strong>epmfs (Most Free Space)</strong>
                      <p>最适合大多数NAS场景，始终写入剩余空间最多的磁盘</p>
                    </div>
                  </div>
                </el-card>

                <el-card shadow="hover" class="strategy-card" :class="{ recommended: form.config.category === 'epff' }">
                  <template #header>
                    <div class="card-header">
                      <span>高性能策略</span>
                      <el-tag type="warning" size="small">SSD优化</el-tag>
                    </div>
                  </template>
                  <div class="strategy-content">
                    <div class="strategy-item">
                      <strong>epff (First Free)</strong>
                      <p>按优先级顺序写入，适合SSD作为主存储的场景</p>
                    </div>
                  </div>
                </el-card>

                <el-card shadow="hover" class="strategy-card" :class="{ recommended: form.config.category === 'epall' }">
                  <template #header>
                    <div class="card-header">
                      <span>负载均衡策略</span>
                      <el-tag type="info" size="small">均衡</el-tag>
                    </div>
                  </template>
                  <div class="strategy-content">
                    <div class="strategy-item">
                      <strong>epall (Round Robin)</strong>
                      <p>轮询方式分配，平衡所有磁盘的I/O负载</p>
                    </div>
                  </div>
                </el-card>
              </div>
            </el-form-item>

            <el-form-item label="最小空闲空间" required>
              <el-input
                v-model="form.config.minfreespace"
                placeholder="例如: 10G, 100M"
              />
              <div class="form-tip">
                分支的最小空闲空间要求，低于此值将不用于写入新文件。建议设置为 10G-50G
              </div>
            </el-form-item>

            <el-divider>高级选项</el-divider>

            <el-form-item label="直接 I/O">
              <el-switch v-model="form.config.direct_io" />
              <div class="form-tip">
                启用直接 I/O 可以提高性能，但可能减少缓存效果。适合大文件传输
              </div>
            </el-form-item>

            <el-form-item label="异步读取">
              <el-switch v-model="form.config.async_read" />
              <div class="form-tip">
                启用异步读取可以提高并发读取性能
              </div>
            </el-form-item>

            <el-form-item label="硬删除">
              <el-switch v-model="form.config.hard_remove" />
              <div class="form-tip">
                立即删除文件而非延迟删除，可以防止空间泄露
              </div>
            </el-form-item>

            <el-form-item label="使用 inode">
              <el-switch v-model="form.config.use_ino" />
              <div class="form-tip">
                使用真实的 inode 值，有助于某些应用正常工作
              </div>
            </el-form-item>

            <el-form-item label="跟随符号链接">
              <el-switch v-model="form.config.follow_symlinks" />
              <div class="form-tip">
                允许符号链接被正确解析
              </div>
            </el-form-item>
          </el-form>
        </div>

        <div v-else class="generic-config">
          <el-alert
            title="高级配置"
            type="info"
            description="此存储池类型使用默认配置，高级选项可在创建后修改"
            :closable="false"
            show-icon
          />
        </div>
      </div>

      <!-- Step 4: Confirmation -->
      <div v-show="currentStep === 3" class="wizard-step">
        <div class="confirmation">
          <h4>确认存储池配置</h4>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="存储池名称">
              {{ form.name }}
            </el-descriptions-item>
            <el-descriptions-item label="存储池类型">
              {{ getPoolTypeLabel(form.type) }}
            </el-descriptions-item>
            <el-descriptions-item label="挂载点">
              {{ form.mountPoint }}
            </el-descriptions-item>
            <el-descriptions-item label="磁盘数量">
              {{ selectedDisks.length }}
            </el-descriptions-item>
            <el-descriptions-item label="总容量">
              {{ formatSize(totalCapacity) }}
            </el-descriptions-item>
            <el-descriptions-item label="描述">
              {{ form.description || '无' }}
            </el-descriptions-item>
          </el-descriptions>

          <div class="disk-list">
            <h5>包含磁盘:</h5>
            <el-table :data="selectedDisks" size="small">
              <el-table-column prop="device" label="设备" />
              <el-table-column prop="size" label="容量">
                <template #default="{ row }">
                  {{ formatSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="mode" label="模式" width="80" />
              <el-table-column prop="priority" label="优先级" width="80" />
            </el-table>
          </div>

          <div class="warning-box">
            <el-alert
              title="重要提示"
              type="warning"
              :closable="false"
              show-icon
            >
              <ul>
                <li>创建存储池将格式化所选磁盘，所有数据将被永久删除</li>
                <li>请确保已备份重要数据</li>
                <li>创建过程可能需要较长时间，请勿关闭浏览器</li>
              </ul>
            </el-alert>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="wizard-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          v-if="currentStep > 0"
          @click="previousStep"
        >
          上一步
        </el-button>
        <el-button
          v-if="currentStep < steps.length - 1"
          type="primary"
          @click="nextStep"
          :disabled="!canNextStep"
        >
          下一步
        </el-button>
        <el-button
          v-else
          type="primary"
          @click="createPool"
          :loading="creating"
        >
          创建存储池
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useStoragePoolStore } from '@/stores/storage_pool'
import { storageApi } from '@/api'
import { ElMessage } from 'element-plus'
import { ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import type { StoragePoolCreateRequest, BranchConfig, MergerFSCategory } from '@/types/storage_pool'
import { MERGERFS_CATEGORIES } from '@/types/storage_pool'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'created', pool: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const storagePoolStore = useStoragePoolStore()
const availableDisks = ref<any[]>([])
const formRef = ref()
const creating = ref(false)
const currentStep = ref(0)

const form = ref({
  name: '',
  type: 'mergerfs',
  mountPoint: '',
  description: '',
  config: {
    branches: [],
    category: 'epmfs',
    minfreespace: '10G',
    direct_io: false,
    async_read: true,
    use_ino: false,
    hard_remove: false,
    auto_unshare: false,
    follow_symlinks: false,
    link_exas: false
  }
})

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// Fetch available disks
const loadAvailableDisks = async () => {
  try {
    const response = await storageApi.getDisks()
    // 过滤掉已挂载的磁盘或者系统磁盘（可选策略）
    // 或者允许用户选择任何磁盘，但在确认步骤给予警告
    availableDisks.value = (response.disks || []).map((disk: any) => ({
      device: disk.name,
      size: disk.size,
      model: disk.label || 'Unknown',
      status: disk.mounted ? 'mounted' : 'available',
      fstype: disk.type
    }))
  } catch (error) {
    console.error('Failed to fetch disks:', error)
    ElMessage.error('获取磁盘列表失败')
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    loadAvailableDisks()
  }
})

const selectedDisks = ref<any[]>([])

const mergerFSCategories = MERGERFS_CATEGORIES

const totalCapacity = computed(() => {
  return selectedDisks.value.reduce((sum, disk) => sum + disk.size, 0)
})

const canNextStep = computed(() => {
  switch (currentStep.value) {
    case 0:
      return form.value.name && form.value.type && form.value.mountPoint
    case 1:
      return selectedDisks.value.length > 0
    case 2:
      return true
    default:
      return false
  }
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getTypeDescription = (type: string) => {
  const descriptions: Record<string, any> = {
    mergerfs: {
      title: 'MergerFS 存储池',
      type: 'success',
      description: 'MergerFS 可以将多个文件系统合并到一个挂载点，提供统一的文件视图。适合需要聚合多个存储设备的场景。'
    },
    btrfs: {
      title: 'Btrfs 存储池',
      type: 'warning',
      description: 'Btrfs 是下一代文件系统，支持快照、压缩、 RAID 和其他高级功能。适合需要高级数据管理的场景。'
    },
    zfs: {
      title: 'ZFS 存储池',
      type: 'info',
      description: 'ZFS 是企业级文件系统，提供数据完整性保护、压缩、快照和存储池功能。适合关键业务数据存储。'
    },
    lvm: {
      title: 'LVM 存储池',
      type: 'info',
      description: 'LVM 提供逻辑卷管理，允许动态调整存储容量和创建快照。适合需要灵活存储管理的场景。'
    }
  }
  return descriptions[type] || descriptions.mergerfs
}

const getPoolTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    mergerfs: 'MergerFS',
    btrfs: 'Btrfs',
    zfs: 'ZFS',
    lvm: 'LVM'
  }
  return labels[type] || type
}

const getCategoryDescription = (category: string): string => {
  const cat = MERGERFS_CATEGORIES.find(c => c.value === category)
  return cat?.description || ''
}

const onTypeChange = () => {
  // Update mount point based on name
  if (form.value.name && !form.value.mountPoint) {
    form.value.mountPoint = `/mnt/${form.value.name}`
  }
}

const handleDiskSelection = (selection: any[]) => {
  selectedDisks.value = selection.map((disk, index) => ({
    ...disk,
    mode: 'rw',
    priority: index === 0 ? 100 : 50,
    role: detectDiskRole(disk)
  }))
}

const detectDiskRole = (disk: any): string => {
  const sizeGB = disk.size / (1024 * 1024 * 1024)
  const isSSD = disk.model?.toLowerCase().includes('ssd') || 
                disk.model?.toLowerCase().includes('nvme') ||
                disk.type === 'ssd'
  
  if (isSSD && sizeGB < 1000) return 'hot'
  if (isSSD && sizeGB >= 1000) return 'warm'
  return 'cold'
}

const moveUp = (index: number) => {
  if (index === 0) return
  const temp = selectedDisks.value[index]
  selectedDisks.value[index] = selectedDisks.value[index - 1]
  selectedDisks.value[index - 1] = temp
}

const moveDown = (index: number) => {
  if (index === selectedDisks.value.length - 1) return
  const temp = selectedDisks.value[index]
  selectedDisks.value[index] = selectedDisks.value[index + 1]
  selectedDisks.value[index + 1] = temp
}

const removeDisk = (index: number) => {
  selectedDisks.value.splice(index, 1)
}

const nextStep = async () => {
  if (currentStep.value === 0) {
    try {
      await formRef.value.validate()
    } catch {
      return
    }
  }
  currentStep.value++
}

const previousStep = () => {
  currentStep.value--
}

const createPool = async () => {
  creating.value = true

  try {
    const poolData: StoragePoolCreateRequest = {
      name: form.value.name,
      type: form.value.type as any,
      mountPoint: form.value.mountPoint,
      description: form.value.description,
      config: form.value.type === 'mergerfs' ? form.value.config as any : undefined,
      disks: selectedDisks.value.map(disk => ({
        device: disk.device,
        size: disk.size,
        priority: disk.priority,
        branchPath: disk.device,
        mode: disk.mode
      })) as any
    }

    const newPool = await storagePoolStore.createPool(poolData)

    ElMessage.success('存储池创建成功')
    emit('created', newPool)
    handleClose()

    // Reset form
    currentStep.value = 0
    form.value = {
      name: '',
      type: 'mergerfs',
      mountPoint: '',
      description: '',
      config: {
        branches: [],
        category: 'epmfs',
        minfreespace: '10G',
        direct_io: false,
        async_read: true,
        use_ino: false,
        hard_remove: false,
        auto_unshare: false,
        follow_symlinks: false,
        link_exas: false
      }
    }
    selectedDisks.value = []
  } catch (error) {
    ElMessage.error('创建存储池失败')
  } finally {
    creating.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

// Watch for name changes to update mount point
watch(() => form.value.name, (newName) => {
  if (newName && currentStep.value === 0) {
    form.value.mountPoint = `/mnt/${newName}`
  }
})
</script>

<style scoped lang="scss">
.wizard-steps {
  margin: 20px 0 30px;
}

.wizard-content {
  min-height: 400px;
  padding: 20px 0;
}

.wizard-step {
  h4 {
    margin-bottom: 20px;
    color: #303133;
  }

  h5 {
    margin: 20px 0 10px;
    color: #606266;
  }
}

.type-description {
  margin-top: 20px;
}

.option-content {
  display: flex;
  flex-direction: column;

  .option-label {
    font-weight: bold;
    color: #303133;
  }

  .option-desc {
    font-size: 12px;
    color: #909399;
    margin-top: 2px;
  }
}

.disk-selection {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;

  .available-disks,
  .selected-disks {
    h4 {
      margin-bottom: 15px;
    }

    .total-capacity {
      margin-top: 15px;
      padding: 10px;
      background: #f5f7fa;
      border-radius: 4px;
      text-align: center;
    }
  }
}

.mergerfs-config {
  .strategy-recommendations {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 15px;
    margin-top: 10px;

    .strategy-card {
      cursor: pointer;
      transition: all 0.3s;

      &.recommended {
        border: 2px solid #67c23a;
      }

      &:hover {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-weight: bold;
      }

      .strategy-content {
        .strategy-item {
          strong {
            display: block;
            margin-bottom: 8px;
            color: #303133;
          }

          p {
            margin: 0;
            font-size: 12px;
            color: #606266;
            line-height: 1.4;
          }
        }
      }
    }
  }
}

.sort-buttons {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.confirmation {
  .disk-list {
    margin: 20px 0;
  }

  .warning-box {
    margin: 20px 0;

    ul {
      margin: 10px 0;
      padding-left: 20px;

      li {
        margin: 5px 0;
      }
    }
  }
}

.wizard-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>