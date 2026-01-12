<script setup lang="ts">
import { ChevronRight, Command, Hash, Variable } from "lucide-vue-next";
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

import CommandMenuItem from "./CommandMenuItem.vue";
import CommandMenuKbd from "./CommandMenuKbd.vue";

import { useCommandMenuData } from "@/api/command-menu";
import { Button } from "@/components/ui/button";
import {
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandList,
	Command as CommandRoot,
} from "@/components/ui/command";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Kbd } from "@/components/ui/kbd";
import { Separator } from "@/components/ui/separator";
import { getFlatNavigationItems } from "@/config/navigation";
import { useIsMac } from "@/composables/useIsMac";

const router = useRouter();
const open = ref(false);
const { t } = useI18n();

const isMac = useIsMac();

const { commands, keywords, variables } = useCommandMenuData();

// Get navigation routes from shared config with translations
const navRoutes = computed(() => {
	return getFlatNavigationItems().map((route) => ({
		...route,
		displayName: route.translationKey ? t(route.translationKey) : route.name || "",
	}));
});

function runCommand(command: () => unknown) {
	open.value = false;
	command();
}
const down = (e: KeyboardEvent) => {
	if ((e.key === "k" && (e.metaKey || e.ctrlKey)) || e.key === "/") {
		if (
			(e.target instanceof HTMLElement && e.target.isContentEditable) ||
			e.target instanceof HTMLInputElement ||
			e.target instanceof HTMLTextAreaElement ||
			e.target instanceof HTMLSelectElement
		) {
			return;
		}
		e.preventDefault();
		open.value = !open.value;
	}
};

onMounted(() => document.addEventListener("keydown", down));
onUnmounted(() => document.removeEventListener("keydown", down));
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<Button
				variant="outline"
				class="relative h-9 w-full justify-start text-sm text-muted-foreground sm:pr-12 md:w-48 lg:w-64"
				@click="open = true"
			>
				<Command class="mr-2 h-4 w-4" />
				<span class="hidden lg:inline-flex">Search...</span>
				<span class="inline-flex lg:hidden">Search...</span>
				<div class="absolute top-1.5 right-1.5 hidden gap-1 sm:flex">
					<Kbd>{{ isMac ? "⌘" : "Ctrl" }}</Kbd>
					<Kbd>K</Kbd>
				</div>
			</Button>
		</DialogTrigger>
		<DialogContent class="gap-0 p-0" :show-close-button="false">
			<DialogHeader class="sr-only">
				<DialogTitle>Command Menu</DialogTitle>
				<DialogDescription>Search for pages and entities</DialogDescription>
			</DialogHeader>
			<CommandRoot class="rounded-lg border-none shadow-md">
				<CommandInput placeholder="Type to search..." />
				<CommandList class="max-h-[400px]">
					<CommandEmpty class="py-6 text-center text-sm text-muted-foreground">
						No results found.
					</CommandEmpty>

					<!-- Navigation Routes -->
					<CommandGroup heading="Pages">
						<CommandMenuItem
							v-for="route in navRoutes"
							:key="route.path"
							:value="`page ${route.displayName}`"
							@select="() => runCommand(() => router.push(route.path))"
						>
							<component :is="route.icon" class="mr-2 h-4 w-4 flex-shrink-0" />
							<span class="truncate">{{ route.displayName }}</span>
						</CommandMenuItem>
					</CommandGroup>

					<!-- Commands -->
					<CommandGroup v-if="commands.length > 0" heading="Commands">
						<CommandMenuItem
							v-for="command in commands.filter((c) => c.enabled)"
							:key="command.id"
							:value="`command ${command.name} ${command.description || ''}`"
							@select="
								() => runCommand(() => router.push(`/dashboard/commands/custom/${command.id}`))
							"
						>
							<Command class="mr-2 h-4 w-4 flex-shrink-0" />
							<span class="truncate">{{ command.name }}</span>
							<span
								v-if="command.description"
								class="ml-auto text-xs text-muted-foreground truncate max-w-[200px]"
							>
								{{ command.description }}
							</span>
						</CommandMenuItem>
					</CommandGroup>

					<!-- Keywords -->
					<CommandGroup v-if="keywords.length > 0" heading="Keywords">
						<CommandMenuItem
							v-for="keyword in keywords.filter((k) => k.enabled)"
							:key="keyword.id"
							:value="`keyword ${keyword.text}`"
							@select="() => runCommand(() => router.push(`/dashboard/keywords`))"
						>
							<Hash class="mr-2 h-4 w-4 flex-shrink-0" />
							<span class="truncate">{{ keyword.text }}</span>
						</CommandMenuItem>
					</CommandGroup>

					<!-- Variables -->
					<CommandGroup v-if="variables.length > 0" heading="Variables">
						<CommandMenuItem
							v-for="variable in variables"
							:key="variable.id"
							:value="`variable ${variable.name} ${variable.description || ''}`"
							@select="() => runCommand(() => router.push(`/dashboard/variables/${variable.id}`))"
						>
							<Variable class="mr-2 h-4 w-4 flex-shrink-0" />
							<span class="truncate">{{ variable.name }}</span>
							<span
								v-if="variable.description"
								class="ml-auto text-xs text-muted-foreground truncate max-w-[200px]"
							>
								{{ variable.description }}
							</span>
						</CommandMenuItem>
					</CommandGroup>
				</CommandList>
			</CommandRoot>
			<div class="flex items-center gap-2 border-t px-4 py-3 text-xs text-muted-foreground">
				<div class="flex items-center gap-1">
					<CommandMenuKbd>
						<ChevronRight class="h-3 w-3" />
					</CommandMenuKbd>
					<span>to navigate</span>
				</div>
				<Separator orientation="vertical" class="h-4" />
				<div class="flex items-center gap-1">
					<CommandMenuKbd>Enter</CommandMenuKbd>
					<span>to select</span>
				</div>
				<Separator orientation="vertical" class="h-4" />
				<div class="flex items-center gap-1">
					<CommandMenuKbd>Esc</CommandMenuKbd>
					<span>to close</span>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
