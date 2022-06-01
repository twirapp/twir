import { config } from '@tsuwari/config';
import { ApiClient } from '@twurple/api';
import { ClientCredentialsAuthProvider } from '@twurple/auth';


const staticProvider = new ClientCredentialsAuthProvider(config.TWITCH_CLIENTID, config.TWITCH_CLIENTSECRET);
export const staticApi = new ApiClient({ authProvider: staticProvider });