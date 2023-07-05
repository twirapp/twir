import { QueryClient } from '@tanstack/react-query';

import { printError } from '@/services/api/error';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
    },
    mutations: {
      onError: (error: any) => {
        // if (error instanceof FetcherError) {
        //   printError(error.messages ?? error.message);
        // } else {
        //   printError(error.message);
        // }
      },
    },
  },
});
