import { Component, createSSRApp, defineComponent, h, markRaw, reactive } from 'vue';

import LandingLayout from '@/components/landing/layout/LandingLayout.vue';
import { setPageContext } from '@/hooks/usePageContext.js';
import type { PageContext } from '@/types/pageContext.js';
import { objectAssign } from '@/utils/objectAssign.js';

export function createApp(pageContext: PageContext) {
  const { Page } = pageContext;

  let rootComponent: Component;

  const PageWithWrapper = defineComponent({
    data: () => ({
      Page: markRaw(Page),
      pageProps: markRaw(pageContext.pageProps || {}),
    }),
    created() {
      // eslint-disable-next-line @typescript-eslint/no-this-alias
      rootComponent = this;
    },
    render() {
      return h(
        LandingLayout,
        {},
        {
          default: () => {
            return h(this.Page as any, this.pageProps);
          },
        },
      );
    },
  });

  const app = createSSRApp(PageWithWrapper);

  const pageContextReactive = reactive(pageContext);
  setPageContext(app, pageContextReactive);

  objectAssign(app, {
    changePage: (pageContext: PageContext) => {
      Object.assign(pageContextReactive, pageContext);
      (rootComponent as { Page: Component }).Page = markRaw(pageContext.Page);
      (rootComponent as { pageProps: any }).pageProps = markRaw(pageContext.pageProps || {});
    },
  });

  return app;
}
