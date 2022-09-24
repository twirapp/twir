import type { AuthUser } from '@tsuwari/shared';
import { createRouter, createWebHistory } from 'vue-router';

import { redirectToLogin } from '@/services/auth.service.js';
import { setUser } from '@/stores/user.js';
import { authFetch } from '@/utils/authFetch.js';

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
      const response = await authFetch('/api/auth/profile');

      if (!response.ok) {
        return redirectToLogin();
      }
      const user = (await response.json()) as AuthUser;
      setUser(user);
    }
    next();
  });

  return router;
}
