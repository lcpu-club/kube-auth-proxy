<template>
  <div>
    <h1>ConfigMap</h1>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchConfigMaps">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 ConfigMap
      </t-button>
    </t-space>

    <!-- ConfigMap 表格 -->
    <t-table :data="configMaps" :columns="columns" row-key="metadata.uid">
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

    <!-- 创建 ConfigMap 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 ConfigMap"
      :on-confirm="handleCreate"
      :on-close="closeCreateDialog"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        ref="createFormRef"
      >
        <!-- ConfigMap 名称 -->
        <t-form-item label="ConfigMap 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 ConfigMap 名称"
          />
        </t-form-item>

        <!-- 数据 -->
        <t-form-item label="数据" name="data">
          <t-textarea
            v-model="createFormData.data"
            placeholder="请输入键值对，例如：key1=value1\nkey2=value2"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </t-form-item>
      </t-form>
    </t-dialog>

    <!--显示 ConfigMap 数据的对话框-->
    <t-dialog
      v-model:visible="viewDataDialogVisible"
      header="ConfigMap 数据"
      :cancel-btn="null"
      @confirm="viewDataDialogVisible = false"
    >
      <pre style="text-wrap-mode: wrap;">{{ viewData }}</pre>
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/api";

// ConfigMap 数据
const configMaps = ref([]);

// 创建 ConfigMap 的对话框状态
const createDialogVisible = ref(false);

const viewDataDialogVisible = ref(false);

const viewData = ref("");

// 创建 ConfigMap 的表单数据
const createFormData = ref({
  name: "",
  data: "",
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "ConfigMap 名称不能为空" }],
  namespace: [{ required: true, message: "命名空间不能为空" }],
  data: [{ required: true, message: "数据不能为空" }],
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

// 获取 ConfigMap 列表
const fetchConfigMaps = async () => {
  try {
    const response = await client.get(`${apiRoot.value}/configmaps`);
    configMaps.value = (await response.json()).items;
    MessagePlugin.success("ConfigMap 列表刷新成功");
  } catch (error) {
    console.error("获取 ConfigMap 失败:", error);
    MessagePlugin.error("获取 ConfigMap 失败");
  }
};

// 删除 ConfigMap
const handleDelete = async (configMapName) => {
  try {
    await client.delete(`${apiRoot.value}/configmaps/${configMapName}`);
    MessagePlugin.success(`ConfigMap ${configMapName} 删除成功`);
    fetchConfigMaps(); // 刷新列表
  } catch (error) {
    console.error("删除 ConfigMap 失败:", error);
    MessagePlugin.error("删除 ConfigMap 失败");
  }
};

// 创建 ConfigMap
const handleCreate = async () => {
  try {
    const data = createFormData.value.data.split("\n").reduce((acc, line) => {
      const [key, value] = line.split("=");
      if (key && value) acc[key.trim()] = value.trim();
      return acc;
    }, {});

    const configMapYAML = {
      apiVersion: "v1",
      kind: "ConfigMap",
      metadata: {
        name: createFormData.value.name,
        namespace: `u-${username.value}`,
      },
      data,
    };

    await client.post(`${apiRoot.value}/configmaps`, configMapYAML);
    MessagePlugin.success("ConfigMap 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchConfigMaps(); // 刷新列表
  } catch (error) {
    console.error("创建 ConfigMap 失败:", error);
    MessagePlugin.error("创建 ConfigMap 失败");
  }
};

// 显示创建 ConfigMap 的对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
};

const handleViewData = (configMap) => {
  viewDataDialogVisible.value = true;
  viewData.value = JSON.stringify(configMap.data, null, 2);
};

// 关闭创建 ConfigMap 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    data: "",
  };
};

// 组件挂载时获取用户信息和 ConfigMap 列表
onMounted(async () => {
  await fetchUserInfo();
  fetchConfigMaps();
});
</script>
