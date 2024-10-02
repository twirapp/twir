<script setup lang="ts">
import RewardsSelector from '@/components/rewardsSelector.vue'
import { Label } from '@/components/ui/label'
import { useForm } from '@/features/chat-alerts/composables/use-form'

const { formValue } = useForm()
</script>

<template>
	<div v-if="'ignoredRewardsIds' in formValue.redemptions" class="grid items-center gap-1.5">
		<Label for="rewardsIgnoredIds">Ignored rewards</Label>
		<RewardsSelector
			id="rewardsIgnoredIds"
			v-model="formValue.redemptions.ignoredRewardsIds"
			multiple
			clearable
			:onlyWithInput="false"
		/>
	</div>
</template>
