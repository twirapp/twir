import { protectedApiClient, unprotectedApiClient } from '@/services/apiClients.js';

export const getProfile = async () => {
	const res = await protectedApiClient.userProfile({});

  return res.response;
};

export const authorizeByTwitchCode = async (code: string) => {
	await unprotectedApiClient.authPostCode({
		code,
	});
};

export const logout = async () => {
  return true;
};
