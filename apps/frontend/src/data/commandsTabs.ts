import { mdiLeadPencil, mdiFormatListBulleted, mdiSword } from '@mdi/js';
import { type CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';

export const tabs: Array<{
  name: string,
  value: keyof typeof CommandModule,
  icon?: string,
}> = [
  {
    name: 'Custom',
    value: 'CUSTOM',
    icon: mdiLeadPencil,
  },
  {
    name: 'Manage Commands',
    value: 'MANAGE',
    icon: mdiFormatListBulleted,
  },
  {
    name: 'Moderation',
    value: 'MODERATION',
    icon: mdiSword,
  },
  {
    name: 'Dota 2',
    value: 'DOTA',
  },
];