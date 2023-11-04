<script lang="ts" setup>
import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_kappagen';
import { useNotification } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const formValue = ref<Settings>({});

const message = useNotification();
const { t } = useI18n();

async function save() {
	if (!formValue.value) return;

	// await updater.mutateAsync(formValue.value);
	message.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}
</script>

<tempalte>
	settings of kappagen
</tempalte>
