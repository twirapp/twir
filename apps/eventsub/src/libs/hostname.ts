import { config } from '@tsuwari/config';
import ngrok from 'ngrok';

export const getHostName = async () => {
  let hostName = '';

  if (config.isDev) {
    const tunnel = await ngrok.connect(3003);
    hostName = tunnel.replace('https://', '');
  } else {
    hostName = `eventsub.${config.HOSTNAME.replace('https://', '')}`;
  }

  return hostName;
};
