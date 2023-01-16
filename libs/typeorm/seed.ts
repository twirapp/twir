import { resolve } from 'path';

import { encrypt } from '@tsuwari/crypto';
import * as dotenv from 'dotenv';

dotenv.config({ path: resolve(process.cwd(), '../../.env') });

import { AppDataSource } from './src';
import { Bot, BotType } from './src/entities/Bot';
import { Token } from './src/entities/Token';

import { config } from '@tsuwari/config';

async function bootstrap() {
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
    accessToken: encrypt(BOT_ACCESS_TOKEN, config.TOKENS_CIPHER_KEY!),
    refreshToken: encrypt(BOT_REFRESH_TOKEN, config.TOKENS_CIPHER_KEY!),
    expiresIn: response.expires_in,
    obtainmentTimestamp: new Date(),
    scopes: response.scopes,
  });

  await typeorm.getRepository(Bot).save({
    id: response.user_id,
    type: BotType.DEFAULT,
    tokenId: token.id,
  });
}

bootstrap();
