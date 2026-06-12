/**
 * Plugin Storage - Provides isolated storage for plugins
 */

import type { PluginStorage } from '../types/plugin';

const STORAGE_PREFIX = 'plugin-storage:';

export function createPluginStorage(pluginId: string): PluginStorage {
  const getStorageKey = (key: string): string => {
    return `${STORAGE_PREFIX}${pluginId}:${key}`;
  };

  return {
    async get(key: string): Promise<any> {
      try {
        const storageKey = getStorageKey(key);
        const value = localStorage.getItem(storageKey);
        if (value === null) {
          return null;
        }
        return JSON.parse(value);
      } catch (error) {
        console.error(`Storage get failed for ${key}:`, error);
        return null;
      }
    },

    async set(key: string, value: any): Promise<void> {
      try {
        const storageKey = getStorageKey(key);
        const serialized = JSON.stringify(value);
        localStorage.setItem(storageKey, serialized);
      } catch (error) {
        console.error(`Storage set failed for ${key}:`, error);
        throw error;
      }
    },

    async remove(key: string): Promise<void> {
      try {
        const storageKey = getStorageKey(key);
        localStorage.removeItem(storageKey);
      } catch (error) {
        console.error(`Storage remove failed for ${key}:`, error);
        throw error;
      }
    },

    async clear(): Promise<void> {
      try {
        const prefix = getStorageKey('');
        const keys = Object.keys(localStorage);
        for (const key of keys) {
          if (key.startsWith(prefix)) {
            localStorage.removeItem(key);
          }
        }
      } catch (error) {
        console.error('Storage clear failed:', error);
        throw error;
      }
    },

    async keys(): Promise<string[]> {
      try {
        const prefix = getStorageKey('');
        const allKeys = Object.keys(localStorage);
        return allKeys
          .filter(key => key.startsWith(prefix))
          .map(key => key.substring(prefix.length));
      } catch (error) {
        console.error('Storage keys failed:', error);
        return [];
      }
    }
  };
}

export function clearPluginStorage(pluginId: string): void {
  const storage = createPluginStorage(pluginId);
  storage.clear();
}
