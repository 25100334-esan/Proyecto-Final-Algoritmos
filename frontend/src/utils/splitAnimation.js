import { calcularPasosNavegacion, calcularLayout } from './treeLayout.js'
import { clavesDeNodo } from './formatKeys.js'

/**
 * Analiza qué ocurrió en un SPLIT comparando árbol antes/después.
 */
export function analizarSplit(anterior, nueva, claveInsertada, orden) {
  if (!anterior?.nodos?.length || !nueva?.nodos?.length) {
    return { huboSplit: false }
  }

  const maxPorNodo = Math.max(orden * 2, 2) // Comer: máximo 2d claves
  const mapaAnt = Object.fromEntries(anterior.nodos.map((n) => [n.id, n]))
  const mapaNue = Object.fromEntries(nueva.nodos.map((n) => [n.id, n]))

  const pasos = calcularPasosNavegacion(anterior, claveInsertada)
  let hojaAntesId = null
  for (let i = pasos.length - 1; i >= 0; i--) {
    if (mapaAnt[pasos[i].nodeId]?.esHoja) {
      hojaAntesId = pasos[i].nodeId
      break
    }
  }
  if (hojaAntesId == null) return { huboSplit: false }

  const hojaAntes = mapaAnt[hojaAntesId]
  const keysConInsert = [...new Set([...clavesDeNodo(hojaAntes), claveInsertada])].sort((a, b) => a - b)

  const hojasAntes = anterior.cadenaHojas?.length ?? 0
  const hojasDespues = nueva.cadenaHojas?.length ?? 0
  const nodosAntes = anterior.nodos.length
  const nodosDespues = nueva.nodos.length

  const huboSplit = hojasDespues > hojasAntes || nodosDespues > nodosAntes
  if (!huboSplit) return { huboSplit: false }

  const medio = Math.floor((keysConInsert.length + 1) / 2)
  const izqKeys = keysConInsert.slice(0, medio)
  const derKeys = keysConInsert.slice(medio)
  const clavePromovida = derKeys[0]

  let hojaIzqId = null
  let hojaDerId = null
  for (const id of nueva.cadenaHojas || []) {
    const h = mapaNue[id]
    if (!h?.esHoja) continue
    if (izqKeys.every((k) => h.indices.includes(k))) hojaIzqId = id
    if (derKeys.every((k) => h.indices.includes(k))) hojaDerId = id
  }
  if (hojaIzqId == null) {
    for (const id of nueva.cadenaHojas || []) {
      if (mapaNue[id]?.indices?.includes(izqKeys[0])) hojaIzqId = id
    }
  }
  if (hojaDerId == null) {
    for (const id of nueva.cadenaHojas || []) {
      if (mapaNue[id]?.indices?.includes(derKeys[0])) hojaDerId = id
    }
  }

  let padreId = null
  for (const n of nueva.nodos) {
    if (!n.esHoja && n.separadores?.includes(clavePromovida)) {
      padreId = n.id
      break
    }
  }
  if (padreId == null) padreId = nueva.raiz

  const layoutNue = calcularLayout(nueva)

  return {
    huboSplit: true,
    esHoja: true,
    maxPorNodo,
    claveInsertada,
    keysConInsert,
    izqKeys,
    derKeys,
    clavePromovida,
    medio,
    hojaIzqId,
    hojaDerId,
    padreId,
    hojaIzqPos: layoutNue.posiciones[hojaIzqId],
    hojaDerPos: layoutNue.posiciones[hojaDerId],
    hojaIzqDim: layoutNue.dimensiones[hojaIzqId],
    hojaDerDim: layoutNue.dimensiones[hojaDerId],
    padrePos: layoutNue.posiciones[padreId],
    padreDim: layoutNue.dimensiones[padreId],
  }
}

export const FASES_SPLIT = {
  SPLITTING: 'splitting',
  PROMOTE: 'promote',
}
