/**
 * Plugin Utilities - Common utilities for plugin development
 */

import type { PluginUtils } from '../types/plugin';

// Event Emitter class
class EventEmitter {
  private events: Map<string, Function[]> = new Map();

  on(event: string, handler: Function): void {
    if (!this.events.has(event)) {
      this.events.set(event, []);
    }
    this.events.get(event)!.push(handler);
  }

  off(event: string, handler: Function): void {
    const handlers = this.events.get(event);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  emit(event: string, ...args: any[]): void {
    const handlers = this.events.get(event);
    if (handlers) {
      handlers.forEach(handler => handler(...args));
    }
  }

  once(event: string, handler: Function): void {
    const onceWrapper = (...args: any[]) => {
      handler(...args);
      this.off(event, onceWrapper);
    };
    this.on(event, onceWrapper);
  }

  removeAllListeners(event?: string): void {
    if (event) {
      this.events.delete(event);
    } else {
      this.events.clear();
    }
  }
}

export function createPluginUtils(): PluginUtils {
  return {
    // Function utilities
    debounce: (func: Function, wait: number): Function => {
      let timeout: any;
      return function executedFunction(...args: any[]) {
        const later = () => {
          clearTimeout(timeout);
          func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
      };
    },

    throttle: (func: Function, limit: number): Function => {
      let inThrottle: boolean;
      return function executedFunction(...args: any[]) {
        if (!inThrottle) {
          func(...args);
          inThrottle = true;
          setTimeout(() => inThrottle = false, limit);
        }
      };
    },

    deepMerge: (target: any, source: any): any => {
      const output = { ...target };
      if (isObject(target) && isObject(source)) {
        Object.keys(source).forEach(key => {
          if (isObject(source[key])) {
            if (!(key in target)) {
              Object.assign(output, { [key]: source[key] });
            } else {
              output[key] = createPluginUtils().deepMerge(target[key], source[key]);
            }
          } else {
            Object.assign(output, { [key]: source[key] });
          }
        });
      }
      return output;
    },

    clone: <T>(obj: T): T => {
      if (obj === null || typeof obj !== 'object') {
        return obj;
      }
      if (obj instanceof Date) {
        return new Date(obj.getTime()) as any;
      }
      if (obj instanceof Array) {
        return obj.map(item => createPluginUtils().clone(item)) as any;
      }
      if (isObject(obj)) {
        const clonedObj = {} as any;
        for (const key in obj) {
          if (obj.hasOwnProperty(key)) {
            clonedObj[key] = createPluginUtils().clone(obj[key]);
          }
        }
        return clonedObj;
      }
      return obj;
    },

    // Event Emitter
    EventEmitter: EventEmitter,

    // Validation utilities
    validate: {
      email: (email: string): boolean => {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
      },

      url: (url: string): boolean => {
        try {
          new URL(url);
          return true;
        } catch {
          return false;
        }
      },

      required: (value: any): boolean => {
        if (value === null || value === undefined) {
          return false;
        }
        if (typeof value === 'string' && value.trim() === '') {
          return false;
        }
        if (Array.isArray(value) && value.length === 0) {
          return false;
        }
        return true;
      }
    }
  };
}

// Helper function
function isObject(item: any): boolean {
  return item && typeof item === 'object' && !Array.isArray(item);
}
