<!-- xmlns attribute not parsed -->
<!-- eslint-disable vue/no-parsing-error -->
<template>
  <component
    :is="icon"
    xmlns="http://www.w3.org/2000/svg"
    :width="width"
    :height="height"
    :stroke="stroke" 
    :fill="fill"
    :stroke-width="strokeWidth"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
</template>

<script lang="ts">
import { defineComponent, FunctionalComponent, PropType, SVGAttributes } from 'vue';

import * as icons from '@/icons';

type Icons = {
  [K in icons.IconName]: FunctionalComponent<SVGAttributes>
};

export default defineComponent({
  name: 'Icon',
  props: {
    name: { type: String as PropType<icons.IconName>, required: true },
    width: { type: String, default: '24px' },
    height: { type: String, default: '24px' },
    fill: { type: String, default: 'none' },
    stroke: { type: String, default: 'white' },
    strokeWidth: { type: Number, default: 1.5 },
  },
  setup: (props) => ({ 
    icon: (icons as Icons)[props.name],
  }),
});
</script>
