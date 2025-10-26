<script setup lang="ts">
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
</script>

<template>
	<DropdownMenu>
		<DropdownMenuTrigger as-child>
			<UiButton
				size="icon"
				variant="outline"
				class="w-8 rounded-lg font-semibold border-transparent bg-transparent hover:bg-[hsl(240,11%,15%)]"
			>
				<Icon name="lucide:ellipsis-vertical" class="size-4" />
			</UiButton>
		</DropdownMenuTrigger>
		<DropdownMenuContent
			class="flex flex-col w-fit bg-[hsl(240,11%,12%)] rounded-xl p-1.5 gap-1"
			align="end"
		>
			<slot name="content" />
		</DropdownMenuContent>
	</DropdownMenu>
</template>
