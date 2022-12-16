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
    {
      path: '/dashboard/moderation',
      component: () => import('../pages/Moderation.vue'),
    },
    {
      path: '/dashboard/keywords',
      component: () => import('../pages/Keywords.vue'),
    },
    {
      path: '/dashboard/variables',
      component: () => import('../pages/Variables.vue'),
    },
    {
      path: '/dashboard/greetings',
      component: () => import('../pages/Greetings.vue'),
    },
    {
      path: '/dashboard/settings',
      component: () => import('../pages/Settings.vue'),
    },
  ],
  history: createWebHistory(),
});