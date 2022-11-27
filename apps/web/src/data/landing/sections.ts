export enum LandingSection {
  'hero',
  'stats',
  'features',
  'integrations',
  'reviews',
  'team',
  'pricing',
}

export const navMenuLinks = [
  LandingSection.features,
  LandingSection.reviews,
  LandingSection.team,
  LandingSection.pricing,
] as const;

export type NavMenuLocale = {
  id: typeof navMenuLinks[number];
  name: string;
};
