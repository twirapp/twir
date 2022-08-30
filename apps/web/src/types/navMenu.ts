export interface NavMenuItem {
  id: number;
  name: () => string;
  href: () => string;
}
