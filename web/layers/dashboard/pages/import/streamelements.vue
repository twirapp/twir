<script setup lang="ts">
import IconStreamElements from '~/assets/integrations/streamelements.svg?use';
import OauthComponent from '#layers/dashboard/components/integrations/variants/oauth.vue';
</script>

<template>
	<oauth-component
		title="StreamElements"
		:icon="IconStreamElements"
		icon-width="200px"
		:data="{}"
		:logout="() => {}"
	>
		<template #description>
			Soon...
		</template>
	</oauth-component>
</template>
