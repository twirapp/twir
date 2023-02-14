import { EventType } from '@tsuwari/typeorm/entities/events/Event';
import { OperationType } from '@tsuwari/typeorm/entities/events/EventOperation';

export const operationMapping: Record<keyof typeof OperationType, {
  description: string
  haveInput?: boolean,
  additionalValues?: Array<string>,
  producedVariables?: Array<string>,
  dependsOnEvents?: Array<EventType>
}> = {
  SEND_MESSAGE: {
    description: 'Send message in chat',
    haveInput: true,
    additionalValues: ['useAnnounce'],
  },
  BAN: {
    'description': 'Ban user',
  },
  UNBAN: {
    description: 'Unban user',
  },
  BAN_RANDOM: {
    description: 'Ban random online user',
    producedVariables: ['bannedUserName'],
  },
  VIP: {
    description: 'Vip user',
  },
  UNVIP: {
    description: 'Unvip user',
  },
  UNVIP_RANDOM: {
    description: 'Unvip random user',
    producedVariables: ['unvipedUserName'],
  },
  MOD: {
    description: 'Give user moderation',
  },
  UNMOD: {
    description: 'Remove moderation from user',
  },
  UNMOD_RANDOM: {
    description: 'Remove moderation from random user',
    producedVariables: ['unmodedUserName'],
  },
  CHANGE_TITLE: {
    description: 'Change title of stream',
    haveInput: true,
  },
  CHANGE_CATEGORY: {
    description: 'Change category of stream',
    haveInput: true,
  },
  FULFILL_REDEMPTION: {
    description: 'Verify fulfillment of the reward',
  },
  CANCEL_REDEMPTION: {
    description: 'Cancel reward and back points to user',
  },
  ENABLE_SUBMODE: {
    description: 'Enable submode',
  },
  DISABLE_SUBMODE: {
    description: 'Disable submode',
  },
  ENABLE_EMOTEONLY: {
    description: 'Enable emoteonly',
  },
  DISABLE_EMOTEONLY: {
    description: 'Disable emoteonly',
  },
  CREATE_GREETING: {
    description: 'Create greeting for user. Available only for rewards event, and requires user input.',
    dependsOnEvents: [EventType.REDEMPTION_CREATED],
  },
  TIMEOUT: {
    description: 'Timeout user',
    haveInput: true,
    // dependsOnEvents: [EventType.REDEMPTION_CREATED],
    additionalValues: ['timeoutTime'],
  },
  TIMEOUT_RANDOM: {
    description: 'Timeout random online user',
    producedVariables: ['bannedUserName'],
    additionalValues: ['timeoutTime'],
  },
};