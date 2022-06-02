<script lang="ts" setup>
import { useStore } from '@nanostores/vue';

import LandingButtons from '@/components/LandingButtons.vue';
import { userStore } from '@/stores/userStore';

const user = useStore(userStore);

function redirectToLogin() {
  window.location.replace(`/api/auth?state=${window.btoa(window.location.origin + '/login')}`);
}
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
          class="collapse navbar-collapse select-none flex-grow items-center"
        >
          <svg
            width="20"
            height="19"
            viewBox="0 0 20 19"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M11.2199 2.83999H13.7999V5.17999C13.7999 6.49999 13.9799 6.93999 14.3599 7.31999C14.7199 7.67999 15.3199 7.83999 15.8599 7.83999H17.0599C17.4399 7.83999 17.9599 7.75999 18.2599 7.59999C18.6399 7.39999 18.8799 7.13999 19.0199 6.73999C19.1799 6.33999 19.2799 5.41999 19.3199 4.57999C18.7599 4.37999 17.9599 3.99999 17.5399 3.63999C17.5199 4.37999 17.4999 5.01999 17.4599 5.29999C17.4199 5.57999 17.3399 5.69999 17.2599 5.73999C17.1999 5.77999 17.0599 5.79999 16.9599 5.79999H16.5399C16.4199 5.79999 16.3199 5.77999 16.2599 5.71999C16.1999 5.63999 16.1999 5.45999 16.1999 5.09999V0.759986H8.9599V2.81999C8.9599 4.17999 8.7399 5.91999 6.9399 7.17999C7.4599 7.43999 8.4999 8.11999 8.8999 8.51999C10.8599 7.03999 11.2199 4.71999 11.2199 2.87999V2.83999ZM4.8199 0.0799866C3.9599 1.41999 2.1599 3.05999 0.599902 4.03999C0.979902 4.53999 1.5199 5.51999 1.7799 6.03999C3.6799 4.81999 5.7399 2.83999 7.0599 0.979986L4.8199 0.0799866ZM5.4399 4.25999C4.2399 6.23999 2.2399 8.23999 0.399902 9.49999C0.799902 10.02 1.4399 11.26 1.6599 11.78C2.2399 11.34 2.8199 10.82 3.4199 10.24V18.72H5.7199V7.71999C6.3999 6.87999 7.0399 5.99999 7.5599 5.13999L5.4399 4.25999ZM14.8999 10.8C14.3599 11.82 13.6399 12.72 12.7999 13.52C11.9599 12.72 11.2799 11.82 10.7799 10.8H14.8999ZM16.5399 8.55999L16.0999 8.65999H7.5799V10.8H10.3999L8.6199 11.34C9.2599 12.7 10.0399 13.88 10.9999 14.9C9.6399 15.76 8.1199 16.4 6.4799 16.8C6.9199 17.28 7.4799 18.18 7.7799 18.8C9.5999 18.26 11.2999 17.5 12.7999 16.5C14.2399 17.54 15.9799 18.32 18.0399 18.8C18.3799 18.16 19.0199 17.2 19.5599 16.68C17.6999 16.34 16.0599 15.74 14.6999 14.94C16.2199 13.48 17.3999 11.62 18.1399 9.27999L16.5399 8.55999Z"
              fill="white"
            />
          </svg>

          <a
            class="text-xl ml-2"
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
