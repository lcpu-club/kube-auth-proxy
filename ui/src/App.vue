<script setup lang="ts">
import { RouterLink, RouterView, useRoute, useRouter } from "vue-router";
import Logo from "@/assets/logo.svg";
const router = useRouter();
const route = useRoute();

const changeHandler = (active: string) => {
  router.push({ path: active });
};
</script>

<template>
  <Suspense>
    <template #fallback>
      <main>
        <h1>Loading...</h1>
      </main>
    </template>
    <div
      style="
        display: grid;
        grid-template-columns: auto 1fr;
        height: 100vh;
        width: 100vw;
        overflow: hidden;
      "
    >
      <div class="box">
        <t-menu
          theme="light"
          @change="changeHandler"
          :style="{
            ...(route.path !== '/canvas/' &&
            route.path !== '/canvas' &&
            route.path !== '/webterm/' &&
            route.path !== '/webterm'
              ? { marginRight: '40px' }
              : {}),
          }"
          :value="route.path"
        >
          <template #logo>
            <img :src="Logo" alt="logo" style="height: 40px" />
            <h2 style="margin: 0.5em">HPCGame</h2>
          </template>
          <t-menu-item value="/">仪表盘</t-menu-item>
          <t-menu-item value="/jobs/">任务（Job）</t-menu-item>
          <t-menu-item value="/pods/">容器（Pods）</t-menu-item>
          <t-menu-item value="/storage/">储存管理</t-menu-item>
          <t-menu-item value="/tokens/">令牌管理</t-menu-item>
          <t-submenu title="配置管理">
            <t-menu-item value="/config/configmap/">ConfigMap</t-menu-item>
            <t-menu-item value="/config/secret/">Secret</t-menu-item>
            <t-menu-item value="/config/sshauthkey/"
              >SSHAuthorizedKeys</t-menu-item
            >
            <t-menu-item value="/config/sshkeypair/">SSHKeyPair</t-menu-item>
          </t-submenu>
          <t-menu-item value="/canvas/">画板</t-menu-item>
        </t-menu>
      </div>
      <div style="position: relative">
        <RouterView />
      </div>
    </div>
  </Suspense>
</template>
