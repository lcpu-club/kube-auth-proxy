<template>
  <h1>容器</h1>
  <div>
    <!-- 操作栏 -->
    <t-space style="margin-bottom: 16px">
      <t-button theme="primary" @click="fetchPods">
        <template #icon><t-icon name="refresh" /></template>
        刷新
      </t-button>
      <t-button theme="success" @click="showCreateDialog">
        <template #icon><t-icon name="add" /></template>
        创建 Pod
      </t-button>
    </t-space>

    <!-- Pods 表格 -->
    <t-table :data="pods" :columns="columns" row-key="metadata.uid">
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

      <!-- 状态列 -->
      <template #status="{ row }">
        <t-tag :theme="getStatusTheme(row.status.phase)">
          {{ row.status.phase }}
        </t-tag>
      </template>

      <!-- 节点列 -->
      <template #node="{ row }">
        {{ row.spec.nodeName }}
      </template>

      <!-- 创建时间列 -->
      <template #creationTimestamp="{ row }">
        {{ formatDate(row.metadata.creationTimestamp) }}
      </template>

      <!-- 操作列 -->
      <template #operation="{ row }">
        <t-button
          variant="text"
          shape="square"
          @click="showLog(row.metadata.name)"
          title="查看日志"
        >
          <t-icon name="file-1" />
        </t-button>
        <t-button
          variant="text"
          shape="square"
          @click="handleMoreTools(row)"
          title="更多工具"
        >
          <t-icon name="more" />
        </t-button>
        <t-button
          variant="text"
          shape="square"
          @click="handleDelete(row.metadata.name)"
          title="删除"
        >
          <t-icon name="delete-1" />
        </t-button>
      </template>
    </t-table>

    <!-- 创建 Pod 的对话框 -->
    <t-dialog
      v-model:visible="createDialogVisible"
      header="创建 Pod"
      :footer="false"
      :on-close="closeCreateDialog"
    >
      <t-form
        :data="createFormData"
        :rules="createFormRules"
        @submit="handleCreate"
        ref="createFormRef"
      >
        <!-- Pod 名称 -->
        <t-form-item label="Pod 名称" name="name">
          <t-input
            v-model="createFormData.name"
            placeholder="请输入 Pod 名称"
          />
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
          <t-textarea
            v-model="createFormData.args"
            placeholder="请输入参数（以换行分割）"
          />
        </t-form-item>

        <!-- PVC 选择 -->
        <t-form-item label="挂载 PVC" name="pvc">
          <t-select
            v-model="createFormData.pvc"
            placeholder="请选择 PVC"
            clearable
          >
            <t-option
              v-for="pvc in pvcs"
              :key="pvc.metadata.name"
              :value="pvc.metadata.name"
              :label="pvc.metadata.name"
            />
          </t-select>
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

        <!-- 注解：LXCFS -->
        <t-form-item label="启用 LXCFS" name="lxcfsEnabled">
          <t-switch v-model="createFormData.lxcfsEnabled" />
        </t-form-item>

        <!-- 注解：SSH Operator -->
        <t-form-item label="启用 SSH 服务器" name="sshEnabled">
          <t-switch v-model="createFormData.sshEnabled" />
        </t-form-item>

        <t-form-item>
          <t-button theme="primary" type="submit">创建</t-button>
        </t-form-item>
      </t-form>
    </t-dialog>

    <!-- 查看日志的对话框 -->
    <t-dialog
      v-model:visible="logDialogVisible"
      header="查看日志"
      :cancel-btn="null"
      @confirm="logDialogVisible = false"
    >
      <pre>{{ logContent }}</pre>
    </t-dialog>

    <!-- 更多操作对话框 -->
    <t-dialog
      v-model:visible="selectContainerDialogVisible"
      :header="`更多操作 - ${selectedPod}`"
    >
      <!-- <t-space style="margin-bottom: 1em" size="small">
        <p>VSCode Server 状态：{{ selectedPodCodeStatus }}</p>
        <t-button
          theme="primary"
          :disabled="attachCodeDisabled"
          @click="attachCode(selectedPod)"
          :loading="selectedPodCodeAttaching"
        >
          启动
        </t-button>
        <t-button
          :disabled="detachCodeDisabled"
          @click="openCode(selectedPod)"
          :loading="selectedPodCodeDetaching"
        >
          打开
        </t-button>
        <t-button
          theme="danger"
          :disabled="detachCodeDisabled"
          :loading="selectedPodCodeDetaching"
          @click="detachCode(selectedPod)"
        >
          关闭
        </t-button>
      </t-space> -->
      <t-select v-model:value="selectedPodContianer" placeholder="请选择容器">
        <t-option
          v-for="container in currentPodContainers"
          :key="container"
          :value="container"
        >
          {{ container }}
        </t-option>
      </t-select>
      <template #confirmBtn>
        <t-button
          theme="primary"
          :disabled="!selectedPodContianer"
          @click="
            selectContainerDialogVisible = false;
            execWebTerm(selectedPod, selectedPodContianer);
          "
          >启动网页终端</t-button
        >
      </template>
    </t-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { MessagePlugin } from "tdesign-vue-next";
import { client } from "@/api/api";
import { useRouter } from "vue-router";

let apiRoot = `/api/v1/namespaces/{!NAMESPACE}`;
let username = "";

const router = useRouter();

// Pods 数据
const pods = ref([]);

// PVC 数据
const pvcs = ref([]);

const currentPodContainers = ref([]);
const selectedPod = ref("");
const selectedPodContianer = ref("");
const selectedPodCodeRunning = ref(false);
const selectedPodCodeStatusLoading = ref(false);
const selectedPodCodeStatus = computed(() => {
  if (selectedPodCodeStatusLoading.value) {
    return "loading...";
  }
  return selectedPodCodeRunning.value ? "running" : "stopped";
});
const selectedPodCodeAttaching = ref(false);
const selectedPodCodeDetaching = ref(false);

const attachCodeDisabled = computed(() => {
  return selectedPodCodeStatusLoading.value || selectedPodCodeRunning.value;
});

const detachCodeDisabled = computed(() => {
  return selectedPodCodeStatusLoading.value || !selectedPodCodeRunning.value;
});

// 查看日志的内容
const logContent = ref("");

// 创建 Pod 的对话框状态
const createDialogVisible = ref(false);

// 查看日志的对话框状态
const logDialogVisible = ref(false);

// 选择容器的对话框状态
const selectContainerDialogVisible = ref(false);

// 创建 Pod 的表单数据
const createFormData = ref({
  name: "",
  image: "",
  command: "",
  args: "",
  pvc: "",
  mountPath: "/root",
  lxcfsEnabled: false,
  sshEnabled: false,
});

// 表单验证规则
const createFormRules = {
  name: [{ required: true, message: "Pod 名称不能为空" }],
  image: [{ required: true, message: "容器镜像不能为空" }],
};

// 表格列配置
const columns = [
  { colKey: "name", title: "名称", cell: "name" },
  { colKey: "namespace", title: "命名空间", cell: "namespace" },
  { colKey: "status", title: "状态", cell: "status" },
  { colKey: "node", title: "节点", cell: "node" },
  { colKey: "creationTimestamp", title: "创建时间", cell: "creationTimestamp" },
  { colKey: "operation", title: "操作", cell: "operation" },
];

// 获取 Pods 数据
const fetchPods = async () => {
  try {
    const response = await client.get(`${apiRoot}/pods`);
    const data = await response.json();
    pods.value = data.items;
    MessagePlugin.success("Pods 列表刷新成功");
  } catch (error) {
    console.error("获取 Pods 失败:", error);
    MessagePlugin.error("获取 Pods 失败");
  }
};

// 获取 PVC 数据
const fetchPVCs = async () => {
  try {
    const response = await client.get(`${apiRoot}/persistentvolumeclaims`);
    const data = await response.json();
    pvcs.value = data.items;
  } catch (error) {
    console.error("获取 PVC 失败:", error);
    MessagePlugin.error("获取 PVC 失败");
  }
};

// 删除 Pod
const handleDelete = async (podName) => {
  try {
    await client.delete(`${apiRoot}/pods/${podName}`);
    MessagePlugin.success(`Pod ${podName} 删除成功`);
    fetchPods(); // 刷新列表
  } catch (error) {
    console.error("删除 Pod 失败:", error);
    MessagePlugin.error("删除 Pod 失败");
  }
};

// 创建 Pod
const handleCreate = async ({ validateResult }) => {
  if (validateResult !== true) {
    return;
  }
  try {
    // 构造 Pod 的 YAML
    const podYAML = {
      apiVersion: "v1",
      kind: "Pod",
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
      },
      spec: {
        nodeSelector: {
          "hpc.lcpu.dev/partition": "x86_amd"
        },
        containers: [
          {
            name: createFormData.value.name,
            image: createFormData.value.image,
            command: createFormData.value.command.split(" "),
            args: createFormData.value.args
              .split("\n")
              .map((x) => x.trim())
              .filter((x) => x),
            volumeMounts: [],
          },
        ],
        volumes: [],
      },
    };

    // 添加 PVC 挂载
    if (createFormData.value.pvc) {
      podYAML.spec.containers[0].volumeMounts.push({
        name: "pvc-volume",
        mountPath: createFormData.value.mountPath,
      });
      podYAML.spec.volumes.push({
        name: "pvc-volume",
        persistentVolumeClaim: {
          claimName: createFormData.value.pvc,
        },
      });
    }

    // 模拟 API 调用
    const response = await client.post(`${apiRoot}/pods`, podYAML);

    if (response.status > 299) {
      throw await response.json();
    }

    MessagePlugin.success("Pod 创建成功");
    closeCreateDialog(); // 关闭对话框
    fetchPods(); // 刷新列表
  } catch (error) {
    console.error("创建 Pod 失败:", error);
    MessagePlugin.error("创建 Pod 失败");
    MessagePlugin.error(error.message);
  }
};

// 显示创建 Pod 的对话框
const showCreateDialog = async () => {
  await fetchPVCs(); // 获取 PVC 列表
  createDialogVisible.value = true;
};

// 关闭创建 Pod 的对话框
const closeCreateDialog = () => {
  createDialogVisible.value = false;
  createFormData.value = {
    name: "",
    image: "",
    command: "",
    args: "",
    pvc: "",
    mountPath: "/root",
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
  switch (status) {
    case "Running":
      return "success";
    case "Pending":
      return "warning";
    case "Failed":
      return "danger";
    default:
      return "default";
  }
};

const showLog = async (podName) => {
  try {
    const response = await client.get(`${apiRoot}/pods/${podName}/log`);
    logContent.value = await response.text();
    logDialogVisible.value = true;
  } catch (error) {
    console.error("获取日志失败:", error);
    MessagePlugin.error("获取日志失败");
  }
};

const execWebTerm = (podName, container) => {
  const href = router.resolve({
    name: "webterm",
    query: { podName, container },
  }).href;
  window.open(href, "_blank");
};

const fetchPodCodeStatus = async (podName) => {
  return (await client.get(`/_/code-server/{!NAMESPACE}/${podName}`)).json();
};

const attachCode = async (podName) => {
  try {
    selectedPodCodeAttaching.value = true;
    const response = (
      await client.post(`/_/code-server/{!NAMESPACE}/${podName}`)
    ).json();
    if (response.status === "success") {
      window.open(
        `/kube/_/code-server/u-${username}/${podName}/proxy`,
        "_blank"
      );
      selectedPodCodeRunning.value = true;
    } else {
      MessagePlugin.error("启动 Code Server 失败");
      MessagePlugin.error(response.error);
    }
    selectedPodCodeAttaching.value = false;
  } catch (error) {
    console.error("启动 Code Server 失败:", error);
    MessagePlugin.error("启动 Code Server 失败");
  }
};

const detachCode = async (podName) => {
  try {
    selectedPodCodeDetaching.value = true;
    const response = (
      await client.delete(`/_/code-server/{!NAMESPACE}/${podName}`)
    ).json();
    if (response.status === "success") {
      MessagePlugin.success("关闭 Code Server 成功");
      selectedPodCodeRunning.value = false;
    } else {
      MessagePlugin.error("关闭 Code Server 失败");
      MessagePlugin.error(response.error);
    }
    selectedPodCodeDetaching.value = false;
  } catch (error) {
    console.error("关闭 Code Server 失败:", error);
    MessagePlugin.error("关闭 Code Server 失败");
  }
};

const openCode = (podName) => {
  window.open(`/kube/_/code-server/u-${username}/${podName}/proxy`, "_blank");
};

const handleMoreTools = async (podInfo) => {
  if (podInfo.status.phase !== "Running") {
    MessagePlugin.error("Pod 不在 Running 状态");
    return;
  }
  // defaults to the first container
  selectedPodContianer.value = podInfo.metadata.name;
  selectedPod.value = podInfo.metadata.name;
  currentPodContainers.value = podInfo.spec.containers.map(
    (container) => container.name
  );
  // selectedPodCodeStatusLoading.value = true;
  selectContainerDialogVisible.value = true;
  // const running = await fetchPodCodeStatus(podInfo.metadata.name);
  // selectedPodCodeStatusLoading.value = false;
  // // TODO: read from backend
  // selectedPodCodeRunning.value = false;
};

// 组件挂载时获取数据
onMounted(async () => {
  await client.ensureUsername();
  // opening new window will require this
  username = client.username;
  fetchPods();
});
</script>
