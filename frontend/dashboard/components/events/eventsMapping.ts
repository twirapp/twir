import { EventType } from '@tsuwari/typeorm/entities/events/Event';

export const eventsMapping: Record<keyof typeof EventType, {
  description: string,
  availableVariables: Array<string>
}> = {
  FOLLOW: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
    ],
  },
  SUBSCRIBE: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
      'subLevel',
    ],
  },
  RESUBSCRIBE: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
      'subLevel',
      'resubMonths',
      'resubStreak',
      'resubMessage',
    ],
  },
  SUB_GIFT: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
      'targetUserName',
      'targetDisplayName',
      'subLevel',
    ],
  },
  REDEMPTION_CREATED: {
    description: 'Reward activated',
    availableVariables: [
      'userName',
      'userDisplayName',
      'rewardName',
      'rewardCost',
      'rewardInput',
    ],
  },
  COMMAND_USED: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
      'commandName',
    ],
  },
  FIRST_USER_MESSAGE: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
    ],
  },
  RAIDED: {
    description: '',
    availableVariables: [
      'userName',
      'userDisplayName',
      'raidViewers',
    ],
  },
  TITLE_OR_CATEGORY_CHANGED: {
    description: '',
    availableVariables: [
      'oldStreamTitle',
      'newStreamTitle',
      'newStreamCategory',
      'newStreamCategory',
    ],
  },
  STREAM_ONLINE: {
    description: '',
    availableVariables: [
      'streamTitle',
      'streamCategory',
    ],
  },
  STREAM_OFFLINE: {
    description: '',
    availableVariables: [],
  },
  ON_CHAT_CLEAR: {
    description: '',
    availableVariables: [],
  },
  DONATE: {
    description: '',
    availableVariables: [
      'userName',
      'donateAmount',
      'donateCurrency',
      'donateMessage',
    ],
  },
};