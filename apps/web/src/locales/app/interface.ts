import type { AppMenu } from '@/pages/app/router.js';

interface IAppLocale {
  hello: string;
  pages: Record<AppMenu, string>;
}

export default IAppLocale;
