import { createRouter, createWebHistory } from 'vue-router';

import { fetchAndSetUser } from '@/functions/fetchAndSetUser.js';
import { userStore } from '@/stores/userStore';

export const router = createRouter({
  routes: [
    {
      path: '/',
      component: () => import('../pages/Landing.vue'),
    },
    {
      path: '/login',
      component: () => import('../pages/Login.vue'),
    },
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
    { path: '/:pathMatch(.*)*', component: () => import('../pages/NotFound.vue') },
  ],
  history: createWebHistory(),
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
    let user = userStore.get();
    if (!user) await fetchAndSetUser();
    user = userStore.get();


    next();
    /* if (!user?.isTester) {
      alert('We are sorry, but currently bot access only via invites.');
      next('/');
    } else {
      next();
    } */

  } else {
    next();
  }
});