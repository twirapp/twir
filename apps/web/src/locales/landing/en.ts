import type ILandingLocale from '@/types/landingLocaleInterface.js';
import { NavMenuTabs } from '@/types/navMenu';

const messages: ILandingLocale = {
  navMenu: [
    { id: NavMenuTabs.features, name: 'Features' },
    { id: NavMenuTabs.pricing, name: 'Pricing' },
    { id: NavMenuTabs.reviews, name: 'Reviews' },
    { id: NavMenuTabs.team, name: 'Team' },
  ],
  buttons: {
    buyPlan: 'Buy plan',
    getStarted: 'Get started',
    learnMore: 'Learn more',
    login: 'Login',
    startForFree: 'Start for free',
    tryFeature: 'Try feature',
  },
  tagline: 'Created by streamers. Made for streamers. Used by streamers. For streamers with love.',
  sections: {
    features: {
      title: 'Bot features',
      featuresInDev: 'Features in development',
      content: [
        {
          name: 'Commands',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Moderation',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Timers',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
        {
          name: 'Greatings',
          description:
            'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
        },
      ],
    },
    firstScreen: {
      title: 'The <span>perfect bot</span> for your stream',
    },
    footer: {
      rights: '© Tsuwari {year}. All rights reserved.',
    },
    integrations: {
      preTitle: 'Integrations',
      title: 'Bot has a built-in API for the most necessary apps',
      description:
        'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    },
    pricing: {
      title: 'We’ve got a plan that’s perfect for you',
      features: 'Features',
      perMonth: 'per month',
      plans: [
        {
          id: 1,
          name: 'Basic plan',
          features: [
            { id: 1, name: 'Unlimited commands' },
            { id: 2, name: '24 hours access' },
            { id: 3, name: '5 integrations' },
            { id: 4, name: 'Unlimited commands' },
            { id: 5, name: 'Maximum 3 users' },
            { id: 6, name: 'Maximum 3 users' },
          ],
        },
        {
          id: 2,
          name: 'Pro plan',
          features: [
            { id: 1, name: 'Unlimited commands' },
            { id: 2, name: '24 hours access' },
            { id: 3, name: '5 integrations' },
            { id: 4, name: 'Unlimited commands' },
            { id: 5, name: 'Maximum 3 users' },
            { id: 6, name: 'Maximum 3 users' },
          ],
        },
      ],
    },
    reviews: {
      title: 'Reviews from streamers and other viewers',
    },
    statLine: {
      statPlaceholder: 'Aliquam nulla',
    },
    subscribeForUpdates: {
      title: 'Subscribe for updates',
      description:
        'Non rhoncus, neque arcu, commodo malesuada sed porttitor dictumst integer. Suscipit dictum quam ut blandit amet.',
      inputPlaceholder: 'Type your email',
    },
    team: {
      title: 'Our team',
      description:
        'Sed eget leo adipiscing lectus nunc laoreet. Scelerisque est justo, pellentesque ut eu sit in. Suspendisse venenatis, odio dui a. Vivamus in fames augue blandit ut non sagittis, sagittis, pretium. Mollis rhoncus, pretium, morbi',
      founder: 'Founder',
      members: [
        {
          id: 1,
          role: 'Backend developer',
        },
        {
          id: 2,
          role: 'Backend developer',
        },
        {
          id: 3,
          role: 'UI-UX Designer Frontend developer',
        },
      ],
    },
  },
};

export default messages;
