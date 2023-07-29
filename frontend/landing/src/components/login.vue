<script lang="ts" setup>
import { onMounted } from 'vue';

import { browserUnProtectedClient } from '../api/twirp-browser.js';

const url = new URL(window.location.href);
const code = url.searchParams.get('code');
const err = url.searchParams.get('error');

onMounted(async () => {
	if (!code) return window.location.replace('/');

	await browserUnProtectedClient.authPostCode({
		code,
	});
	window.location.replace('/dashboard');
});
</script>
