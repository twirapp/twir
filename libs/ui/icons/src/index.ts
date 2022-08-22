/// <reference types="vite-svg-loader" />

import { FunctionalComponent, SVGAttributes } from 'vue';

import ArrowInCircle from '@/ArrowInCircleIcon.svg?component';
import ArrowLarge from '@/ArrowLargeIcon.svg?component';
import ArrowMedium from '@/ArrowMediumIcon.svg?component';
import ArrowNarrow from '@/ArrowNarrowIcon.svg?component';
import Bell from '@/BellIcon.svg?component';
import Check from '@/CheckIcon.svg?component';
import CommandLine from '@/CommandLineIcon.svg?component';
import Cross from '@/CrossIcon.svg?component';
import Danger from '@/DangerIcon.svg?component';
import Eye from '@/EyeIcon.svg?component';
import EyeOff from '@/EyeOffIcon.svg?component';
import Github from '@/GithubIcon.svg?component';
import Home from '@/HomeIcon.svg?component';
import Instagram from '@/InstagramIcon.svg?component';
import Key from '@/KeyIcon.svg?component';
import Layout from '@/LayoutIcon.svg?component';
import Menu from '@/MenuIcon.svg?component';
import Message from '@/MessageIcon.svg?component';
import Minus from '@/MinusIcon.svg?component';
import QuestionMark from '@/QuestionMarkIcon.svg?component';
import Selector from '@/SelectorIcon.svg?component';
import Star from '@/StarIcon.svg?component';
import SuccessCircle from '@/SuccessCircleIcon.svg?component';
import Sword from '@/SwordIcon.svg?component';
import Telegram from '@/TelegramIcon.svg?component';
import Timer from '@/TimerIcon.svg?component';
import Twitch from '@/TwitchIcon.svg?component';
import Users from '@/UsersIcon.svg?component';
import Variable from '@/VariableIcon.svg?component';
import Warning from '@/WarningIcon.svg?component';
import World from '@/WorldIcon.svg?component';

const icons = {
  Star,
  Twitch,
  ArrowNarrow,
  World,
  Telegram,
  Minus,
  Instagram,
  Github,
  Timer,
  Message,
  Variable,
  QuestionMark,
  Selector,
  Menu,
  Key,
  Cross,
  ArrowInCircle,
  Sword,
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
  Eye,
  EyeOff,
};

export type IconName = keyof typeof icons;

export default icons as {
  readonly [K in IconName]: FunctionalComponent<SVGAttributes>;
};
