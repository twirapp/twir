import { ChatBadge } from '@twir/frontend-chat';


export const globalBadges: Map<string, ChatBadge> = new Map([
	['broadcaster', {
		'set_id': 'broadcaster',
		'versions': [
			{
				'id': '1',
				'image_url_1x': 'https://static-cdn.jtvnw.net/badges/v1/5527c58c-fb7d-422d-b71b-f309dcb85cc1/1',
				'image_url_2x': 'https://static-cdn.jtvnw.net/badges/v1/5527c58c-fb7d-422d-b71b-f309dcb85cc1/2',
				'image_url_4x': 'https://static-cdn.jtvnw.net/badges/v1/5527c58c-fb7d-422d-b71b-f309dcb85cc1/3',
			},
		],
	}],
	['moderator', {
		'set_id': 'moderator',
    'versions': [
      {
        'id': '1',
        'image_url_1x': 'https://static-cdn.jtvnw.net/badges/v1/3267646d-33f0-4b17-b3df-f923a41db1d0/1',
        'image_url_2x': 'https://static-cdn.jtvnw.net/badges/v1/3267646d-33f0-4b17-b3df-f923a41db1d0/2',
        'image_url_4x': 'https://static-cdn.jtvnw.net/badges/v1/3267646d-33f0-4b17-b3df-f923a41db1d0/3',
      },
    ],
	}],
	['no_audio', {
    'set_id': 'no_audio',
    'versions': [
      {
        'id': '1',
        'image_url_1x': 'https://static-cdn.jtvnw.net/badges/v1/aef2cd08-f29b-45a1-8c12-d44d7fd5e6f0/1',
        'image_url_2x': 'https://static-cdn.jtvnw.net/badges/v1/aef2cd08-f29b-45a1-8c12-d44d7fd5e6f0/2',
        'image_url_4x': 'https://static-cdn.jtvnw.net/badges/v1/aef2cd08-f29b-45a1-8c12-d44d7fd5e6f0/3',
      },
    ],
  }],
	['vip', {
    'set_id': 'vip',
    'versions': [
      {
        'id': '1',
        'image_url_1x': 'https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/1',
        'image_url_2x': 'https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/2',
        'image_url_4x': 'https://static-cdn.jtvnw.net/badges/v1/b817aba4-fad8-49e2-b88a-7cc744dfa6ec/3',
      },
    ],
  }],
]);
