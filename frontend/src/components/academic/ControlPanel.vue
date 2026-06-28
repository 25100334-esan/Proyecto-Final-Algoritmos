<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { verificarBackend, obtenerConteoBD } from '../../api/bplustreeApi.js'
import { LIMITE_DEMO_INICIAL } from '../../utils/initAnimation.js'
import { LIMITE_VISUALIZACION } from '../../utils/visualizationLimits.js'
import { plantillaCancionSimulada } from '../../utils/cancionPlantilla.js'

const emit = defineEmits(['inicializar', 'buscar', 'insertar', 'eliminar', 'rango', 'cambiar-vista', 'cargar-toda-bd'])

const props = defineProps({
  cargando: Boolean,
  inicializado: Boolean,
  proximoIndice: { type: Number, default: null },
  deshabilitarAcciones: { type: Boolean, default: false },
  tipoVistaArbol: { type: String, default: 'indice' },
  cambiandoVista: { type: Boolean, default: false },
  visualizacionHabilitada: { type: Boolean, default: true },
  totalCanciones: { type: Number, default: 0 },
  limiteVisualizacion: { type: Number, default: LIMITE_VISUALIZACION },
})

const tipoVistaLocal = ref('indice')

watch(
  () => props.tipoVistaArbol,
  (v) => { tipoVistaLocal.value = v || 'indice' },
  { immediate: true },
)

function onCambiarVista() {
  emit('cambiar-vista', tipoVistaLocal.value)
}

const opcionesVista = [
  { id: 'indice', label: 'Indice (id) — int' },
  { id: 'nombre', label: 'TrackName — string' },
  { id: 'popularidad', label: 'Popularidad — int' },
  { id: 'tempo', label: 'Tempo — float64' },
  { id: 'danceability', label: 'Danceability — float64' },
]
const orden = ref(2)
const limite = ref(100)
const backendOk = ref(null)
const totalBD = ref(null)

const demoHabilitado = computed(
  () => limite.value >= 1 && limite.value <= LIMITE_DEMO_INICIAL,
)

const avisoSinCanvas = computed(
  () => limite.value > LIMITE_VISUALIZACION,
)

onMounted(async () => {
  try {
    await verificarBackend()
    backendOk.value = true
    const data = await obtenerConteoBD()
    totalBD.value = data.total
  } catch {
    backendOk.value = false
  }
})

async function onInicializarRapido() {
  if (backendOk.value === false) {
    try {
      await verificarBackend()
      backendOk.value = true
    } catch {
      emit('inicializar', {
        orden: Number(orden.value),
        limite: Number(limite.value),
        modo: 'rapido',
      })
      return
    }
  }
  emit('inicializar', {
    orden: Number(orden.value),
    limite: Number(limite.value),
    modo: 'rapido',
  })
}

function emitirDemo(modo) {
  emit('inicializar', {
    orden: Number(orden.value),
    limite: Number(limite.value),
    modo,
  })
}

function onInicializarCompleto() {
  emitirDemo('completo')
}

function onInicializarPasoAPaso() {
  emitirDemo('pasoAPaso')
}

function onCargarTodaBD() {
  emit('cargar-toda-bd', { orden: Number(orden.value) })
}

const buscarId = ref('')
const eliminarId = ref('')
const rangoInicio = ref('')
const rangoFin = ref('')

const esVistaNombre = computed(() => tipoVistaLocal.value === 'nombre')

const configRango = computed(() => {
  switch (tipoVistaLocal.value) {
    case 'nombre':
      return {
        input: 'text',
        step: null,
        hint: 'Range Scan lexicográfico (TrackName). Letras M–O incluye M*, N*, O*.',
        phInicio: 'Desde (ej: M)',
        phFin: 'Hasta (ej: O)',
        campo: 'nombre',
        tipoClave: 'string',
      }
    case 'popularidad':
      return {
        input: 'number',
        step: 1,
        hint: 'Range Scan int — arbolPorPopularidad (sequence set).',
        phInicio: 'Desde (ej: 80)',
        phFin: 'Hasta (ej: 100)',
        campo: 'popularidad',
        tipoClave: 'int',
      }
    case 'tempo':
      return {
        input: 'number',
        step: 0.1,
        hint: 'Range Scan float64 — arbolPorTempo (BPM).',
        phInicio: 'Desde BPM',
        phFin: 'Hasta BPM',
        campo: 'tempo',
        tipoClave: 'float',
      }
    case 'danceability':
      return {
        input: 'number',
        step: 0.01,
        hint: 'Range Scan float64 — arbolPorDanceability (ej: 0.5–0.8).',
        phInicio: 'Desde (0.5)',
        phFin: 'Hasta (0.8)',
        campo: 'danceability',
        tipoClave: 'float',
      }
    default:
      return {
        input: 'number',
        step: 1,
        hint: 'Range Scan por ID — índice principal (int).',
        phInicio: 'Inicio ID',
        phFin: 'Fin ID',
        campo: 'indice',
        tipoClave: 'int',
      }
  }
})
const mostrarDetallesInsertar = ref(false)

const nuevaCancion = ref(plantillaCancionSimulada(1))

function aplicarPlantillaInsertar() {
  if (props.proximoIndice != null) {
    nuevaCancion.value = plantillaCancionSimulada(props.proximoIndice)
  }
}

watch(
  () => props.proximoIndice,
  (id) => {
    if (id != null) aplicarPlantillaInsertar()
  },
  { immediate: true },
)

function toggleDetallesInsertar() {
  mostrarDetallesInsertar.value = !mostrarDetallesInsertar.value
  if (mostrarDetallesInsertar.value) aplicarPlantillaInsertar()
}

function onBuscarCompleto() {
  emit('buscar', { indice: Number(buscarId.value), modo: 'completo' })
}

function onBuscarPasoAPaso() {
  emit('buscar', { indice: Number(buscarId.value), modo: 'pasoAPaso' })
}

function onInsertarCompleto() {
  emit('insertar', { cancion: { ...nuevaCancion.value }, modo: 'completo' })
}

function onInsertarPasoAPaso() {
  emit('insertar', { cancion: { ...nuevaCancion.value }, modo: 'pasoAPaso' })
}

function onEliminarCompleto() {
  emit('eliminar', { indice: Number(eliminarId.value), modo: 'completo' })
}

function onEliminarPasoAPaso() {
  emit('eliminar', { indice: Number(eliminarId.value), modo: 'pasoAPaso' })
}

function onRangoCompleto() {
  const cfg = configRango.value
  const inicio = cfg.input === 'text' ? rangoInicio.value.trim() : Number(rangoInicio.value)
  const fin = cfg.input === 'text' ? rangoFin.value.trim() : Number(rangoFin.value)
  emit('rango', { inicio, fin, modo: 'completo', campo: cfg.campo, tipoClave: cfg.tipoClave })
}

function onRangoPasoAPaso() {
  const cfg = configRango.value
  const inicio = cfg.input === 'text' ? rangoInicio.value.trim() : Number(rangoInicio.value)
  const fin = cfg.input === 'text' ? rangoFin.value.trim() : Number(rangoFin.value)
  emit('rango', { inicio, fin, modo: 'pasoAPaso', campo: cfg.campo, tipoClave: cfg.tipoClave })
}
</script>

<template>
  <aside class="control-panel">
    <section class="panel-section">
      <h3>Configuración Inicial</h3>
      <label>
        Orden d (Comer)
        <input v-model.number="orden" type="number" min="2" max="10" :disabled="cargando" />
        <span class="hint">d={{ orden }} → máx. {{ orden * 2 }} claves/nodo, split al insertar la {{ orden * 2 + 1 }}ª</span>
      </label>
      <label>
        Límite de canciones
        <input v-model.number="limite" type="number" min="1" max="5000" :disabled="cargando" />
        <span v-if="avisoSinCanvas" class="hint hint--warn">
          Con más de {{ LIMITE_VISUALIZACION }} canciones el canvas se desactiva (solo operaciones API).
        </span>
      </label>

      <template v-if="demoHabilitado">
        <button
          class="btn btn--primary"
          :disabled="cargando || deshabilitarAcciones"
          @click="onInicializarCompleto"
        >
          {{ cargando ? 'Cargando…' : 'Inicializar (secuencia completa)' }}
        </button>
        <button
          class="btn btn--primary-outline"
          :disabled="cargando || deshabilitarAcciones"
          @click="onInicializarPasoAPaso"
        >
          Inicializar (paso a paso)
        </button>
        <button
          class="btn btn--ghost"
          :disabled="cargando || deshabilitarAcciones"
          @click="onInicializarRapido"
        >
          Inicializar (rápido, sin animación)
        </button>
        <p class="hint hint--info">
          Demo animada: ≤ {{ LIMITE_DEMO_INICIAL }} canc. · «Rápido» = POST /api/inicializar
        </p>
      </template>

      <template v-else>
        <button class="btn btn--primary" :disabled="cargando" @click="onInicializarRapido">
          {{ cargando ? 'Cargando…' : 'Inicializar B+-Tree' }}
        </button>
        <p class="hint hint--info">
          Demo animada solo con ≤ {{ LIMITE_DEMO_INICIAL }} canciones.
        </p>
      </template>

      <button
        class="btn btn--bd"
        :disabled="cargando || deshabilitarAcciones"
        @click="onCargarTodaBD"
      >
        {{ cargando ? 'Cargando…' : `Cargar toda la BD${totalBD != null ? ` (${totalBD})` : ''}` }}
      </button>
      <p class="hint hint--info">
        Carga todas las filas de <code>tracks</code>. Sin canvas si superan {{ LIMITE_VISUALIZACION }}.
      </p>

      <p v-if="backendOk === false" class="status status--error">
        Backend Go apagado. En otra terminal: <code>go run .</code>
      </p>
      <p v-else-if="backendOk === true" class="status status--backend">Backend conectado (puerto 8080)</p>
      <p v-if="inicializado" class="status status--ok">Árbol activo en RAM</p>
    </section>

    <section v-if="inicializado" class="panel-section panel-section--vista">
      <h3>Vista del árbol</h3>
      <label class="label-vista">
        Seleccionar índice
        <select
          v-model="tipoVistaLocal"
          :disabled="deshabilitarAcciones || cambiandoVista"
          @change="onCambiarVista"
        >
          <option v-for="op in opcionesVista" :key="op.id" :value="op.id">
            {{ op.label }}
          </option>
        </select>
      </label>
      <p v-if="visualizacionHabilitada" class="hint hint--info">
        Animaciones usan la clave del índice seleccionado (genérico K).
      </p>
      <p v-else class="hint hint--warn">
        Sin canvas — el selector cambia las estadísticas del índice en RAM.
      </p>
    </section>

    <section class="panel-section" :class="{ 'panel-section--disabled': !inicializado || deshabilitarAcciones }">
      <h3>Buscar ID</h3>
      <div class="row">
        <input v-model="buscarId" type="number" placeholder="ID" :disabled="!inicializado" />
      </div>
      <button class="btn btn--search" :disabled="!inicializado" @click="onBuscarCompleto">
        Buscar (secuencia completa)
      </button>
      <button class="btn btn--search-outline" :disabled="!inicializado" @click="onBuscarPasoAPaso">
        Buscar (paso a paso)
      </button>
      <p class="hint hint--delete">Paso a paso: usa «Siguiente» en el canvas para avanzar.</p>
    </section>

    <section class="panel-section" :class="{ 'panel-section--disabled': !inicializado || deshabilitarAcciones }">
      <h3>Insertar Canción (Simulada)</h3>
      <p v-if="proximoIndice != null" class="proximo-id">
        Próximo ID automático: <strong>{{ proximoIndice }}</strong>
      </p>
      <label>
        Nombre
        <input v-model="nuevaCancion.TrackName" :disabled="!inicializado" />
      </label>
      <label>
        Artista
        <input v-model="nuevaCancion.Artists" :disabled="!inicializado" />
      </label>
      <button
        type="button"
        class="btn btn--details"
        :disabled="!inicializado"
        @click="toggleDetallesInsertar"
      >
        {{ mostrarDetallesInsertar ? '▲ Ocultar detalles' : '▼ Detalles' }}
      </button>
      <div v-if="mostrarDetallesInsertar" class="detalles-insertar">
        <button
          type="button"
          class="btn btn--details-outline"
          :disabled="!inicializado || proximoIndice == null"
          @click="aplicarPlantillaInsertar"
        >
          Autocompletar con id {{ proximoIndice }}
        </button>
        <label>
          Track ID (Spotify)
          <input v-model="nuevaCancion.TrackID" :disabled="!inicializado" placeholder="sim_…" />
        </label>
        <label>
          Álbum
          <input v-model="nuevaCancion.AlbumName" :disabled="!inicializado" />
        </label>
        <label>
          Popularidad (0–100)
          <input v-model.number="nuevaCancion.Popularity" type="number" min="0" max="100" :disabled="!inicializado" />
        </label>
        <label>
          Duración (ms)
          <input v-model.number="nuevaCancion.DurationMs" type="number" min="0" :disabled="!inicializado" />
        </label>
        <label class="label--inline">
          <input v-model="nuevaCancion.Explicit" type="checkbox" :disabled="!inicializado" />
          Contenido explícito
        </label>
        <label>
          Género
          <input v-model="nuevaCancion.TrackGenre" :disabled="!inicializado" />
        </label>
        <div class="detalles-grid">
          <label>
            Danceability
            <input v-model.number="nuevaCancion.Danceability" type="number" step="0.01" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Energy
            <input v-model.number="nuevaCancion.Energy" type="number" step="0.01" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Key
            <input v-model.number="nuevaCancion.Key" type="number" min="0" max="11" :disabled="!inicializado" />
          </label>
          <label>
            Loudness
            <input v-model.number="nuevaCancion.Loudness" type="number" step="0.1" :disabled="!inicializado" />
          </label>
          <label>
            Mode
            <input v-model.number="nuevaCancion.Mode" type="number" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Speechiness
            <input v-model.number="nuevaCancion.Speechiness" type="number" step="0.001" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Acousticness
            <input v-model.number="nuevaCancion.Acousticness" type="number" step="0.01" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Instrumentalness
            <input v-model.number="nuevaCancion.Instrumentalness" type="number" step="0.001" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Liveness
            <input v-model.number="nuevaCancion.Liveness" type="number" step="0.01" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Valence
            <input v-model.number="nuevaCancion.Valence" type="number" step="0.01" min="0" max="1" :disabled="!inicializado" />
          </label>
          <label>
            Tempo (BPM)
            <input v-model.number="nuevaCancion.Tempo" type="number" step="0.1" min="0" :disabled="!inicializado" />
          </label>
          <label>
            Compás
            <input v-model.number="nuevaCancion.TimeSignature" type="number" min="1" max="7" :disabled="!inicializado" />
          </label>
        </div>
      </div>
      <button class="btn btn--accent" :disabled="!inicializado" @click="onInsertarCompleto">
        Insertar (secuencia completa)
      </button>
      <button class="btn btn--accent-outline" :disabled="!inicializado" @click="onInsertarPasoAPaso">
        Insertar (paso a paso)
      </button>
      <p class="hint hint--delete">Detalles: campos extra autocompletados según el próximo id.</p>
    </section>

    <section class="panel-section" :class="{ 'panel-section--disabled': !inicializado || deshabilitarAcciones }">
      <h3>Eliminar ID</h3>
      <p class="hint">
        Busca en arbolPorIndice → EliminarDeBD (SQLite) → cascada en los 5 árboles.
        En vista secundaria anima la clave de esa vista (ej. TrackName).
      </p>
      <div class="row">
        <input v-model="eliminarId" type="number" placeholder="ID" :disabled="!inicializado" />
      </div>
      <button class="btn btn--danger" :disabled="!inicializado" @click="onEliminarCompleto">
        Eliminar (secuencia completa)
      </button>
      <button class="btn btn--danger-outline" :disabled="!inicializado" @click="onEliminarPasoAPaso">
        Eliminar (paso a paso)
      </button>
      <p class="hint hint--delete">Paso a paso: usa «Siguiente» en el canvas para avanzar.</p>
    </section>

    <section class="panel-section" :class="{ 'panel-section--disabled': !inicializado || deshabilitarAcciones }">
      <h3>Buscar Rango</h3>
      <p class="hint">{{ configRango.hint }}</p>
      <div class="row">
        <input
          v-model="rangoInicio"
          :type="configRango.input"
          :step="configRango.step ?? undefined"
          :placeholder="configRango.phInicio"
          :disabled="!inicializado"
        />
        <input
          v-model="rangoFin"
          :type="configRango.input"
          :step="configRango.step ?? undefined"
          :placeholder="configRango.phFin"
          :disabled="!inicializado"
        />
      </div>
      <button class="btn btn--search" :disabled="!inicializado" @click="onRangoCompleto">
        Range Scan (secuencia completa)
      </button>
      <button class="btn btn--search-outline" :disabled="!inicializado" @click="onRangoPasoAPaso">
        Range Scan (paso a paso)
      </button>
      <p class="hint hint--delete">Paso a paso: usa «Siguiente» en el canvas para avanzar.</p>
    </section>
  </aside>
</template>

<style scoped>
.control-panel {
  width: 280px;
  flex-shrink: 0;
  background: var(--bg-panel);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  overflow-y: auto;
  max-height: calc(100vh - 140px);
}

.panel-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.panel-section--disabled {
  opacity: 0.45;
  pointer-events: none;
}

.panel-section h3 {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--accent-academic);
  border-bottom: 1px solid var(--border);
  padding-bottom: 0.35rem;
}

label {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.row {
  display: flex;
  gap: 0.5rem;
}

.row input {
  flex: 1;
}

.btn {
  padding: 0.5rem 0.75rem;
  background: var(--bg-card);
  color: var(--text-primary);
  font-size: 0.85rem;
}

.btn--primary {
  background: var(--accent-academic);
  color: white;
  font-weight: 600;
}

.btn--primary-outline {
  background: transparent;
  color: var(--accent-academic);
  border: 1px solid var(--accent-academic);
  font-size: 0.8rem;
}

.btn--bd {
  margin-top: 0.35rem;
  background: rgba(5, 150, 105, 0.18);
  border: 1px solid #059669;
  color: #34d399;
  font-size: 0.82rem;
}

.btn--bd:hover:not(:disabled) {
  background: rgba(5, 150, 105, 0.32);
}

.btn--search {
  background: var(--accent-academic);
  color: white;
  font-weight: 600;
}

.btn--search-outline {
  background: transparent;
  color: var(--accent-academic);
  border: 1px solid var(--accent-academic);
  font-size: 0.8rem;
}

.btn--accent {
  background: var(--leaf-node);
  color: white;
}

.btn--accent-outline {
  background: transparent;
  color: var(--leaf-node);
  border: 1px solid var(--leaf-node);
  font-size: 0.8rem;
}

.btn--danger {
  background: var(--danger);
  color: white;
}

.btn--danger-outline {
  background: transparent;
  color: var(--danger);
  border: 1px solid var(--danger);
  font-size: 0.8rem;
}

.btn--ghost {
  background: var(--bg-card);
  color: var(--text-secondary);
  border: 1px dashed var(--border);
  font-size: 0.8rem;
}

.btn--ghost:hover:not(:disabled) {
  border-color: var(--accent-academic);
  color: var(--text-primary);
}

.btn--details {
  background: rgba(99, 102, 241, 0.12);
  color: var(--accent-academic);
  font-size: 0.8rem;
  text-align: left;
}

.btn--details-outline {
  background: transparent;
  color: var(--accent-academic);
  border: 1px solid var(--border);
  font-size: 0.75rem;
}

.detalles-insertar {
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
  padding: 0.5rem;
  background: rgba(0, 0, 0, 0.15);
  border-radius: 6px;
  border: 1px solid var(--border);
  max-height: 280px;
  overflow-y: auto;
}

.detalles-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.4rem;
}

.detalles-grid label {
  font-size: 0.72rem;
}

.label--inline {
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.78rem;
}

.label--inline input[type='checkbox'] {
  width: auto;
}

.hint {
  margin: 0.25rem 0 0;
  font-size: 0.68rem;
  line-height: 1.35;
  color: var(--text-secondary);
  opacity: 0.85;
}

.hint code {
  font-size: 0.62rem;
  background: rgba(0, 0, 0, 0.2);
  padding: 0.05rem 0.25rem;
  border-radius: 2px;
}

.hint--info {
  color: var(--text-secondary);
}

.hint--delete {
  margin: 0;
}

.hint--warn {
  margin: 0;
  color: #d97706;
  font-size: 0.68rem;
}

.status {
  font-size: 0.75rem;
  padding: 0.35rem 0.5rem;
  border-radius: 4px;
}

.status--ok {
  background: rgba(5, 150, 105, 0.2);
  color: #34d399;
}

.status--backend {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
  font-size: 0.7rem;
}

.status--error {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
  font-size: 0.72rem;
  line-height: 1.4;
}

.status--error code {
  background: rgba(0, 0, 0, 0.2);
  padding: 0.1rem 0.35rem;
  border-radius: 3px;
}

.proximo-id {
  font-size: 0.8rem;
  color: var(--text-secondary);
  background: rgba(99, 102, 241, 0.15);
  padding: 0.4rem 0.6rem;
  border-radius: 4px;
}

.proximo-id strong {
  color: var(--accent-academic);
}

.label-vista select {
  margin-top: 0.25rem;
  font-size: 0.8rem;
  padding: 0.35rem;
  background: var(--bg-card);
  color: var(--text-primary);
  border-radius: 4px;
  border: 1px solid var(--border);
}
</style>
