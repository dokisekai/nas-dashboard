/**
 * Plugin Marketplace API - Handles plugin discovery and search
 */

import type { PluginMarketplaceInfo, PluginManifest } from '../types/plugin';

export interface MarketplaceConfig {
  apiEndpoint: string;
  cacheTimeout?: number;
}

export class PluginMarketplace {
  private config: MarketplaceConfig;
  private cache: Map<string, { data: any; timestamp: number }> = new Map();

  constructor(config: MarketplaceConfig) {
    this.config = config;
  }

  /**
   * Get all available plugins
   */
  async getAllPlugins(): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache('/plugins', 'all-plugins');
  }

  /**
   * Get plugin by ID
   */
  async getPlugin(pluginId: string): Promise<PluginMarketplaceInfo | null> {
    try {
      return await this.fetchWithCache(`/plugins/${pluginId}`, `plugin-${pluginId}`);
    } catch (error) {
      console.error(`Failed to get plugin ${pluginId}:`, error);
      return null;
    }
  }

  /**
   * Search plugins
   */
  async searchPlugins(query: string, filters?: SearchFilters): Promise<PluginMarketplaceInfo[]> {
    const params = new URLSearchParams();
    params.append('q', query);

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (Array.isArray(value)) {
          value.forEach(v => params.append(key, String(v)));
        } else if (value !== undefined) {
          params.append(key, String(value));
        }
      });
    }

    return this.fetchWithCache(`/plugins/search?${params}`, `search-${query}-${JSON.stringify(filters)}`);
  }

  /**
   * Get plugins by category
   */
  async getPluginsByCategory(category: string): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache(`/plugins/category/${category}`, `category-${category}`);
  }

  /**
   * Get plugin categories
   */
  async getCategories(): Promise<string[]> {
    return this.fetchWithCache('/categories', 'categories');
  }

  /**
   * Get featured plugins
   */
  async getFeaturedPlugins(): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache('/plugins/featured', 'featured');
  }

  /**
   * Get popular plugins
   */
  async getPopularPlugins(limit = 10): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache(`/plugins/popular?limit=${limit}`, `popular-${limit}`);
  }

  /**
   * Get recently updated plugins
   */
  async getRecentlyUpdated(limit = 10): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache(`/plugins/recent?limit=${limit}`, `recent-${limit}`);
  }

  /**
   * Get plugin reviews
   */
  async getPluginReviews(pluginId: string): Promise<PluginReview[]> {
    return this.fetchWithCache(`/plugins/${pluginId}/reviews`, `reviews-${pluginId}`);
  }

  /**
   * Submit plugin review
   */
  async submitReview(
    pluginId: string,
    review: Omit<PluginReview, 'id' | 'date'>
  ): Promise<{ success: boolean; error?: string }> {
    try {
      await this.fetch(`/plugins/${pluginId}/reviews`, {
        method: 'POST',
        body: JSON.stringify(review)
      });

      // Clear cache for this plugin's reviews
      this.cache.delete(`reviews-${pluginId}`);

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Rate a plugin
   */
  async ratePlugin(
    pluginId: string,
    rating: number
  ): Promise<{ success: boolean; error?: string }> {
    try {
      await this.fetch(`/plugins/${pluginId}/rate`, {
        method: 'POST',
        body: JSON.stringify({ rating })
      });

      // Clear cache for this plugin
      this.cache.delete(`plugin-${pluginId}`);

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Get plugin download URL
   */
  async getDownloadUrl(pluginId: string, version?: string): Promise<string> {
    const params = version ? `?version=${version}` : '';
    const response = await this.fetch(`/plugins/${pluginId}/download${params}`);
    return response.url;
  }

  /**
   * Get plugin versions
   */
  async getVersions(pluginId: string): Promise<PluginVersion[]> {
    return this.fetchWithCache(`/plugins/${pluginId}/versions`, `versions-${pluginId}`);
  }

  /**
   * Report plugin
   */
  async reportPlugin(
    pluginId: string,
    reason: string,
    details?: string
  ): Promise<{ success: boolean; error?: string }> {
    try {
      await this.fetch(`/plugins/${pluginId}/report`, {
        method: 'POST',
        body: JSON.stringify({ reason, details })
      });

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Get developer info
   */
  async getDeveloperInfo(developerId: string): Promise<DeveloperInfo | null> {
    try {
      return await this.fetchWithCache(`/developers/${developerId}`, `dev-${developerId}`);
    } catch (error) {
      console.error(`Failed to get developer ${developerId}:`, error);
      return null;
    }
  }

  /**
   * Get plugins by developer
   */
  async getPluginsByDeveloper(developerId: string): Promise<PluginMarketplaceInfo[]> {
    return this.fetchWithCache(`/developers/${developerId}/plugins`, `dev-plugins-${developerId}`);
  }

  /**
   * Private methods
   */

  private async fetchWithCache(key: string, cacheKey: string): Promise<any> {
    const cached = this.cache.get(cacheKey);
    const cacheTimeout = this.config.cacheTimeout || 5 * 60 * 1000; // 5 minutes default

    if (cached && Date.now() - cached.timestamp < cacheTimeout) {
      return cached.data;
    }

    const data = await this.fetch(key);
    this.cache.set(cacheKey, { data, timestamp: Date.now() });
    return data;
  }

  private async fetch(endpoint: string, options?: RequestInit): Promise<any> {
    const url = `${this.config.apiEndpoint}${endpoint}`;

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  /**
   * Clear cache
   */
  clearCache(): void {
    this.cache.clear();
  }

  /**
   * Clear specific cache entry
   */
  clearCacheEntry(cacheKey: string): void {
    this.cache.delete(cacheKey);
  }
}

// Types
export interface SearchFilters {
  category?: string;
  author?: string;
  tags?: string[];
  minRating?: number;
  price?: 'free' | 'paid' | 'all';
  sortBy?: 'relevance' | 'rating' | 'downloads' | 'updated';
}

export interface PluginReview {
  id: string;
  pluginId: string;
  userId: string;
  userName: string;
  rating: number;
  title: string;
  content: string;
  date: string;
  helpful: number;
}

export interface PluginVersion {
  version: string;
  downloadUrl: string;
  releaseDate: string;
  changelog: string;
  size: number;
}

export interface DeveloperInfo {
  id: string;
  name: string;
  email?: string;
  website?: string;
  bio?: string;
  avatar?: string;
  plugins: string[];
  joined: string;
  verified: boolean;
}

// Singleton instance with default config
let marketplaceInstance: PluginMarketplace | null = null;

export function getPluginMarketplace(): PluginMarketplace {
  if (!marketplaceInstance) {
    marketplaceInstance = new PluginMarketplace({
      apiEndpoint: '/api/marketplace', // Default endpoint
      cacheTimeout: 5 * 60 * 1000 // 5 minutes
    });
  }
  return marketplaceInstance;
}

export function configurePluginMarketplace(config: MarketplaceConfig): PluginMarketplace {
  marketplaceInstance = new PluginMarketplace(config);
  return marketplaceInstance;
}
