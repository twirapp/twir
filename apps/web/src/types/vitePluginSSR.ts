import type { PageContext } from '@/types/pageContext';

export type PassToClient = (keyof PageContext)[];
export type PrerenderFn = () => { url: string; pageContext: Partial<PageContext> }[];
