const BASE = '/api'

function mensajeErrorRed(status, data) {
  if (status === 404) {
    return (
      data.error ||
      'Backend no encontrado (HTTP 404). Abre otra terminal en la carpeta del proyecto y ejecuta: go run .'
    )
  }
  if (status === 502 || status === 503) {
    return 'Backend Go no responde. Ejecuta en otra terminal: go run . (puerto 8080)'
  }
  return data.error || `Error HTTP ${status}`
}

async function request(url, options = {}) {
  let response
  try {
    response = await fetch(`${BASE}${url}`, {
      headers: { 'Content-Type': 'application/json', ...options.headers },
      ...options,
    })
  } catch {
    throw new Error(
      'No se puede conectar al backend Go (puerto 8080). En otra terminal: cd BPlusTree_Proyecto_Final_Algoritmos && go run .',
    )
  }

  const data = await response.json().catch(() => ({}))

  if (!response.ok) {
    throw new Error(mensajeErrorRed(response.status, data))
  }

  return data
}

export function verificarBackend() {
  return request('/health')
}

export function inicializar(orden, limite) {
  return request('/inicializar', {
    method: 'POST',
    body: JSON.stringify({ orden, limite }),
  })
}

/** Demo: árbol vacío + lista de canciones (máx. 50). */
export function inicializarDemo(orden, limite) {
  return request('/inicializar-demo', {
    method: 'POST',
    body: JSON.stringify({ orden, limite }),
  })
}

/** Inserta en el árbol B+ sin tocar la BD (solo demo de inicialización). */
export function insertarEnArbol(cancion) {
  return request('/insertar-arbol', {
    method: 'POST',
    body: JSON.stringify(cancion),
  })
}

export function obtenerEstructura(tipo = 'indice') {
  const q = tipo && tipo !== 'indice' ? `?tipo=${encodeURIComponent(tipo)}` : ''
  return request(`/estructura${q}`)
}

export function obtenerEstadisticas(tipo = 'indice') {
  const q = tipo && tipo !== 'indice' ? `?tipo=${encodeURIComponent(tipo)}` : ''
  return request(`/estadisticas${q}`)
}

export function obtenerConteoBD() {
  return request('/conteo-bd')
}

export function buscar(indice) {
  return request(`/buscar?indice=${indice}`)
}

export function buscarRango(inicio, fin) {
  return request(`/rango?inicio=${inicio}&fin=${fin}`)
}

/** Range scan unificado: campo = indice | popularidad | tempo | danceability | nombre */
export function buscarRangoCampo(campo, inicio, fin) {
  if (campo === 'nombre') {
    return request(
      `/rango?campo=nombre&inicio=${encodeURIComponent(inicio)}&fin=${encodeURIComponent(fin)}`,
    )
  }
  return request(
    `/rango?campo=${encodeURIComponent(campo)}&inicio=${inicio}&fin=${fin}`,
  )
}

export function buscarRangoNombre(desde, hasta) {
  return request(
    `/rango-nombre?desde=${encodeURIComponent(desde)}&hasta=${encodeURIComponent(hasta)}`,
  )
}

export function buscarPorNombre(nombre) {
  return request(`/buscar-nombre?nombre=${encodeURIComponent(nombre)}`)
}

export function buscarPorPrefijo(prefijo) {
  return request(`/buscar-prefijo?prefijo=${encodeURIComponent(prefijo)}`)
}

export function buscarTempoExacto(tempo) {
  return request(`/buscar-tempo?tempo=${tempo}`)
}

export function buscarRangoPopularidad(inicio, fin) {
  return request(`/rango-popularidad?inicio=${inicio}&fin=${fin}`)
}

export function buscarRangoTempo(inicio, fin) {
  return request(`/rango-tempo?inicio=${inicio}&fin=${fin}`)
}

export function buscarRangoDanceability(inicio, fin) {
  return request(`/rango-danceability?inicio=${inicio}&fin=${fin}`)
}

export function obtenerIndices() {
  return request('/indices')
}

export function obtenerSiguienteIndice() {
  return request('/siguiente-indice')
}

export function insertar(cancion) {
  return request('/insertar', {
    method: 'POST',
    body: JSON.stringify(cancion),
  })
}

export function eliminar(indice) {
  return request(`/eliminar?indice=${indice}`, { method: 'DELETE' })
}
