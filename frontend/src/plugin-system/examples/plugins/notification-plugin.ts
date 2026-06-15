/**
 * Example: Notification Plugin
 * A complete plugin that adds custom notification features
 */

import type { PluginContext } from '../../types/plugin';

export const manifest = {
  id: 'notification-plugin',
  name: 'Notification Plugin',
  version: '1.0.0',
  description: 'Enhanced notification system with custom sounds and filters',
  author: 'Dashboard Team',
  main: 'notification-plugin.ts',
  permissions: ['notifications', 'storage', 'ui'] as any,
  category: 'productivity',
  keywords: ['notifications', 'alerts', 'sounds']
};

export default function createPlugin(context: PluginContext) {
  const { api, storage, logger, utils } = context;

  // Notification settings
  const defaultSettings = {
    soundEnabled: true,
    soundVolume: 0.5,
    desktopEnabled: true,
    filters: {
      priority: ['high', 'urgent'],
      sources: []
    }
  };

  // Sound effects
  const sounds = {
    default: new Audio('/sounds/notification.mp3'),
    urgent: new Audio('/sounds/urgent.mp3'),
    success: new Audio('/sounds/success.mp3'),
    error: new Audio('/sounds/error.mp3')
  };

  // Notification history
  let notificationHistory: any[] = [];

  return {
    async onInstall() {
      logger.info('Installing notification plugin');
      await storage.set('settings', defaultSettings);
      await storage.set('history', []);
    },

    async onEnable() {
      logger.info('Enabling notification plugin');

      // Register notification panel
      api.ui.registerComponent({
        id: 'notification-panel',
        name: 'Notifications',
        render: (container: HTMLElement) => {
          // Create notification panel UI
          const panel = document.createElement('div');
          panel.className = 'notification-panel';
          panel.innerHTML = `
            <div class="notification-header">
              <h3>Notifications</h3>
              <button id="clear-notifications">Clear All</button>
            </div>
            <div id="notification-list"></div>
          `;

          container.appendChild(panel);

          // Add event listeners
          panel.querySelector('#clear-notifications')?.addEventListener('click', () => {
            notificationHistory = [];
            storage.set('history', []);
            renderNotifications();
          });

          function renderNotifications() {
            const list = panel.querySelector('#notification-list');
            if (list) {
              list.innerHTML = notificationHistory.map(n => `
                <div class="notification-item priority-${n.priority}">
                  <span class="notification-time">${new Date(n.timestamp).toLocaleTimeString()}</span>
                  <span class="notification-message">${n.message}</span>
                </div>
              `).join('');
            }
          }

          renderNotifications();

          return () => {
            container.removeChild(panel);
          };
        }
      });
    },

    async onActivate() {
      logger.info('Activating notification plugin');
      notificationHistory = await storage.get('history') || [];
    },

    async onDeactivate() {
      logger.info('Deactivating notification plugin');
    },

    async onDisable() {
      logger.info('Disabling notification plugin');
      api.ui.unregisterComponent('notification-panel');
    },

    async onUninstall() {
      logger.info('Uninstalling notification plugin');
      await storage.clear();
    },

    // Notification methods
    async showNotification(message: string, options: NotificationOptions = {}) {
      const settings = await storage.get('settings') || defaultSettings;

      // Check filters
      if (options.priority && !settings.filters.priority.includes(options.priority)) {
        logger.debug('Notification filtered by priority:', options.priority);
        return;
      }

      // Add to history
      const notification = {
        id: Date.now(),
        message,
        timestamp: new Date().toISOString(),
        priority: options.priority || 'normal',
        type: options.type || 'info'
      };

      notificationHistory.unshift(notification);
      if (notificationHistory.length > 100) {
        notificationHistory.pop();
      }
      await storage.set('history', notificationHistory);

      // Play sound
      if (settings.soundEnabled) {
        const sound = options.priority === 'urgent' ? sounds.urgent :
                      options.type === 'success' ? sounds.success :
                      options.type === 'error' ? sounds.error :
                      sounds.default;

        sound.volume = settings.soundVolume;
        sound.play().catch(e => logger.error('Failed to play sound:', e));
      }

      // Show desktop notification
      if (settings.desktopEnabled && 'Notification' in window) {
        if (Notification.permission === 'granted') {
          new Notification(options.title || 'Notification', {
            body: message,
            icon: options.icon
          });
        } else if (Notification.permission !== 'denied') {
          Notification.requestPermission();
        }
      }

      // Show in-app notification
      api.ui.showNotification(message, options.type || 'info');

      logger.info('Notification shown:', message);
    },

    // Settings management
    async getSettings() {
      return await storage.get('settings') || defaultSettings;
    },

    async updateSettings(newSettings: any) {
      const current = await this.getSettings();
      const merged = utils.deepMerge(current, newSettings);
      await storage.set('settings', merged);
      logger.info('Settings updated:', merged);
    },

    // History management
    async getHistory() {
      return notificationHistory;
    },

    async clearHistory() {
      notificationHistory = [];
      await storage.set('history', []);
    },

    // Filter management
    async addFilter(type: 'priority' | 'source', value: string) {
      const settings = await this.getSettings();
      if (!settings.filters[type].includes(value)) {
        settings.filters[type].push(value);
        await storage.set('settings', settings);
      }
    },

    async removeFilter(type: 'priority' | 'source', value: string) {
      const settings = await this.getSettings();
      const index = settings.filters[type].indexOf(value);
      if (index > -1) {
        settings.filters[type].splice(index, 1);
        await storage.set('settings', settings);
      }
    }
  };
}

// Types
interface NotificationOptions {
  title?: string;
  message?: string;
  type?: 'info' | 'success' | 'warning' | 'error';
  priority?: 'low' | 'normal' | 'high' | 'urgent';
  icon?: string;
  sound?: string;
}
