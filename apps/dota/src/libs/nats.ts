import { config } from '@tsuwari/config';
import * as Dota from '@tsuwari/nats/dota';
import { connect } from 'nats';

import { AppService } from '../app.service.js';
import { app } from '../index.js';

export const nats = await connect({
  servers: [config.NATS_URL],
});

const sub = nats.subscribe('dota.getProfileCard');

(async () => {
  for await (const m of sub) {
    const service = app.get(AppService);
    const data = Dota.GetPlayerCardRequest.fromBinary(m.data);
    const card = await service.getDotaProfileCard(data.accountId);

    const response = Dota.GetPlayerCardResponse.toBinary({
      accountId: data.accountId,
      leaderboardRank: card.leaderboard_rank,
      rankTier: card.rank_tier,
    });

    m.respond(response);
  }
})();
