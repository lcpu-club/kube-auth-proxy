<template>
  <div class="main">
    <div class="xterm" ref="terminal"></div>
  </div>
</template>

<script setup>
import { onDeactivated, onMounted, useTemplateRef } from "vue";
import { useRoute } from "vue-router";
import { client } from "@/api/client";
import { useToken } from "@/api/token";
import { ref } from "vue";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { WebLinksAddon } from "xterm-addon-web-links";
import { AttachAddon } from "@xterm/addon-attach";
import "xterm/css/xterm.css";

const route = useRoute();
const token = useToken();
const terminal = useTemplateRef("terminal");
const xterm = new Terminal();
const fitAddon = new FitAddon();
const webLinksAddon = new WebLinksAddon();
let ws;

const windowResizeHandler = () => {
  fitAddon.fit();
};

xterm.onData((data) => {
  if (ws.readyState != 1) {
    xterm.write(`WebSocket not ready, readyState ${ws.readyState}\n`);
    return;
  }
  const message = new Uint8Array(data.length + 1);
  message[0] = 0; // stdin 通道
  new TextEncoder().encodeInto(data, message.subarray(1));
  try {
    ws.send(message);
  } catch (error) {
    xterm.write(`WebSocket error: ${error.message}\n`);
    console.log("WebSocket error: ", error);
  }
});

onMounted(async () => {
  const userInfo = await (await client.get("/_/whoami")).json();
  const currentHost = window.location.host;
  const apiServer = `${
    currentHost.includes("localhost") ? "ws" : "wss"
  }://${currentHost}/kube`;
  const podName = route.query.podName;
  const namespace = `u-${userInfo.username}`;
  const container = route.query.container;
  const command = "/bin/bash";
  // 构造 WebSocket URL
  const wsURL = `${apiServer}/api/v1/namespaces/${namespace}/pods/${podName}/exec?command=${command}&container=${container}&stdin=true&stdout=true&stderr=true&tty=true&auth=${token.value}`;
  // 创建 WebSocket 连接
  ws = new WebSocket(wsURL, "v4.channel.k8s.io");
  ws.onclose = () => {
    xterm.write("\n\nWebSocket closed\n");
  };
  window.addEventListener("resize", windowResizeHandler);
  xterm.open(terminal.value);
  xterm.loadAddon(fitAddon);
  xterm.loadAddon(webLinksAddon);
  fitAddon.fit();
  const attachAddon = new AttachAddon(ws, {
    bidirectional: false,
  });
  xterm.loadAddon(attachAddon);
});

onDeactivated(() => {
  ws.close();
  window.removeEventListener("resize", windowResizeHandler);
});
</script>

<style scoped>
.main .xterm {
  width: 100%;
  height: 100vh;
  background: #000;
}
</style>
