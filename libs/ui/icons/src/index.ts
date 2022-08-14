/// <reference types="vite-svg-loader" />

import { FunctionalComponent, SVGAttributes } from 'vue';

import ArrowLarge from '@/ArrowLargeIcon.svg?component';
import ArrowMedium from '@/ArrowMediumIcon.svg?component';
import Bell from '@/BellIcon.svg?component';
import Check from '@/CheckIcon.svg?component';
import CommandLine from '@/CommandLineIcon.svg?component';
import Danger from '@/DangerIcon.svg?component';
import Home from '@/HomeIcon.svg?component';
import Layout from '@/LayoutIcon.svg?component';
import SuccessCircle from '@/SuccessCircleIcon.svg?component';
import Users from '@/UsersIcon.svg?component';
import Warning from '@/WarningIcon.svg?component';

export type IconName =
  | 'ArrowLarge'
  | 'ArrowMedium'
  | 'Bell'
  | 'Check'
  | 'CommandLine'
  | 'Danger'
  | 'Home'
  | 'Layout'
  | 'SuccessCircle'
  | 'Users'
  | 'Warning';

type Icons = {
  [K in IconName]: FunctionalComponent<SVGAttributes>;
};

const icons: Icons = {
  ArrowLarge,
  ArrowMedium,
  Bell,
  Check,
  CommandLine,
  Danger,
  Home,
  Layout,
  SuccessCircle,
  Users,
  Warning,
};

export default icons;
