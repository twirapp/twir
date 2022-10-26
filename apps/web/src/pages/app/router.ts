import { createRouter, createWebHistory } from 'vue-router';

import {
  createUserDashboard,
  getProfile,
  redirectToLogin,
  selectedDashboardStore,
} from '@/services/auth';

export function createAppRouter() {
  const router = createRouter({
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
      {
        path: '/:pathMatch(.*)*',
        redirect: '/dashboard',
      },
    ],
  });

  router.beforeEach(async (_to, from, next) => {
    if (from.path === '/') {
      try {
        const user = await getProfile();

        if (!selectedDashboardStore.get()) {
          const userDashboard = createUserDashboard(user);
          selectedDashboardStore.set(userDashboard);
        }
      } catch (error) {
        return redirectToLogin();
      }
    }
    next();
  });

  return router;
}
