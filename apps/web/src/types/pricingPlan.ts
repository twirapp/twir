import type { BasicPlanFeatures, PlanId, ProPlanFeatures } from '@/data/pricingPlans.js';

export type FeatureType = 'accessible' | 'limited';

interface PlanFeatures {
  [PlanId.basic]: BasicPlanFeatures;
  [PlanId.pro]: ProPlanFeatures;
}

export type PricePlanLocale = PricePlansLocale[keyof PricePlansLocale];

export type PricePlansLocale = {
  [P in PlanId]: {
    name: string;
    price: number;
    features: {
      [F in PlanFeatures[P]]: {
        name: string;
        status: FeatureType;
      };
    };
  };
};

export type PlanColorTheme = 'purple' | 'gray';

export type PlanColorThemes = { [K in PlanId]: PlanColorTheme };
