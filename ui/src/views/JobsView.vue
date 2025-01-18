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
      placement="center"
      :footer="false"
    >
      <t-form
        labelWidth="120px"
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
          <t-select
            v-model="createFormData.image"
            placeholder="请输入容器镜像"
            creatable
            filterable
          >
            <t-option
              v-for="image in images"
              :key="image"
              :value="image"
              :label="image"
            />
          </t-select>
        </t-form-item>

        <t-form-item label="架构" name="architecture">
          <t-select
            v-model="createFormData.architecture"
            placeholder="请选择架构"
          >
            <t-option
              v-for="architecture in availableArchitectures"
              :key="architecture"
              :value="architecture"
              :label="architecture"
            />
          </t-select>
        </t-form-item>

        <!-- 命令 -->
        <t-form-item label="命令" name="command">
          <t-input v-model="createFormData.command" placeholder="请输入命令" />
        </t-form-item>

        <!-- 参数 -->
        <t-form-item label="参数" name="args">
          <t-textarea
            v-model="createFormData.args"
            placeholder="请输入参数（换行分割）"
          />
        </t-form-item>

        <!-- 并行数 -->
        <t-form-item label="并行数" name="parallelism">
          <t-input-number v-model="createFormData.parallelism" :min="1" />
        </t-form-item>

        <!-- 重试次数 -->
        <t-form-item label="重试次数" name="backoffLimit">
          <t-input-number v-model="createFormData.backoffLimit" :min="0" />
        </t-form-item>

        <!-- PVC 选择 -->
        <t-form-item label="PVC" name="pvc">
          <t-select
            v-model="createFormData.pvc"
            placeholder="请选择 PVC"
            clearable
            filterable
          >
            <t-option
              v-for="pvc in pvcs"
              :key="pvc.metadata.name"
              :value="pvc.metadata.name"
              :label="pvc.metadata.name"
            />
          </t-select>
        </t-form-item>

        <!-- 标签 -->
        <t-form-item label="标签" name="labels">
          <t-input v-model="createFormData.labels" placeholder="a=b, c=d" />
        </t-form-item>

        <!-- 挂载路径 -->
        <t-form-item label="挂载路径" name="mountPath">
          <t-select
            v-model="createFormData.mountPath"
            placeholder="请选择挂载路径"
            creatable
            filterable
          >
            <t-option value="/root" label="/root" />
            <t-option value="/data" label="/data" />
            <t-option value="/shared" label="/shared" />
          </t-select>
        </t-form-item>

        <!-- lxcfs.lcpu.dev/inject 开关 -->
        <t-form-item label="启用 lxcfs" name="lxcfsEnabled">
          <t-switch v-model="createFormData.lxcfsEnabled" />
        </t-form-item>

        <!-- ssh-operator.lcpu.dev/inject 开关 -->
        <t-form-item label="启用 SSH 服务器" name="sshEnabled">
          <t-switch v-model="createFormData.sshEnabled" />
        </t-form-item>

        <t-form-item label="RDMA" name="rdma">
          <t-switch v-model="createFormData.rdma" />
        </t-form-item>

        <!-- 资源请求和限制 -->
        <t-form-item label="CPU" name="cpuRequest">
          <t-input
            v-model="createFormData.cpuRequest"
            placeholder="例如：1000m"
          />
        </t-form-item>

        <t-form-item label="内存" name="memoryRequest">
          <t-input
            v-model="createFormData.memoryRequest"
            placeholder="例如：1Gi"
          />
        </t-form-item>

        <t-form-item label="Nvidia GPU" name="gpuRequest">
          <t-input-number v-model="createFormData.gpuRequest" :min="0" />
        </t-form-item>

        <t-form-item label="Ascend910" name="ascend910Request">
          <t-input-number v-model="createFormData.ascend910Request" :min="0" />
        </t-form-item>

        <t-form-item label="Ascend310P" name="ascend310PRequest">
          <t-input-number v-model="createFormData.ascend310PRequest" :min="0" />
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
import { ref, onMounted, reactive, computed, watch } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/api";
import { useRouter, useRoute } from "vue-router";

let apiRoot = "/apis/batch/v1/namespaces/{!NAMESPACE}";

const router = useRouter();

// Jobs 数据
const jobs = ref([]);
const localQueues = ref([]);
const pvcs = ref([]);

const x86OnlyImages = ref(["intel", "cuda", "aocc"]);
const armOnlyImages = ref(["hpckit"]);
const commonImages = ref(["full", "llvm", "gcc", "nvhpc", "julia", "base"]);

const images = computed(() => [
  ...commonImages.value,
  ...x86OnlyImages.value,
  ...armOnlyImages.value,
]);

const route = useRoute();
const queryString = (key, init = "") =>
  typeof route.query[key] === "string" ? route.query[key] : init;
const queryBoolean = (key, init = false) =>
  typeof route.query[key] === "string" ? route.query[key] === "true" : init;
const queryNumber = (key, init = 0) =>
  typeof route.query[key] === "string" && isFinite(route.query[key])
    ? +route.query[key]
    : init;

// 创建 Job 的对话框状态
const createDialogVisible = ref("name" in route.query);

// 创建 Job 的表单数据
const createFormDataFactory = () => ({
  name: queryString("name"),
  localQueue: queryString("localQueue"),
  image: queryString("image"),
  architecture: queryString("architecture", "x86_amd"),
  command: queryString("command", "sleep"),
  args: queryString("args", "inf"),
  parallelism: queryNumber("parallelism", 1),
  backoffLimit: queryNumber("backoffLimit", 6),
  pvc: queryString("pvc"),
  mountPath: queryString("mountPath", "/root"),
  labels: queryString("labels"),
  lxcfsEnabled: queryBoolean("lxcfsEnabled", false),
  sshEnabled: queryBoolean("sshEnabled", true),
  rdma: queryBoolean("rdma", true),
  cpuRequest: queryString("cpuRequest"),
  memoryRequest: queryString("memoryRequest"),
  gpuRequest: queryNumber("gpuRequest", 0),
  ascend910Request: queryNumber("ascend910Request", 0),
  ascend310PRequest: queryNumber("ascend310PRequest", 0),
});
const createFormData = ref(createFormDataFactory());

const availableArchitectures = computed(() => {
  if (commonImages.value.includes(createFormData.value.image))
    return ["x86_amd", "arm"];
  if (x86OnlyImages.value.includes(createFormData.value.image))
    return ["x86_amd"];
  if (armOnlyImages.value.includes(createFormData.value.image)) return ["arm"];
  return ["x86_amd", "arm"];
});

watch(availableArchitectures, (newVal) => {
  if (newVal.includes(createFormData.value.architecture)) return;
  createFormData.value.architecture = newVal[0];
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "Job 名称不能为空" }],
  image: [{ required: true, message: "容器镜像不能为空" }],
  architecture: [{ required: true, message: "架构不能为空" }],
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
          "k8s.v1.cni.cncf.io/networks": "wm2-roce",
        },
        labels: createFormData.value.localQueue
          ? {
              "kueue.x-k8s.io/queue-name": createFormData.value.localQueue,
            }
          : {},
      },
      spec: {
        template: {
          metadata: {
            labels: createFormData.value.labels.split(",").reduce((acc, x) => {
              const [k, v] = x.split("=");
              if (k && v) acc[k.trim()] = v.trim();
              return acc;
            }, {}),
          },
          spec: {
            nodeSelector: {
              "hpc.lcpu.dev/partition": createFormData.value.architecture,
            },
            containers: [
              {
                name: createFormData.value.name,
                securityContext: {
                  capabilities: {
                    add: ["IPC_LOCK"],
                  },
                },
                image: images.value.includes(createFormData.value.image)
                  ? `crmirror.lcpu.dev/hpcgame/${createFormData.value.image}:latest`
                  : createFormData.value.image,
                command: createFormData.value.command.split(" "),
                args: createFormData.value.args
                  .split("\n")
                  .map((x) => x.trim())
                  .filter((x) => x),
                resources: {
                  requests: {
                    "nvidia.com/gpu": createFormData.value.gpuRequest,
                    "huawei.com/Ascend910":
                      createFormData.value.ascend910Request,
                    "huawei.com/Ascend310P":
                      createFormData.value.ascend310PRequest,
                  },
                  limits: {
                    "nvidia.com/gpu": createFormData.value.gpuRequest,
                    "huawei.com/Ascend910":
                      createFormData.value.ascend910Request,
                    "huawei.com/Ascend310P":
                      createFormData.value.ascend310PRequest,
                    "rdma.hpc.lcpu.dev/hca_cx5": createFormData.value.rdma
                      ? 1
                      : 0,
                  },
                },
                volumeMounts: [],
              },
            ],
            volumes: [],
            restartPolicy: "Never",
          },
        },
        parallelism: createFormData.value.parallelism,
        backoffLimit: createFormData.value.backoffLimit,
      },
    };

    if (createFormData.value.cpuRequest) {
      jobYAML.spec.template.spec.containers[0].resources.requests.cpu =
        createFormData.value.cpuRequest;
      jobYAML.spec.template.spec.containers[0].resources.limits.cpu =
        createFormData.value.cpuLimit;
    }

    if (createFormData.value.memoryRequest) {
      jobYAML.spec.template.spec.containers[0].resources.requests.memory =
        createFormData.value.memoryRequest;
      jobYAML.spec.template.spec.containers[0].resources.limits.memory =
        createFormData.value.memoryLimit;
    }

    // 如果用户选择了 PVC 和挂载路径，则添加到 YAML 中
    if (createFormData.value.pvc && createFormData.value.mountPath) {
      jobYAML.spec.template.spec.volumes.push({
        name: "pvc-volume",
        persistentVolumeClaim: {
          claimName: createFormData.value.pvc,
        },
      });

      jobYAML.spec.template.spec.containers[0].volumeMounts.push({
        name: "pvc-volume",
        mountPath: createFormData.value.mountPath,
      });
    }

    if (createFormData.value.localQueue) {
      jobYAML.spec.suspend = true;
    }

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
  }
};

// 获取 PVC 数据
const fetchPVCs = async () => {
  try {
    const response = await client.get(
      `/api/v1/namespaces/{!NAMESPACE}/persistentvolumeclaims`
    );
    const data = await response.json();
    pvcs.value = data.items;
  } catch (error) {
    console.error("获取 PVC 失败:", error);
    MessagePlugin.error("获取 PVC 失败");
  }
};

// 显示创建 Job 的对话框
const showCreateDialog = async () => {
  createDialogVisible.value = true;
  await fetchLocalQueues(); // 获取 LocalQueues 列表
  await fetchPVCs();
};

// 关闭创建 Job 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = createFormDataFactory();
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
  router.replace({ ...route, query: {} });
  fetchJobs();
});
</script>
