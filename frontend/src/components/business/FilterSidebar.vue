<script setup>
import { computed } from 'vue'

const props = defineProps({
  popularidadMin: { type: Number, required: true },
  popularidadMax: { type: Number, required: true },
  tempoMin: { type: [String, Number], required: true },
  tempoMax: { type: [String, Number], required: true },
  nombreDesde: { type: String, required: true },
  nombreHasta: { type: String, required: true },
  cargando: { type: Boolean, default: false },
  deshabilitado: { type: Boolean, default: false },
})

const emit = defineEmits([
  'update:popularidadMin',
  'update:popularidadMax',
  'update:tempoMin',
  'update:tempoMax',
  'update:nombreDesde',
  'update:nombreHasta',
  'aplicarPopularidad',
  'aplicarTempo',
  'aplicarRangoNombre',
])

const rangoLo = computed(() => Math.min(props.popularidadMin, props.popularidadMax))
const rangoHi = computed(() => Math.max(props.popularidadMin, props.popularidadMax))

const sliderStyle = computed(() => ({
  '--lo': rangoLo.value,
  '--hi': rangoHi.value,
}))

function onMinChange(event) {
  let valor = Number(event.target.value)
  if (valor > props.popularidadMax) valor = props.popularidadMax
  emit('update:popularidadMin', valor)
}

function onMaxChange(event) {
  let valor = Number(event.target.value)
  if (valor < props.popularidadMin) valor = props.popularidadMin
  emit('update:popularidadMax', valor)
}
</script>

<template>
  <aside class="filter-sidebar">
    <h2 class="sidebar-title">Filtros para DJs y Playlists</h2>
    <p class="sidebar-sub">Range Scans sobre índices secundarios B+</p>

    <section class="filter-block">
      <h3>Rango alfabético (TrackName)</h3>
      <p class="filter-hint">arbolPorNombre · Range Scan string · Prefix B+-Tree</p>
      <div class="tempo-fields">
        <label class="tempo-field">
          <span class="tempo-label">Desde</span>
          <input
            type="text"
            :value="nombreDesde"
            placeholder="M"
            :disabled="deshabilitado"
            @input="$emit('update:nombreDesde', $event.target.value)"
          />
        </label>
        <label class="tempo-field">
          <span class="tempo-label">Hasta</span>
          <input
            type="text"
            :value="nombreHasta"
            placeholder="O"
            :disabled="deshabilitado"
            @input="$emit('update:nombreHasta', $event.target.value)"
          />
        </label>
      </div>
      <button
        class="btn-filter"
        :disabled="cargando || deshabilitado"
        @click="$emit('aplicarRangoNombre')"
      >
        Catálogo A–Z (rango)
      </button>
    </section>

    <section class="filter-block">
      <h3>Popularidad</h3>
      <p class="filter-hint">arbolPorPopularidad · Sequence Set</p>

      <div class="range-display">
        <span class="range-badge">{{ popularidadMin }}</span>
        <span class="range-sep">—</span>
        <span class="range-badge">{{ popularidadMax }}</span>
      </div>

      <div class="dual-slider" :style="sliderStyle">
        <div class="slider-rail">
          <div class="slider-fill" aria-hidden="true" />
          <input
            type="range"
            min="0"
            max="100"
            :value="popularidadMin"
            :disabled="deshabilitado"
            class="slider slider--min"
            aria-label="Popularidad mínima"
            @input="onMinChange"
          />
          <input
            type="range"
            min="0"
            max="100"
            :value="popularidadMax"
            :disabled="deshabilitado"
            class="slider slider--max"
            aria-label="Popularidad máxima"
            @input="onMaxChange"
          />
        </div>
      </div>

      <button
        class="btn-filter"
        :disabled="cargando || deshabilitado"
        @click="$emit('aplicarPopularidad')"
      >
        Top Hits Globales
      </button>
    </section>

    <section class="filter-block">
      <h3>Tempo (BPM)</h3>
      <p class="filter-hint">arbolPorTempo · Range Scan</p>

      <div class="tempo-fields">
        <label class="tempo-field">
          <span class="tempo-label">Mínimo</span>
          <input
            type="number"
            step="0.1"
            min="0"
            :value="tempoMin"
            :disabled="deshabilitado"
            @input="$emit('update:tempoMin', $event.target.value)"
          />
        </label>
        <label class="tempo-field">
          <span class="tempo-label">Máximo</span>
          <input
            type="number"
            step="0.1"
            min="0"
            :value="tempoMax"
            :disabled="deshabilitado"
            @input="$emit('update:tempoMax', $event.target.value)"
          />
        </label>
      </div>

      <button
        class="btn-filter"
        :disabled="cargando || deshabilitado"
        @click="$emit('aplicarTempo')"
      >
        Playlist DJ
      </button>
    </section>

    <div class="sidebar-foot">
      <p>Las hojas enlazadas permiten recorrer el rango en O(k) tras el primer O(log n).</p>
    </div>
  </aside>
</template>

<style scoped>
.filter-sidebar {
  width: 288px;
  flex-shrink: 0;
  background: var(--bg-elevated);
  border-right: 1px solid var(--border);
  padding: 1.25rem 1.1rem;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-title {
  font-size: 0.95rem;
  font-weight: 700;
  margin-bottom: 0.25rem;
  line-height: 1.3;
}

.sidebar-sub {
  font-size: 0.72rem;
  color: var(--text-muted);
  margin-bottom: 1.25rem;
  line-height: 1.4;
}

.filter-block {
  background: var(--bg-card);
  border-radius: 10px;
  padding: 1rem 0.9rem;
  margin-bottom: 1rem;
  overflow: hidden;
}

.filter-block h3 {
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 0.2rem;
}

.filter-hint {
  font-size: 0.65rem;
  color: var(--text-muted);
  margin-bottom: 0.85rem;
  font-family: 'Consolas', 'Courier New', monospace;
  line-height: 1.35;
  word-break: break-word;
}

.range-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 0.65rem;
}

.range-badge {
  min-width: 2.25rem;
  text-align: center;
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--accent);
  background: rgba(29, 185, 84, 0.12);
  padding: 0.2rem 0.55rem;
  border-radius: 6px;
}

.range-sep {
  color: var(--text-muted);
  font-size: 0.75rem;
}

.dual-slider {
  --track-h: 6px;
  --thumb-h: 16px;
  --rail-h: 24px;
  --track-top: calc((var(--rail-h) - var(--track-h)) / 2);
  --thumb-offset: calc((var(--track-h) - var(--thumb-h)) / 2);
  margin-bottom: 1rem;
  padding: 0 2px;
}

.slider-rail {
  position: relative;
  width: 100%;
  height: var(--rail-h);
}

.slider-rail::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: var(--track-top);
  height: var(--track-h);
  background: #3a3a3a;
  border-radius: 3px;
  pointer-events: none;
  z-index: 0;
}

.slider-fill {
  position: absolute;
  top: var(--track-top);
  height: var(--track-h);
  left: calc(var(--lo) * 1%);
  width: calc((var(--hi) - var(--lo)) * 1%);
  background: var(--accent);
  border-radius: 3px;
  pointer-events: none;
  z-index: 1;
}

.slider {
  position: absolute;
  width: 100%;
  height: var(--rail-h);
  top: 0;
  left: 0;
  margin: 0;
  padding: 0;
  -webkit-appearance: none;
  appearance: none;
  background: transparent;
  pointer-events: none;
}

.slider::-webkit-slider-runnable-track {
  height: var(--track-h);
  background: transparent;
  border: none;
}

.slider::-moz-range-track {
  height: var(--track-h);
  background: transparent;
  border: none;
}

.slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: var(--thumb-h);
  height: var(--thumb-h);
  border-radius: 50%;
  background: var(--accent);
  border: 2px solid #121212;
  cursor: pointer;
  pointer-events: auto;
  margin-top: var(--thumb-offset);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}

.slider::-moz-range-thumb {
  width: var(--thumb-h);
  height: var(--thumb-h);
  border-radius: 50%;
  background: var(--accent);
  border: 2px solid #121212;
  cursor: pointer;
  pointer-events: auto;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
  box-sizing: border-box;
}

.slider--min {
  z-index: 3;
}

.slider--max {
  z-index: 4;
}

.tempo-fields {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 0.65rem;
  margin-bottom: 1rem;
}

.tempo-field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  min-width: 0;
}

.tempo-label {
  font-size: 0.72rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.tempo-field input {
  width: 100%;
  min-width: 0;
  box-sizing: border-box;
  padding: 0.45rem 0.5rem;
  background: var(--bg-elevated);
  border-radius: 6px;
  font-size: 0.85rem;
  border: 1px solid var(--border);
  color: var(--text-primary);
}

.tempo-field input:focus {
  outline: none;
  border-color: var(--accent);
}

.btn-filter {
  width: 100%;
  padding: 0.55rem 0.75rem;
  border-radius: 500px;
  background: var(--accent);
  color: #000;
  font-weight: 600;
  font-size: 0.8rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-filter:hover:not(:disabled) {
  background: var(--accent-hover);
}

.btn-filter:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.sidebar-foot {
  font-size: 0.68rem;
  color: var(--text-muted);
  line-height: 1.45;
  padding-top: 0.25rem;
}

@media (max-width: 768px) {
  .filter-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--border);
  }
}
</style>
