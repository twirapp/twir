<script setup lang="ts">
import { toast } from 'vue-sonner'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardAction, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { createMcpClientGuides } from './config.js'

const { t } = useI18n()
const requestUrl = useRequestURL()

const endpoint = computed(() => `${requestUrl.origin}/api/mcp`)
const guides = computed(() => createMcpClientGuides(endpoint.value))

async function copy(value: string) {
	try {
		await navigator.clipboard.writeText(value)
		toast.success(t('sharedTexts.copied'), { duration: 1500 })
	} catch {
		toast.error(t('mcpGuide.copyError'))
	}
}

function copyGuide(id: string) {
	const guide = guides.value.find((item) => item.id === id)
	if (guide) copy(guide.config)
}
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('mcpGuide.title') }}
		</template>

		<template #title-footer>
			<p class="text-muted-foreground">
				{{ t('mcpGuide.description') }}
			</p>
		</template>

		<template #content>
			<div class="flex flex-col gap-6">
				<Card>
					<CardHeader>
						<CardTitle>{{ t('mcpGuide.credentials.title') }}</CardTitle>
						<CardDescription>{{ t('mcpGuide.credentials.description') }}</CardDescription>
					</CardHeader>
					<CardContent class="flex flex-col gap-5">
						<div class="flex flex-col gap-2">
							<p class="text-sm font-medium">{{ t('mcpGuide.credentials.endpoint') }}</p>
							<div class="flex gap-2">
								<Input :model-value="endpoint" readonly />
								<Button variant="outline" size="icon" type="button" @click="copy(endpoint)">
									<Icon name="lucide:copy" />
									<span class="sr-only">{{ t('mcpGuide.copyEndpoint') }}</span>
								</Button>
							</div>
						</div>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>{{ t('mcpGuide.clientsTitle') }}</CardTitle>
						<CardDescription>{{ t('mcpGuide.clientsDescription') }}</CardDescription>
					</CardHeader>
					<CardContent>
						<Tabs default-value="claude" class="flex flex-col gap-4">
							<TabsList class="grid h-auto w-full grid-cols-2 sm:grid-cols-4">
								<TabsTrigger v-for="guide in guides" :key="guide.id" :value="guide.id">
									<Icon :name="guide.icon" />
									{{ guide.name }}
								</TabsTrigger>
							</TabsList>

							<TabsContent v-for="guide in guides" :key="guide.id" :value="guide.id">
								<Card>
									<CardHeader>
										<CardTitle class="flex items-center gap-2">
											<Icon :name="guide.icon" />
											{{ guide.name }}
										</CardTitle>
										<CardDescription>{{ t(guide.descriptionKey) }}</CardDescription>
										<CardAction>
											<Button variant="outline" size="sm" as-child>
												<a :href="guide.docsUrl" target="_blank" rel="noreferrer">
													{{ t('mcpGuide.documentation') }}
													<Icon name="lucide:external-link" data-icon="inline-end" />
												</a>
											</Button>
										</CardAction>
									</CardHeader>
									<CardContent class="flex flex-col gap-5">
										<ol class="flex list-decimal flex-col gap-2 pl-5 text-sm text-muted-foreground">
											<li v-for="stepKey in guide.stepKeys" :key="stepKey">
												{{ t(stepKey) }}
											</li>
										</ol>

										<div class="flex items-center justify-between gap-2">
											<Badge v-if="guide.fileName" variant="secondary">{{ guide.fileName }}</Badge>
											<span v-else />
											<Button type="button" size="sm" @click="copyGuide(guide.id)">
												<Icon name="lucide:copy" data-icon="inline-start" />
												{{ t('mcpGuide.copyConfig') }}
											</Button>
										</div>

										<div class="overflow-hidden rounded-md bg-muted">
											<pre class="overflow-x-auto p-4 text-sm"><code>{{ guide.config }}</code></pre>
										</div>
									</CardContent>
								</Card>
							</TabsContent>
						</Tabs>
					</CardContent>
				</Card>
			</div>
		</template>
	</PageLayout>
</template>
