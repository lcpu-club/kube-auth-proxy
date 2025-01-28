<template>
  <main>
    <h1>统计数据<sup>*</sup></h1>
    <span class="block"
      ><b>参与人数：</b>{{ Object.entries(stats.players).length }}</span
    >
    <span class="block"><b>占有像素数：</b>{{ stats.sum }}</span>
    <span class="block"
      ><b>占有比例：</b
      >{{ ((stats.sum / (600 * 800)) * 100).toFixed(2) }}%</span
    >
    <span class="block" style="margin-top: 1em"><b>最后涂下的像素</b></span>
    <span class="block"
      >位置：({{ stats.lastPixel.x }}, {{ stats.lastPixel.y }})</span
    >
    <span class="block"
      >颜色：<span
        :style="{
          textDecorationColor: stats.lastPixel.color,
          textDecorationLine: 'underline',
          textDecorationStyle: 'solid',
          textDecorationThickness: '0.4em',
          textDecorationSkipInk: 'none',
          textUnderlineOffset: '-0.1em',
        }"
        >{{ stats.lastPixel.color }}</span
      ></span
    >
    <span class="block"
      >时间：{{
        new Date(stats.lastPixel.last_updated * 1000).toLocaleString()
      }}</span
    >
    <span class="block">所有者：{{ stats.lastPixel.owner }}</span>
    <h2>排行榜</h2>
    <table>
      <tr v-for="(stat, player) of stats.players">
        <td
          class="text-right font-sm no-wrap text-ellipsis"
          style="max-width: 8em"
        >
          {{ player }}
        </td>
        <td class="w-full no-wrap font-sm font-mono">
          <div
            class="font-sm font-mono"
            :style="{
              display: 'inline-block',
              height: '1.1em',
              width: `${(stat.sum / max) * 100}%`,
              backgroundColor: stringToColor(player),
              verticalAlign: 'middle',
            }"
          ></div>
          <span style="margin-left: .2em; vertical-align: middle;">{{ stat.sum }}</span>
        </td>
      </tr>
    </table>

    <h2>颜色分布</h2>
    <table>
      <tr v-for="(stat, player) of stats.players">
        <td
          class="text-right font-sm no-wrap text-ellipsis"
          style="max-width: 8em"
        >
          {{ player }}
        </td>
        <td class="w-full">
          <div
            v-for="(sum, color) of stat.fav"
            class="color-block"
            :title="`${color}: ${sum}, ${((sum / stat.sum) * 100).toFixed(2)}%`"
            :style="{
              backgroundColor: color,
              width: `${(sum / stat.sum) * 100}%`,
            }"
          ></div>
        </td>
      </tr>
    </table>

    <span class="font-sm text-gray block" style="margin-top: 2em;">*未整理的数据，存在偏差。</span>
  </main>
</template>

<script setup lang="ts">
import stats from "@/data/canvas-statistics.json";
const max = Math.max(Object.entries(stats.players)[0][1].sum);

function stringToColor(input: string): string {
  const hash = input.split("").reduce((acc, curr) => {
    return acc * 31 + curr.charCodeAt(0);
  }, 0);
  const hue = hash % 360;
  const saturation = 70;
  const maxLightness = 20;
  const lightness = (hash % maxLightness) + 75;

  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}
</script>

<style scoped>
main {
  overflow-y: auto;
  overflow-x: hidden;
  height: 100vh;
  padding-bottom: 40px;
  padding-right: 40px;
  box-sizing: border-box;
}
.block {
  display: block;
}
.font-sm {
  font-size: 0.75rem;
}
.text-gray {
  color: #666;
}
.text-right {
  text-align: right;
}
.w-full {
  width: 100%;
}
.font-mono {
  font-family: monospace;
}
.color-block {
  display: inline-block;
  position: relative;
  transition: all 150ms ease-out;
  height: 1.1em;
}
.color-block:hover {
  transform: scale(1.15);
  z-index: 999;
  box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.2);
}
.no-wrap {
  white-space: nowrap;
}
.text-ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
