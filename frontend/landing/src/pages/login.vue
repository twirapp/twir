<script lang="ts" setup>
// import { unProtectedClient } from "../api/twirp.js";

// const url = new URL(Astro.request.url)
// const code = url.searchParams.get('code')
// let error = url.searchParams.get('error')

// if (code) {
// 	try {
// 		await unProtectedClient.authPostCode({
// 			code,
// 		})
// 	} catch (e) {
// 		console.error(e)
// 		error = 'internal error'
// 	}
// }
</script>

<template>
	{{ window.location.origin }}
</template>
