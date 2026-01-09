<script setup lang="ts">
import RewardsSelector from '#layers/dashboard/components/rewardsSelector.vue'

import { useForm } from '~/features/chat-alerts/composables/use-form'

const { formValue } = useForm()
</script>

<template>
	<div class="grid items-center gap-1.5">
		<UiLabel for="rewardsIgnoredIds">Ignored rewards</UiLabel>
		<RewardsSelector
			id="rewardsIgnoredIds"
			v-model="formValue.redemptions.ignoredRewardsIds"
			multiple
			clearable
			:onlyWithInput="false"
		/>
	</div>
</template>
