/**
 * Example: Analytics Plugin
 * A plugin that adds analytics and reporting features
 */

import type { PluginContext } from '../../types/plugin';

export const manifest = {
  id: 'analytics-plugin',
  name: 'Analytics Plugin',
  version: '1.0.0',
  description: 'Advanced analytics and reporting dashboard',
  author: 'Dashboard Team',
  main: 'analytics-plugin.ts',
  permissions: ['network', 'storage', 'ui'] as any,
  category: 'analytics',
  keywords: ['analytics', 'reporting', 'statistics']
};

export default function createPlugin(context: PluginContext) {
  const { api, storage, logger, utils } = context;

  // Analytics data cache
  let dataCache: Map<string, any> = new Map();

  // Event tracking
  const eventQueue: any[] = [];
  let flushInterval: any = null;

  const startEventTracking = () => {
    flushInterval = setInterval(() => {
      flushEvents();
    }, 30000); // Flush every 30 seconds
  };

  const stopEventTracking = () => {
    if (flushInterval) {
      clearInterval(flushInterval);
      flushInterval = null;
    }
  };

  const flushEvents = async () => {
    if (eventQueue.length === 0) return;

    try {
      const events = [...eventQueue];
      eventQueue.length = 0;

      await api.network.post('/api/analytics/events', {
        events,
        timestamp: new Date().toISOString()
      });

      logger.info(`Flushed ${events.length} events`);
    } catch (error) {
      logger.error('Failed to flush events:', error);
      // Re-queue events on failure
      eventQueue.unshift(...events);
    }
  };

  return {
    async onInstall() {
      logger.info('Installing analytics plugin');

      // Initialize settings
      await storage.set('settings', {
        enabled: true,
        samplingRate: 1.0,
        flushInterval: 30000,
        batchSize: 100
      });

      // Initialize data cache
      await storage.set('cache', {});
    },

    async onEnable() {
      logger.info('Enabling analytics plugin');

      // Register analytics dashboard component
      api.ui.registerComponent({
        id: 'analytics-dashboard',
        name: 'Analytics Dashboard',
        render: (container: HTMLElement) => {
          const dashboard = document.createElement('div');
          dashboard.className = 'analytics-dashboard';
          dashboard.innerHTML = `
            <div class="analytics-header">
              <h2>Analytics Dashboard</h2>
              <div class="date-range-selector">
                <select id="date-range">
                  <option value="24h">Last 24 hours</option>
                  <option value="7d">Last 7 days</option>
                  <option value="30d">Last 30 days</option>
                </select>
              </div>
            </div>
            <div class="analytics-metrics">
              <div class="metric-card">
                <h3>Total Events</h3>
                <div class="metric-value" id="total-events">0</div>
              </div>
              <div class="metric-card">
                <h3>Active Users</h3>
                <div class="metric-value" id="active-users">0</div>
              </div>
              <div class="metric-card">
                <h3>Page Views</h3>
                <div class="metric-value" id="page-views">0</div>
              </div>
              <div class="metric-card">
                <h3>Avg. Session</h3>
                <div class="metric-value" id="avg-session">0s</div>
              </div>
            </div>
            <div class="analytics-charts">
              <canvas id="events-chart"></canvas>
            </div>
          `;

          container.appendChild(dashboard);

          // Load data when component mounts
          loadAnalyticsData('24h');

          return () => {
            container.removeChild(dashboard);
          };

          async function loadAnalyticsData(dateRange: string) {
            try {
              const data = await api.network.get('/api/analytics/summary', {
                params: { dateRange }
              });

              // Update metrics
              dashboard.querySelector('#total-events')!.textContent = data.totalEvents.toLocaleString();
              dashboard.querySelector('#active-users')!.textContent = data.activeUsers.toLocaleString();
              dashboard.querySelector('#page-views')!.textContent = data.pageViews.toLocaleString();
              dashboard.querySelector('#avg-session')!.textContent = `${Math.round(data.avgSession)}s`;

              // Render chart
              renderChart(data.timeline);
            } catch (error) {
              logger.error('Failed to load analytics data:', error);
            }
          }

          function renderChart(timeline: any[]) {
            // Chart rendering logic would go here
            logger.debug('Rendering chart with data:', timeline);
          }
        }
      });
    },

    async onActivate() {
      logger.info('Activating analytics plugin');

      // Load cache
      const cache = await storage.get('cache') || {};
      dataCache = new Map(Object.entries(cache));

      // Start event tracking
      startEventTracking();
    },

    async onDeactivate() {
      logger.info('Deactivating analytics plugin');

      // Stop event tracking
      stopEventTracking();

      // Flush remaining events
      await flushEvents();

      // Save cache
      await storage.set('cache', Object.fromEntries(dataCache));
    },

    async onDisable() {
      logger.info('Disabling analytics plugin');
      api.ui.unregisterComponent('analytics-dashboard');
    },

    async onUninstall() {
      logger.info('Uninstalling analytics plugin');

      // Cleanup
      stopEventTracking();
      await storage.clear();
      dataCache.clear();
    },

    // Event tracking methods
    trackEvent(event: string, properties?: any) {
      const settingsPromise = storage.get('settings');

      settingsPromise.then(settings => {
        if (!settings.enabled) return;

        // Apply sampling
        if (Math.random() > settings.samplingRate) return;

        eventQueue.push({
          event,
          properties,
          timestamp: new Date().toISOString(),
          sessionId: getSessionId()
        });

        // Auto-flush if batch size reached
        if (eventQueue.length >= settings.batchSize) {
          flushEvents();
        }
      });

      logger.debug('Event tracked:', event, properties);
    },

    trackPageView(page: string, properties?: any) {
      this.trackEvent('page_view', {
        page,
        ...properties
      });
    },

    trackUserAction(action: string, properties?: any) {
      this.trackEvent('user_action', {
        action,
        ...properties
      });
    },

    // Analytics queries
    async getEvents(filters?: any) {
      try {
        const response = await api.network.get('/api/analytics/query', {
          params: { filters }
        });
        return response.data;
      } catch (error) {
        logger.error('Failed to query events:', error);
        return [];
      }
    },

    async getMetrics(metrics: string[], dateRange: string) {
      try {
        const response = await api.network.post('/api/analytics/metrics', {
          metrics,
          dateRange
        });
        return response.data;
      } catch (error) {
        logger.error('Failed to get metrics:', error);
        return {};
      }
    },

    async generateReport(reportConfig: any) {
      try {
        const response = await api.network.post('/api/analytics/reports', reportConfig);
        return response.data;
      } catch (error) {
        logger.error('Failed to generate report:', error);
        return null;
      }
    },

    // Settings management
    async getSettings() {
      return await storage.get('settings') || {};
    },

    async updateSettings(newSettings: any) {
      const current = await this.getSettings();
      const merged = utils.deepMerge(current, newSettings);
      await storage.set('settings', merged);

      // Restart event tracking with new settings
      stopEventTracking();
      startEventTracking();

      logger.info('Settings updated:', merged);
    }
  };
}

// Session ID helper
function getSessionId(): string {
  let sessionId = sessionStorage.getItem('analytics_session_id');
  if (!sessionId) {
    sessionId = `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    sessionStorage.setItem('analytics_session_id', sessionId);
  }
  return sessionId;
}
