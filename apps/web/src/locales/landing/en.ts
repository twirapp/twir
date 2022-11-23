import { NavMenuTabs } from '@/data/landing/navMenu.js';
import { BasicPlanFeatures, PlanId, ProPlanFeatures } from '@/data/landing/pricingPlans.js';
import { TeamMemberId } from '@/data/landing/team.js';
import type ILandingLocale from '@/locales/landing/interface.js';

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
            'A powerful command system with aliases, counters, and all sorts of variables for all occasions',
        },
        {
          name: 'Moderation',
          description: `Not enough helpers to moderate the chat room? No problem. You'll find everything you need in our system, including quick nuke message deletion`,
        },
        {
          name: 'Timers',
          description:
            'A simple system, but with verve, has become a popular announcement system from Twitch',
        },
        {
          name: 'Greatings',
          description: 'Do you want to somehow highlight your favorite viewers? Add a greeting!',
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
        'The backend part was written entirely by Satont, early versions of the site too. Later Melkam joined us and drew a new, gorgeous design, and then brought our ideas to life.',
      founder: 'Founder',
      members: {
        [TeamMemberId.Satont]: 'Backend developer',
        [TeamMemberId.Melkam]: 'UI-UX Designer Frontend developer',
      },
    },
  },
};

export default messages;
