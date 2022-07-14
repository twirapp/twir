<script lang="ts" setup>
import LandingButtons from '../LandingButtons.vue';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import LanguageSelector from '../LanguageSelector.vue';

import Logo from '@/assets/logo.svg?component';
</script>

<template>
  <nav class="flex items-center justify-between mt-1.5 p-1 select-none text-white w-full">
    <div>
      <a
        class="flex flex-grow items-center ml-2"
        href="/"
      >
        <Logo />

        <a
          class="font-bold ml-2 text-xl"
          translate="no"
        >Tsuwari</a>
      </a>
    </div>

    <div class="flex hidden md:block space-x-1">
      <a
        href="/dashboard"
        class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-sm text-white transition"
      >Dashboard</a>
      <a
        href="#"
        class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-sm text-white transition"
      >Team</a>
      <a
        href="#"
        class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-sm text-white transition"
      >About</a>
    </div>

    <div class="flex md:space-x-5 space-x-2">
      <LanguageSelector />

      <div class="hidden md:block">
        <LandingButtons />
      </div>

      <div class="block flex justify-center md:hidden">
        <div>
          <div class="dropdown relative">
            <button
              id="dropdownMenuButton1d"
              class="dropdown-toggle
          duration-150
          ease-in-out
          flex
          focus:outline-none
          focus:ring-0
          font-medium
          hover:bg-[#9146FF]
          items-center
          leading-tight
          px-2
          py-2
          rounded
          shadow
          text-white
          text-xs
          transition
          uppercase
          whitespace-nowrap"
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <svg
                class="h-6 w-6"
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
              class="absolute
              bg-[#141414]
              bg-clip-padding
              border-none
              dropdown-menu
              float-left
              hidden
              list-none
              m-0
              min-w-max
              mt-1
              px-1
              py-2
              rounded
              shadow
              space-y-1
              text-base
              text-left
              z-50"
              aria-labelledby="dropdownMenuButton1d"
            >
              <li>
                <a
                  href="/dashboard"
                  class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-center text-sm text-white transition w-full"
                >Dashboard</a>
              </li>
              <li>
                <a
                  href="#"
                  class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-center text-sm text-white transition w-full"
                >Team</a>
              </li>
              <li>
                <a
                  href="#"
                  class="bg-transparent duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-white hover:text-black inline-block leading-tight px-6 py-2.5 rounded text-center text-sm text-white transition w-full"
                >About</a>
              </li>
              <hr class="border border-gray-700 border-solid border-t-0 h-0 opacity-50">
              <LandingButtons type="col" />
            </ul>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>