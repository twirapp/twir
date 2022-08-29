<template>
  <router-view v-slot="{ Component }">
    <Suspense>
      <component :is="Component" />
    </Suspense>
  </router-view>
  <router-link to="/dashboard">
    Dashboard
  </router-link>
  <router-link to="/commands">
    Commands
  </router-link>
</template>