import type { IconName } from '@tsuwari/ui-icons';

import type { FeatureType, PlanColorThemes } from '@/types/pricingPlan.js';

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
