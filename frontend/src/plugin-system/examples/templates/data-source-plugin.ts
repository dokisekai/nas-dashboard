/**
 * Data Source Plugin Template
 * Use this template to create data providers for the dashboard
 */

import type { PluginContext } from '../../types/plugin';

export const manifest = {
  id: 'my-datasource-plugin',
  name: 'My Data Source Plugin',
  version: '1.0.0',
  description: 'A custom data source for the dashboard',
  author: 'Your Name',
  main: 'plugin.ts',
  permissions: ['network', 'storage'] as any,
  category: 'data',
  keywords: ['data', 'source', 'api']
};

export default function createPlugin(context: PluginContext) {
  const { api, storage, logger, utils } = context;

  // Data fetcher with debouncing
  const fetchData = utils.debounce(async () => {
    try {
      const response = await api.network.get('https://api.example.com/data');
      await storage.set('cachedData', response.data);
      return response.data;
    } catch (error) {
      logger.error('Failed to fetch data:', error);
      throw error;
    }
  }, 1000);

  // Real-time data stream
  let ws: any = null;

  const startDataStream = () => {
    try {
      ws = api.websocket.connect('wss://api.example.com/stream');

      api.websocket.on(ws, 'message', (event: any) => {
        const data = JSON.parse(event.data);
        storage.set('realtimeData', data);
        logger.debug('Received real-time data:', data);
      });

      api.websocket.on(ws, 'open', () => {
        logger.info('Data stream connected');
      });

      api.websocket.on(ws, 'close', () => {
        logger.info('Data stream closed');
      });

      api.websocket.on(ws, 'error', (error: any) => {
        logger.error('WebSocket error:', error);
      });

      // Send authentication
      api.websocket.send(ws, {
        type: 'auth',
        token: 'your-token'
      });
    } catch (error) {
      logger.error('Failed to start data stream:', error);
    }
  };

  const stopDataStream = () => {
    if (ws) {
      ws.close();
      ws = null;
      logger.info('Data stream stopped');
    }
  };

  return {
    async onActivate() {
      logger.info('Starting data source');
      await fetchData();
    },

    async onDeactivate() {
      logger.info('Stopping data source');
      stopDataStream();
    },

    // Data source methods
    async getData() {
      // Try cache first
      const cached = await storage.get('cachedData');
      if (cached) {
        return cached;
      }

      // Fetch fresh data
      return await fetchData();
    },

    getRealtimeData() {
      return storage.get('realtimeData');
    },

    startStream: startDataStream,
    stopStream: stopDataStream,

    // Aggregation methods
    async aggregateData(metric: string, timeRange: string) {
      const data = await this.getData();
      // Implement aggregation logic
      return {
        metric,
        timeRange,
        value: 0,
        unit: ''
      };
    },

    // Historical data
    async getHistoricalData(startDate: string, endDate: string) {
      try {
        const response = await api.network.get('/api/historical', {
          params: { startDate, endDate }
        });
        return response.data;
      } catch (error) {
        logger.error('Failed to fetch historical data:', error);
        return [];
      }
    }
  };
}
