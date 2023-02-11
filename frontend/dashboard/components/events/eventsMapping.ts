import { EventType } from '@tsuwari/typeorm/entities/events/Event';

export const eventsMapping: Record<keyof typeof EventType, {
  description: string,
  availableVariables: Array<[string, string]>
}> = {
  FOLLOW: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
    ],
  },
  SUBSCRIBE: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
      ['subLevel', 'Level of subscription (1/2/3/prime)'],
    ],
  },
  RESUBSCRIBE: {
    description: '',
    availableVariables: [
      ['userName', 'Username of user who subscribed channel'],
      ['userDisplayName', 'Username of user who subscribed channel with case'],
      ['subLevel', 'Level of subscription (1/2/3/prime)'],
      ['resubMonths', 'How many months user subscribed'],
      ['resubStreak', 'Cumulative streak of months of subscription'],
      ['resubMessage', 'Message of re-subscription'],
    ],
  },
  SUB_GIFT: {
    description: '',
    availableVariables: [
      ['senderUserName', 'Username of user who gifted subscription'],
      ['senderDisplayName', 'Username of user who gifted subscription with case'],
      ['targetUserName', 'Username of user who receive subscription'],
      ['targetDisplayName', 'Username of user who receive subscription with case'],
      ['subLevel', 'Level of subscription (1/2/3/prime)'],
    ],
  },
  REDEMPTION_CREATED: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
      ['rewardName', 'Name of reward'],
      ['rewardCost', 'Cost of reward'],
      ['rewardInput', 'User input of reward'],
    ],
  },
  COMMAND_USED: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
      ['commandName', 'Command name'],
    ],
  },
  FIRST_USER_MESSAGE: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
    ],
  },
  RAIDED: {
    description: '',
    availableVariables: [
      ['userName', 'Triggered user name'],
      ['userDisplayName', 'Triggered user name with case'],
      ['raidViewers', 'Count of viewers'],
    ],
  },
  TITLE_OR_CATEGORY_CHANGED: {
    description: '',
    availableVariables: [
      ['oldTitle', 'Old title of stream'],
      ['newTitle', 'New title of stream'],
      ['newCategory', 'Old category of stream'],
      ['newCategory', 'Old category of stream'],
    ],
  },
  STREAM_ONLINE: {
    description: '',
    availableVariables: [
      ['streamTitle', 'Title of stream'],
      ['streamCategory', 'Category of stream'],
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
      ['userName', 'Who donated'],
      ['donateAmount', 'Donate amount'],
      ['donateCurrency', 'Donate currency'],
      ['donateMessage', 'Donate message'],
    ],
  },
};