import { createRouter, createWebHistory } from 'vue-router';

export const router = createRouter({
  routes: [
    {
      path: '/dashboard',
      component: () => import('../pages/Home.vue'),
    },
    {
      path: '/dashboard/integrations',
      component: () => import('../pages/Integrations.vue'),
    },
    {
      path: '/dashboard/commands',
      component: () => import('../pages/Commands.vue'),
    },
    {
      path: '/dashboard/timers',
      component: () => import('../pages/Timers.vue'),
    },
  ],
  history: createWebHistory(),
});