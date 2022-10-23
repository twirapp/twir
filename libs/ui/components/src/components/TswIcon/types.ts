import * as icons from './icons';

export type IconName = keyof typeof icons;

export interface Icon {
  style: 'solid' | 'outline';
  width: number;
  height: number;
  path: { d: string }[];
}
