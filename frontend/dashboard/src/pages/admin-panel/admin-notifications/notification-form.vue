<script setup lang="ts">
import { ChevronDownIcon, XIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { SelectIcon } from 'radix-vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useNotificationsForm } from './use-notifications-form';

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

const notificationsForm = useNotificationsForm();
const { data: streamers } = useStreamers();

const streamersOptions = computed(() => {
	if (!streamers.value?.streamers) return [];
	return streamers.value.streamers;
});

function resetFieldUserId(event: Event): void {
	event.stopPropagation();
	notificationsForm.form.resetField('userId');
}
</script>

<template>
	<n-card :title="t('adminPanel.notifications.createNotification')" size="small" bordered>
		<form class="flex flex-col gap-4" @submit="notificationsForm.onSubmit">
			<FormField v-slot="{ componentField }" name="userId">
				<FormItem>
					<FormLabel>{{ t('adminPanel.notifications.userLabel') }}</FormLabel>

					<Select :disabled="notificationsForm.isEditableForm" v-bind="componentField">
						<FormControl>
							<SelectTriggerWithoutChevron>
								<SelectValue :placeholder="t('adminPanel.notifications.userPlaceholder')" />

								<SelectIcon>
									<ChevronDownIcon v-if="!componentField.modelValue" class="w-5 h-5 opacity-50" />
									<XIcon v-else class="w-5 h-5 opacity-50" @pointerdown="resetFieldUserId" />
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
					<FormLabel>{{ t('adminPanel.notifications.messageLabel') }}</FormLabel>
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

			<div class="flex justify-end gap-4">
				<Button type="button" variant="secondary" @click="notificationsForm.onReset">
					{{ t('adminPanel.notifications.resetButton') }}
				</Button>
				<Button type="submit">
					{{ t('adminPanel.notifications.sendButton') }}
				</Button>
			</div>
		</form>
	</n-card>
</template>
