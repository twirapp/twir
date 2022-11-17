import { createRouter, createWebHistory } from 'vue-router';

import { fetchAndSetUser } from '@/functions/fetchAndSetUser';
import { redirectToLogin } from '@/functions/redirectToLogin.js';
import { userStore } from '@/stores/userStore';

export const router = createRouter({
  routes: [
    {
      path: '/',
      component: () => import('../dashboard/DashBoard.vue'),
    },
    {
      path: '/integrations',
      component: () => import('../dashboard/Integrations.vue'),
    },
    {
      path: '/integrations/:integration',
      component: () => import('../dashboard/Integrations.vue'),
    },
    {
      path: '/events',
      component: () => import('../components/Soon.vue'),
    },
    {
      name: 'Commands',
      path: '/commands',
      component: () => import('../dashboard/Commands.vue'),
    },
    {
      path: '/greetings',
      component: () => import('../dashboard/Greetings.vue'),
    },
    {
      path: '/timers',
      component: () => import('../dashboard/Timers.vue'),
    },
    {
      path: '/keywords',
      component: () => import('../dashboard/Keywords.vue'),
    },
    {
      path: '/variables',
      component: () => import('../dashboard/Variables.vue'),
    },
    {
      path: '/moderation',
      component: () => import('../dashboard/Moderation.vue'),
    },
    {
      path: '/settings',
      component: () => import('../dashboard/Settings.vue'),
    },
    {
      path: '/users',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/overlays',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/files',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/quotes',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/admin',
      component: () => import('../admin/Main.vue'),
    },
    { name: '404', path: '/:pathMatch(.*)*', component: () => import('../pages/NotFound.vue') },
  ],
  history: createWebHistory(),
});

router.beforeResolve(async (to, _from, next) => {
  if (to.fullPath.startsWith('/admin')) {
    const user = userStore.get() || (await fetchAndSetUser().then(() => userStore.get()));

    if (user?.isBotAdmin) {
      return next();
    } else return router.push({ name: '404', params: { pagePath: to.fullPath } });
  } else {
    return next();
  }
});

router.beforeEach(async (to, _from, next) => {
  if (to.path === '/') {
    const accessToken = localStorage.getItem('accessToken');
    const refreshToken = localStorage.getItem('refreshToken');
    if (accessToken && refreshToken) {
      fetchAndSetUser();
    }

    return next();
  }

  if (to.path.startsWith('/dashboard')) {
    await fetchAndSetUser();

    if (!userStore.get()) {
      return redirectToLogin();
    }

    next();
  } else {
    next();
  }
});
