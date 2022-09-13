export enum NavMenuTabs {
  'features',
  'reviews',
  'pricing',
  'team',
}

export type NavMenuHrefs = {
  [K in NavMenuTabs]: string;
};

export type NavMenuLocale = {
  id: NavMenuTabs;
  name: string;
};

export const navMenuHrefs: NavMenuHrefs = {
  [NavMenuTabs.features]: 'features',
  [NavMenuTabs.reviews]: 'reviews',
  [NavMenuTabs.pricing]: 'pricing',
  [NavMenuTabs.team]: 'team',
};
