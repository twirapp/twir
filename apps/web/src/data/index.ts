import TextAvatarUrl from '@/assets/avatar.png';
import type { BotFeature } from '@/types/botFeatures.js';
import { NavMenuHrefs, NavMenuTabs } from '@/types/navMenu.js';
import { FeatureType, PlanColorThemes, PlanId, PricePlan } from '@/types/pricingPlan.js';
import type { Review } from '@/types/review.js';
import type { SocialMediaItem } from '@/types/socialMedia.js';
import type { StatInfo } from '@/types/statsLine.js';
import type { TeamMember } from '@/types/teamMember.js';

export const socials: SocialMediaItem[] = [
  { id: 1, media: 'Telegram', href: '#' },
  { id: 2, media: 'Instagram', href: '#' },
];

export const navMenuHrefs: NavMenuHrefs = {
  [NavMenuTabs.features]: 'features',
  [NavMenuTabs.reviews]: 'reviews',
  [NavMenuTabs.pricing]: 'pricing',
  [NavMenuTabs.team]: 'team',
};

export const reviews: Review[] = [
  {
    id: 1,
    username: 'random_usergsdagdsagsadgsda',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 2,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 3,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 2,
  },
  {
    id: 4,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 1,
  },
  {
    id: 5,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 6,
    username: 'random_user',
    comment:
      'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
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
