import type { PageContext } from '@/types/pageContext';

export function getPageTitle(pageContext: PageContext): string {
  return (
    (pageContext.exports.documentProps || {}).title ||
    (pageContext.documentProps || {}).title ||
    'Tsuwari'
  );
}
