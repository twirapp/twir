<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod';
import { ChevronDownIcon, XIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { SelectIcon } from 'radix-vue';
import { useForm } from 'vee-validate';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import * as z from 'zod';

import { useAdminNotifications } from '@/api/notifications';
import { useStreamers } from '@/api/streamers';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
	SelectTriggerWithoutChevron,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';

const { t } = useI18n();

const notifications = useAdminNotifications();

const { data: streamers } = useStreamers();
const streamersOptions = computed(() => {
	if (!streamers.value?.streamers) return [];
	return streamers.value.streamers;
});

const formSchema = toTypedSchema(z.object({
  userId: z.string().optional(),
	message: z.string(),
}));

const form = useForm({
  validationSchema: formSchema,
});

const onSubmit = form.handleSubmit(async (values) => {
  await notifications.create.mutateAsync(values);
	form.resetForm();
});

function onResetUserId(event: Event): void {
	event.stopPropagation();
	form.resetField('userId');
}
</script>

<template>
	<n-card :title="t('adminPanel.notifications.createNotification')" size="small" bordered>
		<form class="flex flex-col gap-4" @submit="onSubmit">
			<FormField v-slot="{ componentField }" name="userId">
				<FormItem>
					<FormLabel>User</FormLabel>

					<Select v-bind="componentField">
						<FormControl>
							<SelectTriggerWithoutChevron>
								<SelectValue placeholder="Select a user" />

								<SelectIcon>
									<ChevronDownIcon v-if="!componentField.modelValue" class="w-5 h-5 opacity-50" />
									<XIcon v-else class="w-5 h-5 opacity-50" @pointerdown="onResetUserId" />
								</SelectIcon>
							</SelectTriggerWithoutChevron>
						</FormControl>
						<SelectContent :hide-when-detached="true">
							<SelectGroup>
								<SelectItem v-for="streamer of streamersOptions" :key="streamer.userId" :value="streamer.userId">
									<div class="flex items-center gap-2">
										<Avatar class="h-6 w-6">
											<AvatarImage :src="streamer.avatar" :alt="streamer.userDisplayName" loading="lazy" />
											<AvatarFallback>{{ streamer.userLogin.charAt(0).toUpperCase() }}</AvatarFallback>
										</Avatar>
										<span>{{ streamer.userDisplayName }}</span>
									</div>
								</SelectItem>
							</SelectGroup>
						</SelectContent>
					</Select>
				</FormItem>
			</FormField>

			<FormField v-slot="{ componentField }" name="message">
				<FormItem>
					<FormLabel>Message</FormLabel>
					<FormControl>
						<Textarea
							placeholder=""
							class="resize-none"
							v-bind="componentField"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<div class="flex justify-end">
				<Button class="" type="submit">
					Submit
				</Button>
			</div>
		</form>
	</n-card>
</template>
