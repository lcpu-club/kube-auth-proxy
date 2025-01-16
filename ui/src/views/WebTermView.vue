<template>
  <div class="main">
    <t-list
      ref="scrolling-list"
      size="small"
      split
      :scroll="{
        type: 'virtual',
        threshold: 1,
      }"
    >
      <t-list-item
        v-for="item in output"
        :content="item.message"
        :class="item.type"
      >
        <span :class="item.type">{{ item.message }}</span>
      </t-list-item>
    </t-list>
    <div class="input-wrapper">
      <t-textarea
        v-model="input"
        placeholder="输入命令"
        class="input"
        @keypress="handleInputKeypress"
      />
      <t-button class="send-btn" :disabled="!input" @click="send">
        <t-icon name="send" />
      </t-button>
    </div>
  </div>
</template>

<script setup>
import { onMounted, useTemplateRef } from "vue";
import { useRoute } from "vue-router";
import { client } from "@/api/api";
import { useToken } from "@/api/token";
import { ref } from "vue";

const route = useRoute();
const token = useToken();
const input = ref("");
const scrollingListRef = useTemplateRef("scrolling-list");

const output = ref([]);

let ws;

output.value.push({
  message: "Establishing WebSocket connection...",
  type: "info",
});

onMounted(async () => {
  const userInfo = await (await client.get("/_/whoami")).json();
  const apiServer = "wss://hpcgame.pku.edu.cn/kube";
  const podName = route.query.podName;
  const namespace = `u-${userInfo.username}`;
  const container = route.query.container;
  const command = "/bin/sh";
  // 构造 WebSocket URL
  const wsURL = `${apiServer}/api/v1/namespaces/${namespace}/pods/${podName}/exec?command=${command}&container=${container}&stdin=true&stdout=true&stderr=true&tty=true`;
  // 创建 WebSocket 连接
  ws = new WebSocket(wsURL, null, {
    headers: {
      Authorization: token.value,
    },
  });

  // WebSocket 连接成功
  ws.onopen = () => {
    appendOutput("WebSocket connection established.", "info");
  };

  // 处理 WebSocket 消息
  ws.onmessage = async (event) => {
    const data = new Uint8Array(await event.data.arrayBuffer());
    const channel = data[0]; // 通道标识符
    const message = new TextDecoder().decode(data.slice(1)); // 数据部分
    switch (channel) {
      case 1: // stdout
        appendOutput(`${message}`, "stdout");
        break;
      case 2: // stderr
        appendOutput(`${message}`, "stderr");
        break;
      case 3: // error
        appendOutput(`${message}`, "error");
        break;
      default:
        appendOutput(`message: ${message}`, `Unknown channel: ${channel}`);
    }
  };

  // WebSocket 连接关闭
  ws.onclose = () => {
    appendOutput("WebSocket connection closed.", "info");
  };

  // WebSocket 错误处理
  ws.onerror = (error) => {
    console.log(error);
    appendOutput(`WebSocket error: ${error.message}`, "error");
  };
});

function send() {
  if (!input.value) return;
  if (ws.CLOSING || ws.CLOSED) {
    appendOutput("WebSocket CLOSING or CLOSED.", "error");
    return;
  }
  const data = new Uint8Array(input.value.length + 1);
  data[0] = 0; // stdin 通道
  new TextEncoder().encodeInto(input.value, data.subarray(1));
  try {
    ws.send(data);
  } catch (error) {
    appendOutput(`WebSocket error: ${error.message}`, "error");
  }
  input.value = "";
}

// 在页面上显示输出
function appendOutput(message, type) {
  output.value.push({ message, type });
  scrollingListRef?.value.scrollTo({ index: output.value.length - 1 });
}

function handleInputKeypress(value, event) {
  if (event.e.key === "Enter" && (event.e.altKey || event.e.ctrlKey)) {
    send();
  }
}
</script>

<style scoped>
.stdout {
  color: #000;
}
.stderr {
  color: #f00;
}
.error {
  color: #f00;
}
.unknown {
  color: rgb(166, 147, 3);
}
.info {
  color: rgb(58, 58, 231);
}

.main {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  display: grid;
  padding: 1em;
  box-sizing: border-box;
  grid-template-rows: 1fr auto;
}

.input-wrapper {
  display: flex;
  gap: 0.5em;
}

.input {
  flex: 1;
}

.send-btn {
  height: 100%;
}
</style>
