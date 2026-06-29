/** Clave de indexación según la vista activa del canvas académico. */
export function truncarClaveVisual(s) {
  const chars = [...String(s ?? '')]
  if (chars.length <= 14) return chars.join('')
  return `${chars.slice(0, 12).join('')}…`
}

export function claveDeCancionPorVista(cancion, tipoVista) {
  switch (tipoVista) {
    case 'nombre':
    case 'trackname':
    case 'name':
      return cancion.TrackName ?? cancion.trackName ?? ''
    case 'popularidad':
    case 'popularity':
      return cancion.Popularity ?? cancion.popularity ?? 0
    case 'tempo':
      return cancion.Tempo ?? cancion.tempo ?? 0
    case 'danceability':
    case 'bailabilidad':
      return cancion.Danceability ?? cancion.danceability ?? 0
    default:
      return cancion.Indice ?? cancion.indice ?? 0
  }
}

/** Tipo de vista inferido del árbol exportado (fuente de verdad para animaciones). */
export function tipoVistaDesdeEstructura(estructura) {
  const t = (estructura?.tipoIndice ?? 'indice').toLowerCase()
  if (t === 'id') return 'indice'
  if (t === 'trackname' || t === 'name') return 'nombre'
  if (t === 'popularity') return 'popularidad'
  if (t === 'bailabilidad') return 'danceability'
  return t
}

/** Clave según el árbol que se está visualizando (evita desincronía con el selector). */
export function claveDeCancionPorEstructura(cancion, estructura) {
  return claveDeCancionPorVista(cancion, tipoVistaDesdeEstructura(estructura))
}

/** Normaliza clave al formato del JSON exportado (truncado string, float64, int). */
export function normalizarClaveArbol(estructura, clave) {
  const tipo = estructura?.tipoClave ?? 'int'
  if (tipo === 'string') return truncarClaveVisual(clave)
  if (tipo === 'float64') {
    const n = Number(clave)
    if (Number.isNaN(n)) return 0
    return Number.parseFloat(n.toPrecision(12))
  }
  return Number(clave)
}

/** Índice primario: claves únicas. Secundarios: duplicados + sequence set (Comer). */
export function clavesUnicasEnEstructura(estructura) {
  if (estructura?.clavesUnicas != null) return Boolean(estructura.clavesUnicas)
  const t = (estructura?.tipoIndice ?? 'indice').toLowerCase()
  return t === 'indice' || t === 'id'
}

export function etiquetaClaveVista(tipoVista) {
  const map = {
    indice: 'ID',
    nombre: 'TrackName',
    popularidad: 'Popularidad',
    tempo: 'Tempo',
    danceability: 'Danceability',
    bailabilidad: 'Danceability',
  }
  return map[tipoVista] ?? map[tipoVistaDesdeEstructura({ tipoIndice: tipoVista })] ?? 'ID'
}
