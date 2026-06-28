import { calcularPasosNavegacion, compararClave } from './treeLayout.js'

function enRango(k, inicio, fin, tipoClave) {
  if (tipoClave === 'string') {
    const a = String(k).localeCompare(String(inicio), undefined, { sensitivity: 'base' })
    const b = String(k).localeCompare(String(fin), undefined, { sensitivity: 'base' })
    return a >= 0 && b <= 0
  }
  const nk = Number(k)
  const ni = Number(inicio)
  const nf = Number(fin)
  return nk >= ni && nk <= nf
}

function antesDeFin(k, fin, tipoClave) {
  if (tipoClave === 'string') {
    return String(k).localeCompare(String(fin), undefined, { sensitivity: 'base' }) <= 0
  }
  return Number(k) <= Number(fin)
}

function despuesDeInicio(k, inicio, tipoClave) {
  if (tipoClave === 'string') {
    return String(k).localeCompare(String(inicio), undefined, { sensitivity: 'base' }) >= 0
  }
  return Number(k) >= Number(inicio)
}

/**
 * Pasos del Range Scan: búsqueda al inicio + recorrido horizontal O(1) por sequence set.
 * Soporta claves numéricas (ID) y string (TrackName).
 */
export function calcularPasosRangeScan(estructura, inicio, fin, tipoClave = 'int') {
  if (!estructura?.nodos?.length) return []
  if (compararClave(inicio, fin) > 0) return []

  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const pasos = []
  const etiqueta = tipoClave === 'string' ? 'TrackName' : 'ID'

  pasos.push({
    nodeId: estructura.raiz,
    tipo: 'range_intro',
    mensaje: `📋 Range Scan [${inicio}, ${fin}] (${etiqueta}): paso 1 — búsqueda O(log n) al límite inferior «${inicio}»`,
    activeEdge: null,
    activeSlot: null,
    activeKey: inicio,
  })

  const nav = calcularPasosNavegacion(estructura, inicio)
  pasos.push(...nav.slice(1))

  let hojaId = null
  for (let i = nav.length - 1; i >= 0; i--) {
    if (mapa[nav[i].nodeId]?.esHoja) {
      hojaId = nav[i].nodeId
      break
    }
  }
  if (hojaId == null) return pasos

  const cadena = estructura.cadenaHojas || []
  let idxCadena = cadena.indexOf(hojaId)

  pasos.push({
    nodeId: hojaId,
    tipo: 'range_hoja_inicial',
    mensaje: `📄 Hoja inicial localizada — recorrido horizontal del sequence set (Comer)`,
    activeKey: inicio,
  })

  while (hojaId != null && idxCadena >= 0) {
    const hoja = mapa[hojaId]
    if (!hoja?.esHoja) break

    let finEnHoja = false
    for (const k of hoja.indices || []) {
      if (!despuesDeInicio(k, inicio, tipoClave)) continue
      if (!enRango(k, inicio, fin, tipoClave)) {
        finEnHoja = true
        break
      }
      pasos.push({
        nodeId: hojaId,
        tipo: 'range_leer',
        mensaje: `→ Leyendo clave «${k}» en rango [${inicio}, ${fin}]`,
        activeKey: k,
        leafChainFrom: null,
        leafChainTo: null,
      })
    }

    if (finEnHoja) break
    const ultima = hoja.indices?.[hoja.indices.length - 1]
    if (ultima != null && !antesDeFin(ultima, fin, tipoClave)) break

    const nextIdx = idxCadena + 1
    if (nextIdx >= cadena.length) {
      pasos.push({
        nodeId: hojaId,
        tipo: 'range_fin',
        mensaje: `✅ Fin del sequence set — rango completado`,
        activeKey: null,
      })
      break
    }

    const nextHojaId = cadena[nextIdx]
    pasos.push({
      nodeId: hojaId,
      tipo: 'range_salto',
      mensaje: `↔ Salto O(1) por enlace horizontal (siguienteHoja) — sin recorrer el índice`,
      activeKey: null,
      leafChainFrom: hojaId,
      leafChainTo: nextHojaId,
    })

    hojaId = nextHojaId
    idxCadena = nextIdx
  }

  return pasos
}

/** Secuencia Range Scan + paso de consulta al backend. */
export function calcularSecuenciaRangeScan(estructura, inicio, fin, tipoClave = 'int') {
  const pasos = calcularPasosRangeScan(estructura, inicio, fin, tipoClave)
  pasos.push({
    nodeId: estructura?.raiz,
    tipo: 'aplicar',
    mensaje: `⚙️ Ejecutando Range Scan en backend [${inicio}, ${fin}]…`,
  })
  return { pasos, inicio, fin }
}

/** Confirmación con resultados del rango. */
export function calcularPasosConfirmacionRango(inicio, fin, resultado, ms, tipoClave = 'int') {
  const keys =
    tipoClave === 'string'
      ? (resultado?.canciones?.map((c) => c.TrackName) ?? [])
      : (resultado?.canciones?.map((c) => c.Indice) ?? [])
  return [
    {
      nodeId: null,
      tipo: 'confirmacion',
      mensaje: `✅ Range Scan [${inicio}, ${fin}]: ${resultado?.total ?? 0} canciones en ${ms.toFixed(1)} ms (O(log_d n + k), sequence set)`,
      highlightedKeys: keys,
    },
  ]
}
