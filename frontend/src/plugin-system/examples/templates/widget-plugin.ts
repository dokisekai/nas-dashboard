/**
 * Widget Plugin Template
 * Use this template to create dashboard widgets
 */

import type { PluginContext } from '../../types/plugin';

export const manifest = {
  id: 'my-widget-plugin',
  name: 'My Widget Plugin',
  version: '1.0.0',
  description: 'A custom dashboard widget',
  author: 'Your Name',
  main: 'plugin.ts',
  permissions: ['ui', 'storage'] as const,
  category: 'widgets',
  keywords: ['widget', 'dashboard']
};

export default function createPlugin(context: PluginContext) {
  const { api, storage, logger } = context;

  // Widget component definition
  const widgetComponent = {
    id: 'my-widget',
    name: 'My Widget',
    description: 'A custom widget',
    icon: '🎯',
    size: 'medium', // small, medium, large
    render: (container: HTMLElement) => {
      // Create widget UI
      const widget = document.createElement('div');
      widget.className = 'my-widget';
      widget.innerHTML = `
        <h3>My Widget</h3>
        <div class="widget-content">
          <p>Widget content here</p>
        </div>
      `;

      // Add styles
      const style = document.createElement('style');
      style.textContent = `
        .my-widget {
          padding: 16px;
          background: white;
          border-radius: 8px;
          box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .my-widget h3 {
          margin: 0 0 12px 0;
          color: #333;
        }
        .widget-content {
          color: #666;
        }
      `;

      container.appendChild(style);
      container.appendChild(widget);

      // Return cleanup function
      return () => {
        container.removeChild(widget);
        container.removeChild(style);
      };
    }
  };

  return {
    async onActivate() {
      logger.info('Registering widget');
      api.ui.registerComponent(widgetComponent);
    },

    async onDeactivate() {
      logger.info('Unregistering widget');
      api.ui.unregisterComponent('my-widget');
    },

    // Widget-specific methods
    getWidgetConfig() {
      return {
        refreshInterval: 5000,
        showHeader: true
      };
    },

    async getWidgetData() {
      return await storage.get('widgetData');
    }
  };
}
