    <script lang="ts">
    import { createPopper } from '@popperjs/core';
    
    export default {
      name: 'RightPinkTooltip',
      props:['text'],

      data() {
        return {
          tooltipShow: false,
        };
      },
      methods: {
        toggleTooltip: function(){
          if(this.tooltipShow){
            this.tooltipShow = false;
          } else {
            this.tooltipShow = true;
            createPopper(this.$refs.btnRef, this.$refs.tooltipRef, {
              placement: 'bottom',
            });
          }
        },
      },
    };
</script>

<template>
  <div
    class="flex flex-wrap"
  >
    <div
      ref="btnRef"
      type="button"
      @mouseenter="toggleTooltip()"
      @mouseleave="toggleTooltip()"
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
        d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
      /></svg>
    </div>

    <div
      ref="tooltipRef"
      :class="{'hidden': !tooltipShow, 'block': tooltipShow}"
      class="bg-[#111111] block break-words px-3 py-1 rounded w-auto z-50"
    >
      <div>
        <p
          class="font-light text-left"
        >
          {{ text }}
        </p>
      </div>
    </div>
  </div>
</template>
