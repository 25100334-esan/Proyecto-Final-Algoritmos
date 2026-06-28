/** Máximo de canciones para dibujar el árbol en canvas (SVG). Por encima: solo operaciones API. */
export const LIMITE_VISUALIZACION = 500

export function visualizacionPermitida(totalCanciones) {
  return totalCanciones > 0 && totalCanciones <= LIMITE_VISUALIZACION
}

/**
 * ¿Mostrar canvas SVG?
 * - Con canciones en RAM: total ≤ 500
 * - Demo animada (árbol vacío pero inicializado): limite planificado ≤ 500
 */
export function debeMostrarCanvas(totalCanciones, limitePlanificado = 0, arbolInicializado = false) {
  if (totalCanciones > LIMITE_VISUALIZACION) return false
  if (totalCanciones > 0) return totalCanciones <= LIMITE_VISUALIZACION
  return (
    arbolInicializado &&
    limitePlanificado > 0 &&
    limitePlanificado <= LIMITE_VISUALIZACION
  )
}
