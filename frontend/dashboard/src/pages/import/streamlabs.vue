<script setup lang="ts">
import IconStreamLabs from '@/assets/integrations/streamlabs.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';
</script>

<template>
	<oauth-component
		title="Streamlabs"
		:icon="IconStreamLabs"
	>
		<template #description>
			Soon...
		</template>
	</oauth-component>
</template>
