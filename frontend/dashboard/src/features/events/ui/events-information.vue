<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { InfoIcon } from 'lucide-vue-next'
</script>

<template>
  <Alert>
    <InfoIcon class="h-4 w-4" />
    <AlertTitle>{{ $t('events.information.title') }}</AlertTitle>
    <AlertDescription>
      {{ $t('events.information.description') }}
    </AlertDescription>
  </Alert>
</template>
