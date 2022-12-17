import {
  TablerIcon,
  IconLink,
  IconLetterCaseUpper,
  IconMoodSmile,
  IconTextWrapDisabled,
  IconPlaylistX,
  IconLambda,
} from '@tabler/icons';
import type { SettingsType } from '@tsuwari/typeorm/entities/ChannelModerationSetting';

export const typesMapping: Record<
  keyof typeof SettingsType,
  {
    icon: TablerIcon;
    iconColor?: string;
    name?: string;
    description: string;
  }
> = {
  links: {
    icon: IconLink,
    iconColor: 'cyan',
    description: `Remove messages containing any links you haven't whitelisted.`,
  },
  caps: {
    icon: IconLetterCaseUpper,
    iconColor: 'orange',
    description: `Remove messages containing excessive amounts of capital letters.`,
  },
  emotes: {
    icon: IconMoodSmile,
    iconColor: 'yellow',
    description: 'Remove messages containing an excessive amount of emotes.',
  },
  longMessage: {
    icon: IconTextWrapDisabled,
    name: 'Long Messages',
    description: `Remove lengthy messages.`,
  },
  blacklists: {
    icon: IconPlaylistX,
    name: 'Deny List',
    description: 'Remove denied words from chat.',
  },
  symbols: {
    icon: IconLambda,
    iconColor: 'green',
    description: `Remove messages containing disruptive or excessive use of symbols.`,
  },
};
