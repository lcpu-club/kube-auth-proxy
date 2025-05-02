<script setup lang="ts">
import type { Partition } from "@/model/partition";
import { MessagePlugin } from "tdesign-vue-next";
import { computed, onMounted, ref, watch } from "vue";

const model = defineModel<{
  image: string;
  partition: string;
  command: string;
  args: string;
  pvc: string;
  mountPath: string;
  labels: string;
  lxcfsEnabled: boolean;
  sshEnabled: boolean;
  rdma: boolean;
  cpuRequest: number;
  memoryRequest: number;
  gpuTag: string | undefined;
  gpuRequest: number | undefined;
}>({
  required: true,
});

defineProps<{
  pvcs: any[];
}>();

let partitions = ref<Partition[]>([]);

onMounted(async () => {
  partitions.value = await (
    await fetch("https://hpcgame.pku.edu.cn/oss/images/public/partitions.json")
  ).json();
  MessagePlugin.success("成功获取分区信息");
});

const selectedPartitionIndex = ref();
const selectedPartition = computed(() => {
  if (selectedPartitionIndex.value < 0) return undefined;
  return partitions.value[selectedPartitionIndex.value];
});
watch(selectedPartition, (newVal) => {
  if (!newVal) return;
  model.value.partition = newVal.Name;
  if (newVal.GPUTag) {
    model.value.gpuTag = newVal.GPUTag;
    model.value.gpuRequest = 0;
  } else {
    model.value.gpuTag = undefined;
    model.value.gpuRequest = undefined;
  }
});

function getImageLabel(image: string) {
  const match = image.match(/^.*\/(.+):(.+)$/);
  if (!match) {
    return image;
  }
  return `${match[1]}:${match[2]}`;
}
</script>

<template>
  <!-- 分区 -->
  <t-form-item label="分区" name="partition">
    <t-select
      v-model="selectedPartitionIndex"
      placeholder="请选择分区"
      :popup-props="{ overlayClassName: 'select__overlay-option' }"
    >
      <t-option
        v-for="(partition, index) in partitions"
        :key="partition.Name"
        :value="index"
        :label="partition.Name"
      >
        <div class="option-wrapper">
          <span class="option-label">{{ partition.Name }}</span>
          <span class="option-description">{{ partition.Description }}</span>
        </div>
      </t-option>
    </t-select>
  </t-form-item>

  <!-- 容器镜像 -->
  <t-form-item label="容器镜像" name="image">
    <t-select
      v-model="model.image"
      placeholder="请输入容器镜像"
      creatable
      filterable
    >
      <t-option
        v-for="image in selectedPartition?.Images"
        :key="image"
        :value="image"
        :label="getImageLabel(image)"
      />
    </t-select>
  </t-form-item>

  <!-- 命令 -->
  <t-form-item label="命令" name="command">
    <t-input v-model="model.command" placeholder="请输入命令" />
  </t-form-item>

  <!-- 参数 -->
  <t-form-item label="参数" name="args">
    <t-textarea v-model="model.args" placeholder="请输入参数（以换行分割）" />
  </t-form-item>

  <!-- PVC 选择 -->
  <t-form-item label="挂载 PVC" name="pvc">
    <t-select v-model="model.pvc" placeholder="请选择 PVC" clearable>
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
      v-model="model.mountPath"
      placeholder="请选择挂载路径"
      creatable
      filterable
    >
      <t-option value="/root" label="/root" />
      <t-option value="/data" label="/data" />
      <t-option value="/shared" label="/shared" />
    </t-select>
  </t-form-item>

  <!-- 标签 -->
  <t-form-item label="标签" name="labels" tips="YAML.metadata.labels">
    <t-input v-model="model.labels" placeholder="a=b, c=d" />
  </t-form-item>

  <!-- 注解：LXCFS -->
  <t-form-item label="启用 LXCFS" name="lxcfsEnabled">
    <t-switch v-model="model.lxcfsEnabled" />
  </t-form-item>

  <!-- 注解：SSH Operator -->
  <t-form-item label="启用 SSH 服务器" name="sshEnabled">
    <t-switch v-model="model.sshEnabled" />
  </t-form-item>

  <!-- RDMA -->
  <t-form-item label="RDMA" name="rdma">
    <t-switch v-model="model.rdma" />
  </t-form-item>

  <!-- 资源请求和限制 -->
  <t-form-item label="CPU" name="cpuRequest">
    <t-input-number
      v-model="model.cpuRequest"
      :min="0"
      :max="selectedPartition?.CPULimit ?? 0"
    />
  </t-form-item>

  <t-form-item label="内存 (GiB)" name="memoryRequest">
    <t-input-number
      v-model="model.memoryRequest"
      :min="0"
      :max="selectedPartition?.MemoryLimit ?? 0"
      placeholder=""
    />
  </t-form-item>

  <t-form-item
    v-if="selectedPartition?.GPUName"
    :label="selectedPartition?.GPUName"
    name="gpu"
  >
    <t-input-number
      :min="0"
      v-model="model.gpuRequest"
      placeholder="请输入 GPU 数量"
    />
  </t-form-item>
</template>

<style lang="css" scoped>
.option-label {
  display: block;
}
.option-description {
  font-size: 90%;
  color: gray;
  display: block;
}
</style>

<style>
.select__overlay-option .t-select-option {
  height: 100%;
  padding: 8px;
}
</style>
