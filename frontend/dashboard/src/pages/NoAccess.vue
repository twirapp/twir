<script setup lang="ts">
import { NResult } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
	<n-result status="403" :title="t('haveNoAccess.message')" />
</template>

<style scoped lang="postcss">

</style>
