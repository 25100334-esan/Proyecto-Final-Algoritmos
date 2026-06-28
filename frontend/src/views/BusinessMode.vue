<script setup>
import { ref, watch, onUnmounted } from 'vue'
import SongGridCard from '../components/business/SongGridCard.vue'
import FilterSidebar from '../components/business/FilterSidebar.vue'
import DebugConsole from '../components/business/DebugConsole.vue'
import SpotifyLogo from '../components/icons/SpotifyLogo.vue'
import { useTreeState } from '../composables/useTreeState.js'
import { useDebugLog } from '../composables/useDebugLog.js'
import * as api from '../api/bplustreeApi.js'
import { calcularRutaBusqueda } from '../utils/treeLayout.js'
import {
  mensajePrefijo,
  mensajeRangoPopularidad,
  mensajeRangoTempo,
  mensajeRangoNombre,
  mensajeBuscarId,
} from '../utils/perfLog.js'

const { estructura, inicializado, config } = useTreeState()
const { agregar } = useDebugLog()

const searchQuery = ref('')
const searchId = ref('')
const mostrarBusquedaId = ref(false)
const popularidadMin = ref(80)
const popularidadMax = ref(100)
const tempoMin = ref('120')
const tempoMax = ref('125')
const nombreDesde = ref('M')
const nombreHasta = ref('O')

const resultados = ref([])
const tituloResultados = ref('')
const cancionReproduciendo = ref(null)
const cargando = ref(false)
const error = ref(null)

const statsNombre = ref(null)
const statsPopularidad = ref(null)
const statsTempo = ref(null)

let debounceTimer = null
let abortController = null

async function cargarStats() {
  if (!inicializado.value) return
  const [n, p, t] = await Promise.all([
    api.obtenerEstadisticas('nombre').catch(() => null),
    api.obtenerEstadisticas('popularidad').catch(() => null),
    api.obtenerEstadisticas('tempo').catch(() => null),
  ])
  statsNombre.value = n
  statsPopularidad.value = p
  statsTempo.value = t
}

watch(inicializado, (ok) => {
  if (ok) cargarStats()
}, { immediate: true })

function limpiarDebounce() {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
    debounceTimer = null
  }
}

async function buscarPorPrefijo(prefijo) {
  if (!inicializado.value) {
    error.value = 'Primero inicializa el árbol en Modo Académico'
    return
  }

  if (abortController) abortController.abort()
  abortController = new AbortController()

  error.value = null
  cargando.value = true
  tituloResultados.value = ''
  const inicio = performance.now()

  try {
    const resultado = await api.buscarPorPrefijo(prefijo)
    const ms = performance.now() - inicio
    resultados.value = resultado.canciones ?? []
    tituloResultados.value =
      resultado.total > 0
        ? `Resultados para «${prefijo}»`
        : `Sin coincidencias para «${prefijo}»`

    if (prefijo.length >= 2) {
      agregar(
        mensajePrefijo({
          prefijo,
          total: resultado.total,
          ms,
          stats: statsNombre.value,
          orden: config.value.orden,
          totalCanciones: config.value.total,
        }),
        resultado.total > 0 ? 'success' : 'info',
      )
    }
  } catch (err) {
    if (err.name !== 'AbortError') {
      error.value = err.message
      resultados.value = []
    }
  } finally {
    cargando.value = false
  }
}

watch(searchQuery, (valor) => {
  limpiarDebounce()
  const prefijo = valor.trim()
  if (!prefijo) {
    resultados.value = []
    tituloResultados.value = ''
    error.value = null
    return
  }
  debounceTimer = setTimeout(() => buscarPorPrefijo(prefijo), 280)
})

async function onBuscarId() {
  if (!inicializado.value) return
  const indice = Number(searchId.value)
  if (Number.isNaN(indice)) {
    error.value = 'ID inválido'
    return
  }

  error.value = null
  cargando.value = true
  searchQuery.value = ''
  const inicio = performance.now()

  try {
    const cancion = await api.buscar(indice)
    const ms = performance.now() - inicio
    resultados.value = [cancion]
    tituloResultados.value = `Acceso directo ID ${indice}`
    cancionReproduciendo.value = cancion

    const nodos = estructura.value
      ? calcularRutaBusqueda(estructura.value, indice).length
      : Math.max(1, Math.ceil(Math.log2(config.value.total + 2)))

    agregar(mensajeBuscarId({ indice, ms, nodosVisitados: nodos }), 'success')
  } catch (err) {
    error.value = err.message
    resultados.value = []
  } finally {
    cargando.value = false
  }
}

async function onAplicarPopularidad() {
  if (!inicializado.value) return
  let min = popularidadMin.value
  let max = popularidadMax.value
  if (min > max) [min, max] = [max, min]

  error.value = null
  cargando.value = true
  searchQuery.value = ''
  const inicio = performance.now()

  try {
    const resultado = await api.buscarRangoCampo('popularidad', min, max)
    const ms = performance.now() - inicio
    resultados.value = resultado.canciones ?? []
    tituloResultados.value = `Top Hits Globales (${min}–${max})`

    agregar(
      mensajeRangoPopularidad({
        inicio: min,
        fin: max,
        total: resultado.total,
        ms,
        stats: statsPopularidad.value,
        orden: config.value.orden,
        totalCanciones: config.value.total,
      }),
      'success',
    )
  } catch (err) {
    error.value = err.message
  } finally {
    cargando.value = false
  }
}

async function onAplicarRangoNombre() {
  if (!inicializado.value) return
  const desde = nombreDesde.value.trim()
  const hasta = nombreHasta.value.trim()
  if (!desde || !hasta) {
    error.value = 'Indica desde y hasta para el rango alfabético'
    return
  }

  error.value = null
  cargando.value = true
  searchQuery.value = ''
  const inicio = performance.now()

  try {
    const resultado = await api.buscarRangoNombre(desde, hasta)
    const ms = performance.now() - inicio
    resultados.value = resultado.canciones ?? []
    tituloResultados.value = `Catálogo ${desde}–${hasta} (TrackName)`

    agregar(
      mensajeRangoNombre({
        desde,
        hasta,
        total: resultado.total,
        ms,
        stats: statsNombre.value,
        orden: config.value.orden,
        totalCanciones: config.value.total,
      }),
      'success',
    )
  } catch (err) {
    error.value = err.message
  } finally {
    cargando.value = false
  }
}

async function onAplicarTempo() {
  if (!inicializado.value) return
  const min = Number(tempoMin.value)
  const max = Number(tempoMax.value)
  if (Number.isNaN(min) || Number.isNaN(max)) {
    error.value = 'Tempo inválido'
    return
  }

  error.value = null
  cargando.value = true
  searchQuery.value = ''
  const inicio = performance.now()

  try {
    const resultado = await api.buscarRangoCampo('tempo', min, max)
    const ms = performance.now() - inicio
    resultados.value = resultado.canciones ?? []
    tituloResultados.value = `Playlist DJ · ${min}–${max} BPM`

    agregar(
      mensajeRangoTempo({
        inicio: min,
        fin: max,
        total: resultado.total,
        ms,
        stats: statsTempo.value,
        orden: config.value.orden,
        totalCanciones: config.value.total,
      }),
      'success',
    )
  } catch (err) {
    error.value = err.message
  } finally {
    cargando.value = false
  }
}

function onPlay(cancion) {
  cancionReproduciendo.value = cancion
  agregar(
    `▶ Reproducción simulada: «${cancion.TrackName}» — puntero Registro en RAM (0 accesos BD).`,
    'info',
  )
}

onUnmounted(() => {
  limpiarDebounce()
  if (abortController) abortController.abort()
})
</script>

<template>
  <div class="business-mode">
    <div class="business-layout">
      <FilterSidebar
        v-model:popularidad-min="popularidadMin"
        v-model:popularidad-max="popularidadMax"
        v-model:tempo-min="tempoMin"
        v-model:tempo-max="tempoMax"
        v-model:nombre-desde="nombreDesde"
        v-model:nombre-hasta="nombreHasta"
        :cargando="cargando"
        :deshabilitado="!inicializado"
        @aplicar-popularidad="onAplicarPopularidad"
        @aplicar-tempo="onAplicarTempo"
        @aplicar-rango-nombre="onAplicarRangoNombre"
      />

      <div class="main-panel">
        <div class="hero-search">
          <div class="search-bar">
            <SpotifyLogo class="search-icon" :size="22" color="#1db954" />
            <input
              v-model="searchQuery"
              type="search"
              placeholder="¿Qué quieres escuchar? Ej: Shape of..."
              autocomplete="off"
              :disabled="!inicializado"
            />
            <span v-if="cargando" class="search-spinner">…</span>
          </div>
          <p class="search-caption">
            Prefix B+-Tree sobre <code>arbolPorNombre</code> — autocompletado en tiempo real
          </p>
        </div>

        <button
          class="id-toggle"
          type="button"
          :title="mostrarBusquedaId ? 'Ocultar búsqueda por ID' : 'Buscar por ID exacto (profesor)'"
          @click="mostrarBusquedaId = !mostrarBusquedaId"
        >
          {{ mostrarBusquedaId ? '✕' : 'ID' }}
        </button>

        <div v-if="mostrarBusquedaId" class="id-search">
          <input
            v-model="searchId"
            type="number"
            placeholder="ID exacto (árbol principal)…"
            @keyup.enter="onBuscarId"
          />
          <button :disabled="cargando || !inicializado" @click="onBuscarId">Buscar ID</button>
        </div>

        <div v-if="!inicializado" class="notice">
          Ve al <strong>Modo Académico</strong> e inicializa el B+-Tree para usar esta vista.
        </div>

        <p v-if="error" class="error-msg">{{ error }}</p>

        <section v-if="tituloResultados || resultados.length" class="results-zone">
          <header class="results-header">
            <h2>{{ tituloResultados }}</h2>
            <span class="results-count">{{ resultados.length }} canción(es)</span>
          </header>

          <div v-if="resultados.length" class="results-grid">
            <SongGridCard
              v-for="cancion in resultados"
              :key="cancion.Indice"
              :cancion="cancion"
              :playing="cancionReproduciendo?.Indice === cancion.Indice"
              @play="onPlay"
            />
          </div>

          <p v-else-if="!cargando" class="empty-results">
            Escribe más caracteres o ajusta los filtros del panel lateral.
          </p>
        </section>

        <div
          v-else-if="inicializado && !cargando"
          class="empty-state"
        >
          <h2>Explora tu catálogo indexado</h2>
          <p>
            Escribe en la barra central para autocompletar por nombre, o usa los filtros DJ
            para range scans sobre popularidad y tempo.
          </p>
        </div>
      </div>
    </div>

    <DebugConsole />
  </div>
</template>

<style scoped>
.business-mode {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-dark);
}

.business-layout {
  flex: 1;
  display: flex;
  min-height: 0;
  overflow: hidden;
}

.main-panel {
  flex: 1;
  overflow-y: auto;
  padding: 2rem 2.5rem 1.5rem;
  position: relative;
}

.hero-search {
  max-width: 680px;
  margin: 0 auto 2rem;
  text-align: center;
}

.search-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: 500px;
  padding: 0.85rem 1.5rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.35);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-bar:focus-within {
  border-color: var(--accent);
  box-shadow: 0 8px 32px rgba(29, 185, 84, 0.12);
}

.search-bar input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 1.05rem;
  color: var(--text-primary);
}

.search-bar input:focus {
  outline: none;
}

.search-bar input::placeholder {
  color: var(--text-muted);
}

.search-icon {
  flex-shrink: 0;
}

.search-spinner {
  color: var(--accent);
  font-weight: 700;
  animation: pulse 0.8s infinite;
}

@keyframes pulse {
  50% { opacity: 0.3; }
}

.search-caption {
  margin-top: 0.65rem;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.search-caption code {
  color: #a5b4fc;
  font-size: 0.72rem;
}

.id-toggle {
  position: absolute;
  top: 1rem;
  right: 1.25rem;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: var(--bg-card);
  color: var(--text-muted);
  font-size: 0.7rem;
  font-weight: 700;
  border: 1px solid var(--border);
  opacity: 0.6;
}

.id-toggle:hover {
  opacity: 1;
  color: var(--text-primary);
}

.id-search {
  position: absolute;
  top: 3rem;
  right: 1.25rem;
  display: flex;
  gap: 0.35rem;
  background: var(--bg-card);
  padding: 0.5rem;
  border-radius: 8px;
  border: 1px solid var(--border);
  z-index: 2;
}

.id-search input {
  width: 140px;
  padding: 0.35rem 0.5rem;
  background: var(--bg-elevated);
  border-radius: 4px;
  font-size: 0.8rem;
}

.id-search button {
  padding: 0.35rem 0.65rem;
  background: var(--accent);
  color: #000;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
}

.notice {
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid var(--warning);
  color: var(--warning);
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.88rem;
  max-width: 680px;
  margin-left: auto;
  margin-right: auto;
}

.error-msg {
  color: #f87171;
  text-align: center;
  margin-bottom: 1rem;
}

.results-zone {
  margin-top: 0.5rem;
}

.results-header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1.25rem;
  flex-wrap: wrap;
}

.results-header h2 {
  font-size: 1.5rem;
  font-weight: 800;
}

.results-count {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 1.25rem;
}

.empty-results,
.empty-state {
  text-align: center;
  color: var(--text-secondary);
  padding: 3rem 1rem;
}

.empty-state h2 {
  font-size: 1.75rem;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

@media (max-width: 768px) {
  .business-layout {
    flex-direction: column;
  }

  .main-panel {
    padding: 1.25rem 1rem;
  }

  .results-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 0.85rem;
  }
}
</style>
