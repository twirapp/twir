<script setup lang="ts">
import { ArrowDownIcon, ArrowUpIcon, EditIcon, TrashIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { z } from 'zod'

import type { UserPublicSettingsQuery } from '@/gql/graphql'

import { useUserSettings } from '@/api/auth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { toast } from 'vue-sonner'

type Link = Omit<
	UserPublicSettingsQuery['userPublicSettings']['socialLinks'][number],
	'__typename'
> & { isEditing?: boolean }

const formSchema = z.object({
	description: z.string().max(1000, 'Description must be less than 1000 characters'),
	socialLinks: z.array(
		z.object({
			title: z
				.string()
				.min(1, 'Title is required')
				.max(30, 'Title must be less than 30 characters'),
			href: z.string().min(1, 'URL is required').max(256, 'URL must be less than 256 characters'),
		})
	),
})

const { t } = useI18n()
const manager = useUserSettings()
const { data } = manager.usePublicQuery()
const updater = manager.usePublicMutation()

const socialLinksWithEditing = ref<Link[]>([])
const newLinkTitle = ref('')
const newLinkHref = ref('')

const form = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		description: '',
		socialLinks: [],
	},
})

const [description] = form.defineField('description')

watch(
	data,
	(v) => {
		if (!v) return
		const rawData = v.userPublicSettings
		form.setValues({
			description: rawData.description ?? '',
			socialLinks: rawData.socialLinks.map((link) => ({ title: link.title, href: link.href })),
		})
		socialLinksWithEditing.value = rawData.socialLinks.map((link) => ({
			...link,
			isEditing: false,
		}))
	},
	{ immediate: true }
)

const onSubmit = form.handleSubmit(async (values) => {
	await updater.executeMutation({
		opts: {
			description: values.description,
			socialLinks: values.socialLinks,
		},
	})

	toast.success(t('sharedTexts.saved'), {
		duration: 1500,
	})
})

const linksLimitReached = computed(() => (form.values.socialLinks?.length ?? 0) >= 10)

function addLink() {
	const title = newLinkTitle.value.trim()
	const href = newLinkHref.value.trim()

	if (!title || !href) return
	if (title.length === 0 || title.length > 30) return
	if (href.length === 0 || href.length > 256) return

	const currentLinks = form.values.socialLinks || []
	form.setFieldValue('socialLinks', [...currentLinks, { title, href }])
	socialLinksWithEditing.value.push({ title, href, isEditing: false })
	newLinkTitle.value = ''
	newLinkHref.value = ''
}

function removeLink(index: number) {
	const currentLinks = form.values.socialLinks || []
	form.setFieldValue(
		'socialLinks',
		currentLinks.filter((_, i) => i !== index)
	)
	socialLinksWithEditing.value = socialLinksWithEditing.value.filter((_, i) => i !== index)
}

function changeSort(from: number, to: number) {
	const currentLinks = [...(form.values.socialLinks || [])]
	const element = currentLinks.splice(from, 1)[0]
	currentLinks.splice(to, 0, element)
	form.setFieldValue('socialLinks', currentLinks)

	const element2 = socialLinksWithEditing.value.splice(from, 1)[0]
	socialLinksWithEditing.value.splice(to, 0, element2)
}

function toggleEdit(index: number) {
	socialLinksWithEditing.value[index].isEditing = !socialLinksWithEditing.value[index].isEditing
}

function updateLinkTitle(index: number, newTitle: string | number) {
	const title = String(newTitle)
	const currentLinks = [...(form.values.socialLinks || [])]
	currentLinks[index] = { ...currentLinks[index], title }
	form.setFieldValue('socialLinks', currentLinks)
	socialLinksWithEditing.value[index].title = title
}

function updateLinkHref(index: number, newHref: string | number) {
	const href = String(newHref)
	const currentLinks = [...(form.values.socialLinks || [])]
	currentLinks[index] = { ...currentLinks[index], href }
	form.setFieldValue('socialLinks', currentLinks)
	socialLinksWithEditing.value[index].href = href
}
</script>

<template>
	<form class="w-full flex flex-wrap gap-6" @submit="onSubmit">
		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.description') }}
			</h4>

			<Card class="p-0">
				<CardContent class="pt-6 pb-2 space-y-2">
					<Label for="description">{{ t('userSettings.public.description') }}</Label>
					<Textarea
						id="description"
						v-model="description"
						maxlength="1000"
						placeholder=""
						class="min-h-20"
					/>
					<p class="text-xs text-muted-foreground">{{ description?.length || 0 }}/1000</p>
				</CardContent>
			</Card>
		</div>

		<div class="flex flex-col w-full gap-6">
			<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
				{{ t('userSettings.public.socialLinks') }}
			</h4>

			<Card>
				<CardContent class="space-y-4">
					<div class="flex flex-col gap-2">
						<Card v-for="(link, idx) of socialLinksWithEditing" :key="idx" class="overflow-hidden">
							<CardHeader class="pb-2">
								<div class="flex items-center justify-between gap-2">
									<div class="flex-1 min-w-0">
										<Input
											v-if="link.isEditing"
											:model-value="link.title"
											:maxlength="30"
											class="w-full h-8"
											@update:model-value="(v) => updateLinkTitle(idx, v)"
										/>
										<template v-else>
											<CardTitle class="text-base truncate">{{ link.title }}</CardTitle>
											<CardDescription>
												<a class="text-white underline" :href="link.href" target="_blank">
													{{ link.href }}
												</a>
											</CardDescription>
										</template>
									</div>
									<div class="flex gap-1 shrink-0">
										<Button
											variant="ghost"
											size="icon"
											class="size-8"
											type="button"
											:disabled="!socialLinksWithEditing[idx + 1]"
											@click="changeSort(idx, idx + 1)"
										>
											<ArrowDownIcon class="size-4" />
										</Button>
										<Button
											variant="ghost"
											size="icon"
											class="size-8"
											type="button"
											:disabled="idx === 0"
											@click="changeSort(idx, idx - 1)"
										>
											<ArrowUpIcon class="size-4" />
										</Button>
										<Button
											variant="ghost"
											size="icon"
											class="size-8"
											type="button"
											@click="toggleEdit(idx)"
										>
											<EditIcon class="size-4" />
										</Button>
										<Button
											variant="ghost"
											size="icon"
											class="size-8 hover:bg-destructive/10 hover:text-destructive"
											type="button"
											@click="removeLink(idx)"
										>
											<TrashIcon class="size-4" />
										</Button>
									</div>
								</div>
							</CardHeader>
							<CardContent v-if="link.isEditing">
								<Textarea
									v-if="link.isEditing"
									:model-value="link.href"
									:maxlength="256"
									class="min-h-15"
									@update:model-value="(v) => updateLinkHref(idx, v)"
								/>
							</CardContent>
						</Card>
					</div>

					<div class="flex flex-wrap gap-2 items-end w-full mt-5">
						<div class="flex-1 space-y-2 min-w-50">
							<Label for="newLinkTitle">{{ t('userSettings.public.linkTitle') }}</Label>
							<Input
								id="newLinkTitle"
								v-model="newLinkTitle"
								:maxlength="30"
								placeholder="Twir"
								:disabled="linksLimitReached"
							/>
						</div>
						<div class="flex-1 space-y-2 min-w-50">
							<Label for="newLinkHref">{{ t('userSettings.public.linkLabel') }}</Label>
							<Input
								id="newLinkHref"
								v-model="newLinkHref"
								:maxlength="256"
								placeholder="https://twir.app"
								:disabled="linksLimitReached"
							/>
						</div>
						<Button
							type="button"
							variant="secondary"
							:disabled="linksLimitReached"
							@click="addLink"
						>
							{{ t('sharedButtons.add') }}
						</Button>
					</div>
				</CardContent>
			</Card>
		</div>

		<div class="flex justify-start w-full">
			<Button type="submit">
				{{ t('sharedButtons.save') }}
			</Button>
		</div>
	</form>
</template>
