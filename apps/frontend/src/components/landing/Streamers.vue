<script lang="ts" setup>
import Verified from '../../assets/icons/verified.svg?component';
</script>

<template>
  <div class="gap-4 grid grid-cols-1 lg:grid-cols-3 md:grid-cols-2">
    <div class="flex justify-center">
      <div class="bg-[#181818] block drop-shadow-xl max-w-sm p-6 rounded text-white">
        <div class="flex flex-col items-center">
          <div class="flex items-center justify-around mb-5">
            <img
              src="https://i.picsum.photos/id/57/512/512.jpg?hmac=tup-ETUL0LHChWZTjAG31HXqECqaxgngmRKP87zOjDs"
              alt="LOGO"
              class="flex justify-center rounded-full w-1/3"
            >

            <h5 class="flex flex-col my-2 text-xl">
              <div class="flex items-center space-x-2">
                <span class="font-bold">LWGerry</span>
                <Verified />
              </div>

              <span class="italic text-lg">4.5k followers</span>
            </h5>
          </div>
          <p class="break-words mb-4">
            Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.
          </p>
        </div>

        <div class="flex justify-end">
          <button
            type="button"
            class="active:bg-blue-800 active:shadow-lg bg-blue-600 duration-150 ease-in-out focus:bg-blue-700 focus:outline-none focus:ring-0 focus:shadow-lg font-medium hover:bg-blue-700 hover:shadow-lg inline-block leading-tight px-6 py-2.5 rounded shadow-md text-white text-xs transition uppercase"
          >
            Button
          </button>
        </div>
      </div>
    </div>
  </div>
</template>