<script setup>
import { ref, computed, nextTick, onMounted } from 'vue'
import ControlPanel from '../components/academic/ControlPanel.vue'
import TreeCanvas from '../components/academic/TreeCanvas.vue'
import TreeStatsPanel from '../components/academic/TreeStatsPanel.vue'
import OperationLog from '../components/academic/OperationLog.vue'
import { claveDeCancionPorVista, claveDeCancionPorEstructura, etiquetaClaveVista, normalizarClaveArbol, tipoVistaDesdeEstructura } from '../utils/claveVista.js'
import { useTreeState } from '../composables/useTreeState.js'
import { useDebugLog } from '../composables/useDebugLog.js'
import * as api from '../api/bplustreeApi.js'
import { calcularLayout } from '../utils/treeLayout.js'
import { calcularSecuenciaBusqueda, calcularPasosConfirmacionBusqueda } from '../utils/searchAnimation.js'
import { calcularSecuenciaRangeScan, calcularPasosConfirmacionRango } from '../utils/rangeAnimation.js'
import { FASES_SPLIT } from '../utils/splitAnimation.js'
import {
  calcularSecuenciaInsercion,
  calcularPasosConfirmacionInsercion,
  estructuraConInsercionVirtual,
} from '../utils/insertAnimation.js'
import {
  analizarEliminacion,
  calcularSecuenciaEliminacion,
  calcularPasosConfirmacion,
  estructuraConBorradoVirtual,
  FASES_DELETE,
} from '../utils/deleteAnimation.js'
import {
  LIMITE_DEMO_INICIAL,
  idDeCancion,
  calcularPasosPreLecturaBD,
  calcularPasosIntroInicializacion,
  calcularPasoRegistroBD,
  calcularPasoEncabezadoCancion,
  calcularPasosFinInicializacion,
} from '../utils/initAnimation.js'

const { estructura, estructuraIndice, estadisticas, ultimaOperacion, tipoVistaArbol, visualizacionHabilitada, limiteVisualizacion, inicializado, cargando, config, refrescarEstructura, refrescarEstadisticas, cambiarVistaArbol, inicializarArbol, prepararDemoInicial, finalizarDemoInicial, ajustarTotal, registrarOperacion, setCargando } = useTreeState()
const { agregar } = useDebugLog()

const highlightedNodes = ref([])
const splitNodes = ref([])
const underflowNodes = ref([])
const highlightedKeys = ref([])
const deletingKeys = ref([])
const rangeHighlightLeaf = ref(null)
const activeLeafChain = ref(null)
const logAcademico = ref([])
const estructuraAnterior = ref(null)
const estructuraAnimacion = ref(null)
const proximoIndice = ref(null)

const modoPasoAPaso = ref(false)
const pasoActual = ref(0)
const pasoTotal = ref(0)
let resolverSiguiente = null

const activeEdge = ref(null)
const activeSlot = ref(null)
const activeKey = ref(null)
const pasoMensaje = ref('')
const animando = ref(false)
const splitAnim = ref(null)
const deleteAnim = ref(null)
const canvasRef = ref(null)
const cambiandoVista = ref(false)

const arbolOperaciones = computed(() => {
  if (estructura.value) return estructura.value
  if (tipoVistaArbol.value === 'indice') return estructuraIndice.value
  return null
})

const estructuraVisual = computed(() => {
  if (animando.value) {
    return estructuraAnimacion.value ?? arbolOperaciones.value
  }
  return estructuraAnimacion.value ?? estructura.value
})

async function onCambiarVista(tipo) {
  if (!inicializado.value || animando.value) return
  cambiandoVista.value = true
  try {
    await cambiarVistaArbol(tipo)
    if (visualizacionHabilitada.value) canvasRef.value?.verTodo()
    logLocal(`Vista: ${tipo}`, 'info')
  } catch (err) {
    logLocal(err.message, 'error')
  } finally {
    cambiandoVista.value = false
  }
}

async function onCargarTodaBD({ orden }) {
  limpiarAnimacion()
  logLocal('Consultando total de canciones en BD…', 'info')
  try {
    const { total } = await api.obtenerConteoBD()
    if (!total) {
      logLocal('La base de datos no tiene canciones', 'error')
      return
    }
    logLocal(`Cargando ${total} canciones (toda la BD)…`, 'info')
    await onInicializar({ orden: orden ?? config.value.orden ?? 2, limite: total, modo: 'rapido' })
  } catch (err) {
    logLocal(
      err.message?.includes('404')
        ? 'Endpoint /api/conteo-bd no encontrado — reinicia el backend: go run .'
        : err.message,
      'error',
    )
  }
}

function mensajeStatsOperacion(accion, cancion, stats) {
  const id = cancion?.Indice ?? '?'
  const nombre = cancion?.TrackName ?? ''
  const vista = stats?.etiquetaIndice ?? tipoVistaArbol.value
  return `${accion}: «${nombre}» (ID ${id}). Índice «${vista}» — ${stats?.nodosInternos ?? '?'} internos, ${stats?.nodosHoja ?? '?'} hojas, ${stats?.totalCanciones ?? '?'} canciones.`
}
const DELAY_RAPIDO = 550
const DELAY_POST_DELETE = 1100

async function cargarProximoIndice() {
  if (!inicializado.value) {
    proximoIndice.value = null
    return
  }
  try {
    const data = await api.obtenerSiguienteIndice()
    proximoIndice.value = data.indice
  } catch (err) {
    proximoIndice.value = null
    logLocal(`No se pudo obtener el próximo ID: ${err.message}`, 'warn')
  }
}

onMounted(() => {
  limpiarAnimacion()
})

function requiereArbolCargado(accion) {
  if (!inicializado.value) {
    logLocal(`${accion}: inicializa el árbol primero`, 'error')
    return false
  }
  if (!visualizacionHabilitada.value) return true
  if (!arbolOperaciones.value) {
    logLocal(
      `${accion}: árbol de la vista «${tipoVistaArbol.value}» no cargado — reinicia el backend (go run .) y cambia de vista`,
      'error',
    )
    return false
  }
  return true
}

function logLocal(texto, tipo = 'info') {
  logAcademico.value.unshift({
    id: Date.now() + Math.random(),
    texto,
    tipo,
    hora: new Date().toLocaleTimeString('es-PE', { hour12: false }),
  })
  if (logAcademico.value.length > 50) logAcademico.value.pop()
}

function limpiarAnimacion() {
  activeEdge.value = null
  activeSlot.value = null
  activeKey.value = null
  pasoMensaje.value = ''
  animando.value = false
  splitAnim.value = null
  deleteAnim.value = null
  splitNodes.value = []
  underflowNodes.value = []
  deletingKeys.value = []
  activeLeafChain.value = null
  rangeHighlightLeaf.value = null
  estructuraAnimacion.value = null
  modoPasoAPaso.value = false
  pasoActual.value = 0
  pasoTotal.value = 0
  resolverSiguiente = null
}

function onSiguientePaso() {
  if (resolverSiguiente) {
    resolverSiguiente()
    resolverSiguiente = null
  }
}

function esperarSiguiente() {
  return new Promise((resolve) => {
    resolverSiguiente = resolve
  })
}

const DELAY_PASO = 900

function esperar(ms = DELAY_PASO) {
  return new Promise((r) => setTimeout(r, ms))
}

/** Aplica estado visual de un paso (eliminación o inserción). */
function aplicarEstadoPaso(paso, animInfo = null) {
  pasoMensaje.value = paso.mensaje
  highlightedNodes.value = paso.highlightedNodes?.length ? paso.highlightedNodes : [paso.nodeId]
  activeSlot.value = paso.activeSlot ?? null
  activeKey.value = paso.activeKey ?? null
  activeEdge.value = paso.activeEdge ?? null
  deletingKeys.value = paso.deletingKeys ?? []

  if (paso.highlightedKeys?.length) {
    highlightedKeys.value = paso.highlightedKeys
  }

  if (paso.underflowNodes !== undefined) {
    underflowNodes.value = paso.underflowNodes
  } else if (paso.underflowPreview || paso.tipo?.includes('underflow')) {
    underflowNodes.value = paso.highlightedNodes?.filter((id) => {
      const n = (estructuraVisual.value ?? estructura.value)?.nodos?.find((x) => x.id === id)
      return n?.esHoja
    }) ?? [paso.nodeId]
  } else if (!paso.deletePhase) {
    underflowNodes.value = []
  }

  if (paso.leafChainFrom != null && paso.leafChainTo != null) {
    activeLeafChain.value = { from: paso.leafChainFrom, to: paso.leafChainTo }
  } else if (!paso.tipo?.startsWith('range_')) {
    activeLeafChain.value = null
  }

  if (paso.borradoVirtual && paso.deletingKeys?.length && estructuraAnterior.value) {
    estructuraAnimacion.value = estructuraConBorradoVirtual(
      estructuraAnterior.value,
      paso.nodeId,
      paso.deletingKeys[0],
      paso.posEliminar ?? animInfo?.posEliminar,
    )
    deletingKeys.value = []
  } else if (paso.borradoVirtual && estructuraAnterior.value && animInfo?.claveEliminada != null) {
    const hId =
      animInfo.hojaId ??
      paso.highlightedNodes?.find((id) => {
        const n = estructuraAnterior.value.nodos.find((x) => x.id === id)
        return n?.esHoja
      })
    if (hId != null) {
      estructuraAnimacion.value = estructuraConBorradoVirtual(
        estructuraAnterior.value,
        hId,
        animInfo.claveEliminada,
        animInfo.posEliminar,
      )
    }
  }

  if (paso.insercionVirtual && estructuraAnterior.value) {
    const hId = paso.hojaId ?? paso.nodeId
    const clave = paso.claveInsertada ?? animInfo?.claveInsertada
    if (hId != null && clave != null) {
      estructuraAnimacion.value = estructuraConInsercionVirtual(
        estructuraAnterior.value,
        hId,
        clave,
      )
    }
  }

  if (paso.splitPhase && animInfo) {
    splitAnim.value = {
      ...animInfo,
      phase: paso.splitPhase,
      promoteProgress: paso.animarPromote ? 0 : 1,
      overflowNodeId: animInfo.hojaDestinoId ?? paso.nodeId,
    }
    if (paso.tipo?.includes('split') || paso.overflowPreview || paso.splitPhase) {
      splitNodes.value = [animInfo.hojaDestinoId ?? paso.nodeId]
    }
  }

  if (paso.deletePhase && animInfo?.claveEliminada != null) {
    const layoutEval =
      paso.hermanoEvaluado && estructuraVisual.value
        ? calcularLayout(estructuraVisual.value).posiciones[paso.hermanoEvaluado]
        : null
    deleteAnim.value = {
      ...animInfo,
      phase: paso.deletePhase,
      progress: paso.animarPrestamo ? 0 : 1,
      underflowNodeId: paso.underflowNodes?.[0] ?? animInfo.hojaId ?? null,
      atenuarIds: paso.atenuarIds ?? [],
      evaluarHermano: paso.evaluarHermano ?? null,
      hermanoRechazado: paso.hermanoRechazado ?? null,
      ...(layoutEval ? { hermanoPos: layoutEval } : {}),
    }
  } else if (paso.atenuarIds?.length && animInfo?.claveEliminada != null) {
    deleteAnim.value = { ...animInfo, atenuarIds: paso.atenuarIds }
  }

  if (paso.animarFusion) {
    splitNodes.value = paso.highlightedNodes?.length ? paso.highlightedNodes : [paso.nodeId]
  } else if (paso.tipo?.startsWith('confirmacion')) {
    splitNodes.value = []
    if (!paso.splitPhase) splitAnim.value = null
  }

  const esWarn =
    paso.tipo?.includes('underflow') ||
    paso.tipo?.includes('fusion') ||
    paso.tipo?.includes('evaluar') ||
    paso.tipo?.includes('propagacion') ||
    paso.tipo?.includes('overflow') ||
    paso.tipo?.includes('split')
  const esOk = paso.tipo?.startsWith('confirmacion')
  logLocal(paso.mensaje, esWarn ? 'warn' : esOk ? 'success' : 'info')
}

async function avanzarPaso(paso, animInfo, opciones) {
  aplicarEstadoPaso(paso, animInfo)

  if (opciones.manual) {
    await esperarSiguiente()
  }

  if (paso.animarPrestamo && deleteAnim.value) {
    deleteAnim.value = { ...deleteAnim.value, phase: FASES_DELETE.PRESTAMO, progress: 0 }
    if (opciones.manual) {
      await esperarSiguiente()
      deleteAnim.value = { ...deleteAnim.value, progress: 1 }
    } else {
      await animarValor(deleteAnim.value, 'progress', 0, 1, 1800)
    }
  } else if (paso.animarPromote && splitAnim.value) {
    splitAnim.value = { ...splitAnim.value, phase: FASES_SPLIT.PROMOTE, promoteProgress: 0 }
    if (opciones.manual) {
      await esperarSiguiente()
      splitAnim.value = { ...splitAnim.value, promoteProgress: 1 }
    } else {
      await animarValor(splitAnim.value, 'promoteProgress', 0, 1, 2000)
    }
  } else if (!opciones.manual && paso.tipo !== 'aplicar') {
    const delay =
      paso.tipo?.includes('evaluar') || paso.tipo?.includes('split') ? 1200 : DELAY_PASO
    await esperar(
      paso.tipo === 'eliminar_fisico' || paso.tipo === 'insertar_virtual' ? 1100 : delay,
    )
  }
}

/** Aplica un paso de lectura (búsqueda / range scan). */
function aplicarPasoGenerico(paso) {
  pasoMensaje.value = paso.mensaje
  highlightedNodes.value = paso.highlightedNodes?.length
    ? paso.highlightedNodes
    : paso.nodeId != null
      ? [paso.nodeId]
      : []
  activeSlot.value = paso.activeSlot ?? null
  activeKey.value = paso.activeKey ?? null
  activeEdge.value = paso.activeEdge ?? null
  deletingKeys.value = []

  if (paso.highlightedKeys?.length) {
    highlightedKeys.value = paso.highlightedKeys
  }

  if (paso.leafChainFrom != null && paso.leafChainTo != null) {
    activeLeafChain.value = { from: paso.leafChainFrom, to: paso.leafChainTo }
    rangeHighlightLeaf.value = paso.leafChainTo
  } else if (paso.tipo?.startsWith('range_')) {
    rangeHighlightLeaf.value = paso.nodeId
    activeLeafChain.value = null
  } else if (paso.tipo !== 'confirmacion') {
    activeLeafChain.value = null
    rangeHighlightLeaf.value = null
  }

  const esWarn = paso.tipo === 'no_encontrado'
  const esError = paso.tipo === 'confirmacion' && paso.mensaje?.startsWith('❌')
  const esOk =
    paso.tipo === 'encontrado' ||
    paso.tipo === 'range_fin' ||
    (paso.tipo === 'confirmacion' && !esError)
  logLocal(paso.mensaje, esWarn ? 'warn' : esError ? 'error' : esOk ? 'success' : 'info')
}

/** Secuencia de pasos de lectura con modo automático o manual. */
async function ejecutarSecuenciaLectura(pasos, opciones = {}) {
  const manual = opciones.manual ?? false
  const delay = opciones.delay ?? DELAY_PASO
  const delayRange = opciones.delayRange ?? DELAY_RAPIDO

  animando.value = true
  modoPasoAPaso.value = manual
  if (!opciones.mantenerKeys) highlightedKeys.value = []
  activeEdge.value = null
  deletingKeys.value = []
  pasoTotal.value = pasos.length

  for (let i = 0; i < pasos.length; i++) {
    pasoActual.value = i + 1
    aplicarPasoGenerico(pasos[i])

    if (manual) {
      await esperarSiguiente()
    } else {
      const paso = pasos[i]
      const ms = paso.tipo?.startsWith('range_')
        ? delayRange
        : paso.tipo === 'aplicar'
          ? 600
          : paso.tipo === 'confirmacion'
            ? 1200
            : delay
      await esperar(ms)
    }
  }

  if (pasos.length && !manual) {
    const ids = pasos.flatMap((p) =>
      p.highlightedNodes?.length ? p.highlightedNodes : p.nodeId != null ? [p.nodeId] : [],
    )
    if (ids.length) highlightedNodes.value = [...new Set(ids)]
    await esperar(350)
  }
}

async function animarValor(obj, prop, desde, hasta, ms) {
  const pasos = Math.max(Math.round(ms / 40), 1)
  const intervalo = ms / pasos
  for (let i = 0; i <= pasos; i++) {
    const t = i / pasos
    obj[prop] = desde + (hasta - desde) * t
    await esperar(intervalo)
  }
}

async function animarInserccionEnDemo(cancion, orden, manual, pasoContador) {
  const id = idDeCancion(cancion)
  estructuraAnterior.value = JSON.parse(JSON.stringify(arbolOperaciones.value))

  const { pasos, insertInfo } = calcularSecuenciaInsercion(arbolOperaciones.value, id, orden)
  const pasoAplicar = pasos.find((p) => p.tipo === 'aplicar')
  const pasosPrevios = pasos.filter((p) => p.tipo !== 'aplicar')

  if (insertInfo?.huboSplit) {
    splitAnim.value = { ...insertInfo, promoteProgress: 0 }
  }

  for (const paso of pasosPrevios) {
    pasoContador.value++
    pasoActual.value = pasoContador.value
    await avanzarPaso(paso, insertInfo, { manual })
  }

  pasoContador.value++
  pasoActual.value = pasoContador.value
  if (pasoAplicar) {
    aplicarEstadoPaso(pasoAplicar, insertInfo)
    if (manual) await esperarSiguiente()
    else await esperar(350)
  }

  await api.insertarEnArbol(cancion)
  estructuraAnimacion.value = null
  await refrescarEstructura()
  splitAnim.value = null
  splitNodes.value = []
}

async function onInicializar({ orden, limite, modo = 'rapido' }) {
  limpiarAnimacion()
  highlightedNodes.value = []

  const esDemo =
    limite <= LIMITE_DEMO_INICIAL && (modo === 'completo' || modo === 'pasoAPaso')

  if (!esDemo) {
    try {
      const resultado = await inicializarArbol(orden, limite)
      const msgViz = visualizacionHabilitada.value
        ? ''
        : ` — canvas desactivado (>${limiteVisualizacion} canciones); operaciones API activas`
      logLocal(`Árbol inicializado: d=${orden}, ${resultado.total} canciones cargadas${msgViz}`, 'success')
      agregar(`🌳 Árbol B+-Tree inicializado (d=${orden}, n=${resultado.total})`, 'info')
      estructuraAnterior.value = null
      await cargarProximoIndice()
      canvasRef.value?.verTodo()
    } catch (err) {
      logLocal(err.message, 'error')
    }
    return
  }

  const manual = modo === 'pasoAPaso'
  modoPasoAPaso.value = manual
  animando.value = true
  pasoTotal.value = 0
  const pasoContador = ref(0)

  try {
    // Preparar árbol vacío primero → habilita canvas (total=0 pero limite≤500)
    const prep = await prepararDemoInicial(orden, limite)
    const canciones = prep.canciones ?? []
    await nextTick()
    canvasRef.value?.verTodo()

    const preLectura = calcularPasosPreLecturaBD(limite)
    for (const paso of preLectura) {
      pasoContador.value++
      pasoActual.value = pasoContador.value
      aplicarPasoGenerico(paso)
      if (manual) await esperarSiguiente()
      else await esperar(700)
    }

    const intro = calcularPasosIntroInicializacion(orden, canciones.length)
    for (const paso of intro) {
      pasoContador.value++
      pasoActual.value = pasoContador.value
      aplicarPasoGenerico(paso)
      if (manual) await esperarSiguiente()
      else await esperar(900)
    }

    setCargando(true)
    const t0 = performance.now()

    for (let i = 0; i < canciones.length; i++) {
      const registro = calcularPasoRegistroBD(i + 1, canciones.length, canciones[i])
      pasoContador.value++
      pasoActual.value = pasoContador.value
      aplicarPasoGenerico(registro)
      if (manual) await esperarSiguiente()
      else await esperar(500)

      const enc = calcularPasoEncabezadoCancion(i + 1, canciones.length, canciones[i])
      pasoContador.value++
      pasoActual.value = pasoContador.value
      aplicarPasoGenerico(enc)
      if (manual) await esperarSiguiente()
      else await esperar(450)

      await animarInserccionEnDemo(canciones[i], orden, manual, pasoContador)
    }

    finalizarDemoInicial(canciones.length)
    const fin = calcularPasosFinInicializacion(canciones.length, performance.now() - t0)
    for (const paso of fin) {
      pasoContador.value++
      pasoActual.value = pasoContador.value
      aplicarPasoGenerico(paso)
      if (manual) await esperarSiguiente()
      else await esperar(1200)
    }

    logLocal(
      `Demo inicialización: ${canciones.length} canciones (d=${orden})`,
      'success',
    )
    agregar(
      `🌳 Demo construcción árbol: ${canciones.length} inserciones (d=${orden})`,
      'info',
    )
    await cargarProximoIndice()
    canvasRef.value?.verTodo()
    if (!manual) await esperar(800)
  } catch (err) {
    logLocal(err.message, 'error')
  } finally {
    setCargando(false)
    limpiarAnimacion()
  }
}

async function onInsertar({ cancion, modo }) {
  const idPrevisto = proximoIndice.value ?? cancion?.Indice ?? cancion?.indice
  if (idPrevisto == null) {
    logLocal(
      'Insertar: no hay próximo ID — inicializa el árbol y verifica que el backend Go esté en marcha (go run .)',
      'error',
    )
    return
  }

  if (!visualizacionHabilitada.value) {
    try {
      const t0 = performance.now()
      const payload = { ...cancion, Indice: cancion.Indice ?? idPrevisto }
      await api.insertar(payload)
      await refrescarEstadisticas(tipoVistaArbol.value)
      await cargarProximoIndice()
      const ms = (performance.now() - t0).toFixed(1)
      const msg = mensajeStatsOperacion('Inserción', payload, estadisticas.value)
      registrarOperacion(`${msg} (${ms} ms)`)
      logLocal(`✅ ${msg}`, 'success')
      agregar(`Inserción sin canvas — ${estadisticas.value?.nodosTotales ?? '?'} nodos`, 'info')
    } catch (err) {
      logLocal(err.message, 'error')
    }
    return
  }

  if (!requiereArbolCargado('Insertar')) return

  const claveAnim = normalizarClaveArbol(
    arbolOperaciones.value,
    claveDeCancionPorVista({ ...cancion, Indice: idPrevisto }, tipoVistaArbol.value),
  )
  const etiquetaClave = etiquetaClaveVista(tipoVistaArbol.value)
  estructuraAnterior.value = JSON.parse(JSON.stringify(arbolOperaciones.value))

  const manual = modo === 'pasoAPaso'
  try {
    limpiarAnimacion()
    modoPasoAPaso.value = manual
    animando.value = true
    highlightedKeys.value = []

    const secuencia = calcularSecuenciaInsercion(
      arbolOperaciones.value,
      claveAnim,
      config.value.orden,
    )
    const { pasos, insertInfo, tipoCorreccion } = secuencia

    const pasoAplicar = pasos.find((p) => p.tipo === 'aplicar')
    const pasosPrevios = pasos.filter((p) => p.tipo !== 'aplicar')
    pasoTotal.value = pasosPrevios.length + 2

    if (insertInfo?.huboSplit) {
      splitAnim.value = { ...insertInfo, promoteProgress: 0 }
    }

    for (let i = 0; i < pasosPrevios.length; i++) {
      pasoActual.value = i + 1
      await avanzarPaso(pasosPrevios[i], insertInfo, { manual })
    }

    pasoActual.value = pasosPrevios.length + 1
    if (pasoAplicar) {
      aplicarEstadoPaso(pasoAplicar, insertInfo)
      if (manual) await esperarSiguiente()
      else await esperar(700)
    }

    const inicio = performance.now()
    const resultado = await api.insertar({ ...cancion, Indice: idPrevisto })
    const idAsignado = resultado.indice
    estructuraAnimacion.value = null
    await refrescarEstructura()
    await cargarProximoIndice()
    const ms = performance.now() - inicio

    const confirmacion = calcularPasosConfirmacionInsercion(
      insertInfo,
      arbolOperaciones.value,
      idAsignado,
    )
    pasoTotal.value = pasosPrevios.length + 1 + confirmacion.length

    splitAnim.value = null
    splitNodes.value = []

    for (let i = 0; i < confirmacion.length; i++) {
      pasoActual.value = pasosPrevios.length + 1 + i
      await avanzarPaso(confirmacion[i], insertInfo, { manual })
    }

    const msgExtra = tipoCorreccion === 'split' ? ' con SPLIT' : ''
    pasoMensaje.value = `✅ ${etiquetaClave}=${claveAnim} insertada${msgExtra} — ${ms.toFixed(1)} ms`
    logLocal(pasoMensaje.value, 'success')
    agregar(`Inserción ${etiquetaClave}=${claveAnim}${msgExtra} (${ms.toFixed(3)} ms)`, 'info')
    if (!manual) await esperar(1000)
  } catch (err) {
    pasoMensaje.value = `❌ ${err.message}`
    logLocal(`Insertar: ${err.message}`, 'error')
  } finally {
    limpiarAnimacion()
  }
}

async function onBuscar({ indice, modo }) {
  if (!visualizacionHabilitada.value) {
    try {
      const t0 = performance.now()
      const cancion = await api.buscar(indice)
      await refrescarEstadisticas(tipoVistaArbol.value)
      const ms = (performance.now() - t0).toFixed(1)
      const msg = `Búsqueda ID=${indice} — «${cancion.TrackName}» encontrada (${ms} ms).`
      registrarOperacion(msg)
      logLocal(`✅ ${msg}`, 'success')
    } catch (err) {
      logLocal(err.message, 'error')
    }
    return
  }

  if (!requiereArbolCargado('Buscar')) return

  const manual = modo === 'pasoAPaso'
  const t0 = performance.now()

  let claveBusqueda = indice
  let cancionPrevia = null
  const vistaActual = tipoVistaDesdeEstructura(arbolOperaciones.value)
  if (vistaActual !== 'indice') {
    try {
      cancionPrevia = await api.buscar(indice)
      claveBusqueda = claveDeCancionPorEstructura(cancionPrevia, arbolOperaciones.value)
    } catch (err) {
      logLocal(err.message, 'error')
      return
    }
  }
  claveBusqueda = normalizarClaveArbol(arbolOperaciones.value, claveBusqueda)

  try {
    limpiarAnimacion()

    const { pasos, encontradoPredicho } = calcularSecuenciaBusqueda(
      arbolOperaciones.value,
      claveBusqueda,
      vistaActual !== 'indice' ? indice : null,
    )
    const pasoAplicar = pasos.find((p) => p.tipo === 'aplicar')
    const pasosPrevios = pasos.filter((p) => p.tipo !== 'aplicar')

    pasoTotal.value = pasosPrevios.length + 2

    await ejecutarSecuenciaLectura(pasosPrevios, { manual })

    pasoActual.value = pasosPrevios.length + 1
    if (pasoAplicar) {
      aplicarPasoGenerico(pasoAplicar)
      if (manual) await esperarSiguiente()
      else await esperar(600)
    }

    const cancion = cancionPrevia ?? (await api.buscar(indice))
    const ms = performance.now() - t0
    const confirmacion = calcularPasosConfirmacionBusqueda(
      claveBusqueda,
      cancion,
      encontradoPredicho,
      ms,
    )
    pasoTotal.value = pasosPrevios.length + 1 + confirmacion.length

    for (let i = 0; i < confirmacion.length; i++) {
      pasoActual.value = pasosPrevios.length + 1 + i
      aplicarPasoGenerico(confirmacion[i])
      if (manual) await esperarSiguiente()
      else await esperar(1500)
    }

    agregar(`Búsqueda ID=${indice} — "${cancion.TrackName}" (${ms.toFixed(3)} ms)`, 'success')
    if (!manual) await esperar(800)
  } catch (err) {
    const confirmacion = calcularPasosConfirmacionBusqueda(
      claveBusqueda,
      null,
      false,
      performance.now() - t0,
      err.message,
    )
    for (const paso of confirmacion) {
      aplicarPasoGenerico(paso)
      if (manual) await esperarSiguiente()
      else await esperar(1500)
    }
  } finally {
    limpiarAnimacion()
  }
}

async function onEliminar({ indice, modo }) {
  if (!visualizacionHabilitada.value) {
    try {
      const t0 = performance.now()
      const prev = await api.buscar(indice).catch(() => null)
      await api.eliminar(indice)
      await refrescarEstadisticas(tipoVistaArbol.value)
      await cargarProximoIndice()
      const ms = (performance.now() - t0).toFixed(1)
      const msg = mensajeStatsOperacion('Eliminación', prev ?? { Indice: indice }, estadisticas.value)
      registrarOperacion(`${msg} (${ms} ms)`)
      logLocal(`✅ ${msg}`, 'success')
    } catch (err) {
      logLocal(err.message, 'error')
    }
    return
  }

  if (!requiereArbolCargado('Eliminar')) return

  // 1) Resolución interna por ID (sin animación de búsqueda en primario)
  let cancion
  try {
    cancion = await api.buscar(indice)
  } catch (err) {
    logLocal(`❌ ID ${indice} no existe en el índice primario — ${err.message}`, 'error')
    return
  }

  const vistaActual = tipoVistaDesdeEstructura(arbolOperaciones.value)
  const claveElim = normalizarClaveArbol(
    arbolOperaciones.value,
    claveDeCancionPorEstructura(cancion, arbolOperaciones.value),
  )
  const etiqueta = etiquetaClaveVista(vistaActual)
  const idElim = Number(indice)

  if (vistaActual !== 'indice') {
    logLocal(
      `🔎 ID ${idElim} («${cancion.TrackName}») verificado en índice primario — animación con clave ${etiqueta}=${claveElim}`,
      'info',
    )
  }

  estructuraAnterior.value = JSON.parse(JSON.stringify(arbolOperaciones.value))
  const manual = modo === 'pasoAPaso'

  try {
    limpiarAnimacion()
    modoPasoAPaso.value = manual
    animando.value = true
    highlightedKeys.value = []

    const secuencia = calcularSecuenciaEliminacion(
      arbolOperaciones.value,
      claveElim,
      config.value.orden,
      idElim,
    )
    const { pasos, encontrado, deleteInfo, tipoCorreccion } = secuencia

    if (!encontrado) {
      pasoTotal.value = pasos.length
      for (let i = 0; i < pasos.length; i++) {
        pasoActual.value = i + 1
        aplicarPasoGenerico(pasos[i])
        if (manual) await esperarSiguiente()
        else await esperar(900)
      }
      pasoMensaje.value = `❌ No se puede eliminar: clave ${claveElim} (ID ${indice}) no está en el árbol ${etiqueta}`
      logLocal(pasoMensaje.value, 'error')
      if (manual) await esperarSiguiente()
      else await esperar(1500)
      return
    }

    const pasoAplicar = pasos.find((p) => p.tipo === 'aplicar')
    const pasosPrevios = pasos.filter((p) => p.tipo !== 'aplicar')
    pasoTotal.value = pasosPrevios.length + 2

    deleteAnim.value = deleteInfo ? { ...deleteInfo, progress: 0 } : null

    for (let i = 0; i < pasosPrevios.length; i++) {
      pasoActual.value = i + 1
      await avanzarPaso(pasosPrevios[i], deleteInfo, { manual })
    }

    pasoActual.value = pasosPrevios.length + 1
    if (pasoAplicar) {
      aplicarEstadoPaso(pasoAplicar, deleteInfo)
      if (manual) await esperarSiguiente()
      else await esperar(700)
    }

    const inicio = performance.now()
    const respuesta = await api.eliminar(indice)
    estructuraAnimacion.value = null
    await refrescarEstructura()
    const ms = performance.now() - inicio

    const cancionResp = respuesta?.cancion
    if (cancionResp?.TrackName) {
      logLocal(
        `🗑 Cascada: «${cancionResp.TrackName}» eliminada de Indice, TrackName, Popularity, Tempo y Danceability`,
        'info',
      )
    }

    const verificado = analizarEliminacion(
      estructuraAnterior.value,
      arbolOperaciones.value,
      claveElim,
      config.value.orden,
      idElim,
    )
    const infoFinal = {
      ...deleteInfo,
      profAntes: verificado.profAntes,
      profDespues: verificado.profDespues,
    }
    const confirmacion = calcularPasosConfirmacion(infoFinal, arbolOperaciones.value)
    pasoTotal.value = pasosPrevios.length + 1 + confirmacion.length

    deleteAnim.value = null
    splitNodes.value = []
    underflowNodes.value = []

    for (let i = 0; i < confirmacion.length; i++) {
      pasoActual.value = pasosPrevios.length + 1 + i
      await avanzarPaso(confirmacion[i], infoFinal, { manual })
    }

    const msgExtra =
      tipoCorreccion === 'fantasma'
        ? ' (separador fantasma)'
        : tipoCorreccion === 'prestamo'
          ? ' (redistribución)'
          : tipoCorreccion === 'fusion'
            ? ' (concatenación)'
            : ''

    pasoMensaje.value = `✅ ${etiqueta} ${claveElim} eliminado (ID ${indice})${msgExtra} — ${ms.toFixed(1)} ms`
    logLocal(pasoMensaje.value, 'success')
    agregar(`Eliminación ID=${indice}${msgExtra} (${ms.toFixed(3)} ms)`, 'info')
    await cargarProximoIndice()
    if (!manual) await esperar(1000)
  } catch (err) {
    pasoMensaje.value = `❌ ${err.message}`
    logLocal(err.message, 'error')
  } finally {
    limpiarAnimacion()
  }
}

async function onRango({ inicio, fin, modo, campo, tipoClave }) {
  const vista = campo || tipoVistaArbol.value
  const clave = tipoClave || (vista === 'nombre' ? 'string' : vista === 'tempo' || vista === 'danceability' ? 'float' : 'int')

  const ejecutarRango = () => {
    switch (vista) {
      case 'nombre':
        return api.buscarRangoNombre(inicio, fin)
      case 'popularidad':
        return api.buscarRangoCampo('popularidad', inicio, fin)
      case 'tempo':
        return api.buscarRangoCampo('tempo', inicio, fin)
      case 'danceability':
        return api.buscarRangoCampo('danceability', inicio, fin)
      default:
        return api.buscarRango(inicio, fin)
    }
  }

  if (!visualizacionHabilitada.value) {
    try {
      const t0 = performance.now()
      const resultado = await ejecutarRango()
      logLocal(
        `✅ Range [${inicio}–${fin}] (${clave}): ${resultado.total} resultados (${(performance.now() - t0).toFixed(1)} ms)`,
        'success',
      )
    } catch (err) {
      logLocal(err.message, 'error')
    }
    return
  }

  if (!requiereArbolCargado('Range Scan')) return

  const manual = modo === 'pasoAPaso'
  const t0 = performance.now()

  try {
    limpiarAnimacion()

    const { pasos } = calcularSecuenciaRangeScan(arbolOperaciones.value, inicio, fin, clave)
    const pasoAplicar = pasos.find((p) => p.tipo === 'aplicar')
    const pasosPrevios = pasos.filter((p) => p.tipo !== 'aplicar')

    pasoTotal.value = pasosPrevios.length + 2

    await ejecutarSecuenciaLectura(pasosPrevios, {
      manual,
      delayRange: DELAY_RAPIDO,
    })

    pasoActual.value = pasosPrevios.length + 1
    if (pasoAplicar) {
      aplicarPasoGenerico(pasoAplicar)
      if (manual) await esperarSiguiente()
      else await esperar(600)
    }

    const resultado = await ejecutarRango()
    const ms = performance.now() - t0

    const confirmacion = calcularPasosConfirmacionRango(inicio, fin, resultado, ms, clave)
    pasoTotal.value = pasosPrevios.length + 1 + confirmacion.length

    for (let i = 0; i < confirmacion.length; i++) {
      pasoActual.value = pasosPrevios.length + 1 + i
      aplicarPasoGenerico(confirmacion[i])
      if (manual) await esperarSiguiente()
      else await esperar(1800)
    }

    agregar(
      `Range Scan [${inicio}–${fin}] (${clave}): ${resultado.total} resultados (${ms.toFixed(3)} ms)`,
      'success',
    )
    if (!manual) await esperar(800)
  } catch (err) {
    logLocal(err.message, 'error')
    pasoMensaje.value = `❌ ${err.message}`
    if (manual) await esperarSiguiente()
    else await esperar(1500)
  } finally {
    limpiarAnimacion()
  }
}

const stats = computed(() => {
  if (!inicializado.value) return null
  const est = estadisticas.value
  const base = {
    canciones: est?.totalCanciones ?? config.value.total,
    orden: est?.orden ?? config.value.orden,
    vista: est?.etiquetaIndice ?? estructura.value?.etiquetaIndice ?? tipoVistaArbol.value,
    tipoClave: est?.tipoClave ?? estructura.value?.tipoClave ?? 'int',
    visualizacion: visualizacionHabilitada.value,
    altura: est?.altura ?? null,
  }
  if (est) {
    return {
      ...base,
      nodos: est.nodosTotales,
      hojas: est.nodosHoja,
      internos: est.nodosInternos,
      minClaves: config.value.orden,
      maxClaves: config.value.orden * 2,
    }
  }
  if (!estructura.value) {
    return {
      ...base,
      nodos: '—',
      hojas: '—',
      internos: '—',
      minClaves: config.value.orden,
      maxClaves: config.value.orden * 2,
    }
  }
  return {
    ...base,
    nodos: estructura.value.nodos.length,
    hojas: estructura.value.cadenaHojas.length,
    internos: estructura.value.nodos.length - estructura.value.cadenaHojas.length,
    minClaves: estructura.value.minClaves ?? config.value.orden,
    maxClaves: estructura.value.maxClaves ?? config.value.orden * 2,
  }
})
</script>

<template>
  <div class="academic-mode">
    <ControlPanel
      :cargando="cargando"
      :inicializado="inicializado"
      :proximo-indice="proximoIndice"
      :deshabilitar-acciones="animando"
      :tipo-vista-arbol="tipoVistaArbol"
      :cambiando-vista="cambiandoVista"
      :visualizacion-habilitada="visualizacionHabilitada"
      :total-canciones="config.total"
      :limite-visualizacion="limiteVisualizacion"
      @inicializar="onInicializar"
      @buscar="onBuscar"
      @insertar="onInsertar"
      @eliminar="onEliminar"
      @rango="onRango"
      @cambiar-vista="onCambiarVista"
      @cargar-toda-bd="onCargarTodaBD"
    />

    <div class="academic-main">
      <div v-if="stats" class="stats-bar">
        <span>Vista: <strong>{{ stats.vista }}</strong> ({{ stats.tipoClave }})</span>
        <span>Orden d: <strong>{{ stats.orden }}</strong> (min {{ stats.minClaves }}, máx {{ stats.maxClaves }})</span>
        <span>Nodos: <strong>{{ stats.nodos }}</strong></span>
        <span>Internos: <strong>{{ stats.internos }}</strong></span>
        <span>Hojas: <strong>{{ stats.hojas }}</strong></span>
        <span v-if="stats.altura != null">Altura: <strong>{{ stats.altura }}</strong></span>
        <span>Canciones: <strong>{{ stats.canciones }}</strong></span>
        <span v-if="!stats.visualizacion" class="stats-warn">
          Canvas off (límite {{ limiteVisualizacion }}) — operaciones activas
        </span>
        <span class="legend">
          <span class="leg leg--ptr">▢</span> Puntero
          <span class="leg leg--sep">▢</span> Separador
          <span class="leg leg--leaf">▢</span> Hoja
        </span>
      </div>

      <TreeCanvas
        v-if="visualizacionHabilitada"
        ref="canvasRef"
        :estructura="estructuraVisual"
        :split-anim="splitAnim"
        :delete-anim="deleteAnim"
        :highlighted-nodes="highlightedNodes"
        :split-nodes="splitNodes"
        :underflow-nodes="underflowNodes"
        :highlighted-keys="highlightedKeys"
        :deleting-keys="deletingKeys"
        :range-highlight-leaf="rangeHighlightLeaf"
        :active-leaf-chain="activeLeafChain"
        :active-edge="activeEdge"
        :active-slot="activeSlot"
        :active-key="activeKey"
        :paso-mensaje="pasoMensaje"
        :animando="animando"
        :paso-actual="pasoActual"
        :paso-total="pasoTotal"
        :modo-paso-a-paso="modoPasoAPaso"
        @siguiente="onSiguientePaso"
      />

      <div v-else-if="inicializado && config.total > limiteVisualizacion">
        <TreeStatsPanel
          :stats="estadisticas"
          :ultima-operacion="ultimaOperacion"
          :tipo-vista="tipoVistaArbol"
          :limite-visualizacion="limiteVisualizacion"
        />
      </div>

      <p v-else class="canvas-placeholder">
        Inicializa el árbol (≤ {{ limiteVisualizacion }} canciones) para ver el canvas interactivo.
      </p>

      <OperationLog :mensajes="logAcademico" />
    </div>
  </div>
</template>

<style scoped>
.academic-mode {
  display: flex;
  gap: 1rem;
  height: 100%;
  padding: 1rem;
}

.academic-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-width: 0;
}

.stats-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  padding: 0.5rem 0.75rem;
  background: var(--bg-panel);
  border-radius: 6px;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.stats-bar strong {
  color: var(--text-primary);
}

.stats-warn {
  color: #e6a817;
  font-weight: 500;
}

.canvas-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 280px;
  padding: 1.5rem;
  text-align: center;
  font-size: 0.85rem;
  color: var(--text-secondary);
  background: var(--bg-panel);
  border-radius: 8px;
  border: 1px dashed var(--border);
}

.canvas-off {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 0.75rem;
  padding: 2rem;
  background: var(--bg-panel);
  border-radius: 8px;
  border: 1px dashed #555;
  color: var(--text-secondary);
  min-height: 320px;
}

.canvas-off h3 {
  margin: 0;
  color: var(--text-primary);
}

.canvas-off ul {
  margin: 0;
  padding-left: 1.25rem;
}

.legend {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.75rem;
}

.leg {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 2px;
  margin-right: 3px;
  vertical-align: middle;
}

.leg--ptr { background: #bfdbfe; border: 1px solid #93c5fd; }
.leg--sep { background: #fff; border: 1px solid #9ca3af; }
.leg--leaf { background: #faf5ff; border: 2px solid #7c3aed; }
</style>
