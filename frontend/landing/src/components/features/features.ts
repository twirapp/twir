import CommandsSvg from '@/assets/features/commands.svg?use';
import ConnectionSvg from '@/assets/features/connection.svg?use';
import EventsSvg from '@/assets/features/events.svg?use';
import GreetingsSvg from '@/assets/features/greetings.svg?use';
import KeywordsSvg from '@/assets/features/keywords.svg?use';
import ModerationSvg from '@/assets/features/moderation.svg?use';
import SongRequestsSvg from '@/assets/features/song-requests.svg?use';
import StatsSvg from '@/assets/features/stats.svg?use';
import TimersSvg from '@/assets/features/timers.svg?use';

interface Feature {
	title: string;
	description: string;
	icon: any;
}

export const features: Feature[] = [
  {
    title: 'Commands',
    description:
      'A powerful command system with aliases, counters, and all sorts of variables for all occasions',
    icon: CommandsSvg,
  },
  {
    title: 'Timers',
    description:
      'A simple system, but with verve, has become a popular announcement system from Twitch',
    icon: TimersSvg,
  },
  {
    title: 'Greetings',
    description: 'Do you want to somehow highlight your favorite viewers? Add a greeting!',
    icon: GreetingsSvg,
  },
  {
    title: 'Song requests',
    description:
      'Viewers request songs via chat commands. Bot manages queue, plays songs, and offers controls. Enhances stream with viewer engagement',
    icon: SongRequestsSvg,
  },
  {
    title: 'Keywords',
    description:
      'Identifies specified keywords in chat, triggers automated messages for engagement or information. Enhances interaction and delivers targeted content during live stream',
    icon: KeywordsSvg,
  },
  {
    title: 'Events',
    description:
      'With this powerful tool, you can easily set up customized listeners to keep track of specific events happening in the chat, or even trigger actions based on system events',
    icon: EventsSvg,
  },
	{
		title: 'Moderation',
		description: 'Create and manage chat filters to keep safe and kind communication',
		icon: ModerationSvg,
	},
  {
    title: 'OBS Websockets',
    description:
      'Highly integrate with obs studio via websockets. Change scenes, mute audio, toggle source visibility via bot',
    icon: ConnectionSvg,
  },
	{
		title: 'Stats tracking',
		description: 'Track users watch time, messages, used channel points',
		icon: StatsSvg,
	},
];
