import { createRouter, createWebHistory } from 'vue-router';

import { fetchAndSetUser } from '@/functions/fetchAndSetUser';
import { userStore } from '@/stores/userStore';

export const router = createRouter({
  routes: [
    {
      path: '/dashboard',
      component: () => import('../dashboard/DashBoard.vue'),
    },
    {
      path: '/dashboard/integrations',
      component: () => import('../dashboard/Integrations.vue'),
    },
    {
      path: '/dashboard/integrations/:integration',
      component: () => import('../dashboard/Integrations.vue'),
    },
    {
      path: '/dashboard/events',
      component: () => import('../components/Soon.vue'),
    },
    {
      name: 'Commands',
      path: '/dashboard/commands',
      component: () => import('../dashboard/Commands.vue'),
    },
    {
      path: '/dashboard/greetings',
      component: () => import('../dashboard/Greetings.vue'),
    },
    {
      path: '/dashboard/timers',
      component: () => import('../dashboard/Timers.vue'),
    },
    {
      path: '/dashboard/keywords',
      component: () => import('../dashboard/Keywords.vue'),
    },
    {
      path: '/dashboard/wordscounters',
      component: () => import('../dashboard/WordsCounters.vue'),
    },
    {
      path: '/dashboard/variables',
      component: () => import('../dashboard/Variables.vue'),
    },
    {
      path: '/dashboard/moderation',
      component: () => import('../dashboard/Moderation.vue'),
    },
    {
      path: '/dashboard/settings',
      component: () => import('../dashboard/Settings.vue'),
    },
    {
      path: '/dashboard/users',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/dashboard/overlays',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/dashboard/files',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/dashboard/quotes',
      component: () => import('../components/Soon.vue'),
    },
    {
      path: '/dashboard/admin',
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

