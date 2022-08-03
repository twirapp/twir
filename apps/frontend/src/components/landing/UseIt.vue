<script lang="ts" setup>
// import Kappa from '@/assets/kappa_start.gif?component';
</script>

<template>
  <div class="bg-[#101010] mb-32 mx-auto py-10 text-center text-white">
    <h2 class="font-bold mb-5 text-3xl text-center text-white">
      Ready to use?
    </h2>

    <div>
      <div class="relative">
        <div class="kappa__button">
          <button
            type="button"
            class="bg-purple-600 ease-in-out file:duration-300 focus:outline-none focus:ring-0 hover:bg-purple-700 hover:scale-105 inline-block kappa__button leading-tight px-14 py-3 relative rounded shadow text-lg text-white transition z-10"
          >
            Start!
          </button>

          <img
            class="kappa"
            src="/src/assets/kappa_start.gif"
          >
        </div>

        

        <p class="text-gray-300">
          It's free and secure
        </p>
      </div>
    </div>
  </div>
</template>