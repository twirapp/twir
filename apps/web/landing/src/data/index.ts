import type { BotFeature } from '@/types/botFeatures.js';
import type { NavMenuItem } from '@/types/navMenu.js';
import { FeatureType, PlanColorThemes, PlanId, PricePlan } from '@/types/pricingPlan.js';
import type { SocialMediaItem } from '@/types/socialMedia.js';
import type { StatInfo } from '@/types/statsLine.js';
import type { TeamMember } from '@/types/teamMember.js';

export const navMenuItems: NavMenuItem[] = [
  { id: 1, name: 'Features', href: '#' },
  { id: 2, name: 'Reviews', href: '#' },
  { id: 3, name: 'Pricing', href: '#' },
  { id: 4, name: 'Team', href: '#' },
];

export const socials: SocialMediaItem[] = [
  { id: 1, media: 'Telegram', href: '#' },
  { id: 2, media: 'Instagram', href: '#' },
];

export const features: BotFeature[] = [
  {
    id: 1,
    title: 'Commands',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 2,
    title: 'Moderation',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 3,
    title: 'Timers',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 4,
    title: 'Greatings',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
];

export const pricePlans: PricePlan[] = [
  {
    id: PlanId.basic,
    name: 'Basic plan',
    price: 0,
    features: [
      { status: FeatureType.accessibly, feature: 'Unlimited commands' },
      { status: FeatureType.accessibly, feature: '24 hours access' },
      { status: FeatureType.accessibly, feature: '5 integrations' },
      { status: FeatureType.accessibly, feature: 'Unlimited commands' },
      { status: FeatureType.limited, feature: 'Maximum 3 users' },
      { status: FeatureType.limited, feature: 'Maximum 3 users' },
    ],
  },
  {
    id: PlanId.pro,
    name: 'Pro plan',
    price: 10,
    features: [
      { status: FeatureType.accessibly, feature: 'Unlimited commands' },
      { status: FeatureType.accessibly, feature: '24 hours access' },
      { status: FeatureType.accessibly, feature: '5 integrations' },
      { status: FeatureType.accessibly, feature: 'Unlimited commands' },
      { status: FeatureType.accessibly, feature: 'Maximum 3 users' },
      { status: FeatureType.accessibly, feature: 'Maximum 3 users' },
    ],
  },
];

export const planColorThemes: PlanColorThemes = {
  [PlanId.basic]: 'gray',
  [PlanId.pro]: 'purple',
};

export const teamMembers: TeamMember[] = [
  {
    id: 1,
    name: 'Satont',
    role: 'Backend developer',
    isFounder: true,
    socials: [
      { media: 'Github', link: 'https://github.com/satont' },
      { media: 'Telegram', link: 'https://t.me/satont' },
      { media: 'Twitch', link: 'https://www.twitch.tv/sadisnamenya' },
      { media: 'Website', link: 'https://satont.dev/' },
    ],
  },
  {
    id: 2,
    name: 'LwGerry',
    role: 'Backend developer',
    socials: [
      { media: 'Telegram', link: 'https://t.me/LWJerri' },
      { media: 'Website', link: 'https://lwjerri.js.org/' },
      { media: 'Github', link: 'https://github.com/LWJerri' },
      { media: 'Twitch', link: 'https://www.twitch.tv/lwgerry' },
    ],
  },
  {
    id: 3,
    name: 'Melkam',
    role: 'UI-UX Designer Frontend developer',
    socials: [
      { media: 'Github', link: 'https://github.com/MellKam' },
      { media: 'Instagram', link: 'https://www.instagram.com/mel._.kam/' },
      { media: 'Telegram', link: 'https://t.me/mellkam' },
    ],
  },
];

export const stats: StatInfo[] = [
  { id: 1, stat: '1053+', description: 'Streamers' },
  { id: 2, stat: '51251', description: 'Reviews' },
  { id: 3, stat: '50%', description: 'Aliquam nulla' },
  { id: 4, stat: '64531', description: 'Commads created' },
];
