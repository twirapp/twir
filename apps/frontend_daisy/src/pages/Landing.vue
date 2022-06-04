<script lang="ts" setup>
import { useStore } from '@nanostores/vue';

import { userStore } from '@/stores/userStore';

const user = useStore(userStore);

function redirectToLogin() {
  window.location.replace(`/api/auth?state=${window.btoa(window.location.origin)}/login`);
}
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <div class="navbar text-white">
      <div class="navbar-start">
        <svg
          class="ml-2"
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
          class="ml-2 text-xl font-bold select-none"
          href="/"
          translate="no"
        >Tsuwari</a>
      </div>

      <div class="navbar-center hidden lg:flex">
        <ul class="menu menu-horizontal p-0 space-x-1">
          <!-- @TODO -->
          <li>
            <a
              href="/"
              class="btn-ghost btn"
            >Home</a>
          </li>
          <li>
            <a
              href="/"
              class="btn-ghost btn"
            >Product</a>
          </li>
          <li>
            <a
              href="/"
              class="btn-ghost btn"
            >Features</a>
          </li>
          <li>
            <a
              href="/"
              class="btn-ghost btn"
            >Contact</a>
          </li>
        </ul>
      </div>

      <div class="navbar-end">
        <div class="dropdown lg:hidden dropdown-end">
          <!-- svelte-ignore a11y-label-has-associated-control -->
          <label
            tabindex="0"
            class="btn btn-ghost btn-circle"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h7"
              />
            </svg>
          </label>

          <ul
            tabindex="0"
            class="menu menu-compact dropdown-content mt-3 p-2 bg-base-200 rounded-box w-52"
          >
            <!-- @TODO -->
            <li>
              <a
                href="/"
                class="btn-ghost btn"
              >Home</a>
            </li>
            <li>
              <a
                href="/"
                class="btn-ghost btn"
              >Product</a>
            </li>
            <li>
              <a
                href="/"
                class="btn-ghost btn"
              >Features</a>
            </li>
            <li>
              <a
                href="/"
                class="btn-ghost btn"
              >Contact</a>
            </li>

            <div class="divider" />

            <router-link
              v-if="user"
              to="/dashboard"
              class="btn btn-outline rounded mb-2"
            >
              Dashboard
            </router-link>
            <a
              v-else
              class="btn btn-outline rounded mb-2"
              @click="redirectToLogin"
            >Login</a>
          </ul>
        </div>

        <div class="hidden lg:block">
          <router-link
            v-if="user"
            to="/dashboard"
            class="btn btn-outline rounded"
          >
            Dashboard
          </router-link>
          <a
            v-else
            class="btn btn-primary rounded"
            @click="redirectToLogin"
          >Login</a>
        </div>
      </div>
    </div>

    <div class="hero h-fit mb-14 mt-32 text-white">
      <div class="hero-content text-center">
        <div class="md:w-2/3">
          <h1 class="text-7xl outline-none font-bold">
            The perfect bot for your stream
          </h1>

          <p class="py-6 text-xl">
            Created by streamers. Made
            <span style="color: #9146ff">for streamers</span>.<br>Used
            <span style="color: #9146ff">by streamers</span>. For streamers with love.
          </p>

          <!-- @TODO -->
          <a
            class="btn btn-outline mr-2 rounded"
            href="/"
          >Learn More</a>
          <a
            class="btn btn-primary rounded"
            href="/"
          >Try for free</a>
        </div>
      </div>
    </div>

    <img
      src="../assets/bottom.png"
      class="select-none"
      alt="PREVIEW"
    >
  </div>
</template>
