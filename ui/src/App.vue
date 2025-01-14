<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <Suspense>
    <RouterView />
    <template #fallback>
      <main>
        <h1>Loading...</h1>
      </main>
    </template>
  </Suspense>
</template>
