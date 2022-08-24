import { FunctionalComponent, SVGAttributes } from 'vue';

import ArrowInCircle from '@/ArrowInCircle.svg?component';
import ArrowNarrow from '@/ArrowNarrow.svg?component';
import ArrowNarrowShort from '@/ArrowNarrowShort.svg?component';
import ArrowTriangleLarge from '@/ArrowTriangleLarge.svg?component';
import ArrowTriangleMedium from '@/ArrowTriangleMedium.svg?component';
import Bell from '@/Bell.svg?component';
import Check from '@/Check.svg?component';
import CommandLine from '@/CommandLine.svg?component';
import Cross from '@/Cross.svg?component';
import Danger from '@/Danger.svg?component';
import Eye from '@/Eye.svg?component';
import EyeOff from '@/EyeOff.svg?component';
import Github from '@/Github.svg?component';
import Home from '@/Home.svg?component';
import Instagram from '@/Instagram.svg?component';
import Key from '@/Key.svg?component';
import Layout from '@/Layout.svg?component';
import Menu from '@/Menu.svg?component';
import Message from '@/Message.svg?component';
import Minus from '@/Minus.svg?component';
import QuestionMark from '@/QuestionMark.svg?component';
import Selector from '@/Selector.svg?component';
import Star from '@/Star.svg?component';
import SuccessCircle from '@/SuccessCircle.svg?component';
import Sword from '@/Sword.svg?component';
import Telegram from '@/Telegram.svg?component';
import Timer from '@/Timer.svg?component';
import Twitch from '@/Twitch.svg?component';
import Users from '@/Users.svg?component';
import Variable from '@/Variable.svg?component';
import Warning from '@/Warning.svg?component';
import Website from '@/Website.svg?component';

const icons = {
  Star,
  Twitch,
  ArrowNarrow,
  Website,
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
  ArrowTriangleLarge,
  ArrowTriangleMedium,
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
  ArrowNarrowShort,
};

export type IconName = keyof typeof icons;

export default icons as {
  readonly [K in IconName]: FunctionalComponent<SVGAttributes>;
};
