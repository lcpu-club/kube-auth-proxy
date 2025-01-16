<template>
  <div>
    <h1>SSHKeyPair</h1>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchSSHKeyPairs">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 SSHKeyPair
      </t-button>
    </t-space>

    <!-- SSHKeyPair 表格 -->
    <t-table :data="sshKeyPairs" :columns="columns" row-key="metadata.uid">
      <!-- 名称列 -->
      <template #name="{ row }">
        <span>
          {{ row.metadata.name }}
        </span>
      </template>

      <!-- 命名空间列 -->
      <template #namespace="{ row }">
        <t-tag>{{ row.metadata.namespace }}</t-tag>
      </template>

      <!-- 操作列 -->
      <template #operation="{ row }">
        <t-button variant="text" @click="handleDelete(row.metadata.name)">
          删除
        </t-button>
      </template>
    </t-table>

    <!-- 创建 SSHKeyPair 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 SSHKeyPair"
      :footer="false"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        ref="createFormRef"
        @submit="handleCreate"
      >
        <!-- SSHKeyPair 名称 -->
        <t-form-item label="SSHKeyPair 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 SSHKeyPair 名称"
          />
        </t-form-item>

        <!-- 私钥 -->
        <t-form-item label="私钥" name="privateKey">
          <t-textarea
            v-model="createFormData.privateKey"
            placeholder="请输入 SSH 私钥"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </t-form-item>

        <!-- 公钥 -->
        <t-form-item label="公钥" name="publicKey">
          <t-textarea
            v-model="createFormData.publicKey"
            placeholder="请输入 SSH 公钥"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </t-form-item>

        <!-- 提交按钮 -->
        <t-form-item>
          <t-space>
            <t-button theme="primary" type="submit">提交</t-button>
            <t-button theme="default" @click="closeCreateDialog">取消</t-button>
          </t-space>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/api";

// SSHKeyPair 数据
const sshKeyPairs = ref([]);

// 创建 SSHKeyPair 的对话框状态
const createDialogVisible = ref(false);

// 创建 SSHKeyPair 的表单数据
const createFormData = ref({
  name: "",
  privateKey: "",
  publicKey: "",
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "SSHKeyPair 名称不能为空" }],
  privateKey: [{ required: true, message: "私钥不能为空" }],
  publicKey: [{ required: true, message: "公钥不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "namespace", title: "命名空间", cell: "namespace" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

// 用户信息和 API 根路径
let apiRoot = ref("");
let username = ref("");

// 获取用户信息并设置 API 根路径
const fetchUserInfo = async () => {
  try {
    const userInfo = await (await client.get("/_/whoami")).json();
    apiRoot.value = `/api/v1/namespaces/u-${userInfo.username}`;
    username.value = userInfo.username;
  } catch (error) {
    console.error("获取用户信息失败:", error);
    MessagePlugin.error("获取用户信息失败");
  }
};

// 获取 SSHKeyPair 列表
const fetchSSHKeyPairs = async () => {
  try {
    const response = await client.get(`${apiRoot.value}/sshkeypairs`);
    sshKeyPairs.value = response.items;
    MessagePlugin.success("SSHKeyPair 列表刷新成功");
  } catch (error) {
    console.error("获取 SSHKeyPair 失败:", error);
    MessagePlugin.error("获取 SSHKeyPair 失败");
  }
};

// 删除 SSHKeyPair
const handleDelete = async (sshKeyPairName) => {
  try {
    await client.delete(`${apiRoot.value}/sshkeypairs/${sshKeyPairName}`);
    MessagePlugin.success(`SSHKeyPair ${sshKeyPairName} 删除成功`);
    fetchSSHKeyPairs(); // 刷新列表
  } catch (error) {
    console.error("删除 SSHKeyPair 失败:", error);
    MessagePlugin.error("删除 SSHKeyPair 失败");
  }
};

// 创建 SSHKeyPair
const handleCreate = async () => {
  try {
    const sshKeyPairYAML = {
      apiVersion: "v1",
      kind: "SSHKeyPair",
      metadata: {
        name: createFormData.value.name,
        namespace: `u-${username.value}`, // 自动使用 u-${username} 作为 namespace
      },
      privateKey: createFormData.value.privateKey,
      publicKey: createFormData.value.publicKey,
    };

    await client.post(`${apiRoot.value}/sshkeypairs`, sshKeyPairYAML);
    MessagePlugin.success("SSHKeyPair 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchSSHKeyPairs(); // 刷新列表
  } catch (error) {
    console.error("创建 SSHKeyPair 失败:", error);
    MessagePlugin.error("创建 SSHKeyPair 失败");
  }
};

// 显示创建 SSHKeyPair 的对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
};

// 关闭创建 SSHKeyPair 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    privateKey: "",
    publicKey: "",
  };
};

// 组件挂载时获取用户信息和 SSHKeyPair 列表
onMounted(async () => {
  await fetchUserInfo();
  fetchSSHKeyPairs();
});
</script>
