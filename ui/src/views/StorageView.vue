<template>
  <h1>储存管理</h1>
  <div>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchPVCs">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 PVC
      </t-button>
    </t-space>

    <!-- PVC 表格 -->
    <t-table :data="pvcs" :columns="columns" row-key="metadata.uid">
      <!-- 名称列 -->
      <template #name="{ row }">
        <span>
          {{ row.metadata.name }}
        </span>
      </template>

      <template #storageClass="{ row }">
        <span>
          {{ row.spec.storageClassName }}
        </span>
      </template>

      <template #accessMode="{ row }">
        <span>
          {{ row.spec.accessModes.join(", ") }}
        </span>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <t-tag :theme="getStatusTheme(row.status.phase)">
          {{ row.status.phase }}
        </t-tag>
      </template>

      <!-- 容量列 -->
      <template #capacity="{ row }">
        {{ row.spec.resources.requests.storage }}
      </template>

      <!-- 操作列 -->
      <template #operation="{ row }">
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

    <!-- 创建 PVC 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 PVC"
      :footer="false"
      :on-close="closeCreateDialog"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        @submit="handleCreate"
        ref="createFormRef"
      >
        <!-- PVC 名称 -->
        <t-form-item label="PVC 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 PVC 名称"
          />
        </t-form-item>

        <!-- 存储大小 -->
        <t-form-item label="存储大小" name="storage">
          <t-input v-model="createFormData.storage" placeholder="例如 10Gi" />
        </t-form-item>

        <!-- 访问模式 -->
        <t-form-item label="访问模式" name="accessMode">
          <t-select v-model="createFormData.accessMode">
            <t-option value="ReadWriteOnce">ReadWriteOnce</t-option>
            <t-option value="ReadOnlyMany">ReadOnlyMany</t-option>
            <t-option value="ReadWriteMany">ReadWriteMany</t-option>
          </t-select>
        </t-form-item>

        <!-- StorageClass -->
        <t-form-item label="StorageClass" name="storageClass">
          <t-select
            v-model="createFormData.storageClass"
            placeholder="请输入 StorageClass"
            filterable
          >
            <t-option
              v-for="item in availableStorageClasses"
              :key="item"
              :value="item"
            >
              {{ item }}
            </t-option>
          </t-select>
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit">创建</t-button>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/client";

let apiRoot = "/api/v1/namespaces/{!NAMESPACE}";

// PVC 数据
const pvcs = ref([]);

// 创建 PVC 的对话框状态
const createDialogVisible = ref(false);

// 创建 PVC 的表单数据
const createFormData = ref({
  name: "",
  storage: "10Gi",
  accessMode: "ReadWriteOnce",
  storageClass: "",
});

const rwOnceStorageClasses = ref([
  "npu-local-data",
  "openebs-hostpath",
  "x86-amd-local-hostpath",
]);

const commmonStorageClasses = ref(["juicefs", "wm2-nfs", "yanyuan-nfs"]);

const availableStorageClasses = computed(() => {
  if (createFormData.value.accessMode === "ReadWriteOnce")
    return [...rwOnceStorageClasses.value, ...commmonStorageClasses.value];
  else return commmonStorageClasses.value;
});

watch(availableStorageClasses, (newVal) => {
  if (
    !availableStorageClasses.value.includes(createFormData.value.storageClass)
  ) {
    createFormData.value.storageClass = "";
  }
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "PVC 名称不能为空" }],
  namespace: [{ required: true, message: "命名空间不能为空" }],
  storage: [{ required: true, message: "存储大小不能为空" }],
  accessMode: [{ required: true, message: "访问模式不能为空" }],
  storageClass: [{ required: true, message: "StorageClass 不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "storageClass", title: "StorageClass", cell: "storageClass" },
  { colKey: "accessMode", title: "访问模式", cell: "accessMode" },
  { colKey: "status", title: "状态", cell: "status" },
  { colKey: "capacity", title: "容量", cell: "capacity" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

// 获取 PVC 列表
const fetchPVCs = async () => {
  try {
    const response = await client.get(`${apiRoot}/persistentvolumeclaims`);
    pvcs.value = (await response.json()).items;
    MessagePlugin.success("PVC 列表刷新成功");
  } catch (error) {
    console.error("获取 PVC 失败:", error);
    MessagePlugin.error("获取 PVC 失败");
  }
};

// 删除 PVC
const handleDelete = async (pvcName) => {
  try {
    await client.delete(`${apiRoot}/persistentvolumeclaims/${pvcName}`);
    MessagePlugin.success(`PVC ${pvcName} 删除成功`);
    fetchPVCs(); // 刷新列表
  } catch (error) {
    console.error("删除 PVC 失败:", error);
    MessagePlugin.error("删除 PVC 失败");
  }
};

// 创建 PVC
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) {
    return;
  }
  try {
    const pvcYAML = {
      apiVersion: "v1",
      kind: "PersistentVolumeClaim",
      metadata: {
        name: createFormData.value.name,
        namespace: "{!NAMESPACE}",
      },
      spec: {
        accessModes: [createFormData.value.accessMode],
        resources: {
          requests: {
            storage: createFormData.value.storage,
          },
        },
        storageClassName: createFormData.value.storageClass,
      },
    };

    await client.post(`${apiRoot}/persistentvolumeclaims`, pvcYAML);
    MessagePlugin.success("PVC 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchPVCs(); // 刷新列表
  } catch (error) {
    console.error("创建 PVC 失败:", error);
    MessagePlugin.error("创建 PVC 失败");
  }
};

// 显示创建 PVC 的对话框
const showCreateDialog = async () => {
  createDialogVisible.value = true;
};

// 关闭创建 PVC 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    storage: "10Gi",
    accessMode: "ReadWriteOnce",
    storageClass: "",
  };
};

// 根据状态获取 Tag 主题
const getStatusTheme = (phase) => {
  switch (phase) {
    case "Bound":
      return "success";
    case "Pending":
      return "warning";
    case "Lost":
      return "danger";
    default:
      return "default";
  }
};

// 组件挂载时获取数据
onMounted(async () => {
  if (!(await client.ensureUsername())) return;
  fetchPVCs();
});
</script>
