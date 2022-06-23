<script lang="ts" setup>
import { useStore } from '@nanostores/vue';

//@ts-ignore
import Logo from '@/assets/logo.svg?component';
import LandingButtons from '@/components/LandingButtons.vue';
import { userStore } from '@/stores/userStore';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <nav class="select-none mt-1.5 p-1 flex items-center w-full justify-between text-white">
      <div>
        <a
          class="ml-2 flex flex-grow items-center"
          href="/"
        >
          <Logo />

          <a
            class="text-xl font-bold ml-2"
            translate="no"
          >Tsuwari</a>
        </a>
      </div>

      <div class="flex space-x-1 hidden md:block">
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
        >Dashboard</a>
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
        >Team</a>
        <a
          href="/dashboard"
          class="inline-block text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
        >About</a>
      </div>

      <div class="hidden md:block">
        <LandingButtons />
      </div>

      <div class="block md:hidden flex justify-center">
        <div>
          <div class="dropdown relative">
            <button
              id="dropdownMenuButton1d"
              class="
          dropdown-toggle
          px-2
          py-2
          text-white
          font-medium
          text-xs
          leading-tight
          uppercase
          rounded
          shadow
          hover:bg-[#9146FF] 
          focus:outline-none focus:ring-0
          transition
          duration-150
          ease-in-out
          flex
          items-center
          whitespace-nowrap
        "
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <svg
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              ><path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16m-7 6h7"
              /></svg>
            </button>
            <ul
              class="
              space-y-1
              bg-[#141414]
              px-1
          dropdown-menu
          min-w-max
          absolute
          text-base
          z-50
          float-left
          py-2
          list-none
          text-left
          rounded
          shadow
          mt-1
          hidden
          m-0
          bg-clip-padding
          border-none
        "
              aria-labelledby="dropdownMenuButton1d"
            >
              <li>
                <a
                  href="/dashboard"
                  class="inline-block w-full text-center text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
                >Dashboard</a>
              </li>
              <li>
                <a
                  href="/dashboard"
                  class="inline-block w-full text-center text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
                >Team</a>
              </li>
              <li>
                <a
                  href="/dashboard"
                  class="inline-block w-full text-center text-white px-6 py-2.5 bg-transparent text-sm leading-tight rounded hover:text-black hover:bg-gray-100  focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
                >About</a>
              </li>
              <hr class="h-0 border border-solid border-t-0 border-gray-700 opacity-50">
              <LandingButtons type="col" />
            </ul>
          </div>
        </div>
      </div>
    </nav>

    <div class="flex text-center h-screen">
      <div class="m-auto text-white">
        <h1 class="text-5xl md:text-7xl break-words outline-none font-bold mb-5">
          The perfect bot for<br> your stream
        </h1>
        <p class="py-6 text-xl">
          Created by streamers. Made
          <span style="color: #772CE8">for streamers</span>.<br>Used
          <span style="color: #772CE8">by streamers</span>. For streamers with love.
        </p>
   
        <LandingButtons size="large" />
      </div>
    </div>

    <img
      src="../assets/bottom.png"
      class="select-none"
      alt="PREVIEW"
    >
  </div>
</template>
