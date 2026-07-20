<script setup lang="ts">
import TwirLogo from '~/components/twir-logo.vue'

import {
	type CompareBot,
	type CompareCell,
	compareBots,
	compareFeatureRows,
	compareTextRows,
} from './compare.data'

const { t } = useI18n()

const supportMeta: Record<CompareCell['support'], { icon: string; class: string; label: string }> =
	{
		yes: { icon: 'lucide:check', class: 'text-emerald-400', label: 'compare.support.yes' },
		partial: { icon: 'lucide:minus', class: 'text-amber-400', label: 'compare.support.partial' },
		no: { icon: 'lucide:x', class: 'text-red-400/70', label: 'compare.support.no' },
	}

function botCellClass(bot: CompareBot): string {
	return bot.id === 'twir' ? 'bg-[#5D58F5]/10 border-x border-[#5D58F5]/25' : ''
}

function botHeaderCellClass(bot: CompareBot): string {
	return bot.id === 'twir'
		? 'bg-[#171629] border-x border-[#5D58F5]/25'
		: 'bg-[#101014]'
}

const wrapperRef = ref<HTMLElement>()

function updateHeaderHeightVar() {
	const siteHeader = document.querySelector('header#top')
	if (!siteHeader || !wrapperRef.value) return

	const { height } = siteHeader.getBoundingClientRect()
	wrapperRef.value.style.setProperty('--site-header-h', `${height - 1}px`)
}

onMounted(() => {
	updateHeaderHeightVar()
	window.addEventListener('resize', updateHeaderHeightVar)
})

onUnmounted(() => {
	window.removeEventListener('resize', updateHeaderHeightVar)
})
</script>

<template>
	<div
		ref="wrapperRef"
		class="overflow-clip rounded-2xl border border-[#72757d26] bg-[#101014]/80 shadow-[0_0_80px_rgba(93,88,245,0.08)]"
	>
		<table class="compare-table w-full border-separate border-spacing-0 text-left">
			<caption class="sr-only">
				{{
					t('compare.table.caption')
				}}
			</caption>
			<thead>
				<tr>
					<th
						scope="col"
						class="sticky top-[var(--site-header-h,67px)] z-10 border-b border-[#72757d26] bg-[#101014] p-2.5 align-bottom text-sm font-medium tracking-wider text-[#ADB0B8] uppercase sm:p-5"
					>
						{{ t('compare.table.feature') }}
					</th>
					<th
						v-for="bot of compareBots"
						:key="bot.id"
						scope="col"
						class="sticky top-[var(--site-header-h,67px)] z-10 border-b border-[#72757d26] p-2.5 align-bottom sm:p-5"
						:class="botHeaderCellClass(bot)"
					>
						<a
							:href="bot.siteUrl"
							target="_blank"
							rel="noopener noreferrer"
							class="flex flex-col items-center justify-center gap-2"
						>
							<span class="flex h-9 items-center justify-center">
								<TwirLogo
									v-if="bot.logo === 'twir-logo'"
									alt=""
									:class="bot.logoClass"
								/>
								<Icon
									v-else-if="bot.logoType === 'icon'"
									:name="bot.logo"
									:class="bot.logoClass"
								/>
								<img
									v-else
									:src="bot.logo"
									:alt="bot.name"
									:class="bot.logoClass"
									loading="lazy"
								/>
							</span>
							<span class="text-center text-xs font-semibold text-white sm:text-sm">
								{{ bot.name }}
							</span>
						</a>
					</th>
				</tr>
			</thead>
			<tbody>
				<tr
					v-for="row of compareFeatureRows"
					:key="row.labelKey"
					class="transition-colors hover:bg-white/[0.02]"
				>
					<th
						scope="row"
						class="border-b border-[#72757d26]/60 p-2.5 text-sm font-medium text-white sm:p-5 sm:text-base"
					>
						{{ t(row.labelKey) }}
					</th>
					<td
						v-for="bot of compareBots"
						:key="bot.id"
						class="border-b border-[#72757d26]/60 p-2.5 sm:p-5"
						:class="botCellClass(bot)"
					>
						<div class="flex items-center justify-center gap-1.5">
							<Icon
								:name="supportMeta[row.cells[bot.id].support].icon"
								class="h-5 w-5"
								:class="supportMeta[row.cells[bot.id].support].class"
								:aria-label="t(supportMeta[row.cells[bot.id].support].label)"
							/>
							<span
								v-if="row.cells[bot.id].noteKey"
								class="cursor-help text-[#ADB0B8] transition-colors hover:text-white"
								:title="t(row.cells[bot.id].noteKey!)"
							>
								<Icon
									name="lucide:info"
									class="h-4 w-4"
									aria-hidden="true"
								/>
								<span class="sr-only">{{ t(row.cells[bot.id].noteKey!) }}</span>
							</span>
						</div>
					</td>
				</tr>
				<tr
					v-for="row of compareTextRows"
					:key="row.labelKey"
					class="transition-colors hover:bg-white/[0.02]"
				>
					<th
						scope="row"
						class="border-b border-[#72757d26]/60 p-2.5 text-sm font-medium text-white sm:p-5 sm:text-base"
					>
						{{ t(row.labelKey) }}
					</th>
					<td
						v-for="bot of compareBots"
						:key="bot.id"
						class="border-b border-[#72757d26]/60 p-2.5 text-center text-sm text-[#ADB0B8] sm:p-5"
						:class="botCellClass(bot)"
					>
						{{ t(row.cells[bot.id]) }}
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<style scoped>
.compare-table tbody tr:last-child > th,
.compare-table tbody tr:last-child > td {
	border-bottom: none;
}
</style>
