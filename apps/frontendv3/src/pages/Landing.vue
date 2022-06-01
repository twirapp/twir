<script lang="ts" setup>
import LandingButtons from '@/components/LandingButtons.vue';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <nav
      class="
        relative
        w-full
        flex flex-wrap
        items-center
        justify-between
        py-4
        navbar navbar-expand-lg navbar-light
        text-white"
    >
      <div class="container-fluid w-full flex flex-wrap items-center justify-between px-6">
        <button
          class="
      navbar-toggler
      text-gray-500
      border-0
      hover:shadow-none hover:no-underline
      py-2
      px-2.5
      bg-transparent
      focus:outline-none focus:ring-0 focus:shadow-none focus:no-underline
    "
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <svg
            aria-hidden="true"
            focusable="false"
            data-prefix="fas"
            data-icon="bars"
            class="w-6"
            role="img"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 448 512"
          >
            <path
              fill="currentColor"
              d="M16 132h416c8.837 0 16-7.163 16-16V76c0-8.837-7.163-16-16-16H16C7.163 60 0 67.163 0 76v40c0 8.837 7.163 16 16 16zm0 160h416c8.837 0 16-7.163 16-16v-40c0-8.837-7.163-16-16-16H16c-8.837 0-16 7.163-16 16v40c0 8.837 7.163 16 16 16zm0 160h416c8.837 0 16-7.163 16-16v-40c0-8.837-7.163-16-16-16H16c-8.837 0-16 7.163-16 16v40c0 8.837 7.163 16 16 16z"
            />
          </svg>
        </button>
        
        <div
          id="navbarSupportedContent"
          class="collapse navbar-collapse flex-grow items-center"
        >
          <a
            class="text-xl"
            href="#"
          >Tsuwari</a>
          <!-- Left links -->
          <ul class="navbar-nav flex flex-col pl-0 list-style-none mr-auto ml-auto">
            <li class="nav-item p-2">
              <a
                class="nav-link hover:text-gray-700 focus:text-gray-700 p-0"
                href="#"
              >Dashboard</a>
            </li>
            <li class="nav-item p-2">
              <a
                class="nav-link hover:text-gray-700 focus:text-gray-700 p-0"
                href="#"
              >Team</a>
            </li>
            <li class="nav-item p-2">
              <a
                class="nav-link hover:text-gray-700 focus:text-gray-700 p-0"
                href="#"
              >Projects</a>
            </li>
          </ul>
          <!-- Left links -->
        </div>
        <!-- Collapsible wrapper -->

        <!-- Right elements -->
        <div class="flex items-center relative">
          <LandingButtons />
        </div>
        <!-- Right elements -->
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
