import dotenv from 'dotenv';

import { Bot, BotType } from './dist/entities/Bot.js';
import { Token } from './dist/entities/Token.js';
import { AppDataSource } from './dist/index.js';

dotenv.config({ path: '../../.env' });

const typeorm = await AppDataSource.initialize();

const entity = await typeorm.getRepository(Bot).findOneBy({
  type: BotType.DEFAULT,
});

if (entity) {
  console.info('✅ Bot already exists, skipping...');
  process.exit(0);
}

const { BOT_ACCESS_TOKEN, BOT_REFRESH_TOKEN } = process.env;

if (!BOT_ACCESS_TOKEN || !BOT_REFRESH_TOKEN) {
  console.error('🚨 Missed bot access token or bot refresh token');
  process.exit(1);
}

const request = await fetch('https://id.twitch.tv/oauth2/validate', {
  headers: {
    Authorization: `OAuth ${BOT_ACCESS_TOKEN}`,
  },
});
const response = await request.json();

if (!request.ok) {
  console.error(response);
  process.exit(1);
}

const token = await typeorm.getRepository(Token).save({
  accessToken: BOT_ACCESS_TOKEN,
  refreshToken: BOT_REFRESH_TOKEN,
  expiresIn: response.expires_in,
  obtainmentTimestamp: new Date(),
});

await typeorm.getRepository(Bot).save({
  id: response.user_id,
  type: BotType.DEFAULT,
  tokenId: token.id,
});
