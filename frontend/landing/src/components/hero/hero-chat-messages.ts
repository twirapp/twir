export type ChatMessage = {
	type: 'message' | 'redemption';
	sender: 'bot' | 'user';
	text: string;
	user?: string
}

export const initialChatMessages: ChatMessage[] = [
	{
    sender: 'user',
    text: 'Hello, World!',
    type: 'message',
  },
  {
    sender: 'bot',
    text: 'Message from timer: follow to my socials!',
    type: 'message',
  },
  {
    sender: 'user',
    text: '!title Playling League of Legends with my friend',
    type: 'message',
  },
  {
    sender: 'bot',
    text: '✅ Title succesfully changed.',
    type: 'message',
  },
  {
    type: 'redemption',
		sender: 'user',
    text: '<b>melkam</b> activated channel reward: timeout chatter (1000 🪙)',
    user: 'Satont',
  },
  {
    sender: 'bot',
    text: 'melkam disabled chat for <b>Satont</b> for 5 minutes',
    type: 'message',
  },
  {
    sender: 'user',
    text: '!song',
    type: 'message',
  },
  {
    sender: 'bot',
    text: 'Linkin Park — Numb',
    type: 'message',
  },
  {
    sender: 'user',
    text: '!category LOL',
    type: 'message',
  },
  {
    sender: 'bot',
    text: '✅ Category changed to League of Legends.',
    type: 'message',
  },
];
