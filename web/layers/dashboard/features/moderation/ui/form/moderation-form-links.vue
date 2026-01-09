<script setup lang="ts">


</script>

<template>
	<div>
		<UiFormField v-slot="{ field }" name="checkClips">
			<UiFormItem class="flex flex-row items-center justify-between rounded-lg border p-4">
				<div class="space-y-0.5">
					<UiFormLabel class="text-base"> Check twitch clips </UiFormLabel>
					<UiFormDescription> If enabled bot will check clips as links. </UiFormDescription>
				</div>
				<UiFormControl>
					<UiSwitch
						:model-value="field.value"
						default-checked
						@update:model-value="field['onUpdate:modelValue']"
					/>
				</UiFormControl>
			</UiFormItem>
		</UiFormField>
	</div>
</template>
