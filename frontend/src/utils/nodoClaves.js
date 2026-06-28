/** Claves visibles de un nodo exportado (int, string o float como texto). */
export function clavesVisuales(nodo) {
  if (nodo?.claves?.length) return nodo.claves
  if (nodo?.indices?.length) return nodo.indices.map(String)
  return []
}

export function separadoresVisuales(nodo) {
  return nodo?.separadores ?? []
}

/** Para animaciones numéricas (solo índice por id). */
export function indicesNumericos(nodo) {
  return nodo?.indices ?? []
}
