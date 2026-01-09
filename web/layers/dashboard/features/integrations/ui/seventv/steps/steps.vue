<script setup lang="ts">
import { CheckIcon, CircleIcon } from 'lucide-vue-next'



import { useSeventvSteps } from '~/features/integrations/composables/seventv/use-seventv-steps.ts'

const {
	currentStep,
	steps,
} = useSeventvSteps()
</script>

<template>
	<UiStepper
		:model-value="currentStep"
		orientation="vertical"
		class="flex w-full flex-col justify-start gap-10"
	>
		<UiStepperItem
			v-for="(step) in steps"
			:key="step.step"
			class="relative flex w-full items-start gap-6"
			:step="step.step"
			:completed="step.completed"
		>
			<UiStepperSeparator
				v-if="step.step !== steps[steps.length - 1].step"
				class="absolute left-[18px] top-[38px] block h-[105%] w-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
			/>

			<UiStepperTrigger as-child>
				<UiButton
					:variant="step.completed ? 'default' : 'outline-solid'"
					size="icon"
					class="z-10 rounded-full shrink-0"
					:class="[step.completed && 'ring-2 ring-ring ring-offset-2 ring-offset-background']"
				>
					<CheckIcon v-if="step.completed" class="size-5" />
					<CircleIcon v-else />
				</UiButton>
			</UiStepperTrigger>

			<div class="flex flex-col gap-1">
				<UiStepperTitle
					:class="[step.completed && 'text-primary']"
					class="text-sm font-semibold transition lg:text-base"
				>
					{{ step.title }}
				</UiStepperTitle>
				<UiStepperDescription
					:class="[step.completed && 'text-primary']"
					class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
					as="div"
				>
					<component :is="step.description" />
				</UiStepperDescription>
			</div>
		</UiStepperItem>
	</UiStepper>
</template>
