import type { BotFeature } from '@/types/botFeatures.js';
import type { NavMenuLocale } from '@/types/navMenu.js';
import type { PricePlansLocale } from '@/types/pricingPlan.js';

interface ILandingLocale {
  navMenu: NavMenuLocale[];
  tagline: string;
  buttons: {
    startForFree: string;
    getStarted: string;
    buyPlan: string;
    learnMore: string;
    login: string;
    tryFeature: string;
  };
  sections: {
    firstScreen: {
      title: string;
    };
    statLine: {
      statPlaceholder: string;
    };
    features: {
      title: string;
      featuresInDev: string;
      content: BotFeature[];
    };
    integrations: {
      preTitle: string;
      title: string;
      description: string;
    };
    reviews: {
      title: string;
    };
    team: {
      title: string;
      description: string;
      founder: string;
      members: {
        id: number;
        role: string;
      }[];
    };
    pricing: {
      title: string;
      features: string;
      perMonth: string;
      plans: PricePlansLocale;
    };
    subscribeForUpdates: {
      title: string;
      description: string;
      inputPlaceholder: string;
    };
    footer: {
      rights: string;
    };
  };
}

export default ILandingLocale;
