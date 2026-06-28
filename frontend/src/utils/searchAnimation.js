import { compararClave } from './treeLayout.js'
import { clavesVisuales } from './nodoClaves.js'
import { localizarEntradaEnHojas } from './leafLocate.js'

/**
 * Pasos de búsqueda exacta: navegación + escaneo celda a celda en la hoja.
 * idCancion: en vistas secundarias (TrackName, etc.) distingue duplicados al buscar por ID.
 */
export function calcularPasosBusquedaCompleta(estructura, clave, idCancion = null) {
  if (!estructura?.nodos?.length) return { pasos: [], encontrado: false }

  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const loc = localizarEntradaEnHojas(estructura, clave, idCancion)
  const pasos = [...loc.nav.slice(0, -1)]
  const { hojaId, claveNorm, startIdx, cadenaIdx } = loc
  const hoja = mapa[hojaId]
  const cadena = estructura.cadenaHojas ?? []
  const hojaNavId = loc.nav[loc.nav.length - 1]?.nodeId

  if (hojaNavId !== hojaId && startIdx != null && cadenaIdx > startIdx) {
    for (let ci = startIdx + 1; ci <= cadenaIdx; ci++) {
      pasos.push({
        nodeId: cadena[ci],
        tipo: 'sequence_set',
        mensaje: `↪ Sequence Set: clave ${claveNorm} duplicada — escaneo en hoja ${cadena[ci]}`,
        highlightedNodes: [cadena[ci - 1], cadena[ci]],
      })
    }
  }

  if (!hoja?.esHoja) {
    return { pasos, encontrado: false }
  }

  const claves = loc.claves?.length ? loc.claves : clavesVisuales(hoja)

  pasos.push({
    nodeId: hojaId,
    tipo: 'hoja_entrada',
    mensaje: `📄 Hoja alcanzada — regla B+: solo aquí vive el dato real (el índice es mapa de ruta)`,
    activeEdge: null,
    activeSlot: null,
    activeKey: null,
  })

  if (claves.length === 0) {
    pasos.push({
      nodeId: hojaId,
      tipo: 'no_encontrado',
      mensaje: `❌ Hoja vacía — clave ${claveNorm} no encontrada`,
      activeKey: claveNorm,
      encontrado: false,
    })
    return { pasos, encontrado: false }
  }

  if (!loc.encontrado) {
    for (let i = 0; i < claves.length; i++) {
      const k = claves[i]
      if (compararClave(k, claveNorm) >= 0) {
        pasos.push({
          nodeId: hojaId,
          tipo: 'no_encontrado',
          mensaje: `❌ ${k} ≥ ${claveNorm} — clave no existe en esta hoja`,
          activeKey: k,
          activeSlot: i,
          encontrado: false,
        })
        break
      }
      pasos.push({
        nodeId: hojaId,
        tipo: 'hoja_comparar',
        mensaje: `🔎 ${k} < ${claveNorm} → continuar escaneo lineal en la hoja`,
        activeKey: k,
        activeSlot: i,
      })
    }
    if (pasos[pasos.length - 1]?.tipo !== 'no_encontrado') {
      pasos.push({
        nodeId: hojaId,
        tipo: 'no_encontrado',
        mensaje: `❌ Clave ${claveNorm}${idCancion != null ? ` (ID ${idCancion})` : ''} no encontrada`,
        activeKey: claveNorm,
        encontrado: false,
      })
    }
    return { pasos, encontrado: false }
  }

  const pos = loc.pos
  for (let i = 0; i < pos; i++) {
    const idEtq =
      idCancion != null && hoja?.idsRegistro?.[i] != null ? ` [ID ${hoja.idsRegistro[i]}]` : ''
    pasos.push({
      nodeId: hojaId,
      tipo: 'hoja_comparar',
      mensaje: `🔎 ${claves[i]}${idEtq} < ${claveNorm} → continuar escaneo`,
      activeKey: claves[i],
      activeSlot: i,
    })
  }

  const idEtqObj =
    idCancion != null && hoja?.idsRegistro?.[pos] != null ? ` / ID ${hoja.idsRegistro[pos]}` : ''
  pasos.push({
    nodeId: hojaId,
    tipo: 'encontrado',
    mensaje: `✅ ¡Coincidencia! Clave ${claveNorm}${idEtqObj} en posición ${pos} de la hoja`,
    activeKey: claves[pos],
    activeSlot: pos,
    encontrado: true,
  })

  return { pasos, encontrado: true }
}

/** Secuencia de búsqueda: navegación + escaneo + paso de consulta al backend. */
export function calcularSecuenciaBusqueda(estructura, clave, idCancion = null) {
  const { pasos, encontrado } = calcularPasosBusquedaCompleta(estructura, clave, idCancion)
  const ultimoNodo = pasos[pasos.length - 1]?.nodeId ?? estructura?.raiz

  pasos.push({
    nodeId: ultimoNodo,
    tipo: 'aplicar',
    mensaje: `⚙️ Consultando backend — verificación en árbol B+ y BD…`,
  })

  return { pasos, encontradoPredicho: encontrado }
}

/** Confirmación tras la respuesta del API. */
export function calcularPasosConfirmacionBusqueda(indice, cancion, predicho, ms, errorMsg = null) {
  if (cancion) {
    return [
      {
        nodeId: null,
        tipo: 'confirmacion',
        mensaje: `✅ Encontrado: "${cancion.TrackName}" — ${cancion.Artists} (${ms.toFixed(1)} ms, O(log_d n))`,
        highlightedKeys: [indice],
        activeKey: indice,
      },
    ]
  }

  const msg = errorMsg
    ? `❌ Error inesperado: ${errorMsg}`
    : predicho
      ? `❌ Clave ${indice} predicha en árbol pero no en BD`
      : `❌ Clave ${indice} no encontrada — escaneo en hoja sin coincidencia`

  return [{ nodeId: null, tipo: 'confirmacion', mensaje: msg }]
}
