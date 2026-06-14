// Sync & Backup API Client
import api from './client'

export const syncApi = {
  getJobs: () => api.get('/api/storage/sync/jobs'),
  createJob: (data: any) => api.post('/api/storage/sync/jobs', data),
  runJob: (id: number) => api.post(`/api/storage/sync/jobs/${id}/run`),
}

export const backupApi = {
  getRepos: () => api.get('/api/storage/backup/repos'),
  createRepo: (data: any) => api.post('/api/storage/backup/repos', data),
  getTasks: () => api.get('/api/storage/backup/tasks'),
  createTask: (data: any) => api.post('/api/storage/backup/tasks', data),
}
