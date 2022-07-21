<script lang="ts" setup>
  import { createPopper } from '@popperjs/core';

  import Question from '../assets/icons/question.svg?component';
</script>

<script lang="ts">
  export default {
    name: 'RightPinkTooltip',
    // eslint-disable-next-line vue/require-prop-types
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
          createPopper(this.$refs.btnRef as any, this.$refs.tooltipRef as any, {
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
      <Question />
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
