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
