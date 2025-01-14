<script setup lang="ts">
import { client } from '@/api/api';
import { getToken, hasToken } from '@/api/token';
import { RouterLink } from 'vue-router';

if (!getToken()) {
  window.location.href = "../oauth/redirect";
}

const userInfo = await (await client.get('/_/whoami')).json();
</script>

<template>
  <main>
    <h1>Home</h1>
    <ul>
      <li>
        <RouterLink :to="{ 'name': 'token' }">Token Management</RouterLink>
      </li>
      <li>
        <RouterLink :to="{ 'name': 'logout' }">Logout</RouterLink>
      </li>
    </ul>
    <p>Welcome to the home page!</p>
    <pre>{{ userInfo }}</pre>
  </main>
</template>
