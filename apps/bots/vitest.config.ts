import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    include: ['./tests/**/*.ts'],
    exclude: ['./tests/**/*.mock.ts', './tests/**/helpers.ts'],
    watch: false,
    coverage: {
      enabled: false,
    },
    reporters: ['verbose'],
  },
});