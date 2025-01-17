<template>
  <h1>任务</h1>
  <div>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchJobs">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 Job
      </t-button>
    </t-space>

    <!-- Jobs 表格 -->
    <t-table :data="jobs" :columns="columns" row-key="metadata.uid">
      <!-- 名称列 -->
      <template #name="{ row }">
        <span>
          {{ row.metadata.name }}
        </span>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <t-tag :theme="getStatusTheme(row.status)">
          {{
            row.status.conditions ? row.status.conditions[0].type : "Unknown"
          }}
        </t-tag>
      </template>

      <!-- 创建时间列 -->
      <template #creationTimestamp="{ row }">
        {{ formatDate(row.metadata.creationTimestamp) }}
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

    <!-- 创建 Job 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 Job"
      :footer="false"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        @submit="handleCreate"
        ref="createFormRef"
      >
        <!-- Job 名称 -->
        <t-form-item label="Job 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 Job 名称"
          />
        </t-form-item>

        <!-- LocalQueue 选择 -->
        <t-form-item label="LocalQueue" name="localQueue">
          <t-select
            v-model="createFormData.localQueue"
            placeholder="请选择 LocalQueue"
            clearable
          >
            <t-option
              v-for="queue in localQueues"
              :key="queue.metadata.name"
              :value="queue.metadata.name"
              :label="queue.metadata.name"
            />
          </t-select>
        </t-form-item>

        <!-- 容器镜像 -->
        <t-form-item label="容器镜像" name="image">
          <t-input
            v-model="createFormData.image"
            placeholder="请输入容器镜像"
          />
        </t-form-item>

        <!-- 命令 -->
        <t-form-item label="命令" name="command">
          <t-input v-model="createFormData.command" placeholder="请输入命令" />
        </t-form-item>

        <!-- 参数 -->
        <t-form-item label="参数" name="args">
          <t-input v-model="createFormData.args" placeholder="请输入参数" />
        </t-form-item>

        <!-- 并行数 -->
        <t-form-item label="并行数" name="parallelism">
          <t-input-number v-model="createFormData.parallelism" :min="1" />
        </t-form-item>

        <!-- 重试次数 -->
        <t-form-item label="重试次数" name="backoffLimit">
          <t-input-number v-model="createFormData.backoffLimit" :min="0" />
        </t-form-item>

        <!-- lxcfs.lcpu.dev/inject 开关 -->
        <t-form-item label="启用 lxcfs" name="lxcfsEnabled">
          <t-switch v-model="createFormData.lxcfsEnabled" />
        </t-form-item>

        <!-- ssh-operator.lcpu.dev/inject 开关 -->
        <t-form-item label="启用 SSH 服务器" name="sshEnabled">
          <t-switch v-model="createFormData.sshEnabled" />
        </t-form-item>

        <t-form-item>
          <t-space size="small">
            <t-button theme="default" @click="closeCreateDialog">取消</t-button>
            <t-button theme="primary" type="submit">创建</t-button>
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

let apiRoot = "/apis/batch/v1/namespaces/{!NAMESPACE}";

// Jobs 数据
const jobs = ref([]);
const localQueues = ref([]);

// 创建 Job 的对话框状态
const createDialogVisible = ref(false);

// 创建 Job 的表单数据
const createFormData = ref({
  name: "",
  localQueue: "",
  image: "",
  command: "",
  args: "",
  parallelism: 1,
  backoffLimit: 6,
  lxcfsEnabled: false,
  sshEnabled: false,
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "Job 名称不能为空" }],
  image: [{ required: true, message: "容器镜像不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "status", title: "状态", cell: "status" },
  { colKey: "creationTimestamp", title: "创建时间", cell: "creationTimestamp" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

const fetchLocalQueues = async () => {
  try {
    const response = await client.get(
      "/apis/kueue.x-k8s.io/v1beta1/namespaces/{!NAMESPACE}/localqueues"
    );
    const data = await response.json();
    localQueues.value = data.items;
    MessagePlugin.success("Local Queues 列表刷新成功");
  } catch (error) {
    console.error("获取 Local Queues 失败:", error);
    MessagePlugin.error("获取 Local Queues 失败");
  }
};

// 获取 Jobs 数据
const fetchJobs = async () => {
  try {
    // 模拟 API 调用
    const response = await client.get(`${apiRoot}/jobs`);
    const data = await response.json();
    jobs.value = data.items;
    MessagePlugin.success("Jobs 列表刷新成功");
  } catch (error) {
    console.error("获取 Jobs 失败:", error);
    MessagePlugin.error("获取 Jobs 失败");
  }
};

// 删除 Job
const handleDelete = async (jobName) => {
  try {
    // 模拟 API 调用
    await client.delete(`${apiRoot}/jobs/${jobName}`, {
      method: "DELETE",
    });
    MessagePlugin.success(`Job ${jobName} 删除成功`);
    fetchJobs(); // 刷新列表
  } catch (error) {
    console.error("删除 Job 失败:", error);
    MessagePlugin.error("删除 Job 失败");
  }
};

// 创建 Job
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) {
    return;
  }
  try {
    // 构造 Job 的 YAML
    const jobYAML = {
      apiVersion: "batch/v1",
      kind: "Job",
      metadata: {
        name: createFormData.value.name,
        namespace: "{!NAMESPACE}",
        annotations: {
          "lxcfs.lcpu.dev/inject": createFormData.value.lxcfsEnabled
            ? "enabled"
            : "disabled",
          "ssh-operator.lcpu.dev/inject": createFormData.value.sshEnabled
            ? "enabled"
            : "disabled",
        },
        labels: createFormData.value.localQueue
          ? {
              "kueue.x-k8s.io/queue-name": createFormData.value.localQueue,
            }
          : {},
      },
      spec: {
        template: {
          spec: {
            containers: [
              {
                name: createFormData.value.name,
                image: createFormData.value.image,
                command: createFormData.value.command.split(" "),
                args: createFormData.value.args.split(" "),
              },
            ],
            restartPolicy: "Never",
          },
        },
        parallelism: createFormData.value.parallelism,
        backoffLimit: createFormData.value.backoffLimit,
      },
    };

    const response = await client.post(`${apiRoot}/jobs`, jobYAML);
    if (response.status > 299) {
      throw await response.json();
    }

    MessagePlugin.success("Job 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchJobs(); // 刷新列表
  } catch (error) {
    console.error("创建 Job 失败:", error);
    MessagePlugin.error("创建 Job 失败");
    MessagePlugin.error(error.message);
  }
};

// 显示创建 Job 的对话框
const showCreateDialog = () => {
  createDialogVisible.value = true;
  fetchLocalQueues(); // 获取 LocalQueues 列表
};

// 关闭创建 Job 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    localQueue: "",
    image: "",
    command: "",
    args: "",
    parallelism: 1,
    backoffLimit: 6,
    lxcfsEnabled: false,
    sshEnabled: false,
  };
};

// 格式化日期
const formatDate = (timestamp) => {
  return new Date(timestamp).toLocaleString();
};

// 根据状态获取 Tag 主题
const getStatusTheme = (status) => {
  if (status.succeeded) return "success";
  if (status.failed) return "danger";
  if (status.active) return "warning";
  return "default";
};

// 组件挂载时获取数据
onMounted(async () => {
  await client.ensureUsername();
  fetchJobs();
});
</script>
