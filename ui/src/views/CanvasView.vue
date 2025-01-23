<script setup lang="ts">
import { getToken } from "@/api/token";
import {
  onMounted,
  useTemplateRef,
  ref,
  computed,
  watch,
  type CSSProperties,
  onUnmounted,
} from "vue";

const currentHost = window.location.host;
const baseUrl = `${
  currentHost.includes("localhost") ? "http" : "https"
}://${currentHost}`;

const WIDTH_PX = ref(800);
const HEIGHT_PX = ref(600);
const SCALING_RATIO = ref(2);

const canvasHeight = computed(() => HEIGHT_PX.value * SCALING_RATIO.value);
const canvasWidth = computed(() => WIDTH_PX.value * SCALING_RATIO.value);

const showDetails = ref(false);
const showDetailsCloseHint = ref(true);
const details = ref<{
  owner: string;
  color: string;
  last_updated: number;
}>();

const detailsTimestamp = computed(() => {
  if (!details.value) return "";
  return new Date(details.value.last_updated).toLocaleString("zh-cn", {
    year: undefined,
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
  });
});

const canvas = useTemplateRef("canvasRef");
const canvasWrapper = useTemplateRef("canvasWrapperRef");
const detailsElement = useTemplateRef("detailsRef");
const detailsMeasure = useTemplateRef("detailsMeasuringRef");
const detailsCoords = ref({ x: 0, y: 0 });
const detailsColorUnderlineStyle = computed<CSSProperties>(() => {
  return {
    textDecorationColor: details.value?.color,
    textDecorationLine: "underline",
    textDecorationStyle: "solid",
    textDecorationThickness: "0.4em",
    textDecorationSkipInk: "none",
    textUnderlineOffset: "-0.1em",
  };
});
let context: CanvasRenderingContext2D | null | undefined;

let imageBitmap: ImageBitmap | null = null;

const updateCanvas = () => {
  if (!context) {
    alert("Canvas not supported");
    return;
  }
  if (!imageBitmap) return;

  context.imageSmoothingEnabled = false;
  context.clearRect(0, 0, canvas.value?.width ?? 0, canvas.value?.height ?? 0);
  context.drawImage(
    imageBitmap,
    0,
    0,
    WIDTH_PX.value * SCALING_RATIO.value,
    HEIGHT_PX.value * SCALING_RATIO.value
  );
};

const handleCanvasScroll = () => {
  showDetails.value = false;
  showDetailsCloseHint.value = false;
};

const handleKeyPress = (e: KeyboardEvent) => {
  if (e.key === "Escape") {
    showDetails.value = false;
    showDetailsCloseHint.value = false;
  }
};

const placeDetailsElement = (mouseX: number, mouseY: number) => {
  if (!detailsElement.value || !detailsMeasure.value) return;
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight;
  if (20 + mouseX + (detailsMeasure.value?.clientWidth ?? 0) > windowWidth) {
    detailsElement.value.style.right = `${windowWidth - mouseX}px`;
    detailsElement.value.style.left = "auto";
  } else {
    detailsElement.value.style.left = `${mouseX}px`;
    detailsElement.value.style.right = "auto";
  }
  if (20 + mouseY + (detailsMeasure.value?.clientHeight ?? 0) > windowHeight) {
    detailsElement.value.style.bottom = `${windowHeight - mouseY}px`;
    detailsElement.value.style.top = "auto";
  } else {
    detailsElement.value.style.top = `${mouseY}px`;
    detailsElement.value.style.bottom = "auto";
  }
};

const handleCanvasClicked = async (e: MouseEvent) => {
  if (!detailsElement.value) return;
  const rect = canvas.value?.getBoundingClientRect();
  const x = Math.floor((e.clientX - (rect?.left ?? 0)) / SCALING_RATIO.value);
  const y = Math.floor((e.clientY - (rect?.top ?? 0)) / SCALING_RATIO.value);
  detailsCoords.value = { x, y };
  details.value = undefined;
  placeDetailsElement(e.clientX, e.clientY);
  showDetails.value = true;
  const pixel = await (
    await fetch(`${baseUrl}/kube/_/painter/board/pixel?x=${x}&y=${y}`)
  ).json();
  if (pixel.owner === "HPCGame") {
    details.value = {
      owner: "HPCGame",
      color: pixel.color,
      last_updated: pixel.last_updated,
    };
    return;
  }
  const token = getToken();
  const userProfile = await (
    await fetch(`${baseUrl}/api/user/${pixel.owner}/profile`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
  ).json();
  details.value = {
    owner: userProfile.name,
    color: pixel.color,
    last_updated: pixel.last_updated,
  };

  placeDetailsElement(e.clientX, e.clientY);
};

onMounted(async () => {
  try {
    context = canvas.value?.getContext("2d");
    if (!context) {
      alert("Canvas not supported");
      return;
    }
    const image = await (
      await fetch(`${baseUrl}/kube/_/painter/board/`)
    ).blob();
    imageBitmap = await createImageBitmap(image);
    HEIGHT_PX.value = imageBitmap.height;
    WIDTH_PX.value = imageBitmap.width;
    updateCanvas();
    document.addEventListener("keydown", handleKeyPress);
  } catch (e) {
    console.error(e);
    alert(
      "Unknown error occurred. Refer to console for more information.\nRefresh the page to reload, contact administrator if the problem persists."
    );
  }
});

onUnmounted(() => {
  document.removeEventListener("keydown", handleKeyPress);
});

watch(
  SCALING_RATIO,
  (newVal, oldVal) => {
    showDetails.value = false;
    const rect = canvas.value?.getBoundingClientRect();
    const wrapperRect = canvasWrapper.value?.getBoundingClientRect();
    if (!rect || !wrapperRect) updateCanvas();
    else {
      const currentCenterX =
        (-rect.left + wrapperRect.left + wrapperRect.width / 2) / oldVal;
      const currentCenterY =
        (-rect.top + wrapperRect.top + wrapperRect.height / 2) / oldVal;
      const newRectX = currentCenterX * newVal - wrapperRect.width / 2;
      const newRectY = currentCenterY * newVal - wrapperRect.height / 2;
      canvasWrapper.value?.scrollTo(newRectX, newRectY);
      updateCanvas();
    }
  },
  { flush: "post" } // so that the canvas is updated after canvas has been resized
);
watch(HEIGHT_PX, updateCanvas, { flush: "post" });
watch(WIDTH_PX, updateCanvas, { flush: "post" });
</script>

<template>
  <div id="canvas-wrapper" @scroll="handleCanvasScroll" ref="canvasWrapperRef">
    <canvas
      ref="canvasRef"
      :height="canvasHeight"
      :width="canvasWidth"
      @click="handleCanvasClicked"
    ></canvas>
  </div>

  <div id="scaling-ratio-slider-wrapper">
    <input
      id="scaling-ratio-slider"
      type="range"
      v-model="SCALING_RATIO"
      min="1"
      max="32"
    />
    <label id="scaling-ratio-label">Scaling Ratio: {{ SCALING_RATIO }}x</label>
  </div>

  <Transition>
    <div id="details" v-show="showDetails" ref="detailsRef">
      <div v-if="!details">
        <div class="loader-wrapper">
          <div class="loader"></div>
          <span>Loading...</span>
        </div>
      </div>
      <div v-else>
        <div>
          <span class="details-title">Coordinates</span>
          <span class="details-content"
            >({{ detailsCoords.y }}, {{ detailsCoords.x }})</span
          >
        </div>
        <div>
          <span class="details-title">Color</span>
          <span class="details-content" :style="detailsColorUnderlineStyle">{{
            details.color
          }}</span>
        </div>
        <div>
          <span class="details-title">Owner</span>
          <span class="details-content">{{ details.owner }}</span>
        </div>
        <div>
          <span class="details-title">Last updated</span>
          <span class="details-content">{{ detailsTimestamp }}</span>
        </div>
        <span class="details-close-hint" v-if="showDetailsCloseHint"
          >Press escape to omit.</span
        >
      </div>
    </div>
  </Transition>

  <div id="details-measuring-element" ref="detailsMeasuringRef">
    <div v-if="!details">
      <div class="loader-wrapper">
        <div class="loader"></div>
        <span>Loading...</span>
      </div>
    </div>
    <div v-else>
      <div>
        <span class="details-title">Coordinates</span>
        <span class="details-content"
          >({{ detailsCoords.y }}, {{ detailsCoords.x }})</span
        >
      </div>
      <div>
        <div>
          <span class="details-title">Color</span>
          <span class="details-content" :style="detailsColorUnderlineStyle">{{
            details.color
          }}</span>
        </div>
        <div>
          <span class="details-title">Owner</span>
          <span class="details-content">{{ details.owner }}</span>
        </div>
        <span class="details-title">Last updated</span>
        <span class="details-content">detailsTimestamp</span>
      </div>
      <span class="details-close-hint" v-if="showDetailsCloseHint"
        >Press escape to omit.</span
      >
    </div>
  </div>
</template>

<style scoped>
#canvas-wrapper {
  height: 100vh;
  width: 100%;
  overflow: auto;
}

#scaling-ratio-slider-wrapper {
  position: fixed;
  left: 0.5em;
  bottom: 0.5em;
  padding: 0.5em;
  border-radius: 0.5em;
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(24px);
  display: flex;
  border: 2px solid rgba(0, 0, 0, 0.1);
  box-shadow: rgba(0, 0, 0, 0.1) 0px 20px 25px -5px,
    rgba(0, 0, 0, 0.04) 0px 10px 10px -5px;
  cursor: grab;
}

#scaling-ratio-label {
  margin-left: 0.5em;
  pointer-events: none;
  user-select: none;
}

#details,
#details-measuring-element {
  position: fixed;
  padding: 0.5em;
  min-width: 8em;
  border-radius: 0.5em;
  border: 2px solid rgba(0, 0, 0, 0.1);
}

#details {
  background-color: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(24px);
  box-shadow: rgba(0, 0, 0, 0.1) 0px 20px 25px -5px,
    rgba(0, 0, 0, 0.04) 0px 10px 10px -5px;
}

#details-measuring-element {
  top: 0;
  left: 0;
  pointer-events: none;
  opacity: 0;
}

.details-title {
  font-weight: bold;
  text-transform: uppercase;
  font-size: 0.5em;
  letter-spacing: 0.1em;
  color: gray;
  display: block;
}

.details-content {
  font-size: 1.25em;
  display: block;
}

.details-close-hint {
  font-size: 0.75em;
  color: gray;
  display: block;
  float: right;
}

#details-mesauring-element > div,
#details > div {
  display: flex;
  gap: 0.5em;
  flex-direction: column;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 150ms;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.loader {
  width: 30px;
  aspect-ratio: 1;
  --g: conic-gradient(from -90deg at 5px 5px, #000 90deg, #0000 0);
  background: var(--g), var(--g), var(--g);
  background-size: 50% 50%;
  animation: l18 1s infinite;
}

.loader-wrapper {
  display: flex;
  gap: 0.5em;
  align-items: center;
  justify-content: center;
}

@keyframes l18 {
  0% {
    background-position: 0 0, 5px 5px, 10px 10px;
  }
  33% {
    background-position: -15px 0, 5px 5px, 10px 10px;
  }
  66% {
    background-position: -15px 0, 5px 20px, 10px 10px;
  }
  100% {
    background-position: -15px 0, 5px 20px, 25px 10px;
  }
}
</style>
