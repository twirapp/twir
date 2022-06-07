import { FaceitDbData } from '../../integrations/faceit.js';
import { Module, State } from '../index.js';

const getFaceitIntegrationData = async (state: State) => {
  const integrations = await state.cache.getEnabledIntegrations();
  const integration = integrations.find((i) => i.integration.service === 'FACEIT');

  if (!integration) return 'Integration with faceit not enabled';
  const integrationData = integration.data as unknown as FaceitDbData;

  const data = await state.cache.getFaceitData(integrationData.username, integrationData.game);

  return data;
};

export const faceit: Module[] = [
  {
    key: 'faceit.elo',
    description: 'Faceit elo',
    handler: async (_key, state) => {
      const data = await getFaceitIntegrationData(state);
      if (!data) {
        return 'Cannot fetch faceit data.';
      }
      if (typeof data === 'string') {
        return data;
      }

      return data.elo;
    },
  },
  {
    key: 'faceit.lvl',
    description: 'Faceit Lvl',
    handler: async (_key, state) => {
      const data = await getFaceitIntegrationData(state);
      if (!data) {
        return 'Cannot fetch faceit data.';
      }
      if (typeof data === 'string') {
        return data;
      }

      return data.lvl;
    },
  },
  {
    key: 'faceit.todayEloDiff',
    description: 'Faceit today elo earned',
    handler: async (_key, state) => {
      const data = await getFaceitIntegrationData(state);
      if (!data) {
        return 'Cannot fetch faceit data.';
      }
      if (typeof data === 'string') {
        return data;
      }

      return data.todayEloDiff;
    },
  },
  {
    key: 'faceit.latestMatchesTrend.simple',
    description: 'Faceit matches trend',
    handler: async (_key, state) => {
      const data = await getFaceitIntegrationData(state);
      if (!data) {
        return 'Cannot fetch faceit data.';
      }
      if (typeof data === 'string') {
        return data;
      }

      return data.latestMatchesTrend.simple;
    },
  },
  {
    key: 'faceit.latestMatchesTrend.extended',
    description: 'Faceit matches trend with elo diff',
    handler: async (_key, state) => {
      const data = await getFaceitIntegrationData(state);
      if (!data) {
        return 'Cannot fetch faceit data.';
      }
      if (typeof data === 'string') {
        return data;
      }

      return data.latestMatchesTrend.simple;
    },
  },
];
