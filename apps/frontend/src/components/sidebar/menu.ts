import { 
  mdiCalendarEdit,
  mdiClockOutline,
  mdiCog,
  mdiConsoleLine,
  mdiHome,
  mdiPackageVariantClosed,
  mdiSword,
  mdiAccountSupervisor,
  mdiScriptTextKeyOutline,
  mdiPulse,
  mdiEmoticon,
  mdiAnchor,
  mdiFolder,
} from '@mdi/js';

export const publicRoutes = [
  {
    name: 'dashboard',
    path: '/dashboard',
    icon: mdiHome,
  },
  // {
  //   name: 'events',
  //   path: '/dashboard/events',
  //   icon: mdiCalendarEdit,
  // },
  {
    name: 'integrations',
    path: '/dashboard/integrations',
    icon: mdiPackageVariantClosed,
  },
  {
    name: 'settings',
    path: '/dashboard/settings',
    icon: mdiCog,
  },
  {
    name: 'commands',
    path: '/dashboard/commands',
    icon: mdiConsoleLine,
  },
  {
    name: 'timers',
    path: '/dashboard/timers',
    icon: mdiClockOutline,
  },
  {
    name: 'moderation',
    path: '/dashboard/moderation',
    icon: mdiSword,
  },
  // {
  //   name: 'users',
  //   path: '/dashboard/users',
  //   icon: mdiAccountSupervisor,
  // },
  {
    name: 'keywords',
    path: '/dashboard/keywords',
    icon: mdiScriptTextKeyOutline,
  },
  {
    name: 'variables',
    path: '/dashboard/variables',
    icon: mdiPulse,
  },
  {
    name: 'greetings',
    path: '/dashboard/greetings',
    icon: mdiEmoticon,
  },
  // {
  //   name: 'overlays',
  //   path: '/dashboard/overlays',
  //   icon: mdiAnchor,
  // },
  // {
  //   name: 'files',
  //   path: '/dashboard/files',
  //   icon: mdiFolder
  // },
  // {
  //   name: 'quotes',
  //   path: '/dashboard/quotes',
  // },
];
