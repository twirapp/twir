import { useAxios } from '@vueuse/integrations/useAxios';
import { watch } from 'vue';
import { useToast } from 'vue-toastification';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const toast = useToast();

type Opts = {
  immediate?: boolean;
};

export function useUpdatingDataFixed<T = any>(
  url: string,
  opts: Opts | undefined = { immediate: false },
) {
  const buildUrl = () => url.replace('{dashboardId}', selectedDashboardStore.get()?.channelId);

  const { execute, data, error, isFinished, isLoading } = useAxios<T>(buildUrl(), api, {
    immediate: opts.immediate,
  });

  selectedDashboardStore.listen(() => {
    execute(buildUrl());
  });

  watch(error, (data: any) => {
    if (data?.response?.data.message) {
      toast(data.response.data.message);
    }
  });

  return {
    data,
    execute,
    error,
    isFinished,
    isLoading,
  };
}
