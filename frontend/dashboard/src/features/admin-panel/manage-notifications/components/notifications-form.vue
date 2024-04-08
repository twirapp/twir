<script setup lang="ts">
import { ChevronDownIcon, XIcon } from 'lucide-vue-next';
import { NCard } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { SelectIcon } from 'radix-vue';
import { useI18n } from 'vue-i18n';

import { useNotificationsForm } from '../composables/use-notifications-form.js';
import { useStreamers } from '../composables/use-streamers.js';
import { useTextarea, textareaButtons } from '../composables/use-textarea.js';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
	SelectTriggerWithoutChevron,
  SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';

const { t } = useI18n();

const { streamers } = storeToRefs(useStreamers());
const notificationsForm = useNotificationsForm();

const textarea = useTextarea();
const { textareaRef } = storeToRefs(textarea);
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.notifications.formTitle') }}
	</h4>
	<n-card size="small" bordered>
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
									<XIcon v-else class="w-5 h-5 opacity-50" @pointerdown="notificationsForm.resetFieldUserId" />
								</SelectIcon>
							</SelectTriggerWithoutChevron>
						</FormControl>
						<SelectContent :hide-when-detached="true">
							<SelectGroup>
								<SelectItem v-for="streamer of streamers" :key="streamer.userId" :value="streamer.userId">
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
					<div class="flex flex-col gap-2">
						<div class="flex gap-2 flex-wrap">
							<TooltipProvider>
								<Tooltip v-for="button in textareaButtons" :key="button.name">
									<TooltipTrigger as-child>
										<Button
											type="button"
											variant="secondary"
											size="icon"
											@click="textarea.applyModifier(button.name)"
										>
											<component :is="button.icon" class="h-4 w-4" />
										</Button>
									</TooltipTrigger>
									<TooltipContent>
										<p>{{ button.title }}</p>
									</TooltipContent>
								</Tooltip>
							</TooltipProvider>
						</div>

						<FormControl>
							<Textarea
								ref="textareaRef"
								rows="8"
								v-bind="componentField"
							/>
						</FormControl>
					</div>
					<FormMessage />
				</FormItem>
			</FormField>

			<template v-if="notificationsForm.message">
				<Label>{{ t('adminPanel.notifications.messagePreview') }}</Label>
				<div class="border rounded-md p-2" v-html="notificationsForm.message"></div>
			</template>

			<div class="flex justify-end gap-4">
				<Button
					:disabled="!notificationsForm.isFormDirty && !notificationsForm.editableMessageId"
					type="button"
					variant="secondary"
					@click="notificationsForm.onReset"
				>
					<template v-if="notificationsForm.editableMessageId">
						{{ t('sharedButtons.cancel') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.reset') }}
					</template>
				</Button>
				<Button type="submit">
					<template v-if="notificationsForm.editableMessageId">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.send') }}
					</template>
				</Button>
			</div>
		</form>
	</n-card>
</template>
