type ChatItemMessage = {
	type: 'message';
	sender: 'bot' | 'user';
	message: string;
}

type ChatItemRedemption = {
	type: 'redemption';
	title: string;
	input: string;
}

export type ChatItem = ChatItemMessage | ChatItemRedemption

export const messages: ChatItem[] = [
  {
    sender: 'bot',
    message: 'Message from timer: follow to my socials!',
    type: 'message',
  },
  {
    sender: 'user',
    message: '!title Playling League of Legends with my friend',
    type: 'message',
  },
  {
    sender: 'bot',
    message: '✅ Title succesfully changed.',
    type: 'message',
  },
  {
    type: 'redemption',
    title: '<b>melkam</b> activated channel reward: timeout chatter (1000 🪙)',
    input: 'Satont',
  },
  {
    sender: 'bot',
    message: 'melkam disabled chat for <b>Satont</b> for 5 minutes',
    type: 'message',
  },
  {
    sender: 'user',
    message: '!song',
    type: 'message',
  },
  {
    sender: 'bot',
    message: 'Linkin Park — Numb',
    type: 'message',
  },
  {
    sender: 'user',
    message: '!category LOL',
    type: 'message',
  },
  {
    sender: 'bot',
    message: '✅ Category changed to League of Legends.',
    type: 'message',
  },
];
