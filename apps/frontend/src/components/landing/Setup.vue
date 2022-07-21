<script lang="ts" setup>
import Love from '../../assets/love.png?url';

import { redirectToLogin } from '@/functions/redirectToLogin';
</script>

<template>
  <section
    id="install"
    class="mb-32 mx-4 text-center text-white"
  >
    <h2 class="font-bold mb-20 text-3xl">
      Easy install & use 🤯
    </h2>

    <div class="gap-4 grid grid-cols-1 lg:grid-cols-3 md:grid-cols-2">
      <div class="duration-700 hover:-translate-y-2">
        <div class="bg-[#181818] block drop-shadow-xl rounded">
          <div class="p-6">
            <h3 class="font-bold mb-4 text-2xl">
              Log in 🚀
            </h3>
            <h5 class="font-medium text-left text-lg text-white">
              After success login you get access to your personal channel dashboard with all features!
            </h5>

            <div class="flex justify-end mt-5">
              <button
                type="button"
                class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-4 py-2.5 rounded shadow text-white text-xs transition uppercase"
                @click="redirectToLogin"
              >
                Login
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="duration-700 hover:-translate-y-2">
        <div class="bg-[#181818] block drop-shadow-xl rounded">
          <div class="p-6">
            <h3 class="font-bold mb-4 text-2xl">
              Setup settings
            </h3>
            <h5 class="font-medium text-left text-lg">
              Use the sidebar to navigate between settings, enable or disable specific functions and click on the Save button.
            </h5>
          </div>
        </div>
      </div>

      <div class="duration-700 hover:-translate-y-2">
        <div class="bg-[#181818] block drop-shadow-xl rounded">
          <div class="p-6">
            <h3 class="font-bold mb-4 text-2xl">
              Use!
            </h3>

            <div class="items-center md:flex">
              <h5 class="font-medium text-left text-lg">
                After setup you just need to start a stream and... wow, it's work! ✨
              </h5>

              <img
                :src="Love"
                alt="LOVE"
                class="h-24"
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>