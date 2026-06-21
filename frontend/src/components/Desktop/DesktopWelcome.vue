<template>
  <div class="welcome-message">
    <h1>NAS Dashboard</h1>
    <p>欢迎使用 NAS 管理系统</p>
    <div class="status">
      <span class="status-indicator"></span>
      <span>系统运行正常</span>
    </div>
  </div>
</template>

<style scoped>
.welcome-message {
  position: absolute;
  top: calc(15% + 48px);
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  color: white;
  z-index: 1;
  pointer-events: none;
}

.welcome-message h1 {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 16px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  letter-spacing: 1px;
}

.welcome-message p {
  font-size: 24px;
  opacity: 0.9;
  margin-bottom: 24px;
  font-weight: 300;
}

.status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
}

.status-indicator {
  width: 12px;
  height: 12px;
  background: #10b981;
  border-radius: 50%;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
