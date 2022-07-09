import AdminMain from '@/admin/Main.vue';
import Commands from '@/assets/sidebar/commands.svg?url';
import Dashboard from '@/assets/sidebar/dashboard.svg?url';
import Events from '@/assets/sidebar/events.svg?url';
import Folder from '@/assets/sidebar/folder.svg?url';
import Greetings from '@/assets/sidebar/greetings.svg?url';
import Integrations from '@/assets/sidebar/integrations.svg?url';
import Keywords from '@/assets/sidebar/keywords.svg?url';
import Overlays from '@/assets/sidebar/overlays.svg?url';
import Quotes from '@/assets/sidebar/quotes.svg?url';
import Settings from '@/assets/sidebar/settings.svg?url';
import Sword from '@/assets/sidebar/sword.svg?url';
import Timers from '@/assets/sidebar/timers.svg?url';
import Users from '@/assets/sidebar/users.svg?url';
import Variables from '@/assets/sidebar/variables.svg?url';

export const publicRoutes = [
  {
    name: 'dashboard',
    icon: Dashboard,
    path: '/dashboard',
  },
  {
    name: 'events',
    icon: Events,
    path: '/dashboard/events',
  },
  {
    name: 'integrations',
    icon: Integrations,
    path: '/dashboard/integrations',
  },
  {
    name: 'settings',
    icon: Settings,
    path: '/dashboard/settings',
  },
  {
    name: 'commands',
    icon: Commands,
    path: '/dashboard/commands',
  },
  {
    name: 'timers',
    icon: Timers,
    path: '/dashboard/timers',
  },
  {
    name: 'moderation',
    icon: Sword,
    path: '/dashboard/moderation',
  },
  {
    name: 'users',
    icon: Users,
    path: '/dashboard/users',
  },
  {
    name: 'keywords',
    icon: Keywords,
    path: '/dashboard/keywords',
  },
  {
    name: 'variables',
    icon: Variables,
    path: '/dashboard/variables',
  },
  {
    name: 'greetings',
    icon: Greetings,
    path: '/dashboard/greetings',
  },
  {
    name: 'overlays',
    icon: Overlays,
    path: '/dashboard/overlays',
  },
  {
    name: 'files',
    icon: Folder,
    path: '/dashboard/files',
  },
  {
    name: 'quotes',
    icon: Quotes,
    path: '/dashboard/quotes',
  },
];

export const adminRoutes = [
  {
    name: 'admin',
    icon: Dashboard,
    path: '/admin',
  },
];