<script setup>
import { computed, ref, onUnmounted } from 'vue'
import TreeNode from './TreeNode.vue'
import StepBanner from './StepBanner.vue'
import SplitOverlay from './SplitOverlay.vue'
import DeleteOverlay from './DeleteOverlay.vue'
import { calcularLayout, anclaPuntero } from '../../utils/treeLayout.js'

const props = defineProps({
  estructura: { type: Object, default: null },
  splitAnim: { type: Object, default: null },
  deleteAnim: { type: Object, default: null },
  highlightedNodes: { type: Array, default: () => [] },
  splitNodes: { type: Array, default: () => [] },
  underflowNodes: { type: Array, default: () => [] },
  highlightedKeys: { type: Array, default: () => [] },
  deletingKeys: { type: Array, default: () => [] },
  rangeHighlightLeaf: { type: Number, default: null },
  activeLeafChain: { type: Object, default: null },
  activeEdge: { type: String, default: null },
  activeSlot: { type: Number, default: null },
  activeKey: { type: [Number, String], default: null },
  pasoMensaje: { type: String, default: '' },
  animando: { type: Boolean, default: false },
  pasoActual: { type: Number, default: 0 },
  pasoTotal: { type: Number, default: 0 },
  modoPasoAPaso: { type: Boolean, default: false },
})

const emit = defineEmits(['siguiente'])

const viewportRef = ref(null)
const zoom = ref(1)
const panX = ref(0)
const panY = ref(0)
const arrastrando = ref(false)
const ultimoMouse = { x: 0, y: 0 }

const MIN_ZOOM = 0.08
const MAX_ZOOM = 48
const ZOOM_STEP = 1.18
/** Zoom base: encaja todo el árbol en el viewport (100% = ver todo). */
const ZOOM_FIT = 1

const layout = computed(() => calcularLayout(props.estructura))

const nodos = computed(() => {
  if (!props.estructura?.nodos) return []
  return props.estructura.nodos.map((n) => ({
    ...n,
    indices: n.indices?.length ? n.indices : (n.claves ?? []),
  }))
})

const overflowNodeId = computed(() => props.splitAnim?.overflowNodeId ?? props.deleteAnim?.underflowNodeId ?? null)

function leafChainActive(link) {
  if (props.activeLeafChain) {
    return (
      props.activeLeafChain.from === link.from && props.activeLeafChain.to === link.to
    )
  }
  return (
    props.rangeHighlightLeaf === link.from || props.rangeHighlightLeaf === link.to
  )
}

const bounds = computed(() => {
  const l = layout.value
  if (!l.width) return { x: -400, y: -10, w: 800, h: 420 }
  return { x: -l.width / 2, y: -10, w: l.width, h: l.height + 20 }
})

const viewBox = computed(() => {
  const b = bounds.value
  const w = b.w / zoom.value
  const h = b.h / zoom.value
  const cx = b.x + b.w / 2
  const cy = b.y + b.h / 2
  return `${cx - w / 2 + panX.value} ${cy - h / 2 + panY.value} ${w} ${h}`
})

const zoomPorcentaje = computed(() => Math.round(zoom.value * 100))

/** Etiqueta legible cuando el zoom supera 1000%. */
const zoomEtiqueta = computed(() => {
  const p = zoomPorcentaje.value
  return p >= 1000 ? `${(p / 100).toFixed(1)}×` : `${p}%`
})

function clampZoom(z) {
  return Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, z))
}

function zoomIn(factor = ZOOM_STEP) {
  zoom.value = clampZoom(zoom.value * factor)
}

function zoomOut(factor = ZOOM_STEP) {
  zoom.value = clampZoom(zoom.value / factor)
}

function resetZoom() {
  zoom.value = ZOOM_FIT
  panX.value = 0
  panY.value = 0
}

/** Encaja todo el árbol en pantalla (zoom 100%). */
function verTodo() {
  resetZoom()
}

/** Acercamiento rápido para árboles grandes (+2× por clic). */
function zoomInRapido() {
  zoomIn(2)
}

function escalaPixelASvg() {
  const b = bounds.value
  const vp = viewportRef.value
  if (!vp?.clientWidth) return 1
  return b.w / zoom.value / vp.clientWidth
}

function onWheel(e) {
  if (!props.estructura) return
  e.preventDefault()
  // Rueda más fina al acercar mucho; sigue siendo O(1), sin lag.
  const factor = e.deltaY > 0 ? 1 / ZOOM_STEP : ZOOM_STEP
  zoom.value = clampZoom(zoom.value * factor)
}

function onMouseDown(e) {
  if (!props.estructura || e.button !== 0) return
  arrastrando.value = true
  ultimoMouse.x = e.clientX
  ultimoMouse.y = e.clientY
  viewportRef.value?.classList.add('canvas-viewport--dragging')
}

function onMouseMove(e) {
  if (!arrastrando.value) return
  const escala = escalaPixelASvg()
  panX.value += (ultimoMouse.x - e.clientX) * escala
  panY.value += (ultimoMouse.y - e.clientY) * escala
  ultimoMouse.x = e.clientX
  ultimoMouse.y = e.clientY
}

function onMouseUp() {
  arrastrando.value = false
  viewportRef.value?.classList.remove('canvas-viewport--dragging')
}

onUnmounted(() => onMouseUp())

defineExpose({ verTodo, resetZoom, zoomIn, zoomOut, zoomInRapido })

function edgePath(edge) {
  const from = anclaPuntero(edge.from, edge.ptrIndex, layout.value.posiciones, layout.value.dimensiones)
  const toPos = layout.value.posiciones[edge.to]
  const toDim = layout.value.dimensiones[edge.to]
  if (!from || !toPos || !toDim) return ''

  const x1 = from.x
  const y1 = from.y
  const x2 = toPos.x
  const y2 = toPos.y - toDim.height / 2
  const midY = y1 + (y2 - y1) * 0.45

  return `M ${x1} ${y1} L ${x1} ${midY} L ${x2} ${midY} L ${x2} ${y2}`
}

function leafChainPath(fromId, toId) {
  const from = layout.value.posiciones[fromId]
  const to = layout.value.posiciones[toId]
  const dimFrom = layout.value.dimensiones[fromId]
  const dimTo = layout.value.dimensiones[toId]
  if (!from || !to || !dimFrom || !dimTo) return ''

  const y = from.y + dimFrom.height / 2 + 14
  const x1 = from.x + dimFrom.width / 2
  const x2 = to.x - dimTo.width / 2

  return `M ${x1} ${y} L ${x2} ${y}`
}

function nodoResaltado(id) {
  return (
    props.highlightedNodes.includes(id) ||
    props.rangeHighlightLeaf === id ||
    props.splitNodes.includes(id) ||
    props.underflowNodes.includes(id)
  )
}
</script>

<template>
  <div class="canvas-wrapper">
    <StepBanner
      :mensaje="pasoMensaje"
      :visible="animando"
      :paso-actual="pasoActual"
      :paso-total="pasoTotal"
      :mostrar-siguiente="modoPasoAPaso"
      @siguiente="emit('siguiente')"
    />

    <div v-if="estructura" class="zoom-toolbar">
      <button type="button" class="zoom-btn" title="Alejar" @click="zoomOut()">−</button>
      <span class="zoom-label" :title="`Zoom ${zoomPorcentaje}% (máx. ${MAX_ZOOM * 100}%)`">
        {{ zoomEtiqueta }}
      </span>
      <button type="button" class="zoom-btn" title="Acercar" @click="zoomIn()">+</button>
      <button type="button" class="zoom-btn" title="Acercar rápido (×2)" @click="zoomInRapido">++</button>
      <span class="zoom-sep" />
      <button type="button" class="zoom-btn zoom-btn--text" title="Ver árbol completo" @click="verTodo">
        ⊞ Ver todo
      </button>
      <button type="button" class="zoom-btn zoom-btn--text" title="Restablecer zoom 100%" @click="resetZoom">
        100%
      </button>
      <span class="zoom-hint">Rueda = zoom · Arrastrar = mover · ++ = acercar rápido</span>
    </div>

    <div v-if="!estructura" class="canvas-empty">
      <p>Inicializa el árbol para ver la visualización</p>
    </div>

    <div
      v-else
      ref="viewportRef"
      class="canvas-viewport"
      @wheel="onWheel"
      @mousedown="onMouseDown"
      @mousemove="onMouseMove"
      @mouseup="onMouseUp"
      @mouseleave="onMouseUp"
    >
      <svg
        class="tree-canvas"
        :viewBox="viewBox"
        preserveAspectRatio="xMidYMid meet"
      >
        <defs>
          <marker id="arrow-down" markerWidth="7" markerHeight="7" refX="5" refY="3.5" orient="auto">
            <path d="M0,0 L0,7 L7,3.5 Z" fill="#4b5563" />
          </marker>
          <marker id="arrow-down-active" markerWidth="7" markerHeight="7" refX="5" refY="3.5" orient="auto">
            <path d="M0,0 L0,7 L7,3.5 Z" fill="#f59e0b" />
          </marker>
        <marker id="arrow-leaf" markerWidth="7" markerHeight="7" refX="6" refY="3.5" orient="auto">
          <path d="M0,0 L0,7 L7,3.5 Z" fill="#7c3aed" />
        </marker>
        <marker id="arrow-promote" markerWidth="7" markerHeight="7" refX="5" refY="3.5" orient="auto">
          <path d="M0,0 L0,7 L7,3.5 Z" fill="#f59e0b" />
        </marker>
      </defs>

        <g class="edges">
          <path
            v-for="edge in layout.edges"
            :key="edge.key"
            :d="edgePath(edge)"
            class="edge"
            :class="{ 'edge--active': activeEdge === edge.key }"
            :marker-end="activeEdge === edge.key ? 'url(#arrow-down-active)' : 'url(#arrow-down)'"
          />
        </g>

        <g class="leaf-chain">
          <path
            v-for="(link, i) in layout.leafChain"
            :key="'lc' + i"
            :d="leafChainPath(link.from, link.to)"
            class="edge edge--leaf-link"
            :class="{
              'edge--range-active': leafChainActive(link),
            }"
            marker-end="url(#arrow-leaf)"
          />
        </g>

        <text
          v-if="layout.leafChain.length"
          x="0"
          :y="layout.height - 4"
          class="chain-label"
          text-anchor="middle"
        >
          ← Sequence Set: hojas enlazadas horizontalmente →
        </text>

        <TreeNode
          v-for="nodo in nodos"
          :key="nodo.id"
          :nodo="nodo"
          :x="layout.posiciones[nodo.id]?.x ?? 0"
          :y="layout.posiciones[nodo.id]?.y ?? 0"
          :dimensiones="layout.dimensiones[nodo.id]"
          :highlighted="nodoResaltado(nodo.id)"
          :split-animating="splitNodes.includes(nodo.id)"
          :overflow-full="overflowNodeId === nodo.id"
          :highlighted-keys="highlightedKeys"
          :deleting-keys="deletingKeys"
          :active-slot="highlightedNodes.includes(nodo.id) ? activeSlot : null"
          :active-key="highlightedNodes.includes(nodo.id) ? activeKey : null"
          :atenuado="deleteAnim?.atenuarIds?.includes(nodo.id)"
        />

        <SplitOverlay :anim="splitAnim" />
        <DeleteOverlay :anim="deleteAnim" />
      </svg>
    </div>
  </div>
</template>

<style scoped>
.canvas-wrapper {
  position: relative;
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  min-height: 420px;
  overflow: hidden;
}

.zoom-toolbar {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.4rem 0.75rem;
  background: #e2e8f0;
  border-bottom: 1px solid #cbd5e1;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.zoom-btn {
  min-width: 32px;
  height: 28px;
  padding: 0 0.5rem;
  background: #fff;
  border: 1px solid #94a3b8;
  border-radius: 4px;
  color: #334155;
  font-size: 1rem;
  font-weight: 600;
  line-height: 1;
}

.zoom-btn:hover {
  background: #f1f5f9;
  border-color: #6366f1;
}

.zoom-btn--text {
  font-size: 0.75rem;
  font-weight: 500;
  min-width: auto;
}

.zoom-label {
  min-width: 44px;
  text-align: center;
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
}

.zoom-sep {
  width: 1px;
  height: 20px;
  background: #94a3b8;
  margin: 0 0.25rem;
}

.zoom-hint {
  margin-left: auto;
  font-size: 0.7rem;
  color: #64748b;
}

.canvas-viewport {
  flex: 1;
  overflow: hidden;
  cursor: grab;
  padding: 0.5rem;
  min-height: 360px;
}

.canvas-viewport--dragging {
  cursor: grabbing;
  user-select: none;
}

.canvas-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  color: #64748b;
}

.tree-canvas {
  width: 100%;
  height: 100%;
  min-height: 360px;
  display: block;
}

.edge {
  fill: none;
  stroke: #6b7280;
  stroke-width: 1.5;
  transition: stroke 0.25s, stroke-width 0.25s;
}

.edge--active {
  stroke: #f59e0b;
  stroke-width: 2.5;
}

.edge--leaf-link {
  stroke: #7c3aed;
  stroke-width: 2;
}

.edge--range-active {
  stroke: #f59e0b;
  stroke-width: 3;
}

.chain-label {
  fill: #64748b;
  font-size: 10px;
}
</style>
