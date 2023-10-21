<script setup lang="ts">
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { onMounted, ref, toValue } from 'vue';

const props = defineProps<{
	item?: ItemWithId
}>();
const formValue = ref<ItemWithId>();

onMounted(() => {
	formValue.value = toValue(props.item);
});

</script>

<template>
	{{ formValue }}
</template>
