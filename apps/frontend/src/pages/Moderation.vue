<script lang="ts" setup>
import { mdiLinkBoxVariant, mdiCapsLock, mdiClipboard, mdiCalculatorVariant, mdiEmoticon, mdiWrapDisabled } from '@mdi/js';

import Card from '@/components/moderation/card.vue';

</script>

<template>
  <div>
    <div class="d-flex justify-center flex-column flex-sm-row flex-md-row flex-lg-row flex-xl-row">
      <Card type="links" :icon="mdiLinkBoxVariant" icon-color="blue" />
      <Card type="caps" :icon="mdiCapsLock" icon-color="orange" />
    </div>
    <div class="d-flex justify-center flex-column flex-sm-row flex-md-row flex-lg-row flex-xl-row">
      <Card type="symbols" :icon="mdiCalculatorVariant" icon-color="green" />
      <Card type="emotes" :icon="mdiEmoticon" icon-color="yellow" />
    </div>
    <div class="d-flex justify-center flex-column flex-sm-row flex-md-row flex-lg-row flex-xl-row">
      <Card type="blacklists" :icon="mdiClipboard" icon-color="red" />
      <Card type="longMessage" :icon="mdiWrapDisabled" icon-color="cyan" />
    </div>
  </div>
</template>