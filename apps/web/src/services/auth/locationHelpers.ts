import { isClient } from '@vueuse/core';

import { API_LOGIN_ROUTE } from './api.js';

import type { Locale } from '@/locales';

export const LOGIN_ROUTE_STATE = isClient ? window.btoa(window.location.origin + '/login') : '';
export const ORIGIN_STATE = isClient ? window.btoa(window.location.origin) : '';

export const redirectToLogin = () => {
  if (isClient) {
    window.location.replace(API_LOGIN_ROUTE);
  }
};

export const redirectToDashboard = () => {
  window.location.replace('/app/dashboard');
};

export const redirectToLanding = (locale?: Locale) => {
  window.location.replace(locale ? `/${locale}` : '/');
};
