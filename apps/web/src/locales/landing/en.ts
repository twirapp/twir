import { BasicPlanFeatures, PlanId, ProPlanFeatures } from '@/data/pricingPlans.js';
import type ILandingLocale from '@/types/landingLocaleInterface.js';
import { NavMenuTabs } from '@/types/navMenu';

const messages: ILandingLocale = {
  navMenu: [
    { id: NavMenuTabs.features, name: 'Features' },
    { id: NavMenuTabs.reviews, name: 'Reviews' },
    { id: NavMenuTabs.team, name: 'Team' },
    { id: NavMenuTabs.pricing, name: 'Pricing' },
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
      plans: {
        [PlanId.basic]: {
          name: 'Basic plan',
          price: 0,
          features: {
            [BasicPlanFeatures.first]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [BasicPlanFeatures.second]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [BasicPlanFeatures.last]: {
              name: 'Unlimited commands',
              status: 'limited',
            },
          },
        },
        [PlanId.pro]: {
          name: 'Pro plan',
          price: 10,
          features: {
            [ProPlanFeatures.first]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [ProPlanFeatures.second]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
            [ProPlanFeatures.last]: {
              name: 'Unlimited commands',
              status: 'accessible',
            },
          },
        },
      },
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
