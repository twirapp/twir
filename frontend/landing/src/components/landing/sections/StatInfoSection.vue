<template>
  <section class="bg-black-15 border-b border-b-black-25">
    <div class="container">
      <div class="flex min-md:w-[740px] mx-auto py-[18px]">
        <ClientOnly>
          <template #default>
            <StatsList />
          </template>
          <template #server>
            <div class="w-full h-[70px]"></div>
          </template>
        </ClientOnly>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import ClientOnly from '@/components/ClientOnly.vue';
import StatsList from '@/components/landing/StatsList.vue';
</script>
