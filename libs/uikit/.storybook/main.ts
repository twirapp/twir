import type { StorybookViteConfig } from "@storybook/builder-vite"
import { mergeConfig, UserConfig } from 'vite';
import { resolve } from "node:path"

const config: StorybookViteConfig = {
  stories: [
    "../src/**/*.stories.mdx",
    "../src/**/*.stories.@(js|jsx|ts|tsx)"
  ],
  addons: [
    "@storybook/addon-links",
    "@storybook/addon-essentials",
    "@storybook/addon-interactions"
  ],
  framework: "@storybook/vue3",
  core: {
    builder: "@storybook/builder-vite"
  },
  features: {
    storyStoreV7: true
  },
  typescript: {
    check: false,
  },
  async viteFinal(config) {
    return mergeConfig(config, {
      resolve: {
        alias: [
          { find: '@', replacement: resolve(__dirname, '..', 'src') }
        ]
      },
    } as UserConfig);
  },
}

module.exports = config;