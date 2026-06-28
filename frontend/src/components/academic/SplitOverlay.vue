<script setup>
import { computed } from 'vue'
import { FASES_SPLIT } from '../../utils/splitAnimation.js'

const props = defineProps({
  anim: { type: Object, default: null },
})

const activo = computed(() => {
  const p = props.anim?.phase
  return p === FASES_SPLIT.SPLITTING || p === FASES_SPLIT.PROMOTE
})

const phase = computed(() => props.anim?.phase)
const promoteProgress = computed(() => props.anim?.promoteProgress ?? 0)

const izq = computed(() => props.anim?.hojaIzqPos ?? { x: 0, y: 0 })
const der = computed(() => props.anim?.hojaDerPos ?? { x: 0, y: 0 })
const padre = computed(() => props.anim?.padrePos ?? { x: der.value.x, y: der.value.y - 90 })

/** Punto superior de la hoja derecha (origen de la promoción). */
const origen = computed(() => ({
  x: der.value.x,
  y: der.value.y - (props.anim?.hojaDerDim?.height ?? 54) / 2 - 4,
}))

const destino = computed(() => ({
  x: padre.value.x,
  y: padre.value.y + (props.anim?.padreDim?.height ?? 54) / 2 + 2,
}))

const claveX = computed(
  () => origen.value.x + (destino.value.x - origen.value.x) * promoteProgress.value,
)
const claveY = computed(
  () => origen.value.y + (destino.value.y - origen.value.y) * promoteProgress.value,
)

const linkPath = computed(() => {
  const y = izq.value.y + (props.anim?.hojaIzqDim?.height ?? 54) / 2 + 12
  return `M ${izq.value.x} ${y} L ${der.value.x} ${y}`
})
</script>

<template>
  <g v-if="activo" class="split-annotations">
    <!-- Enlace SPLIT entre las dos hojas reales del árbol nuevo -->
    <g v-if="phase === 'splitting' || phase === 'promote'">
      <path :d="linkPath" class="split-link" />
      <text
        :x="(izq.x + der.x) / 2"
        :y="izq.y + (anim?.hojaIzqDim?.height ?? 54) / 2 + 28"
        class="split-link-label"
        text-anchor="middle"
      >
        ✂️ SPLIT — una hoja dividida en dos
      </text>
    </g>

    <!-- Flecha de promoción sobre el árbol real -->
    <g v-if="phase === 'promote' && anim.clavePromovida != null">
      <line
        class="promote-line"
        :x1="origen.x"
        :y1="origen.y"
        :x2="destino.x"
        :y2="destino.y"
        pathLength="100"
        :stroke-dasharray="100"
        :stroke-dashoffset="100 * (1 - promoteProgress)"
      />

      <g :transform="`translate(${claveX}, ${claveY})`">
        <rect class="promoted-key" x="-34" y="-14" width="68" height="28" rx="4" />
        <text class="promoted-key__text" y="5">{{ anim.clavePromovida }}</text>
      </g>

      <text
        v-if="promoteProgress > 0.85"
        :x="padre.x"
        :y="padre.y - (anim?.padreDim?.height ?? 54) / 2 - 8"
        class="promote-label"
        text-anchor="middle"
      >
        ↑ separador en índice
      </text>
    </g>
  </g>
</template>

<style scoped>
.split-link {
  fill: none;
  stroke: #f97316;
  stroke-width: 2.5;
  stroke-dasharray: 8 5;
  animation: dashMove 0.8s linear infinite;
}

.split-link-label {
  fill: #ea580c;
  font-size: 9px;
  font-weight: 700;
}

@keyframes dashMove {
  to { stroke-dashoffset: -13; }
}

.promote-line {
  fill: none;
  stroke: #f59e0b;
  stroke-width: 2.5;
  marker-end: url(#arrow-promote);
}

.promoted-key {
  fill: #fef08a;
  stroke: #f59e0b;
  stroke-width: 2;
  filter: drop-shadow(0 0 10px rgba(251, 191, 36, 0.8));
}

.promoted-key__text {
  fill: #92400e;
  font-size: 11px;
  font-weight: 800;
  text-anchor: middle;
}

.promote-label {
  fill: #2563eb;
  font-size: 8px;
  font-weight: 700;
}
</style>
