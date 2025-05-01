<script setup lang="ts">
const route = useRoute()
const api = useOapi()

const { data, pending, error } = await useAsyncData(
  'hastebin-raw',
  async () => {
    const req = await api.v1.pastebinGetById(route.params.id as string)
    if (req.error) {
      throw req.error
    }

    return req.data
  },
)

// Set content type to plain text
useHead({
  title: `Raw Paste ${route.params.id} - TWIR`,
  meta: [
    { 'http-equiv': 'Content-Type', content: 'text/plain; charset=utf-8' }
  ]
})
</script>

<template>
  <div v-if="pending" class="text-white">Loading...</div>
  <div v-else-if="error" class="text-white">Error: {{ error.message }}</div>
  <pre v-else class="text-white font-mono whitespace-pre-wrap">{{ data?.content }}</pre>
</template>

<style>
body {
  background-color: #1E1E1E;
  margin: 0;
  padding: 16px;
}
</style>
