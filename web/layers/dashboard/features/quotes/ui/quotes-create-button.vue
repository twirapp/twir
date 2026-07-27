<script setup lang="ts">
import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth.js'

import { Button } from '@/components/ui/button'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

import QuotesDialog from './quotes-dialog.vue'

const { t } = useI18n()
const canManageQuotes = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageQuotes)
</script>

<template>
	<QuotesDialog>
		<template #dialog-trigger>
			<Button :disabled="!canManageQuotes">
				<Icon name="lucide:plus" class="mr-2 size-4" />
				{{ t('quotes.create') }}
			</Button>
		</template>
	</QuotesDialog>
</template>
