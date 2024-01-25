<script setup lang="ts">
import { useRoute } from 'vue-router';

import {
	NavigationMenu,
	NavigationMenuItem,
	NavigationMenuLink,
	NavigationMenuList,
	navigationMenuTriggerStyle,
} from '@/components/ui/navigation-menu';
import { publicRouteNames } from '@/routes';

const currentRoute = useRoute();
</script>

<template>
	<NavigationMenu>
		<NavigationMenuList>
			<NavigationMenuItem>
				<NavigationMenuLink
					v-for="routeName of Object.values(publicRouteNames)"
					:key="routeName"
					:class="navigationMenuTriggerStyle()"
					as-child
					:active="currentRoute.name === routeName"
				>
					<RouterLink :to="{ name: routeName }">
						{{ routeName }}
					</RouterLink>
				</NavigationMenuLink>
			</NavigationMenuItem>
		</NavigationMenuList>
	</NavigationMenu>
</template>

<style scoped>

</style>
