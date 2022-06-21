<script lang="ts" setup>
import { useStore } from '@nanostores/vue';

//@ts-ignore
import Logo from '@/assets/logo.svg?component';
import LandingButtons from '@/components/LandingButtons.vue';
import { userStore } from '@/stores/userStore';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <nav class="p-1 flex items-center w-full justify-between text-white">
      <div>
        <a
          class="ml-2 flex flex-grow items-center"
          href="/"
        >
          <Logo />

          <a
            class="text-xl font-bold ml-2"
          >Tsuwari</a>
        </a>
      </div>

      <div class="flex space-x-1">
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-blue-700 hover:bg-gray-100 focus:bg-gray-100 focus:outline-none focus:ring-0 active:bg-gray-200 transition duration-150 ease-in-out"
        >Dashabord</a>
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-blue-700 hover:bg-gray-100 focus:bg-gray-100 focus:outline-none focus:ring-0 active:bg-gray-200 transition duration-150 ease-in-out"
        >Team</a>
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-blue-700 hover:bg-gray-100 focus:bg-gray-100 focus:outline-none focus:ring-0 active:bg-gray-200 transition duration-150 ease-in-out"
        >About</a>
      </div>

      <div>
        <LandingButtons />
      </div>
    </nav>

    <!-- Jumbotron -->
    <div class="p-12 mb-14 mt-32 text-center text-white">
      <h1 class="text-7xl outline-none font-bold">
        The perfect bot for your stream
      </h1>
      <p class="py-6 text-xl">
        Created by streamers. Made
        <span style="color: #9146ff">for streamers</span>.<br>Used
        <span style="color: #9146ff">by streamers</span>. For streamers with love.
      </p>
   
      <LandingButtons size="large" />
    </div>
    <!-- Jumbotron -->

    <img
      src="../assets/bottom.png"
      class="select-none"
      alt="PREVIEW"
    >
  </div>
</template>
