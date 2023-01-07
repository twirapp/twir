import { VueQueryPlugin } from '@tanstack/vue-query';
import { Component, createSSRApp, defineComponent, h, markRaw, reactive } from 'vue';

import LandingLayout from '@/components/landing/layout/LandingLayout.vue';
import { objectAssign } from '@/utils/objectAssign.js';
import { PageContext, setPageContext } from '@/utils/pageContext.js';

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

  app.use(VueQueryPlugin);

  objectAssign(app, {
    changePage: (pageContext: PageContext) => {
      Object.assign(pageContextReactive, pageContext);
      (rootComponent as { Page: Component }).Page = markRaw(pageContext.Page);
      (rootComponent as { pageProps: any }).pageProps = markRaw(pageContext.pageProps || {});
    },
  });

  return app;
}
