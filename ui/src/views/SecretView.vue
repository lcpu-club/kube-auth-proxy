<template>
  <div>
    <h1>Secret</h1>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchSecrets">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 Secret
      </t-button>
    </t-space>

    <!-- Secret 表格 -->
    <t-table :data="secrets" :columns="columns" row-key="metadata.uid">
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

      <!-- 类型列 -->
      <template #type="{ row }">
        <t-tag>{{ row.type }}</t-tag>
      </template>

      <!-- 操作列 -->
      <template #operation="{ row }">
        <t-button
          variant="text"
          @click="handleViewData(row)"
          shape="square"
          title="显示数据"
        >
          <t-icon name="data-search" />
        </t-button>
        <t-button
          variant="text"
          @click="handleDelete(row.metadata.name)"
          shape="square"
          title="删除"
        >
          <t-icon name="delete-1" />
        </t-button>
      </template>
    </t-table>

    <!-- 创建 Secret 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 Secret"
      :footer="false"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        ref="createFormRef"
        @submit="handleCreate"
      >
        <!-- Secret 名称 -->
        <t-form-item label="Secret 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 Secret 名称"
          />
        </t-form-item>

        <!-- 类型 -->
        <t-form-item label="类型" name="type">
          <t-select v-model="createFormData.type">
            <t-option value="Opaque">Opaque</t-option>
            <t-option value="kubernetes.io/tls">TLS</t-option>
            <t-option value="kubernetes.io/dockerconfigjson"
              >Docker Config</t-option
            >
          </t-select>
        </t-form-item>

        <!-- 数据 -->
        <t-form-item label="数据" name="data">
          <t-textarea
            v-model="createFormData.data"
            placeholder="请输入键值对，例如：key1=value1\nkey2=value2"
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

    <!--展示 Secret 数据的对话框-->
    <t-dialog
      v-model:visible="viewDataDialogVisible"
      header="Secret 数据（未经 Base 64 解码）"
      @confirm="viewDataDialogVisible = false"
      :cancel-btn="null"
    >
      {{ viewData }}
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/client";

// Secret 数据
const secrets = ref([]);
const viewData = ref("");

// 创建 Secret 的对话框状态
const createDialogVisible = ref(false);

const viewDataDialogVisible = ref(false);

// 创建 Secret 的表单数据
const createFormData = ref({
  name: "",
  type: "Opaque",
  data: "",
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "Secret 名称不能为空" }],
  type: [{ required: true, message: "类型不能为空" }],
  data: [{ required: true, message: "数据不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "namespace", title: "命名空间", cell: "namespace" },
  { colKey: "type", title: "类型", cell: "type" },
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

// 获取 Secret 列表
const fetchSecrets = async () => {
  try {
    const response = await client.get(`${apiRoot.value}/secrets`);
    secrets.value = (await response.json()).items;
    MessagePlugin.success("Secret 列表刷新成功");
  } catch (error) {
    console.error("获取 Secret 失败:", error);
    MessagePlugin.error("获取 Secret 失败");
  }
};

// 删除 Secret
const handleDelete = async (secretName) => {
  try {
    await client.delete(`${apiRoot.value}/secrets/${secretName}`);
    MessagePlugin.success(`Secret ${secretName} 删除成功`);
    fetchSecrets(); // 刷新列表
  } catch (error) {
    console.error("删除 Secret 失败:", error);
    MessagePlugin.error("删除 Secret 失败");
  }
};

// 创建 Secret
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) return;
  try {
    const data = createFormData.value.data.split("\n").reduce((acc, line) => {
      const [key, value] = line.split("=");
      if (key && value) acc[key.trim()] = btoa(value.trim());
      return acc;
    }, {});

    const secretYAML = {
      apiVersion: "v1",
      kind: "Secret",
      metadata: {
        name: createFormData.value.name,
        namespace: "{!NAMESPACE}", // 自动使用 u-${username} 作为 namespace
      },
      type: createFormData.value.type,
      data,
    };

    const response = await client.post(`${apiRoot.value}/secrets`, secretYAML);
    if (response.status > 299) {
      MessagePlugin.error("创建 Secret 失败");
      MessagePlugin.error((await response.json()).message);
    } else {
      MessagePlugin.success("Secret 创建成功");
      closeCreateDialog(); // 关闭对话框
      fetchSecrets(); // 刷新列表
    }
  } catch (error) {
    console.error("创建 Secret 失败:", error);
    MessagePlugin.error("创建 Secret 失败");
  }
};

// 显示创建 Secret 的对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
};

// 关闭创建 Secret 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    type: "Opaque",
    data: "",
  };
};

const handleViewData = (secret) => {
  viewData.value = secret.data;
  viewDataDialogVisible.value = true;
};

// 组件挂载时获取用户信息和 Secret 列表
onMounted(async () => {
  await fetchUserInfo();
  fetchSecrets();
});
</script>
