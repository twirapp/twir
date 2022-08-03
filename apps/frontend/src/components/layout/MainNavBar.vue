<script lang="ts" setup>
import LandingButtons from '../LandingButtons.vue';
// eslint-disable-next-line import/order
import MenuSelector from '../../assets/buttons/menu_selector.svg?component';

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import LanguageSelector from '../LanguageSelector.vue';

import Logo from '@/assets/logo.svg?component';
</script>

<template>
  <nav class="bg-[#131313] drop-shadow-lg fixed flex items-center justify-between p-1 pt-1.5 select-none text-white w-full z-10">
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

    <!-- <div class="flex hidden md:block space-x-1">
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
    </div> -->

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
              <MenuSelector />
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
              <!-- <li>
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
              </li> -->
              <hr class="border border-gray-700 border-solid border-t-0 h-0 opacity-50">
              <LandingButtons type="col" />
            </ul>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>