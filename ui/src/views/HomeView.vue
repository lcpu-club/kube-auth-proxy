<script setup lang="ts">
import { client } from '@/api/api';
import { getToken, hasToken } from '@/api/token';
import { RouterLink } from 'vue-router';
import { ArrowLeftStartOnRectangleIcon, KeyIcon } from '@heroicons/vue/24/outline';

if (!getToken()) {
  window.location.href = "../oauth/redirect";
}

const userInfo = await (await client.get('/_/whoami')).json();
</script>

<template>
  <main>
    <h1 class="page-title">主页</h1>
    <h2 class="m-0">{{ userInfo.extra.name }}</h2>
    <div class="flex gap-2 text-gray text-sm">
      <span>{{ userInfo.extra.realname }}</span>
      <span>{{ userInfo.extra.school }}</span>
      <span>{{ userInfo.extra.studentGrade }}</span>
    </div>
    <div class="flex col m-t-4 gap-2">
      <RouterLink :to="{ 'name': 'token' }" class="link">
        <KeyIcon class="icon" />
        <span>凭据管理</span>
      </RouterLink>
      <RouterLink :to="{ 'name': 'logout' }" class="link">
        <ArrowLeftStartOnRectangleIcon class="icon" />
        <span>退出登录</span>
      </RouterLink>
    </div>
  </main>
</template>

<style scoped>
.m-0 {
  margin: 0;
}

.flex {
  display: flex;
}

.gap-2 {
  gap: 0.5rem;
}

.text-gray {
  color: #666;
}

.text-sm {
  font-size: 0.875rem;
}

main {
  padding: 3rem;
}

.col {
  flex-direction: column;
}

.m-t-4 {
  margin-top: 1rem;
}

.icon {
  width: 1em;
  height: 1em;
  margin-right: .3em;
}

.link {
  color: unset;
  text-decoration: none;
  display: flex;
  align-items: center;
  border-bottom: #666 solid 1px;
  width: fit-content;
}

@media screen and (max-width: 768px) {
  main {
    padding: 1rem;
  }
}
</style>
