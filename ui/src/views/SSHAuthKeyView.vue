<template>
  <div>
    <h1>SSHAuthorizedKeys</h1>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchSSHAuthorizedKeys">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 SSHAuthorizedKey
      </t-button>
    </t-space>

    <!-- SSHAuthorizedKey 表格 -->
    <t-table
      :data="sshAuthorizedKeys"
      :columns="columns"
      row-key="metadata.uid"
    >
      <!-- 名称列 -->
      <template #name="{ row }">
        <span>
          {{ row.metadata.name }}
        </span>
      </template>

      <template #publicKey="{ row }">
        <span>
          {{ row.spec.key }}
        </span>
      </template>

      <!-- 操作列 -->
      <template #operation="{ row }">
        <t-button variant="text" @click="handleDelete(row.metadata.name)">
          删除
        </t-button>
      </template>
    </t-table>

    <!-- 创建 SSHAuthorizedKey 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 SSHAuthorizedKey"
      :footer="false"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        ref="createFormRef"
        @submit="handleCreate"
      >
        <!-- SSHAuthorizedKey 名称 -->
        <t-form-item label="名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 SSHAuthorizedKey 名称"
          />
        </t-form-item>

        <!-- 公钥 -->
        <t-form-item label="公钥" name="publicKey">
          <t-textarea
            v-model="createFormData.publicKey"
            placeholder="ssh-ed25519 ..."
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

// SSHAuthorizedKey 数据
const sshAuthorizedKeys = ref([]);

// 创建 SSHAuthorizedKey 的对话框状态
const createDialogVisible = ref(false);

// 创建 SSHAuthorizedKey 的表单数据
const createFormData = ref({
  name: "",
  publicKey: "",
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "SSHAuthorizedKey 名称不能为空" }],
  publicKey: [{ required: true, message: "公钥不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "publicKey", title: "公钥", cell: "publicKey" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

// 用户信息和 API 根路径
let apiRoot = "/apis/ssh-operator.lcpu.dev/v1alpha1/namespaces/{!NAMESPACE}";

// 获取 SSHAuthorizedKey 列表
const fetchSSHAuthorizedKeys = async () => {
  try {
    const response = await client.get(`${apiRoot}/sshauthorizedkeys`);
    sshAuthorizedKeys.value = (await response.json()).items;
    MessagePlugin.success("SSHAuthorizedKey 列表刷新成功");
  } catch (error) {
    console.error("获取 SSHAuthorizedKey 失败:", error);
    MessagePlugin.error("获取 SSHAuthorizedKey 失败");
  }
};

// 删除 SSHAuthorizedKey
const handleDelete = async (sshAuthorizedKeyName) => {
  try {
    await client.delete(`${apiRoot}/sshauthorizedkeys/${sshAuthorizedKeyName}`);
    MessagePlugin.success(`SSHAuthorizedKey ${sshAuthorizedKeyName} 删除成功`);
    fetchSSHAuthorizedKeys(); // 刷新列表
  } catch (error) {
    console.error("删除 SSHAuthorizedKey 失败:", error);
    MessagePlugin.error("删除 SSHAuthorizedKey 失败");
  }
};

// 创建 SSHAuthorizedKey
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) return;
  try {
    const sshAuthorizedKeyYAML = {
      apiVersion: "ssh-operator.lcpu.dev/v1alpha1",
      kind: "SSHAuthorizedKey",
      metadata: {
        name: createFormData.value.name,
        namespace: `{!NAMESPACE}`, // 自动使用 u-${username} 作为 namespace
      },
      spec: {
        key: createFormData.value.publicKey,
      },
    };

    await client.post(`${apiRoot}/sshauthorizedkeys`, sshAuthorizedKeyYAML);
    MessagePlugin.success("SSHAuthorizedKey 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchSSHAuthorizedKeys(); // 刷新列表
  } catch (error) {
    console.error("创建 SSHAuthorizedKey 失败:", error);
    MessagePlugin.error("创建 SSHAuthorizedKey 失败");
  }
};

// 显示创建 SSHAuthorizedKey 的对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
};

// 关闭创建 SSHAuthorizedKey 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    publicKey: "",
  };
};

// 组件挂载时获取用户信息和 SSHAuthorizedKey 列表
onMounted(async () => {
  await client.ensureUsername();
  fetchSSHAuthorizedKeys();
});
</script>
