import { ChildProcess, exec } from 'node:child_process';
import { readFile, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';

function updateTwitchCliConfig(proc: ChildProcess) {
  return new Promise<boolean>((resolve) => {
    if (proc.stdin === null || proc.stdout === null) {
      return resolve(false);
    }

    proc.stdin.write(
      `doppler run --command='twitch-cli configure -i $TWITCH_CLIENTID -s $TWITCH_CLIENTSECRET'\n`,
    );

    proc.stdout.on('error', (error) => {
      console.error(error);
      resolve(false);
    });
    proc.stdout.on('data', (message) => {
      if (typeof message === 'string' && message.trim() === 'Updated configuration.') {
        return resolve(true);
      }
      resolve(false);
    });
  });
}

type TwitchTokens = { BOT_ACCESS_TOKEN: string; BOT_REFRESH_TOKEN: string };

function getTwitchTokens(proc: ChildProcess) {
  return new Promise<TwitchTokens>((resolve, reject) => {
    if (!proc.stdin || !proc.stdout || !proc.stderr) {
      return reject('Cannot get stdin or stdout from process');
    }

    const tokenRegex = /Token:\s([0-9a-z]+)/gi;

    proc.stdin.write(
      `doppler run --command='twitch-cli token --user-token --scopes="${scopes}"'\n`,
    );

    proc.stdout.on('data', (data) => {
      console.log(data);
    });
    proc.stderr.on('data', (data) => {
      if (typeof data !== 'string') return;

      const result = data.match(tokenRegex);
      if (result === null) {
        return console.log(data);
      }
      if (result.length !== 2) {
        return reject('Cannot find two tokens');
      }

      return resolve({
        BOT_ACCESS_TOKEN: result[0].split(':')[1].trim(),
        BOT_REFRESH_TOKEN: result[1].split(':')[1].trim(),
      });
    });
  });
}

async function writeTokensToEnv(tokens: TwitchTokens) {
  const envPath = resolve(process.cwd(), '.env');
  const data = [
    `BOT_ACCESS_TOKEN=${tokens.BOT_ACCESS_TOKEN}`,
    `BOT_REFRESH_TOKEN=${tokens.BOT_REFRESH_TOKEN}`,
  ].join('\n');

  try {
    let content = await readFile(envPath, 'utf-8');
    content = content.trim();
    if (content === '') {
      return await writeToFile(data, true);
    }

    const tokensRegex = /BOT_.+_TOKEN *= *([0-9a-z]+)/gi;
    const tokensRegexResult = content.match(tokensRegex);
    console.log(tokensRegexResult);
    console.log(content);

    if (tokensRegexResult === null) {
      content += '\n' + data;
      return await writeToFile(content, true);
    }

    tokensRegexResult.forEach((res) => {
      content = content.replace(res, '');
    });
    content += '\n' + data;
    return await writeToFile(content.trim(), true);
  } catch (e) {
    await writeToFile(data, false);
  }

  async function writeToFile(content: string, exists = false) {
    return await writeFile(envPath, content, { flag: !exists ? 'wx' : undefined });
  }
}

(async () => {
  const childProcess = exec('bash');

  const isConfigured = await updateTwitchCliConfig(childProcess);
  if (!isConfigured) throw new Error('Cannot configure twitch cli');
  const tokens = await getTwitchTokens(childProcess);
  await writeTokensToEnv(tokens);

  childProcess.kill(0);
  process.exit(0);
})();

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
