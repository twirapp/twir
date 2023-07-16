<script setup lang='ts'>
import {
	type DataTableCreateSummary,
	NDataTable,
	NTag,
	NSpin,
	NSpace,
	NText,
	NCard,
} from 'naive-ui';
import type { TableColumn, SummaryRowData } from 'naive-ui/es/data-table/src/interface';
import { h } from 'vue';

import { timeAgo, convertMillisToTime } from '@/components/songRequests/helpers.js';
import { Video } from '@/components/songRequests/hook.js';

const props = defineProps<{
	queue: Video[]
}>();
defineEmits<{
	deleteVideo: [id: string]
	deleteAllVideos: []
}>();

const columns: TableColumn<Video>[] = [
	{
		title: '#',
		key: 'position',
		width: 25,
		render(row, index) {
			return index+1;
		},
	},
	{
		title: 'Title',
		key: 'title',
		render(row) {
			return h('a', {
				class: 'queue-song-link',
				href: `https://youtu.be/${row.videoId}`,
			}, {
				default: () => row.title,
			});
		},
	},
	{
		title: 'Added',
		key: 'createdAt',
		render(row) {
			return timeAgo(row.createdAt);
		},
	},
	{
		title: 'Author',
		key: 'author',
		render(row) {
			return h(NTag, { bordered: false, type: 'info' }, { default: () => row.orderedByDisplayName || row.orderedByName });
		},
	},
	{
		title: 'Duration',
		key: 'duration',
		render(row) {
			return convertMillisToTime(row.duration);
		},
	},
];

const createSummary: DataTableCreateSummary<Video> = (pageData) => {
	return{
		position: {
			value: h(
				'span',
				{ },
				pageData.length,
			),
			colSpan: 4,
		},
		duration: {
			value: h(
				'span',
				{ style: 'font-weight: bold;' },
				convertMillisToTime(pageData.reduce((acc, cur) => acc + cur.duration, 0)),
			),
			colSpan: 1,
		},
	};
};
</script>

<template>
  <n-card
    title="Current Song"
    content-style="padding: 0;"
    header-style="padding: 10px;"
    segmented
  >
    <n-data-table
      :columns="columns"
      :data="queue"
      :loading="!queue.length"
      :bordered="false"
      :summary="createSummary"
    >
      <template #loading>
        <n-space vertical align="center" style="margin-top: 50px;">
          <n-spin :rotate="false" stroke="#959596">
            <template #description>
              <n-text>Waiting for songs</n-text>
            </template>
          </n-spin>
        </n-space>
      </template>
    </n-data-table>

    <template #footer>
      footer
    </template>
  </n-card>
</template>

<style>
.queue-song-link {
	color: #63e2b7;
	text-decoration: none
}
</style>
