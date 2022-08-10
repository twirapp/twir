import { Module } from '../index.js';

type FFZSet = {
  emoticons: Array<{ name: string }>;
};
type FFZResponse = {
  sets: Record<string, FFZSet>;
};

type BTTVEmote = {
  code: string;
};
type BTTVResponse = {
  sharedEmotes: Array<BTTVEmote>;
  channelEmotes: Array<BTTVEmote>;
};

export const emotes: Module[] = [
  {
    key: 'emotes.7tv',
    description: 'Emotes of channel from https://7tv.app',
    visible: true,
    async handler(_key, state) {
      if (!state.channelId) return;

      const request = await fetch(`https://api.7tv.app/v2/users/${state.channelId}/emotes`);
      const data: { name: string }[] = await request.json();

      return data.map((s) => s.name).join(' ');
    },
  },
  {
    key: 'emotes.ffz',
    description: 'Emotes of channel from https://frankerfacez.com',
    async handler(_key, state) {
      if (!state.channelId) return;

      const request = await fetch(`https://api.frankerfacez.com/v1/room/id/${state.channelId}`);
      const data: FFZResponse = await request.json();

      const emoticons = Object.entries(data.sets).reduce((acc, set) => {
        return [...acc, ...set[1].emoticons.map((e) => e.name)];
      }, [] as string[]);

      return emoticons.join(' ');
    },
  },
  {
    key: 'emotes.bttv',
    description: 'Emotes of channel from https://betterttv.com/',
    async handler(_key, state) {
      if (!state.channelId) return;

      const request = await fetch(
        `https://api.betterttv.net/3/cached/users/twitch/${state.channelId}`,
      );
      const data: BTTVResponse = await request.json();

      const emoticons = [...data.channelEmotes, ...data.sharedEmotes].reduce((acc, emote) => {
        return [...acc, emote.code];
      }, [] as string[]);

      return emoticons.join(' ');
    },
  },
];
