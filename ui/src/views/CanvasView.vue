<script setup lang="ts">
import { getToken } from "@/api/token";
import { ScaleIcon } from "@heroicons/vue/24/outline";
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

const canvasWidth = ref(800);
const canvasHeight = ref(600);
const canvasScalingRatio = ref(1);

const scrollingCofficient = navigator.userAgent.includes("Macintosh") || navigator.userAgent.includes("Mac OS X") ? 1 : -1;

const showDetails = ref(false);
const showDetailsCloseHint = ref(true);
const details = ref<{
  owner: string;
  color: string;
  last_updated: number;
}>();

const detailsTimestamp = computed(() => {
  if (!details.value) return "";
  return new Date(details.value.last_updated * 1000).toLocaleString("zh-cn", {
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

const canvasTranslateX = ref(0);
const canvasTranslateY = ref(0);

let context: CanvasRenderingContext2D | null | undefined;

let imageBitmap: ImageBitmap | null = null;
let cameraMove = false;
let isCameraMoving = ref(false);
let initialCanvasTranslateX = 0;
let initialCanvasTranslateY = 0;

const updateCanvas = () => {
  if (!context) {
    alert("Canvas not supported");
    return;
  }
  if (!imageBitmap) return;

  context.imageSmoothingEnabled = false;
  context.clearRect(0, 0, canvas.value?.width ?? 0, canvas.value?.height ?? 0);
  context.drawImage(imageBitmap, 0, 0, canvasWidth.value, canvasHeight.value);
};

const handleCanvasScroll = (e: WheelEvent) => {
  const initialScalingRatio = canvasScalingRatio.value;
  canvasScalingRatio.value = Math.min(
    32,
    Math.max(
      0.5,
      canvasScalingRatio.value +
        ((e.deltaY / 100) * canvasScalingRatio.value) / 10 * scrollingCofficient
    )
  );
  showDetails.value = false;
  showDetailsCloseHint.value = false;
  const rect = canvas.value?.getBoundingClientRect();
  if (!rect) return;
  const centerX = rect.x + rect.width / 2;
  const centerY = rect.y + rect.height / 2;
  const offsetX = (centerX - e.clientX) / initialScalingRatio;
  const offsetY = (centerY - e.clientY) / initialScalingRatio;
  canvasTranslateX.value +=
    offsetX * (canvasScalingRatio.value - initialScalingRatio);
  canvasTranslateY.value +=
    offsetY * (canvasScalingRatio.value - initialScalingRatio);
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
  if (isCameraMoving.value || !detailsElement.value) return;
  const rect = canvas.value?.getBoundingClientRect();
  const x = Math.floor(
    (e.clientX - (rect?.left ?? 0)) / canvasScalingRatio.value
  );
  const y = Math.floor(
    (e.clientY - (rect?.top ?? 0)) / canvasScalingRatio.value
  );
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
    canvasHeight.value = imageBitmap.height;
    canvasWidth.value = imageBitmap.width;
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

let beginMoveEvent: MouseEvent | null = null;

const handleCameraMove = (e: MouseEvent) => {
  if (!cameraMove || !canvasWrapper.value || !beginMoveEvent) return;
  isCameraMoving.value = true;
  showDetails.value = false;
  canvasTranslateX.value =
    initialCanvasTranslateX + (e.clientX - beginMoveEvent.clientX);
  canvasTranslateY.value =
    initialCanvasTranslateY + (e.clientY - beginMoveEvent.clientY);
};

const handleCameraMoveBegin = (e: MouseEvent) => {
  cameraMove = true;
  beginMoveEvent = e;
  initialCanvasTranslateX = canvasTranslateX.value;
  initialCanvasTranslateY = canvasTranslateY.value;
};

const handleCameraMoveEnd = (e: MouseEvent) => {
  cameraMove = false;
  beginMoveEvent = null;
  setTimeout(() => {
    isCameraMoving.value = false;
  }, 100);
};

const handleCameraMoveBeginTouch = (e: TouchEvent) => {
  if (e.touches.length !== 1) return;
  const touch = e.touches[0];
  cameraMove = true;
  beginMoveEvent = new MouseEvent("mousedown", {
    clientX: touch.clientX,
    clientY: touch.clientY,
  });
  initialCanvasTranslateX = canvasTranslateX.value;
  initialCanvasTranslateY = canvasTranslateY.value;
};

const handleCameraMoveTouch = (e: TouchEvent) => {
  if (!cameraMove || !canvasWrapper.value || !beginMoveEvent) return;
  isCameraMoving.value = true;
  showDetails.value = false;
  const touch = e.touches[0];
  canvasTranslateX.value =
    initialCanvasTranslateX + (touch.clientX - beginMoveEvent.clientX);
  canvasTranslateY.value =
    initialCanvasTranslateY + (touch.clientY - beginMoveEvent.clientY);
};

const handleCameraMoveEndTouch = (e: TouchEvent) => {
  cameraMove = false;
  beginMoveEvent = null;
  setTimeout(() => {
    isCameraMoving.value = false;
  }, 100);
};

watch(canvasHeight, updateCanvas, { flush: "post" });
watch(canvasWidth, updateCanvas, { flush: "post" });
</script>

<template>
  <div
    id="canvas-wrapper"
    @wheel.prevent="handleCanvasScroll"
    @mousemove="handleCameraMove"
    @mousedown="handleCameraMoveBegin"
    @mouseup="handleCameraMoveEnd"
    @mouseleave="handleCameraMoveEnd"
    @touchstart.prevent="handleCameraMoveBeginTouch"
    @touchmove.prevent="handleCameraMoveTouch"
    @touchcancel.prevent="handleCameraMoveEndTouch"
    @touchend.prevent="handleCameraMoveEndTouch"
    ref="canvasWrapperRef"
    :style="{
      cursor: isCameraMoving ? 'grabbing' : 'default',
    }"
  >
    <canvas
      ref="canvasRef"
      :height="canvasHeight"
      :width="canvasWidth"
      @click="handleCanvasClicked"
      :style="{
        transform: `translate(${canvasTranslateX}px, ${canvasTranslateY}px) scale(${canvasScalingRatio})`,
      }"
    ></canvas>
  </div>

  <div id="scaling-ratio-slider-wrapper">
    <input
      id="scaling-ratio-slider"
      type="range"
      v-model.number="canvasScalingRatio"
      min="0.5"
      max="32"
      step="0.5"
    />
    <label id="scaling-ratio-label">
      {{ (canvasScalingRatio * 100).toFixed(2) }}%</label
    >
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
        <div v-if="details.last_updated">
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
        <div v-if="details.last_updated">
          <span class="details-title">Last updated</span>
          <span class="details-content">detailsTimestamp</span>
        </div>
      </div>
      <span class="details-close-hint" v-if="showDetailsCloseHint"
        >Press escape to omit.</span
      >
    </div>
  </div>
</template>

<style scoped>
canvas {
  image-rendering: pixelated;
  box-shadow: rgba(0, 0, 0, 0.25) 0px 25px 50px -12px;
}

#canvas-wrapper {
  position: relative;
  height: 100vh;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
  background-image: linear-gradient(45deg, #f3f3f3 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #f3f3f3 75%),
    linear-gradient(135deg, #f3f3f3 25%, transparent 25%),
    linear-gradient(135deg, transparent 75%, #f3f3f3 75%);
  background-size: 20px 20px;
  background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
}

#scaling-ratio-slider-wrapper {
  position: absolute;
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
