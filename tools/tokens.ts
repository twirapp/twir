import { exec } from 'node:child_process';
import { readFile, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';
import { promisify } from 'node:util';

const promisedExec = promisify(exec);

const scopes = [
  'analytics:read:extensions',
  'analytics:read:games',
  'bits:read',
  'channel:edit:commercial',
  'channel:manage:broadcast',
  'channel:read:charity',
  'channel:manage:extensions',
  'channel:manage:moderators',
  'channel:manage:polls',
  'channel:manage:predictions',
  'channel:manage:raids',
  'channel:manage:redemptions',
  'channel:manage:schedule',
  'channel:manage:videos',
  'channel:read:editors',
  'channel:read:goals',
  'channel:read:hype_train',
  'channel:read:polls',
  'channel:read:predictions',
  'channel:read:redemptions',
  'channel:read:stream_key',
  'channel:read:subscriptions',
  'channel:read:vips',
  'channel:manage:vips',
  'clips:edit',
  'moderation:read',
  'moderator:manage:announcements',
  'moderator:manage:automod',
  'moderator:read:automod_settings',
  'moderator:manage:automod_settings',
  'moderator:manage:banned_users',
  'moderator:read:blocked_terms',
  'moderator:manage:blocked_terms',
  'moderator:manage:chat_messages',
  'moderator:read:chat_settings',
  'moderator:manage:chat_settings',
  'user:edit',
  'user:edit:follows',
  'user:manage:blocked_users',
  'user:read:blocked_users',
  'user:read:broadcast',
  'user:manage:chat_color',
  'user:read:email',
  'user:read:follows',
  'user:read:subscriptions',
  'user:manage:whispers',
  'channel:moderate',
  'chat:edit',
  'chat:read',
  'whispers:read',
  'whispers:edit',
].join(' ');

const result = await promisedExec(
  `doppler run --command='twitch-cli token --user-token --scopes "${scopes}" --client-id $TWITCH_CLIENTID'`,
);

const regex = /([\sa-z]+)token:\s([0-9a-z]+)/gi;
const firstregexResult = regex.exec(result.stderr);
if (firstregexResult == null) {
  console.error(result.stderr);
  process.exit(1);
}
const secondRegexResult = regex.exec(result.stderr);

const refreshToken = firstregexResult[2];
const accessToken = secondRegexResult![2];

const envPath = resolve(process.cwd(), '.env');

const data = [`BOT_ACCESS_TOKEN=${accessToken}`, `BOT_REFRESH_TOKEN=${refreshToken}`].join('\n');

const tokensRegex = new RegExp('.+TOKEN=([0-9a-z]+)', 'ig');
try {
  let content = await readFile(envPath, 'utf-8');
  const tokensRegexResult = tokensRegex.exec(content);
  if (tokensRegexResult != null) {
    content = content.replaceAll(tokensRegex, '');
  }
  await writeToFile(`${content.replaceAll(/^\s*\n/gm, '')}\n${data}`, true);
} catch {
  await writeToFile(data);
}

function writeToFile(content: string, exists = false) {
  return writeFile(envPath, content, { flag: !exists ? 'wx' : undefined });
}
