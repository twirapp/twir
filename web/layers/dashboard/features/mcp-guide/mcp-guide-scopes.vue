<script setup lang="ts">
import { Badge } from '@/components/ui/badge'

import type { McpGuideScope } from './api.js'

type ScopesState = 'loading' | 'ready' | 'error'

interface Props {
	readonly state: ScopesState
	readonly scopes: readonly McpGuideScope[]
}

defineProps<Props>()

const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-3">
		<div class="flex flex-col gap-1">
			<p class="text-sm font-medium">{{ t('mcpGuide.oauth.scopesTitle') }}</p>
			<p class="text-sm text-muted-foreground">{{ t('mcpGuide.oauth.scopesDescription') }}</p>
		</div>

		<p v-if="state === 'loading'" class="text-sm text-muted-foreground" data-test="scopes-loading">
			{{ t('mcpGuide.oauth.scopesLoading') }}
		</p>
		<p
			v-else-if="state === 'error'"
			class="text-sm text-muted-foreground"
			data-test="scopes-unavailable"
		>
			{{ t('mcpGuide.oauth.scopesUnavailable') }}
		</p>
		<div v-else class="flex flex-col divide-y rounded-md border">
			<div
				v-for="scope in scopes"
				:key="scope.group"
				class="flex flex-col gap-1 p-3"
				data-test="scope-group"
				:data-group="scope.group"
			>
				<div class="flex flex-wrap items-center gap-2">
					<span class="text-sm font-medium">{{ scope.name }}</span>
					<Badge
						v-for="action in scope.actions"
						:key="action"
						variant="secondary"
						class="font-mono"
						data-test="scope-action"
					>
						{{ action }}
					</Badge>
				</div>
				<p class="text-sm text-muted-foreground">{{ scope.description }}</p>
			</div>
		</div>

		<p class="text-sm text-muted-foreground">{{ t('mcpGuide.oauth.scopesEditIncludesRead') }}</p>
		<p class="text-sm text-muted-foreground">{{ t('mcpGuide.oauth.scopesLegacyAliases') }}</p>
		<p class="text-sm text-muted-foreground">{{ t('mcpGuide.oauth.tokenLifetimes') }}</p>
	</div>
</template>
