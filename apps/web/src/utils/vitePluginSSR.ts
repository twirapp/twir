import type { PageContext } from '@/utils/pageContext.js';

export type PassToClient = (keyof PageContext)[];
export type PrerenderFn = () => { url: string; pageContext: Partial<PageContext> }[];

