<template>
  <div class="inline-flex items-center relative w-fit dropdown">
    <div class="select-none absolute inline-block top-0 right-0 bottom-auto left-auto translate-x-2/4 -translate-y-1/2 rotate-0 skew-x-0 skew-y-0 scale-x-100 scale-y-100 py-0.5 px-1.5 text-xs leading-none text-center whitespace-nowrap align-baseline font-bold bg-[#772CE8] text-white rounded z-10">
      99+
    </div>
    <button
      id="alertDropdown"
      data-bs-toggle="dropdown"
      aria-expanded="false"
      class="hover:text-slate-300"
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
        d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
      /></svg>
    </button>
    <ul
      class="
              space-y-0.5
          dropdown-menu
          w-96
         scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-600 scrollbar-track-gray-500
          max-h-[55vh]
          absolute
          hidden
          bg-gray-800
          text-base
          z-50
          float-left
          py-2
          text-left
          rounded
          mt-1
          m-0
          bg-clip-padding
          border-none
        "
      aria-labelledby="alertDropdown"
    >
      <li>
        <p
          class="    
              text-sm
              py-2
              px-2
              font-normal
              w-full
              hover:bg-gray-600
            "
        >
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Neque, facere. Quia obcaecati, veniam voluptatem minima iste deserunt excepturi odio asperiores tenetur ipsum iure nemo officiis, corporis quam ea odit nesciunt?
        </p>
      </li>
    </ul>
  </div>
</template>