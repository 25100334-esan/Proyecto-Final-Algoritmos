import { ref, computed } from 'vue'
import * as api from '../api/bplustreeApi.js'
import { visualizacionPermitida, debeMostrarCanvas, LIMITE_VISUALIZACION } from '../utils/visualizationLimits.js'

const estructura = ref(null)
const estructuraIndice = ref(null)
const estadisticas = ref(null)
const tipoVistaArbol = ref('indice')
const inicializado = ref(false)
const cargando = ref(false)
const error = ref(null)
const config = ref({ orden: 2, limite: 100, total: 0 })
const ultimaOperacion = ref('')

export function useTreeState() {
  const visualizacionHabilitada = computed(() =>
    debeMostrarCanvas(config.value.total, config.value.limite, inicializado.value),
  )

  function canvasActivo() {
    return debeMostrarCanvas(config.value.total, config.value.limite, inicializado.value)
  }

  async function refrescarEstadisticas(tipo = tipoVistaArbol.value) {
    if (!inicializado.value) {
      estadisticas.value = null
      return null
    }
    try {
      const data = await api.obtenerEstadisticas(tipo)
      estadisticas.value = data
      if (data.totalCanciones != null) {
        config.value = { ...config.value, total: data.totalCanciones }
      }
      return data
    } catch (err) {
      error.value = err.message
      return null
    }
  }

  async function refrescarEstructura(tipo = tipoVistaArbol.value) {
    await refrescarEstadisticas(tipo)
    if (!canvasActivo()) {
      estructura.value = null
      estructuraIndice.value = null
      return null
    }
    const data = await api.obtenerEstructura(tipo)
    estructura.value = data
    if (tipo === 'indice') {
      estructuraIndice.value = data
    } else {
      estructuraIndice.value = await api.obtenerEstructura('indice')
    }
    return data
  }

  async function cambiarVistaArbol(tipo) {
    tipoVistaArbol.value = tipo
    if (!inicializado.value) return null
    if (!canvasActivo()) {
      return refrescarEstadisticas(tipo)
    }
    return refrescarEstructura(tipo)
  }

  async function inicializarArbol(orden, limite) {
    cargando.value = true
    error.value = null
    try {
      const resultado = await api.inicializar(orden, limite)
      config.value = { orden, limite, total: resultado.total }
      tipoVistaArbol.value = 'indice'
      ultimaOperacion.value = `Inicializado: ${resultado.total} canciones cargadas desde BD (d=${orden}).`
      inicializado.value = true
      if (visualizacionPermitida(resultado.total)) {
        await refrescarEstructura('indice')
      } else {
        estructura.value = null
        estructuraIndice.value = null
        const s = await refrescarEstadisticas('indice')
        if (s) {
          ultimaOperacion.value += ` Árbol: ${s.nodosInternos} internos, ${s.nodosHoja} hojas, altura ${s.altura}, ${s.nodosTotales} nodos totales.`
        }
      }
      return resultado
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      cargando.value = false
    }
  }

  async function prepararDemoInicial(orden, limite) {
    cargando.value = true
    error.value = null
    try {
      const resultado = await api.inicializarDemo(orden, limite)
      config.value = { orden, limite, total: 0 }
      tipoVistaArbol.value = 'indice'
      inicializado.value = true
      if (debeMostrarCanvas(0, limite, true)) {
        await refrescarEstructura('indice')
      } else {
        estructura.value = null
        estructuraIndice.value = null
      }
      return resultado
    } catch (err) {
      error.value = err.message
      inicializado.value = false
      estructura.value = null
      estructuraIndice.value = null
      estadisticas.value = null
      throw err
    } finally {
      cargando.value = false
    }
  }

  function finalizarDemoInicial(total) {
    config.value = { ...config.value, total }
    if (!canvasActivo()) {
      estructura.value = null
      estructuraIndice.value = null
      refrescarEstadisticas(tipoVistaArbol.value)
    }
  }

  function ajustarTotal(delta) {
    config.value = { ...config.value, total: Math.max(0, config.value.total + delta) }
    if (!canvasActivo()) {
      estructura.value = null
      estructuraIndice.value = null
    }
  }

  function registrarOperacion(texto) {
    ultimaOperacion.value = texto
  }

  return {
    estructura: computed(() => estructura.value),
    estructuraIndice: computed(() => estructuraIndice.value),
    estadisticas: computed(() => estadisticas.value),
    ultimaOperacion: computed(() => ultimaOperacion.value),
    tipoVistaArbol: computed(() => tipoVistaArbol.value),
    visualizacionHabilitada,
    limiteVisualizacion: LIMITE_VISUALIZACION,
    inicializado: computed(() => inicializado.value),
    cargando: computed(() => cargando.value),
    error: computed(() => error.value),
    config: computed(() => config.value),
    refrescarEstructura,
    refrescarEstadisticas,
    cambiarVistaArbol,
    inicializarArbol,
    prepararDemoInicial,
    finalizarDemoInicial,
    ajustarTotal,
    registrarOperacion,
    setCargando: (v) => { cargando.value = v },
  }
}
