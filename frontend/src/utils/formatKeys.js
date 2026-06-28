/** Formatea claves de hoja para mensajes (hoja vacía → ∅). */
export function clavesDeNodo(nodo) {
  if (nodo?.indices?.length) return nodo.indices
  const vis = nodo?.claves
  if (!vis?.length) return []
  return vis.map((c) => {
    const n = Number(c)
    return Number.isNaN(n) ? c : n
  })
}

export function formatearClaves(keys) {
  const lista = keys ?? []
  return lista.length ? lista.join(', ') : '∅'
}
