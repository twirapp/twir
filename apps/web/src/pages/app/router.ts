import { createRouter as _createRouter, createWebHistory } from 'vue-router';

export function createAppRouter() {
  return _createRouter({
    history: createWebHistory('app'),
    routes: [
      {
        path: '/dashboard',
        component: () => import('@/pages/app/Dashboard.vue'),
      },
      {
        path: '/commands',
        component: () => import('@/pages/app/Commands.vue'),
      },
    ],
  });
}
