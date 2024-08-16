<script lang="ts" setup>
import ButtonToggleTheme from '../buttons/button-toggle-theme.vue'
import DropdownLanguage from '../dropdowns/dropdown-language.vue'
import DropdownProfile from '../dropdowns/dropdown-profile'
import NotificationsButton from '../notifications/notifications-button.vue'
import StatsScrollable from '../stats/stats-scrollable.vue'
</script>

<template>
	<div class="drawerHeader">
		<StatsScrollable />

		<div class="drawerSlot">
			<div class="drawerSlotBlocks">
				<div style="display: flex; align-items: center; column-gap: 12px">
					<DropdownLanguage />
					<ButtonToggleTheme />
				</div>
			</div>

			<div style="display: flex; align-items: center; column-gap: 8px">
				<NotificationsButton />
				<DropdownProfile />
			</div>
		</div>
	</div>
</template>

<style scoped>
.drawerHeader {
	display: flex;
	flex-direction: column;
	row-gap: 6px;
}

.drawerSlot {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 8px;
}

.drawerSlotBlocks {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-right: 6px;
}

@media screen and (max-width: 520px) {
	.drawerSlot {
		align-items: flex-start;
	}

	.drawerSlotBlocks {
		flex-direction: column;
	}
}
</style>
