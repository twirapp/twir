<script setup lang="ts">
import { type ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NFormItem, NInput, NInputNumber, useThemeVars } from 'naive-ui';
import { ref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';

const theme = useThemeVars();
const { t } = useI18n();

const props = defineProps<{
	item: ItemWithId
}>();
const formValue = ref(toRaw(props.item));

</script>

<template>
	<div>
		<div class="form-block">
			<n-form-item label="Timeout message">
				<n-input
					v-model:value="formValue.data!.banMessage"
					type="textarea"
					:maxLength="500"
				/>
			</n-form-item>

			<n-form-item label="Ban time">
				<n-input-number
					v-model:value="formValue.data!.banTime"
					:min="0"
					:max="86400"
					style="width: 100%"
				/>
			</n-form-item>
		</div>
	</div>
	{{ formValue }}
</template>

<style scoped>
.form-block {
	display: flex;
	flex-direction: column;
	padding: 8px;
	background-color: v-bind('theme.buttonColor2');
	border-radius: 8px;
	width: 100%;
}
</style>
