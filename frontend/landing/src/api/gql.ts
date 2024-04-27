import { config } from '@twir/config';

const host = config.SITE_BASE_URL ?? 'localhost:3005';
const isDev = config.NODE_ENV === 'development';
const apiAddr = isDev ? `${host}/api-new` : 'api-gql:3009';

export const gqlUrl = `http://${apiAddr}/query`;
