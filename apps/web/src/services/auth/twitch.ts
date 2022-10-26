import { authorizeByTwitchCode } from './api.js';

export const handleTwitchLoginCallback = async () => {
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');
  
  if (params.get('error') !== null || !code) {
    throw new Error('Cannot find code or got an error from twitch!');
  }

  return await authorizeByTwitchCode(code);
};
