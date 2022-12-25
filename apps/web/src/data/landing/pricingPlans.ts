// List of pricing plans
export enum PricingPlans {
  basic,
  pro,
}

/**
 * Enums for plan features are created to index a unique key for each
 * feature. We need it for mapping translations with right feature.
 */
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

// Mapping of plan enum with plan features
interface PlanFeatures {
  [PricingPlans.basic]: BasicPlanFeatures;
  [PricingPlans.pro]: ProPlanFeatures;
}

type PlanFeature = { name: string; isAvaible: boolean };

// General info it's information that will not change depends on translation
export type PlanGeneral<Plan extends PricingPlans> = {
  [Feature in PlanFeatures[Plan]]: Pick<PlanFeature, 'isAvaible'>;
};

type PlanFeaturesGeneral = {
  [Plan in PricingPlans]: PlanGeneral<Plan>;
};

export const planFeaturesGeneral: PlanFeaturesGeneral = {
  [PricingPlans.basic]: {
    [BasicPlanFeatures.first]: {
      isAvaible: true,
    },
    [BasicPlanFeatures.second]: {
      isAvaible: true,
    },
    [BasicPlanFeatures.last]: {
      isAvaible: false,
    },
  },
  [PricingPlans.pro]: {
    [ProPlanFeatures.first]: {
      isAvaible: true,
    },
    [ProPlanFeatures.second]: {
      isAvaible: true,
    },
    [ProPlanFeatures.last]: {
      isAvaible: true,
    },
  },
};

// Represents a single plan type for translation
export type PricingPlanLocale<Plan extends PricingPlans> = {
  name: string;
  price: number;
  features: Record<PlanFeatures[Plan], Pick<PlanFeature, 'name'>>;
};

// Translate object of all plans.
export type PricingPlansLocale = {
  [Plan in PricingPlans]: PricingPlanLocale<Plan>;
};

// Every plan has own color theme
export type PlanColorTheme = 'purple' | 'gray';

export const planColorThemes: Record<PricingPlans, PlanColorTheme> = {
  [PricingPlans.basic]: 'gray',
  [PricingPlans.pro]: 'purple',
};
