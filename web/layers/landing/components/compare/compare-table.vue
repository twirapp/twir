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
</script>

<template>
	<div
		class="overflow-x-auto rounded-2xl border border-[#72757d26] bg-[#101014]/80 shadow-[0_0_80px_rgba(93,88,245,0.08)]"
	>
		<table class="w-full min-w-[820px] border-collapse text-left">
			<caption class="sr-only">
				{{
					t('compare.table.caption')
				}}
			</caption>
			<thead>
				<tr class="border-b border-[#72757d26]">
					<th
						scope="col"
						class="p-5 align-bottom text-sm font-medium tracking-wider text-[#ADB0B8] uppercase"
					>
						{{ t('compare.table.feature') }}
					</th>
					<th
						v-for="bot of compareBots"
						:key="bot.id"
						scope="col"
						class="p-5 align-bottom"
						:class="botCellClass(bot)"
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
							<span class="text-sm font-semibold text-white">
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
					class="border-b border-[#72757d26]/60 transition-colors last:border-b-0 hover:bg-white/[0.02]"
				>
					<th
						scope="row"
						class="p-5 text-base font-medium text-white"
					>
						{{ t(row.labelKey) }}
					</th>
					<td
						v-for="bot of compareBots"
						:key="bot.id"
						class="p-5"
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
					class="border-b border-[#72757d26]/60 transition-colors last:border-b-0 hover:bg-white/[0.02]"
				>
					<th
						scope="row"
						class="p-5 text-base font-medium text-white"
					>
						{{ t(row.labelKey) }}
					</th>
					<td
						v-for="bot of compareBots"
						:key="bot.id"
						class="p-5 text-center text-sm text-[#ADB0B8]"
						:class="botCellClass(bot)"
					>
						{{ t(row.cells[bot.id]) }}
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
