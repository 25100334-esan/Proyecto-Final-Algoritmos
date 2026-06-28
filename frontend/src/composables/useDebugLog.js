import { ref, computed } from 'vue'

const logs = ref([])
const MAX_LOGS = 100

export function useDebugLog() {
  function agregar(mensaje, tipo = 'info') {
    const entrada = {
      id: Date.now() + Math.random(),
      mensaje,
      tipo,
      hora: new Date().toLocaleTimeString('es-PE', { hour12: false }),
    }
    logs.value.unshift(entrada)
    if (logs.value.length > MAX_LOGS) {
      logs.value.pop()
    }
  }

  function limpiar() {
    logs.value = []
  }

  function registrarOperacion(tipo, indice, ms, nodosVisitados, accesosDisco = 0) {
    const iconos = {
      buscar: '🔍',
      rango: '📋',
      insertar: '➕',
      eliminar: '🗑️',
      play: '▶️',
    }
    const icono = iconos[tipo] || '✅'
    agregar(
      `${icono} ${tipo.charAt(0).toUpperCase() + tipo.slice(1)} ID=${indice} — B+-Tree en RAM (${ms.toFixed(3)} ms). Nodos visitados: ${nodosVisitados}. Accesos a disco: ${accesosDisco}`,
      'success',
    )
  }

  return {
    logs: computed(() => logs.value),
    agregar,
    limpiar,
    registrarOperacion,
  }
}
