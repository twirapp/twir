import type { Component } from 'vue';
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

import {
  createUserDashboard,
  getProfile,
  redirectToLogin,
  selectedDashboardStore,
} from '@/services/auth';

export enum AppMenu {
  'dashboard',
  'commands',
  'timers',
  'moderation',
  'greetings',
  'keywords',
  'variables',
  'settings',
}

type AppMenuItems = Record<AppMenu, RouteRecordRaw>;

export const appMenu: AppMenuItems = {
  [AppMenu.dashboard]: {
    path: '/dashboard',
    component: () => import('@/components/app/pages/Dashboard.vue'),
  },
  [AppMenu.commands]: {
    path: '/commands',
    component: () => import('@/components/app/pages/Commands.vue'),
  },
  [AppMenu.greetings]: {
    path: '/greetings',
    component: () => import('@/components/app/pages/Greetings.vue'),
  },
  [AppMenu.keywords]: {
    path: '/keywords',
    component: () => import('@/components/app/pages/Keywords.vue'),
  },
  [AppMenu.moderation]: {
    path: '/moderation',
    component: () => import('@/components/app/pages/Moderation.vue'),
  },
  [AppMenu.timers]: {
    path: '/timers',
    component: () => import('@/components/app/pages/Timers.vue'),
  },
  [AppMenu.variables]: {
    path: '/variables',
    component: () => import('@/components/app/pages/Variables.vue'),
  },
  [AppMenu.settings]: {
    path: '/settings',
    component: () => import('@/components/app/pages/Settings.vue'),
  },
};

export function createAppRouter() {
  const router = createRouter({
    history: createWebHistory('app'),
    routes: [
      ...Object.values(appMenu),
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
