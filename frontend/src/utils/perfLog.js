/** Estima nodos internos visitados en búsqueda O(log n) con altura del árbol. */
export function estimarNodosInternos(altura, totalCanciones, orden = 2) {
  if (altura != null && altura > 0) {
    return Math.max(1, altura - 1)
  }
  if (totalCanciones <= 1) return 1
  const d = Math.max(2, orden)
  const h = Math.ceil(Math.log(totalCanciones) / Math.log(d))
  return Math.max(1, h - 1)
}

export function formatearMs(ms) {
  if (ms < 1) return `${(ms * 1000).toFixed(0)} microsegundos`
  if (ms < 1000) return `${ms.toFixed(1)} milisegundos`
  return `${(ms / 1000).toFixed(2)} segundos`
}

export function mensajePrefijo({ prefijo, total, ms, stats, orden, totalCanciones }) {
  const nodos = estimarNodosInternos(stats?.altura, totalCanciones, orden)
  const saltos = Math.max(0, total - 1)
  return (
    `✅ Autocompletado Prefix B+-Tree «${prefijo}» completado. ` +
    `Primer acceso O(log n) visitando ${nodos} nodo(s) interno(s) en árbolPorNombre. ` +
    `Recolección de ${total} canción(es) vía Sequence Set (costo O(1) por salto, ${saltos} salto(s)). ` +
    `Tiempo total: ${formatearMs(ms)}. Accesos a base de datos externa: 0.`
  )
}

export function mensajeRangoPopularidad({ inicio, fin, total, ms, stats, orden, totalCanciones }) {
  const nodos = estimarNodosInternos(stats?.altura, totalCanciones, orden)
  const saltos = Math.max(0, total - 1)
  return (
    `✅ Filtro de Popularidad (${inicio}–${fin}) completado. ` +
    `Primer acceso O(log n) visitando ${nodos} nodo(s) interno(s) en árbolPorPopularidad. ` +
    `Recolección de ${total} canción(es) vía Sequence Set (costo O(1) por salto, ${saltos} salto(s)). ` +
    `Tiempo total: ${formatearMs(ms)}. Accesos a base de datos externa: 0.`
  )
}

export function mensajeRangoTempo({ inicio, fin, total, ms, stats, orden, totalCanciones }) {
  const nodos = estimarNodosInternos(stats?.altura, totalCanciones, orden)
  const saltos = Math.max(0, total - 1)
  return (
    `✅ Filtro Tempo BPM (${inicio}–${fin}) completado. ` +
    `Primer acceso O(log n) visitando ${nodos} nodo(s) interno(s) en árbolPorTempo. ` +
    `Recolección de ${total} canción(es) vía Sequence Set (costo O(1) por salto, ${saltos} salto(s)). ` +
    `Tiempo total: ${formatearMs(ms)}. Accesos a base de datos externa: 0.`
  )
}

export function mensajeRangoNombre({ desde, hasta, total, ms, stats, orden, totalCanciones }) {
  const nodos = estimarNodosInternos(stats?.altura, totalCanciones, orden)
  const saltos = Math.max(0, total - 1)
  return (
    `✅ Range Scan TrackName [${desde}–${hasta}] completado (Prefix B+-Tree / Comer). ` +
    `Primer acceso O(log n) visitando ${nodos} nodo(s) interno(s) en árbolPorNombre. ` +
    `Recolección de ${total} canción(es) vía Sequence Set (costo O(1) por salto, ${saltos} salto(s)). ` +
    `Tiempo total: ${formatearMs(ms)}. Accesos a base de datos externa: 0.`
  )
}

export function mensajeBuscarId({ indice, ms, nodosVisitados }) {
  return (
    `✅ Búsqueda exacta por ID=${indice} completada. ` +
    `Recorrido O(log n) visitando ${nodosVisitados} nodo(s) en árbolPorIndice. ` +
    `Tiempo total: ${formatearMs(ms)}. Accesos a base de datos externa: 0.`
  )
}
