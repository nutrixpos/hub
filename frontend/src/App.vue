<template>
  <Toast group="br" position="top-center" />
  <RouterView />
</template>

<script setup lang="ts">
import Toast from 'primevue/toast';
</script>

<style lang="scss">
@use '@/style.scss';
</style>
