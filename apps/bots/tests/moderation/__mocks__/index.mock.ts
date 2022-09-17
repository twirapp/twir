/* import { vi } from 'vitest';

vi.mock('tlds', () => ({
  default: ['com'],
}));

vi.mock('../../../src/libs/redis.js', () => ({
  redis: {
    del: vi.fn(() => true),
    set: vi.fn(() => true),
    get: vi.fn(() => true),
  },
}));

vi.mock('../../../src/libs/prisma.js', () => ({
  prisma: {
    permit: {
      findFirst: vi.fn(() => null),
      delete: vi.fn(() => null),
    },
  },
})); */
