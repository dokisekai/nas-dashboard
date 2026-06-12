/**
 * Plugin Router Integration
 * Example of how to integrate plugins with Vue Router
 */

import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';
import type { PluginContext } from '../types/plugin';

/**
 * Plugin Route Helper
 * Provides utilities for plugins to register routes
 */
export class PluginRouter {
  private routes: Map<string, RouteRecordRaw[]> = new Map();

  /**
   * Add routes from a plugin
   */
  addRoutes(pluginId: string, routes: RouteRecordRaw[]): void {
    if (!this.routes.has(pluginId)) {
      this.routes.set(pluginId, []);
    }
    this.routes.get(pluginId)!.push(...routes);

    console.log(`📍 Routes added for plugin: ${pluginId}`, routes);
  }

  /**
   * Remove routes from a plugin
   */
  removeRoutes(pluginId: string): void {
    this.routes.delete(pluginId);
    console.log(`📍 Routes removed for plugin: ${pluginId}`);
  }

  /**
   * Get all plugin routes
   */
  getAllRoutes(): RouteRecordRaw[] {
    const allRoutes: RouteRecordRaw[] = [];
    for (const routes of this.routes.values()) {
      allRoutes.push(...routes);
    }
    return allRoutes;
  }

  /**
   * Get routes for a specific plugin
   */
  getPluginRoutes(pluginId: string): RouteRecordRaw[] {
    return this.routes.get(pluginId) || [];
  }
}

// Singleton instance
let routerInstance: PluginRouter | null = null;

export function getPluginRouter(): PluginRouter {
  if (!routerInstance) {
    routerInstance = new PluginRouter();
  }
  return routerInstance;
}

/**
 * Create router with plugin support
 */
export function createPluginRouter() {
  const pluginRouter = getPluginRouter();

  // Base routes
  const baseRoutes: RouteRecordRaw[] = [
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/Dashboard.vue')
    },
    {
      path: '/plugins',
      name: 'Plugins',
      component: () => import('../views/Plugins.vue')
    }
  ];

  const router = createRouter({
    history: createWebHistory(),
    routes: baseRoutes
  });

  // Watch for plugin route changes
  const updateRoutes = () => {
    const pluginRoutes = pluginRouter.getAllRoutes();
    const currentRoutes = router.getRoutes();

    // Add new plugin routes
    for (const route of pluginRoutes) {
      if (!currentRoutes.find(r => r.path === route.path)) {
        router.addRoute(route);
      }
    }
  };

  // Initial update
  updateRoutes();

  // Return router instance
  return router;
}

/**
 * Plugin Route Component Example
 */
export const PluginRouteComponent = {
  name: 'PluginRoute',

  props: {
    pluginId: {
      type: String,
      required: true
    }
  },

  setup(props: any) {
    const { loader } = usePluginSystem();

    const plugin = computed(() => {
      return loader.getPluginState(props.pluginId);
    });

    const pluginRoutes = computed(() => {
      const pluginRouter = getPluginRouter();
      return pluginRouter.getPluginRoutes(props.pluginId);
    });

    return {
      plugin,
      pluginRoutes
    };
  },

  template: `
    <div class="plugin-routes">
      <h3>{{ plugin?.manifest.name }} Routes</h3>
      <ul>
        <li v-for="route in pluginRoutes" :key="route.path">
          <router-link :to="route.path">
            {{ route.name || route.path }}
          </router-link>
        </li>
      </ul>
    </div>
  `
};

/**
 * Usage in Plugin
 */
export function createPluginWithRoutes(context: PluginContext) {
  const { api, pluginId, logger } = context;
  const pluginRouter = getPluginRouter();

  return {
    async onActivate() {
      // Register routes when plugin activates
      pluginRouter.addRoutes(pluginId, [
        {
          path: '/my-plugin',
          name: 'MyPlugin',
          component: () => import('./MyPluginView.vue')
        },
        {
          path: '/my-plugin/settings',
          name: 'MyPluginSettings',
          component: () => import('./MyPluginSettings.vue')
        }
      ]);

      logger.info('Routes registered');
    },

    async onDeactivate() {
      // Remove routes when plugin deactivates
      pluginRouter.removeRoutes(pluginId);

      logger.info('Routes removed');
    }
  };
}

/**
 * Navigation Guard for Plugins
 */
export function setupPluginNavigationGuards(router: any) {
  // Before each route change
  router.beforeEach((to: any, from: any, next: any) => {
    // Check if route belongs to a plugin
    const pluginId = to.meta.pluginId;

    if (pluginId) {
      const { loader } = usePluginSystem();
      const plugin = loader.getPluginState(pluginId);

      // Check if plugin is active
      if (!plugin?.active) {
        console.warn(`Plugin ${pluginId} is not active`);
        next({ name: 'Plugins' });
        return;
      }
    }

    next();
  });

  // After each route change
  router.afterEach((to: any, from: any) => {
    console.log(`Navigated from ${from.path} to ${to.path}`);
  });
}

/**
 * Route Meta Types
 */
export interface PluginRouteMeta {
  pluginId: string;
  permission?: string;
  settings?: boolean;
  requiresActive?: boolean;
}

/**
 * Example: Plugin with Route Meta
 */
export const pluginRouteExample: RouteRecordRaw = {
  path: '/analytics-plugin',
  name: 'AnalyticsPlugin',
  component: () => import('./views/AnalyticsPlugin.vue'),
  meta: {
    pluginId: 'analytics-plugin',
    permission: 'network',
    settings: false,
    requiresActive: true
  } as PluginRouteMeta,
  children: [
    {
      path: 'dashboard',
      name: 'AnalyticsDashboard',
      component: () => import('./views/AnalyticsDashboard.vue')
    },
    {
      path: 'reports',
      name: 'AnalyticsReports',
      component: () => import('./views/AnalyticsReports.vue')
    },
    {
      path: 'settings',
      name: 'AnalyticsSettings',
      component: () => import('./views/AnalyticsSettings.vue'),
      meta: {
        settings: true
      } as PluginRouteMeta
    }
  ]
};

/**
 * Dynamic Route Loading
 */
export async function loadPluginRoutes(pluginId: string): Promise<void> {
  const { loader } = usePluginSystem();
  const plugin = loader.getPluginState(pluginId);

  if (!plugin) {
    throw new Error(`Plugin ${pluginId} not found`);
  }

  // Check if plugin has routes
  if (plugin.instance?.getRoutes) {
    const routes = await plugin.instance.getRoutes();
    const pluginRouter = getPluginRouter();
    pluginRouter.addRoutes(pluginId, routes);
  }
}

/**
 * Route Update Watcher
 */
export function watchPluginRoutes() {
  const { loader } = usePluginSystem();
  const pluginRouter = getPluginRouter();

  // Watch for plugin changes
  const checkRoutes = () => {
    const activePlugins = loader.getActivePlugins();

    for (const plugin of activePlugins) {
      if (plugin.instance?.getRoutes) {
        const existingRoutes = pluginRouter.getPluginRoutes(plugin.id);
        if (existingRoutes.length === 0) {
          loadPluginRoutes(plugin.id);
        }
      }
    }
  };

  // Check every second
  setInterval(checkRoutes, 1000);
}

export default createPluginRouter;
