<script setup lang="ts">
import { client } from "@/api/api";
import { ref } from "vue";
import { RouterLink } from "vue-router";
import {
  ChevronLeftIcon,
  EyeIcon,
  TrashIcon,
  PlusIcon,
  EyeSlashIcon,
} from "@heroicons/vue/24/outline";

const tokens = ref(
  (await (await client.get("/_/tokens")).json()).tokens as string[]
);

const add_token = async () => {
  const token = await (await client.post("/_/tokens", {})).json();
  tokens.value.push(token.token);
};

const delete_token = async (token: string) => {
  await client.delete(`/_/tokens/${token}`);
  tokens.value = tokens.value.filter((t) => t !== token);
};

const kubeconfig_for = ref("");
const kubeconfig = ref("");

const show_kubeconfig = async (token: string) => {
  const kc = await (await client.get(`/_/tokens/${token}/kubeconfig`)).text();
  kubeconfig_for.value = token;
  kubeconfig.value = kc;
};
</script>

<template>
  <h1 class="flex items-center">
    <span class="flex-grow">令牌管理</span>
    <button @click="add_token()" class="icon-button">
      <PlusIcon class="icon" />
      <span>添加令牌</span>
    </button>
  </h1>
  <div class="token-wrapper" v-for="token in tokens" :key="token">
    <pre class="token-cell m-0">{{ token }}</pre>
    <button @click="show_kubeconfig(token)" class="icon-button">
      <EyeIcon class="icon" />
      <span>显示 Kubeconfig</span>
    </button>
    <button @click="delete_token(token)" class="icon-button">
      <TrashIcon class="icon" />
      <span>删除</span>
    </button>
  </div>
  <div v-if="kubeconfig" class="m-t-2">
    <div class="flex gap-2 items-center text-no-wrap">
      <span class="font-bold">Kubeconfig for</span>
      <pre class="inline flex-grow flex-shrink overflow-auto font-bold">{{
        kubeconfig_for
      }}</pre>
      <button @click="kubeconfig = ''" class="icon-button">
        <EyeSlashIcon class="icon" />
        <span>隐藏</span>
      </button>
    </div>
    <textarea class="m-t-2" rows="16" cols="80" disabled>{{
      kubeconfig
    }}</textarea>
  </div>
</template>

<style scoped>
.m-0 {
  margin: 0;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.gap-2 {
  gap: 0.5rem;
}

.col {
  flex-direction: column;
}

.m-t-2 {
  margin-top: 0.5rem;
}

.icon {
  width: 1em;
  height: 1em;
  margin-right: 0.3em;
}

.title-back {
  text-decoration: none;
  display: flex;
  align-items: center;
  color: unset;
}

button {
  padding: 0.4em 0.8em;
  border: none;
  background-color: transparent;
  border-radius: 0.3em;
  width: max-content;
  transition: background-color 0.1s;
}

button:hover {
  background-color: rgba(100, 100, 100, 0.1);
}

.icon-button {
  display: flex;
  align-items: center;
  padding: 1em;
  text-wrap-mode: nowrap;
}

.token-cell {
  overflow: auto;
  text-overflow: unset;
  text-wrap-mode: nowrap;
  flex-grow: 1;
  padding: 1em 1em;
}

.token-wrapper {
  display: flex;
  align-items: stretch;
}

.token-wrapper:nth-child(even) {
  background-color: rgba(200, 200, 200, 0.1);
}

.flex-grow {
  flex-grow: 1;
}

.flex-shrink {
  flex-shrink: 1;
}

.inline {
  display: inline;
}

.text-no-wrap {
  text-wrap-mode: nowrap;
}

.overflow-auto {
  overflow: auto;
}

.font-bold {
  font-weight: bold;
}

textarea {
  max-width: 100%;
}
</style>
