<script setup lang="ts">
import { setToken } from "@/api/token";
import { useRoute, useRouter } from "vue-router";

const route = useRoute();
const router = useRouter();

if (route.params.token) {
  setToken(route.params.token as string);
  router.push("/");
}
</script>

<template>
  <main>
    <h1>Auth Token</h1>
    <p>Use this token to authenticate with the API:</p>
    <pre>{{ $route.params.token }}</pre>
    <p>Redirecting to home...</p>
  </main>
</template>
