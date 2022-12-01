import 'reflect-metadata';

import { config } from '@tsuwari/config';
import * as DotaServer from '@tsuwari/grpc/generated/dota/dota';
import { PORTS } from '@tsuwari/grpc/servers/constants';
import { ChannelDotaAccount } from '@tsuwari/typeorm/entities/ChannelDotaAccount';
import _ from 'lodash';
import { createServer } from 'nice-grpc';

import { Dota } from './libs/dota.js';
import { typeorm } from './libs/typeorm.js';

const dota = await new Dota().init();

const dotaServer: DotaServer.DotaServiceImplementation = {
  async getPlayerCard(request: DotaServer.GetPlayerCardRequest) {
    const result = await dota.getDotaProfileCard(request.accountId);
    return result;
  },
};

const server = createServer();

server.add(DotaServer.DotaDefinition, dotaServer);

await server.listen(`0.0.0.0:${PORTS.DOTA_SERVER_PORT}`);

setInterval(
  async () => {
    const accounts = await typeorm.getRepository(ChannelDotaAccount).findBy({
      channel: {
        isEnabled: true,
      },
    });

    const chunks = _.chunk(
      accounts.map((a) => a.id),
      50,
    );

    for (const chunk of chunks) {
      dota.getPresences(chunk);
    }
  },
  config.isDev ? 5000 : 1 * 60 * 1000,
);
