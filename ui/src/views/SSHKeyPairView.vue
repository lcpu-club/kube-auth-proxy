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

      <template #privateKey="{ row }">
        <t-button
          variant="text"
          @click="displayKey(row.spec.privateKey, 'private')"
        >
          查看
        </t-button>
      </template>

      <template #publicKey="{ row }">
        <t-button
          variant="text"
          @click="displayKey(row.spec.publicKey, 'public')"
        >
          查看
        </t-button>
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
        <t-form-item label="名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 SSHKeyPair 名称"
          />
        </t-form-item>

        <!-- 类型 -->
        <t-form-item label="类型" name="type">
          <t-select
            v-model="createFormData.type"
            placeholder="请选择类型"
            filterable
            creatable
          >
            <t-option
              v-for="type in predefinedKeyTypes"
              :key="type"
              :value="type"
            >
              {{ type }}
            </t-option>
          </t-select>
        </t-form-item>

        <!-- 提交按钮 -->
        <t-form-item>
          <t-space>
            <t-button theme="primary" type="submit" :loading="creatingKeyPair"
              >提交</t-button
            >
            <t-button theme="default" @click="closeCreateDialog">取消</t-button>
          </t-space>
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog
      v-model:visible="viewKeyDialogVisible"
      :header="'查看' + (displayingKeyType == 'private' ? '私钥' : '公钥')"
      width="50em"
      @confirm="viewKeyDialogVisible = false"
      :cancel-btn="null"
    >
      <t-textarea
        readonly
        autosize
        v-model="displayingKey"
        style="font-family: monospace"
      />
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
const viewKeyDialogVisible = ref(false);
const displayingKey = ref("");
const displayingKeyType = ref("");
const creatingKeyPair = ref(false);

// 创建 SSHKeyPair 的表单数据
const createFormData = ref({
  name: "",
  type: "",
  privateKey: "",
  publicKey: "",
});

const predefinedKeyTypes = ["ssh-ed25519", "ssh-rsa"];

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "SSHKeyPair 名称不能为空" }],
  type: [{ required: true, message: "类型不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "privateKey", title: "私钥", cell: "privateKey" },
  { colKey: "publicKey", title: "公钥", cell: "publicKey" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

// 用户信息和 API 根路径
let apiRoot = "/apis/ssh-operator.lcpu.dev/v1alpha1/namespaces/{!NAMESPACE}";

// 获取 SSHKeyPair 列表
const fetchSSHKeyPairs = async () => {
  try {
    const response = await client.get(`${apiRoot}/sshkeypairs`);
    sshKeyPairs.value = (await response.json()).items;
    MessagePlugin.success("SSHKeyPair 列表刷新成功");
  } catch (error) {
    console.error("获取 SSHKeyPair 失败:", error);
    MessagePlugin.error("获取 SSHKeyPair 失败");
  }
};

// 删除 SSHKeyPair
const handleDelete = async (sshKeyPairName) => {
  try {
    await client.delete(`${apiRoot}/sshkeypairs/${sshKeyPairName}`);
    MessagePlugin.success(`SSHKeyPair ${sshKeyPairName} 删除成功`);
    fetchSSHKeyPairs(); // 刷新列表
  } catch (error) {
    console.error("删除 SSHKeyPair 失败:", error);
    MessagePlugin.error("删除 SSHKeyPair 失败");
  }
};

// 创建 SSHKeyPair
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) return;
  creatingKeyPair.value = true;
  try {
    const sshKeyPairYAML = {
      apiVersion: "ssh-operator.lcpu.dev/v1alpha1",
      kind: "SSHKeyPair",
      metadata: {
        name: createFormData.value.name,
        namespace: `{!NAMESPACE}`, // 自动使用 u-${username} 作为 namespace
      },
      spec: {
        type: createFormData.value.type,
      },
    };

    await client.post(`${apiRoot}/sshkeypairs`, sshKeyPairYAML);
    MessagePlugin.success("SSHKeyPair 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchSSHKeyPairs(); // 刷新列表
  } catch (error) {
    console.error("创建 SSHKeyPair 失败:", error);
    MessagePlugin.error("创建 SSHKeyPair 失败");
  } finally {
    creatingKeyPair.value = false;
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
    type: "",
  };
};

const displayKey = (key, type) => {
  displayingKey.value = key;
  displayingKeyType.value = type;
  viewKeyDialogVisible.value = true;
};

// 组件挂载时获取用户信息和 SSHKeyPair 列表
onMounted(async () => {
  await client.ensureUsername();
  fetchSSHKeyPairs();
});
</script>
