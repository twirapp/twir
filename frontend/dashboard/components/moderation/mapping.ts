import {
  TablerIcon,
  IconLink,
  IconLetterCaseUpper,
  IconMoodSmile,
  IconTextWrapDisabled,
  IconPlaylistX,
  IconLambda,
} from '@tabler/icons';
import type { SettingsType } from '@twir/typeorm/entities/ChannelModerationSetting';

export const typesMapping: Record<
  keyof typeof SettingsType,
  {
    icon: TablerIcon;
    iconColor?: string;
    name?: string;
  }
> = {
  links: {
    icon: IconLink,
    iconColor: 'cyan',
  },
  caps: {
    icon: IconLetterCaseUpper,
    iconColor: 'orange',
  },
  emotes: {
    icon: IconMoodSmile,
    iconColor: 'yellow',
  },
  longMessage: {
    icon: IconTextWrapDisabled,
    name: 'Long Messages',
  },
  blacklists: {
    icon: IconPlaylistX,
    name: 'Deny List',
  },
  symbols: {
    icon: IconLambda,
    iconColor: 'green',
  },
};
