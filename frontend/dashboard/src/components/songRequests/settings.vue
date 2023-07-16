<script setup lang='ts'>
import {
	NTab,
	NTabPane,
} from 'naive-ui';
</script>

<template>
  <n-tab
    class="card-tabs"
    default-value="general"
    size="large"
    animated
    pane-wrapper-style="margin: 0 -4px"
    pane-style="padding-left: 4px; padding-right: 4px; box-sizing: border-box;"
  >
    <n-tab-pane name="general" tab="General">
      first
    </n-tab-pane>

    <n-tab-pane name="users" tab="Users">
      s
    </n-tab-pane>

    <n-tab-pane name="songs" tab="Songs">
      t
    </n-tab-pane>

    <n-tab-pane name="translations" tab="translations">
      t
    </n-tab-pane>
  </n-tab>
</template>

