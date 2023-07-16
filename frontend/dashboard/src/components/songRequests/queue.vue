<script setup lang='ts'>
import { IconDots } from '@tabler/icons-vue';
import { NDataTable, NTag, NSpin, NSpace, NText } from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { h } from 'vue';

import { timeAgo, convertMillisToTime } from '@/components/songRequests/helpers.js';
import { Video } from '@/components/songRequests/hook.js';

const props = defineProps<{
	queue: Video[]
}>();
defineEmits<{
	deleteVideo: [id: string]
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
</script>

<template>
  <n-data-table
    :columns="columns"
    :data="queue"
    bordered
    bottom-bordered
    :loading="!queue.length"
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
<!--  <n-button v-for="video of queue" :key="video.id" @click="$emit('deleteVideo', video.id)">-->
<!--    Skip {{ video.id }}-->
<!--  </n-button>-->
</template>

<style>
.queue-song-link {
	color: #63e2b7;
	text-decoration: none
}
</style>
