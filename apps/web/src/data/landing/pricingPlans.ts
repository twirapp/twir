import type { IconName } from '@tsuwari/ui-components';

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

export enum PlanId {
  basic,
  pro,
}

export enum BasicPlanFeatures {
  first,
  second,
  last,
}

export enum ProPlanFeatures {
  first,
  second,
  last,
}

export const planColorThemes: PlanColorThemes = {
  [PlanId.basic]: 'gray',
  [PlanId.pro]: 'purple',
};

export const featureTypeIcons: Record<FeatureType, IconName> = {
  accessible: 'Check',
  limited: 'Minus',
};
