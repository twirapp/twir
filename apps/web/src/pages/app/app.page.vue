<template>
  <span>{{ t('tabs') }}</span>
  <router-view v-slot="{ Component }">
    <Suspense>
      <component :is="Component" />
    </Suspense>
  </router-view>
  <router-link to="/dashboard">
    Dashboard
  </router-link>
  <router-link to="/commands">
    Commands
  </router-link>
</template>

<script lang="ts" setup>
import useTranslation from '@/hooks/useTranslation';

const t = useTranslation<'app'>();
</script>
