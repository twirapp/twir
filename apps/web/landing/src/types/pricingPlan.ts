export enum PlanId {
  basic,
  pro,
}

export enum FeatureType {
  accessibly,
  limited,
}

export interface PricePlanFeature {
  status: FeatureType;
  feature: string;
}

export interface PricePlan {
  id: PlanId;
  name: string;
  price: number;
  features: PricePlanFeature[];
}

export type PlanColorTheme = 'purple' | 'gray';

export type PlanColorThemes = { [K in PlanId]: PlanColorTheme };
