<script setup lang="ts">
import { NCard } from 'naive-ui';
import { onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';

import BadgesPreview from './badges-preview.vue';
import { useBadgesForm } from '../composables/use-badges-form';

import { Button } from '@/components/ui/button';
import {
	FormMessage,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

const { t } = useI18n();
const badgesForm = useBadgesForm();

onUnmounted(() => {
	badgesForm.onReset();
});
</script>

<template>
	<n-card size="small" bordered>
		<form class="flex flex-col gap-4" @submit="badgesForm.onSubmit">
			<FormField v-slot="{ componentField }" name="name">
				<FormItem>
					<FormLabel>{{ t('adminPanel.manageBadges.name') }}</FormLabel>
					<FormControl>
						<Input type="text" placeholder="" v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField name="file">
				<FormItem>
					<FormLabel>
						{{ t('adminPanel.manageBadges.image') }}
					</FormLabel>
					<FormControl>
						<div className="grid w-full items-center gap-1.5">
							<Input
								:required="!badgesForm.editableBadgeId"
								accept="image/*"
								type="file"
								@change="badgesForm.setImageField"
							/>
						</div>
					</FormControl>
				</FormItem>
			</FormField>

			<div v-if="badgesForm.image">
				<Label>
					{{ t('adminPanel.manageBadges.preview') }}
				</Label>
				<badges-preview :image="badgesForm.image" />
			</div>

			<div class="flex justify-end gap-4">
				<Button
					type="button"
					variant="secondary"
					:disabled="!badgesForm.isFormDirty" @click="badgesForm.onReset"
				>
					<template v-if="badgesForm.editableBadgeId">
						{{ t('sharedButtons.cancel') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.reset') }}
					</template>
				</Button>
				<Button type="submit">
					<template v-if="badgesForm.editableBadgeId">
						{{ t('sharedButtons.edit') }}
					</template>
					<template v-else>
						{{ t('sharedButtons.create') }}
					</template>
				</Button>
			</div>
		</form>
	</n-card>
</template>
