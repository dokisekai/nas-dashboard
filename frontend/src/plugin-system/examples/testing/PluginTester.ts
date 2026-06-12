/**
 * Plugin Testing Tool
 * Comprehensive testing utilities for plugin development
 */

import type { PluginManifest } from '../types/plugin';
import { getPluginLoader, getPluginManager } from '../index';

export interface TestResult {
  name: string;
  passed: boolean;
  error?: string;
  duration: number;
}

export interface TestSuite {
  name: string;
  tests: TestResult[];
  duration: number;
}

export class PluginTester {
  private results: TestResult[] = [];
  private startTime = 0;

  /**
   * Test plugin loading
   */
  async testLoad(
    source: string | object,
    manifest: PluginManifest
  ): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      const result = await loader.loadPlugin(source, manifest);

      if (!result.success) {
        return {
          name: 'Load Plugin',
          passed: false,
          error: result.error?.message,
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Load Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Load Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin enable
   */
  async testEnable(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      await loader.enablePlugin(pluginId);
      const state = loader.getPluginState(pluginId);

      if (!state?.enabled) {
        return {
          name: 'Enable Plugin',
          passed: false,
          error: 'Plugin not enabled after call',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Enable Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Enable Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin activation
   */
  async testActivate(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      await loader.activatePlugin(pluginId);
      const state = loader.getPluginState(pluginId);

      if (!state?.active) {
        return {
          name: 'Activate Plugin',
          passed: false,
          error: 'Plugin not active after call',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Activate Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Activate Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin deactivation
   */
  async testDeactivate(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      await loader.deactivatePlugin(pluginId);
      const state = loader.getPluginState(pluginId);

      if (state?.active) {
        return {
          name: 'Deactivate Plugin',
          passed: false,
          error: 'Plugin still active after deactivation',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Deactivate Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Deactivate Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin disable
   */
  async testDisable(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      await loader.disablePlugin(pluginId);
      const state = loader.getPluginState(pluginId);

      if (state?.enabled) {
        return {
          name: 'Disable Plugin',
          passed: false,
          error: 'Plugin still enabled after call',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Disable Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Disable Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin unload
   */
  async testUnload(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const loader = getPluginLoader();

    try {
      await loader.unloadPlugin(pluginId);
      const state = loader.getPluginState(pluginId);

      if (state) {
        return {
          name: 'Unload Plugin',
          passed: false,
          error: 'Plugin still registered after unload',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Unload Plugin',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Unload Plugin',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin storage
   */
  async testStorage(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();

    try {
      const { createPluginStorage } = await import('../sdk/storage');
      const storage = createPluginStorage(pluginId);

      // Test set/get
      await storage.set('test-key', { value: 'test-data' });
      const value = await storage.get('test-key');

      if (value?.value !== 'test-data') {
        return {
          name: 'Storage Operations',
          passed: false,
          error: 'Storage get returned unexpected value',
          duration: performance.now() - startTime
        };
      }

      // Test remove
      await storage.remove('test-key');
      const removed = await storage.get('test-key');

      if (removed !== null) {
        return {
          name: 'Storage Operations',
          passed: false,
          error: 'Storage remove did not delete key',
          duration: performance.now() - startTime
        };
      }

      // Test keys
      await storage.set('key1', 'value1');
      await storage.set('key2', 'value2');
      const keys = await storage.keys();

      if (keys.length !== 2 || !keys.includes('key1') || !keys.includes('key2')) {
        return {
          name: 'Storage Operations',
          passed: false,
          error: 'Storage keys returned unexpected values',
          duration: performance.now() - startTime
        };
      }

      // Cleanup
      await storage.clear();

      return {
        name: 'Storage Operations',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Storage Operations',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Test plugin configuration
   */
  async testConfig(pluginId: string): Promise<TestResult> {
    const startTime = performance.now();
    const manager = getPluginManager();

    try {
      // Test set
      await manager.updatePluginConfig(pluginId, { test: 'value' });
      const config = await manager.getPluginConfig(pluginId);

      if (config?.test !== 'value') {
        return {
          name: 'Configuration Management',
          passed: false,
          error: 'Configuration not set correctly',
          duration: performance.now() - startTime
        };
      }

      // Test reset
      await manager.resetPluginConfig(pluginId);
      const reset = await manager.getPluginConfig(pluginId);

      if (reset?.test === 'value') {
        return {
          name: 'Configuration Management',
          passed: false,
          error: 'Configuration not reset correctly',
          duration: performance.now() - startTime
        };
      }

      return {
        name: 'Configuration Management',
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: 'Configuration Management',
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Run all lifecycle tests
   */
  async testLifecycle(
    source: string | object,
    manifest: PluginManifest
  ): Promise<TestSuite> {
    this.startTime = performance.now();
    this.results = [];

    console.log(`\n🧪 Testing Plugin: ${manifest.id}`);
    console.log('─'.repeat(50));

    // Load test
    const loadResult = await this.testLoad(source, manifest);
    this.results.push(loadResult);
    this.logResult(loadResult);

    if (!loadResult.passed) {
      return {
        name: 'Lifecycle Tests',
        tests: this.results,
        duration: performance.now() - this.startTime
      };
    }

    // Enable test
    const enableResult = await this.testEnable(manifest.id);
    this.results.push(enableResult);
    this.logResult(enableResult);

    if (!enableResult.passed) {
      await this.cleanup(manifest.id);
      return {
        name: 'Lifecycle Tests',
        tests: this.results,
        duration: performance.now() - this.startTime
      };
    }

    // Activate test
    const activateResult = await this.testActivate(manifest.id);
    this.results.push(activateResult);
    this.logResult(activateResult);

    // Storage test (only if active)
    if (activateResult.passed) {
      const storageResult = await this.testStorage(manifest.id);
      this.results.push(storageResult);
      this.logResult(storageResult);
    }

    // Config test
    const configResult = await this.testConfig(manifest.id);
    this.results.push(configResult);
    this.logResult(configResult);

    // Deactivate test
    const deactivateResult = await this.testDeactivate(manifest.id);
    this.results.push(deactivateResult);
    this.logResult(deactivateResult);

    // Disable test
    const disableResult = await this.testDisable(manifest.id);
    this.results.push(disableResult);
    this.logResult(disableResult);

    // Cleanup
    await this.cleanup(manifest.id);

    return {
      name: 'Lifecycle Tests',
      tests: this.results,
      duration: performance.now() - this.startTime
    };
  }

  /**
   * Test custom functionality
   */
  async testCustomFunction(
    pluginId: string,
    functionName: string,
    ...args: any[]
  ): Promise<TestResult> {
    const startTime = performance.now();

    try {
      const loader = getPluginLoader();
      const state = loader.getPluginState(pluginId);

      if (!state?.instance) {
        return {
          name: `Custom: ${functionName}`,
          passed: false,
          error: 'Plugin instance not found',
          duration: performance.now() - startTime
        };
      }

      const func = state.instance[functionName];
      if (typeof func !== 'function') {
        return {
          name: `Custom: ${functionName}`,
          passed: false,
          error: `Function ${functionName} not found`,
          duration: performance.now() - startTime
        };
      }

      const result = await func.apply(state.instance, args);

      return {
        name: `Custom: ${functionName}`,
        passed: true,
        duration: performance.now() - startTime
      };
    } catch (error) {
      return {
        name: `Custom: ${functionName}`,
        passed: false,
        error: (error as Error).message,
        duration: performance.now() - startTime
      };
    }
  }

  /**
   * Get test summary
   */
  getSummary(): void {
    const passed = this.results.filter(r => r.passed).length;
    const failed = this.results.filter(r => !r.passed).length;
    const total = this.results.length;
    const duration = performance.now() - this.startTime;

    console.log('\n📊 Test Summary');
    console.log('─'.repeat(50));
    console.log(`Total: ${total}`);
    console.log(`✅ Passed: ${passed}`);
    console.log(`❌ Failed: ${failed}`);
    console.log(`⏱️  Duration: ${duration.toFixed(2)}ms`);

    if (failed > 0) {
      console.log('\n❌ Failed Tests:');
      this.results
        .filter(r => !r.passed)
        .forEach(r => console.log(`  - ${r.name}: ${r.error}`));
    }

    console.log('─'.repeat(50));
  }

  /**
   * Cleanup after tests
   */
  private async cleanup(pluginId: string): Promise<void> {
    try {
      const loader = getPluginLoader();
      const state = loader.getPluginState(pluginId);

      if (state?.enabled) {
        await loader.disablePlugin(pluginId);
      }

      await loader.unloadPlugin(pluginId);
    } catch (error) {
      console.error('Cleanup failed:', error);
    }
  }

  /**
   * Log individual test result
   */
  private logResult(result: TestResult): void {
    const status = result.passed ? '✅' : '❌';
    const duration = result.duration.toFixed(2) + 'ms';
    console.log(`${status} ${result.name} (${duration})`);

    if (!result.passed && result.error) {
      console.log(`   Error: ${result.error}`);
    }
  }

  /**
   * Export test results as JSON
   */
  exportResults(): string {
    return JSON.stringify({
      timestamp: new Date().toISOString(),
      duration: performance.now() - this.startTime,
      results: this.results
    }, null, 2);
  }
}

/**
 * Run complete plugin test suite
 */
export async function runPluginTests(
  source: string | object,
  manifest: PluginManifest
): Promise<void> {
  const tester = new PluginTester();
  await tester.testLifecycle(source, manifest);
  tester.getSummary();
}

/**
 * Quick test function
 */
export async function quickTest(
  source: string | object,
  manifest: PluginManifest
): Promise<boolean> {
  const loader = getPluginLoader();

  try {
    // Load
    const loadResult = await loader.loadPlugin(source, manifest);
    if (!loadResult.success) return false;

    // Enable
    await loader.enablePlugin(manifest.id);

    // Activate
    await loader.activatePlugin(manifest.id);

    // Cleanup
    await loader.unloadPlugin(manifest.id);

    return true;
  } catch (error) {
    console.error('Quick test failed:', error);
    return false;
  }
}
