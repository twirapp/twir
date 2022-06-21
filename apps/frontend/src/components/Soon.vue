<template>
  <div class="flex h-full">
    <div class="m-auto text-center space-y-2">
      <h1 class="font-bold text-6xl">
        Oops!
      </h1>
      <p
        class="text-lg"
      >
        This feature is coming soon!
      </p>

      <button
        class="inline-block px-10 py-2.5 bg-[#9146FF] text-white text-lg leading-tight uppercase rounded shadow-md hover:bg-[#772CE8] hover:shadow-lg focus:bg-[#772CE8] focus:shadow-lg focus:outline-none focus:ring-0 active:bg-[#772CE8] active:shadow-lg transition duration-150 ease-in-out"
        @click="$router.go(-1)"
      >
        Back
      </button>
    </div>
  </div>
</template>