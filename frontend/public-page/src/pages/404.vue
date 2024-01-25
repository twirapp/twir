<script setup lang="ts">
import { IconError404 } from '@tabler/icons-vue';

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
</script>

<template>
	<div class="container h-screen flex items-center justify-center">
		<Alert variant="destructive">
			<IconError404 />
			<AlertTitle>Heads up!</AlertTitle>
			<AlertDescription>
				Page or channel you're looking for not found
			</AlertDescription>
		</Alert>
	</div>
</template>
