/**
 * Plugin Logger - Provides logging functionality for plugins
 */

import type { PluginLogger } from '../types/plugin';

export function createPluginLogger(pluginId: string): PluginLogger {
  const formatMessage = (...args: any[]): string => {
    const timestamp = new Date().toISOString();
    const prefix = `[${timestamp}] [Plugin: ${pluginId}]`;
    const message = args.map(arg =>
      typeof arg === 'object' ? JSON.stringify(arg, null, 2) : String(arg)
    ).join(' ');
    return `${prefix} ${message}`;
  };

  return {
    debug: (...args: any[]) => {
      if (import.meta.env.DEV) {
        console.debug(formatMessage(...args));
      }
    },

    info: (...args: any[]) => {
      console.info(formatMessage(...args));
    },

    warn: (...args: any[]) => {
      console.warn(formatMessage(...args));
    },

    error: (...args: any[]) => {
      console.error(formatMessage(...args));
    }
  };
}
