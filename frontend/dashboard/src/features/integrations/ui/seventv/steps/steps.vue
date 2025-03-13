<script setup lang="ts">
import { CheckIcon, CircleIcon } from 'lucide-vue-next'

import { Button } from '@/components/ui/button'
import {
	Stepper,
	StepperDescription,
	StepperItem,
	StepperSeparator,
	StepperTitle,
	StepperTrigger,
} from '@/components/ui/stepper'
import { useSeventvSteps } from '@/features/integrations/composables/seventv/use-seventv-steps.ts'

const {
	currentStep,
	steps,
} = useSeventvSteps()
</script>

<template>
	<Stepper
		:model-value="currentStep"
		orientation="vertical"
		class="flex w-full flex-col justify-start gap-10"
	>
		<StepperItem
			v-for="(step) in steps"
			:key="step.step"
			class="relative flex w-full items-start gap-6"
			:step="step.step"
			:completed="step.completed"
		>
			<StepperSeparator
				v-if="step.step !== steps[steps.length - 1].step"
				class="absolute left-[18px] top-[38px] block h-[105%] w-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
			/>

			<StepperTrigger as-child>
				<Button
					:variant="step.completed ? 'default' : 'outline'"
					size="icon"
					class="z-10 rounded-full shrink-0"
					:class="[step.completed && 'ring-2 ring-ring ring-offset-2 ring-offset-background']"
				>
					<CheckIcon v-if="step.completed" class="size-5" />
					<CircleIcon v-else />
				</Button>
			</StepperTrigger>

			<div class="flex flex-col gap-1">
				<StepperTitle
					:class="[step.completed && 'text-primary']"
					class="text-sm font-semibold transition lg:text-base"
				>
					{{ step.title }}
				</StepperTitle>
				<StepperDescription
					:class="[step.completed && 'text-primary']"
					class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
					as="div"
				>
					<component :is="step.description" />
				</StepperDescription>
			</div>
		</StepperItem>
	</Stepper>
</template>
