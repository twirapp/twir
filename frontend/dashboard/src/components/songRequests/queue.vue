<script setup lang="ts">
import {
	IconTrash,
	IconChevronUp,
	IconChevronDown,
	IconBan,
} from '@tabler/icons-vue';
import {
	type DataTableCreateSummary,
	NDataTable,
	NSpin,
	NSpace,
	NText,
	NCard,
	NButton,
	NTime,
	NPopconfirm,
} from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { storeToRefs } from 'pinia';
import { h, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useYoutubeSocket, Video } from '@/components/songRequests/hook.js';
import { convertMillisToTime } from '@/helpers/convertMillisToTime.js';

const socket = useYoutubeSocket();
const { videos } = storeToRefs(socket);

const { t } = useI18n();

const columns = computed<TableColumn<Video>[]>(() => [
	{
		title: '#',
		key: 'position',
		width: 50,
		render(_, index) {
			return index + 1;
		},
	},
	{
		title: t('sharedTexts.name'),
		key: 'title',
		ellipsis: true,
		render(row) {
			const banButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => socket.banSong(row.videoId),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						text: true,
						size: 'small',
					}, {
						default: () => h(IconBan),
					}),
					default: () => t('songRequests.ban.songConfirm'),
				},
			);

			return h(
				'div',
				{
					style: 'display: flex; align-items: center; gap: 4px;',
				},
				[
					banButton,
					h(NButton, {
						tag: 'a',
						type: 'primary',
						text: true,
						target: '_blank',
						href: row.songLink,
					}, {
						default: () => row.title,
					}),
				],
			);
		},
	},
	{
		title: t('songRequests.table.columns.author'),
		key: 'author',
		width: 300,
		render(row) {
			const banButton = h(
				NPopconfirm,
				{
					onPositiveClick: () => socket.banUser(row.orderedById),
					positiveText: t('deleteConfirmation.confirm'),
					negativeText: t('deleteConfirmation.cancel'),
				},
				{
					trigger: () => h(NButton, {
						text: true,
						size: 'small',
					}, {
						default: () => h(IconBan),
					}),
					default: () => t('songRequests.ban.userConfirm'),
				},
			);

			return h(
				'div',
				{
					style: 'display: flex; align-items: center; gap: 4px;',
				},
				[
					banButton,
					row.orderedByDisplayName || row.orderedByName,
				],
			);
		},
	},
	{
		title: t('songRequests.table.columns.added'),
		key: 'createdAt',
		width: 150,
		render(row) {
			return h(NTime, {
				time: 0,
				to: Date.now() - new Date(row.createdAt).getTime(),
				type: 'relative',
			});
		},
	},
	{
		title: t('songRequests.table.columns.duration'),
		key: 'duration',
		width: 100,
		render(row) {
			return convertMillisToTime(row.duration * 1000);
		},
	},
	{
		title: '',
		key: 'actions',
		width: 150,
		render(row, index) {
			const deleteButton = h(
				NButton,
				{
					size: 'tiny',
					type: 'error',
					text: true,
					onClick: () => {
						socket.deleteVideo(row.id);
					},
				}, {
					default: () => h(IconTrash),
				},
			);

			const moveUpButton = h(NButton, {
				size: 'tiny',
				type: 'primary',
				text: true,
				disabled: index === 0,
				onClick: () => {
					socket.moveVideo(row.id, index - 1);
				},
			}, {
				default: () => h(IconChevronUp),
			});

			const moveDownButton = h(NButton, {
				size: 'tiny',
				type: 'primary',
				text: true,
				disabled: index + 1 === videos.value.length,
				onClick: () => {
					socket.moveVideo(row.id, index + 1);
				},
			}, {
				default: () => h(IconChevronDown),
			});

			return h(NSpace, {
				justify: 'center',
				align: 'center',
			}, {
				default: () => [
					deleteButton,
					moveUpButton,
					moveDownButton,
				],
			});
		},
	},
]);

const createSummary: DataTableCreateSummary<Video> = (pageData) => {
	return {
		position: {
			value: h(
				'span',
				{ style: 'font-weight: bold;' },
				pageData.length,
			),
			colSpan: 4,
		},
		duration: {
			value: h(
				'span',
				{ style: 'font-weight: bold;' },
				convertMillisToTime(pageData.reduce((acc, cur) => acc + cur.duration * 1000, 0)),
			),
			colSpan: 2,
		},
	};
};
</script>

<template>
	<n-card
		:title="t('songRequests.table.title')"
		content-style="padding: 0;"
		header-style="padding: 10px;"
		segmented
	>
		<template #header-extra>
			<n-button tertiary size="small" @click="$emit('deleteAllVideos')">
				<IconTrash />
			</n-button>
		</template>
		<n-data-table
			:columns="columns"
			:data="videos"
			:loading="!videos.length"
			:bordered="false"
			:summary="createSummary"
		>
			<template #loading>
				<n-space vertical align="center" style="margin-top: 50px;">
					<n-spin :rotate="false" stroke="#959596">
						<template #description>
							<n-text>{{ t('songRequests.waiting') }}</n-text>
						</template>
					</n-spin>
				</n-space>
			</template>
		</n-data-table>
	</n-card>
</template>
