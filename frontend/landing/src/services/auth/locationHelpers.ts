import { isClient } from '@vueuse/core';


import type { Locale } from '@/locales';
import { unprotectedApiClient } from '@/services/apiClients.js';

export const ORIGIN_STATE = isClient ? window.btoa(window.location.origin) : '';

export const redirectToLogin = async () => {
  if (isClient) {
		const res = await unprotectedApiClient.authGetLink({
			state: ORIGIN_STATE,
		});

    window.location.replace(res.response.link);
  }
};

export const redirectToDashboard = () => {
  window.location.replace(`/dashboard`);
};

export const redirectToLanding = (locale?: Locale) => {
  window.location.replace(locale ? `/${locale}` : '/');
};
