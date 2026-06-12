/**
 * Plugin API - Provides controlled access to application features
 */

import type { PluginAPI, PluginManifest } from '../types/plugin';
import { useRouter } from 'vue-router';
import { useAppStore } from '../../stores/app';

export function createPluginAPI(manifest: PluginManifest): PluginAPI {
  const router = useRouter();
  const appStore = useAppStore();

  // Check permissions
  const hasPermission = (permission: string): boolean => {
    return manifest.permissions?.includes(permission as any) || false;
  };

  return {
    app: {
      navigate: (path: string) => {
        if (!hasPermission('ui')) {
          throw new Error('Missing UI permission');
        }
        router.push(path);
      },

      getState: () => {
        return appStore.$state;
      },

      setState: (key: string, value: any) => {
        if (!hasPermission('settings')) {
          throw new Error('Missing settings permission');
        }
        (appStore as any)[key] = value;
      }
    },

    ui: {
      registerComponent: (component: any) => {
        if (!hasPermission('ui')) {
          throw new Error('Missing UI permission');
        }
        // Component registration logic
        console.log('Component registered:', component);
      },

      unregisterComponent: (componentId: string) => {
        if (!hasPermission('ui')) {
          throw new Error('Missing UI permission');
        }
        // Component unregistration logic
        console.log('Component unregistered:', componentId);
      },

      showNotification: (message: string, type = 'info') => {
        if (!hasPermission('notifications')) {
          throw new Error('Missing notifications permission');
        }
        // Notification logic
        console.log(`[${type.toUpperCase()}] ${message}`);
      }
    },

    network: {
      request: async (config: any) => {
        if (!hasPermission('network')) {
          throw new Error('Missing network permission');
        }
        const axios = (await import('axios')).default;
        return axios(config);
      },

      get: async (url: string, config?: any) => {
        if (!hasPermission('network')) {
          throw new Error('Missing network permission');
        }
        const axios = (await import('axios')).default;
        return axios.get(url, config);
      },

      post: async (url: string, data?: any, config?: any) => {
        if (!hasPermission('network')) {
          throw new Error('Missing network permission');
        }
        const axios = (await import('axios')).default;
        return axios.post(url, data, config);
      },

      put: async (url: string, data?: any, config?: any) => {
        if (!hasPermission('network')) {
          throw new Error('Missing network permission');
        }
        const axios = (await import('axios')).default;
        return axios.put(url, data, config);
      },

      delete: async (url: string, config?: any) => {
        if (!hasPermission('network')) {
          throw new Error('Missing network permission');
        }
        const axios = (await import('axios')).default;
        return axios.delete(url, config);
      }
    },

    websocket: {
      connect: (url: string) => {
        if (!hasPermission('websocket')) {
          throw new Error('Missing websocket permission');
        }
        // WebSocket connection logic
        const WebSocket = (window as any).WebSocket;
        return new WebSocket(url);
      },

      send: (data: any) => {
        if (!hasPermission('websocket')) {
          throw new Error('Missing websocket permission');
        }
        const ws = (window as any).currentWebSocket;
        if (ws && ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify(data));
        }
      },

      on: (event: string, handler: Function) => {
        if (!hasPermission('websocket')) {
          throw new Error('Missing websocket permission');
        }
        const ws = (window as any).currentWebSocket;
        ws.addEventListener(event, handler);
      },

      off: (event: string, handler: Function) => {
        if (!hasPermission('websocket')) {
          throw new Error('Missing websocket permission');
        }
        const ws = (window as any).currentWebSocket;
        ws.removeEventListener(event, handler);
      }
    }
  };
}
