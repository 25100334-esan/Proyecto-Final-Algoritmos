<script setup>
import { computed } from 'vue'
import { FASES_DELETE } from '../../utils/deleteAnimation.js'

const props = defineProps({
  anim: { type: Object, default: null },
})

const phase = computed(() => props.anim?.phase)
const progress = computed(() => props.anim?.progress ?? 0)

const activo = computed(() => {
  const p = phase.value
  return (
    p === FASES_DELETE.PRESTAMO ||
    p === FASES_DELETE.FUSION ||
    p === FASES_DELETE.FANTASMA ||
    p === FASES_DELETE.UNDERFLOW ||
    p === FASES_DELETE.EVALUAR ||
    p === FASES_DELETE.RAIZ
  )
})

const origen = computed(() => props.anim?.hermanoPos ?? { x: 0, y: 0 })
const destino = computed(() => props.anim?.hojaPos ?? { x: 0, y: 0 })
const padre = computed(() => props.anim?.padrePos ?? { x: 0, y: 0 })
const merged = computed(() => props.anim?.mergedPos ?? { x: 0, y: 0 })

const claveX = computed(
  () => origen.value.x + (destino.value.x - origen.value.x) * progress.value,
)
const claveY = computed(
  () => origen.value.y + (destino.value.y - origen.value.y) * progress.value,
)

const evalLabel = computed(() => {
  if (phase.value !== FASES_DELETE.EVALUAR) return ''
  const dir = props.anim?.evaluarHermano
  if (dir === 'derecha') return '🔎 Evaluando hermano DERECHO'
  if (dir === 'izquierda') return '🔎 Evaluando hermano IZQUIERDO'
  return '🔎 Evaluando hermano'
})
</script>

<template>
  <g v-if="activo" class="delete-annotations">
    <g v-if="phase === 'underflow'">
      <text
        :x="destino.x || merged.x"
        :y="(destino.y || merged.y) - 50"
        class="underflow-label"
        text-anchor="middle"
      >
        ⚠ UNDERFLOW — claves &lt; d (mínimo Comer)
      </text>
    </g>

    <g v-if="phase === 'evaluar'">
      <text
        :x="origen.x || destino.x"
        :y="(origen.y || destino.y) - 48"
        class="eval-label"
        text-anchor="middle"
      >
        {{ evalLabel }}
      </text>
    </g>

    <g v-if="phase === 'prestamo' && anim.clavePrestada != null">
      <line
        class="borrow-line"
        :x1="origen.x"
        :y1="origen.y"
        :x2="destino.x"
        :y2="destino.y"
      />
      <g :transform="`translate(${claveX}, ${claveY})`">
        <rect class="borrow-key" x="-34" y="-14" width="68" height="28" rx="4" />
        <text class="borrow-key__text" y="5">{{ anim.clavePrestada }}</text>
      </g>
      <text
        :x="(origen.x + destino.x) / 2"
        :y="(origen.y + destino.y) / 2 - 14"
        class="borrow-label"
        text-anchor="middle"
      >
        ↔ prestarCancion{{ anim.direccion === 'derecha' ? 'Derecha' : 'Izquierda' }}()
      </text>
      <text
        v-if="anim.nuevoSeparador != null && progress > 0.75"
        :x="padre.x"
        :y="padre.y - 44"
        class="sep-label"
        text-anchor="middle"
      >
        Separador padre → {{ anim.nuevoSeparador }}
      </text>
    </g>

    <g v-if="phase === 'fusion'">
      <text
        :x="merged.x"
        :y="merged.y - (anim.mergedDim?.height ?? 54) / 2 - 14"
        class="merge-label"
        text-anchor="middle"
      >
        🔗 fusionarHoja — concatenación + re-enlace sequence set
      </text>
    </g>

    <g v-if="phase === 'fantasma' && anim.separadorFantasma != null">
      <text
        :x="padre.x"
        :y="padre.y - (anim.padreDim?.height ?? 54) / 2 - 10"
        class="ghost-label"
        text-anchor="middle"
      >
        👻 Separador {{ anim.separadorFantasma }} permanece (fantasma)
      </text>
    </g>

    <g v-if="phase === 'raiz'">
      <text
        :x="padre.x || merged.x"
        :y="(padre.y || merged.y) - 58"
        class="raiz-label"
        text-anchor="middle"
      >
        ⬇ ajustarRaiz() — altura del árbol −1
      </text>
    </g>
  </g>
</template>

<style scoped>
.underflow-label {
  fill: #dc2626;
  font-size: 9px;
  font-weight: 700;
}

.eval-label {
  fill: #d97706;
  font-size: 9px;
  font-weight: 700;
}

.borrow-line {
  fill: none;
  stroke: #3b82f6;
  stroke-width: 2.5;
  stroke-dasharray: 6 4;
}

.borrow-key {
  fill: #dbeafe;
  stroke: #2563eb;
  stroke-width: 2;
  filter: drop-shadow(0 0 8px rgba(59, 130, 246, 0.6));
}

.borrow-key__text {
  fill: #1e40af;
  font-size: 11px;
  font-weight: 800;
  text-anchor: middle;
}

.borrow-label {
  fill: #2563eb;
  font-size: 9px;
  font-weight: 700;
}

.sep-label {
  fill: #6366f1;
  font-size: 8px;
  font-weight: 700;
}

.merge-label {
  fill: #dc2626;
  font-size: 9px;
  font-weight: 700;
}

.ghost-label {
  fill: #7c3aed;
  font-size: 9px;
  font-weight: 700;
}

.raiz-label {
  fill: #059669;
  font-size: 9px;
  font-weight: 700;
}
</style>
