<script setup lang="ts">
import Table from '~~/layers/dashboard/components/table.vue'

import { useSecretsTable } from './composables/use-secrets-table'
import SecretCreateButton from './ui/secret-create-button.vue'

const secretsTable = useSecretsTable()
const { t } = useI18n()
</script>

<template>
	<div class="flex flex-col gap-4">
		<div class="flex justify-between items-center">
			<div>
				<p class="text-sm text-muted-foreground">
					Store sensitive values like API tokens and passwords. Secrets are encrypted and can be accessed in scripts via
					<code class="bg-muted px-1 py-0.5 rounded text-xs">twir.secrets.get('SECRET_NAME')</code>
				</p>
			</div>
			<SecretCreateButton />
		</div>
		<Table
			:table="secretsTable.table"
			:is-loading="secretsTable.isLoading.value"
		/>
	</div>
</template>
