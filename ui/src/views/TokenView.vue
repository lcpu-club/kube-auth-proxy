<script setup lang="ts">
import { client } from "@/api/client";
import { onMounted, ref } from "vue";
import { MessagePlugin } from "tdesign-vue-next";

const tokens = ref<string[]>([]);

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
const kubeconfigDialogVisible = ref(false);

const show_kubeconfig = async (token: string) => {
  const kc = await (await client.get(`/_/tokens/${token}/kubeconfig`)).text();
  kubeconfig_for.value = token;
  kubeconfig.value = kc;
  kubeconfigDialogVisible.value = true;
  console.log(kubeconfigDialogVisible.value);
};

const refreshTokens = async () => {
  tokens.value = (await (await client.get("/_/tokens")).json()).tokens;
  MessagePlugin.success("令牌列表刷新成功");
};

onMounted(async () => {
  await refreshTokens();
});
</script>

<template>
  <h1>令牌管理</h1>
  <t-space style="margin-bottom: 16px">
    <t-button theme="primary" @click="refreshTokens">
      <template #icon><t-icon name="refresh" /></template>
      刷新
    </t-button>
    <t-button theme="success" @click="add_token">
      <template #icon><t-icon name="add" /></template>
      创建令牌
    </t-button>
  </t-space>
  <t-list size="small" split>
    <t-list-item v-for="token in tokens" :key="token">
      <pre style="margin: 0">{{ token }}</pre>
      <template #action>
        <t-space size="small">
          <t-button
            theme="primary"
            variant="text"
            @click="show_kubeconfig(token)"
          >
            <span>Kubeconfig</span>
          </t-button>
          <t-button theme="danger" variant="text" @click="delete_token(token)">
            <span>删除</span>
          </t-button>
        </t-space>
      </template>
    </t-list-item>
  </t-list>

  <t-dialog
    v-model:visible="kubeconfigDialogVisible"
    @confirm="kubeconfigDialogVisible = false"
    :cancel-btn="null"
  >
    <t-textarea
      v-model="kubeconfig"
      autosize
      readonly
      style="font-family: monospace"
    ></t-textarea>
  </t-dialog>
</template>
