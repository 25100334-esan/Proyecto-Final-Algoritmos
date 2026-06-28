<script setup>
import { computed } from 'vue'

const props = defineProps({
  nodo: { type: Object, required: true },
  x: { type: Number, required: true },
  y: { type: Number, required: true },
  dimensiones: { type: Object, required: true },
  highlighted: { type: Boolean, default: false },
  splitAnimating: { type: Boolean, default: false },
  overflowFull: { type: Boolean, default: false },
  oculto: { type: Boolean, default: false },
  atenuado: { type: Boolean, default: false },
  highlightedKeys: { type: Array, default: () => [] },
  deletingKeys: { type: Array, default: () => [] },
  activeSlot: { type: Number, default: null },
  activeKey: { type: [Number, String], default: null },
})

const esHoja = computed(() => props.nodo.esHoja)
const width = computed(() => props.dimensiones.width)
const height = computed(() => props.dimensiones.height)
const celdas = computed(() => props.dimensiones.celdas)

function celdaActiva(celda) {
  if (celda.tipo === 'key') {
    return props.highlightedKeys.includes(celda.valor) || props.activeKey === celda.valor
  }
  if (celda.tipo === 'sep') {
    return props.activeKey === celda.valor
  }
  if (celda.tipo === 'ptr') {
    return props.activeSlot != null && celda.indice === props.activeSlot && props.activeKey == null
  }
  return false
}

function celdaEliminando(celda) {
  return celda.tipo === 'key' && props.deletingKeys.includes(celda.valor)
}

function offsetX(idx) {
  let o = 0
  for (let i = 0; i < idx; i++) o += celdas.value[i].ancho
  return o
}
</script>

<template>
  <g
    v-show="!oculto"
    class="tree-node"
    :class="{
      'tree-node--index': !esHoja,
      'tree-node--leaf': esHoja,
      'tree-node--highlight': highlighted,
      'tree-node--split': splitAnimating,
      'tree-node--overflow': overflowFull,
      'tree-node--atenuado': atenuado,
    }"
    :opacity="atenuado ? 0.12 : 1"
    :transform="`translate(${x - width / 2}, ${y - height / 2})`"
  >
    <!-- Borde doble en hojas (Sequence Set) -->
    <rect
      v-if="esHoja"
      class="tree-node__outer-leaf"
      :width="width + 6"
      :height="height + 6"
      x="-3"
      y="-3"
      rx="4"
    />

    <rect
      class="tree-node__box"
      :width="width"
      :height="height"
      rx="3"
    />

    <text class="tree-node__tipo" :x="width / 2" y="11">
      {{ esHoja ? 'HOJA' : 'ÍNDICE' }}
    </text>

    <g v-for="(celda, idx) in celdas" :key="idx">
      <rect
        class="celda"
        :class="{
          'celda--ptr': celda.tipo === 'ptr',
          'celda--sep': celda.tipo === 'sep',
          'celda--key': celda.tipo === 'key',
          'celda--active': celdaActiva(celda),
          'celda--deleting': celdaEliminando(celda),
        }"
        :x="offsetX(idx)"
        y="16"
        :width="celda.ancho"
        :height="height - 20"
      />

      <text
        v-if="celda.tipo === 'sep' || celda.tipo === 'key'"
        class="celda__texto"
        :class="{
          'celda__texto--leaf': esHoja,
          'celda__texto--deleting': celdaEliminando(celda),
        }"
        :x="offsetX(idx) + celda.ancho / 2"
        :y="height / 2 + 6"
      >
        {{ celda.valor }}
      </text>

      <!-- Subrayado en claves de hoja -->
      <line
        v-if="esHoja && celda.tipo === 'key'"
        class="celda__underline"
        :class="{ 'celda__underline--active': celdaActiva(celda) }"
        :x1="offsetX(idx) + 4"
        :x2="offsetX(idx) + celda.ancho - 4"
        :y1="height / 2 + 10"
        :y2="height / 2 + 10"
      />

      <!-- Divisor vertical entre celdas -->
      <line
        v-if="idx < celdas.length - 1"
        class="celda__divider"
        :x1="offsetX(idx) + celda.ancho"
        :x2="offsetX(idx) + celda.ancho"
        y1="18"
        :y2="height - 2"
      />
    </g>
  </g>
</template>

<style scoped>
.tree-node__outer-leaf {
  fill: none;
  stroke: #7c3aed;
  stroke-width: 2.5;
}

.tree-node__box {
  fill: #fafafa;
  stroke: #374151;
  stroke-width: 1.5;
  transition: filter 0.3s;
}

.tree-node--index .tree-node__box {
  stroke: #2563eb;
}

.tree-node--leaf .tree-node__box {
  fill: #fdf4ff;
  stroke: #7c3aed;
  stroke-width: 2;
}

.tree-node--highlight .tree-node__box {
  filter: drop-shadow(0 0 10px #fbbf24);
  stroke: #f59e0b;
  stroke-width: 2.5;
}

.tree-node--split .tree-node__box {
  animation: splitPulse 0.7s ease-in-out infinite;
}

.tree-node--overflow .tree-node__box {
  fill: #fee2e2;
  stroke: #ef4444;
  stroke-width: 3;
  animation: overflowShake 0.5s ease-in-out infinite;
}

@keyframes overflowShake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-2px); }
  75% { transform: translateX(2px); }
}

@keyframes splitPulse {
  0%, 100% { filter: drop-shadow(0 0 4px #f97316); }
  50% { filter: drop-shadow(0 0 18px #f97316); }
}

.tree-node__tipo {
  fill: #6b7280;
  font-size: 7px;
  text-anchor: middle;
  letter-spacing: 0.08em;
  font-weight: 700;
}

.celda {
  fill: transparent;
  stroke: none;
  transition: fill 0.2s;
}

.celda--ptr {
  fill: #bfdbfe;
  stroke: #93c5fd;
  stroke-width: 0.5;
}

.celda--sep {
  fill: #ffffff;
}

.celda--key {
  fill: #faf5ff;
}

.celda--active {
  fill: #fef08a !important;
  stroke: #f59e0b;
  stroke-width: 1.5;
}

.celda--deleting {
  fill: #fecaca !important;
  stroke: #ef4444;
  stroke-width: 2;
  animation: deleteBlink 0.6s ease-in-out infinite;
}

@keyframes deleteBlink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.35; }
}

.celda__texto--deleting {
  text-decoration: line-through;
  fill: #b91c1c;
}

.celda__texto {
  fill: #111827;
  font-size: 11px;
  font-weight: 600;
  text-anchor: middle;
  dominant-baseline: middle;
}

.celda__texto--leaf {
  text-decoration: none;
}

.celda__underline {
  stroke: #7c3aed;
  stroke-width: 1;
}

.celda__underline--active {
  stroke: #f59e0b;
  stroke-width: 2;
}

.celda__divider {
  stroke: #9ca3af;
  stroke-width: 0.5;
}
</style>
