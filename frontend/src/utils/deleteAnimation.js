import { calcularPasosNavegacion, calcularLayout, estimarProfundidad, compararClave } from './treeLayout.js'
import { clavesVisuales } from './nodoClaves.js'
import { normalizarClaveArbol } from './claveVista.js'
import { localizarEntradaEnHojas } from './leafLocate.js'

export const FASES_DELETE = {
  BORRADO: 'borrado',
  FANTASMA: 'fantasma',
  UNDERFLOW: 'underflow',
  EVALUAR: 'evaluar',
  PRESTAMO: 'prestamo',
  FUSION: 'fusion',
  RAIZ: 'raiz',
}

function normalizarClave(estructura, clave) {
  return normalizarClaveArbol(estructura, clave)
}

function clavesEnHoja(hoja) {
  return clavesVisuales(hoja)
}

function ordenarClaves(arr) {
  return [...arr].sort((a, b) => compararClave(a, b))
}

function findHojaPorClaves(estructura, claves) {
  if (!estructura?.cadenaHojas?.length || !claves?.length) return null
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const sorted = ordenarClaves(claves)
  for (const id of estructura.cadenaHojas) {
    const h = mapa[id]
    if (!h?.esHoja) continue
    const hKeys = ordenarClaves(clavesEnHoja(h))
    if (hKeys.length === sorted.length && hKeys.every((k, i) => compararClave(k, sorted[i]) === 0)) {
      return id
    }
  }
  return null
}

function quitarUnaClave(arr, clave) {
  const copy = [...arr]
  const idx = copy.findIndex((k) => compararClave(k, clave) === 0)
  if (idx >= 0) copy.splice(idx, 1)
  return copy
}

function findPadreDeHoja(estructura, hojaId) {
  for (const n of estructura?.nodos || []) {
    if (!n.esHoja && n.hijos?.includes(hojaId)) return n
  }
  return null
}

function hojaOrigen(anterior, clave, idCancion = null) {
  const loc = localizarEntradaEnHojas(anterior, clave, idCancion)
  if (!loc.encontrado) return { hojaId: null, hoja: null, pos: -1 }
  const mapa = Object.fromEntries(anterior.nodos.map((n) => [n.id, n]))
  return { hojaId: loc.hojaId, hoja: mapa[loc.hojaId], pos: loc.pos }
}

export function estructuraConBorradoVirtual(estructura, hojaId, claveEliminada, posicion = null) {
  const clone = JSON.parse(JSON.stringify(estructura))
  const nodo = clone.nodos.find((n) => n.id === hojaId)
  if (!nodo) return clone
  const claves = clavesEnHoja(nodo)
  const pos =
    posicion ??
    claves.findIndex((k) => compararClave(k, normalizarClave(estructura, claveEliminada)) === 0)
  if (pos < 0) return clone
  if (nodo.claves?.length) nodo.claves.splice(pos, 1)
  if (nodo.indices?.length) nodo.indices.splice(pos, 1)
  if (nodo.idsRegistro?.length) nodo.idsRegistro.splice(pos, 1)
  return clone
}

function buscarSeparadorFantasma(estructura, clave) {
  for (const n of estructura?.nodos || []) {
    if (!n.esHoja && n.separadores?.includes(clave)) {
      return { padreId: n.id, separadorFantasma: clave }
    }
  }
  return { padreId: null, separadorFantasma: null }
}

function layoutInfo(estructura) {
  const layout = calcularLayout(estructura)
  const pos = (id) => (id != null ? layout.posiciones[id] : null)
  const dim = (id) => (id != null ? layout.dimensiones[id] : null)
  return { layout, pos, dim }
}

function siguienteEnCadena(estructura, hojaId) {
  const idx = estructura.cadenaHojas?.indexOf(hojaId) ?? -1
  if (idx < 0 || idx >= estructura.cadenaHojas.length - 1) return hojaId
  return estructura.cadenaHojas[idx + 1]
}

/**
 * Simula corregirUnderflowHoja sobre el árbol ANTES del API (misma lógica que Go).
 * Genera pasos narrativos en tiempo presente: "revisando hermano…", "no puede prestar…".
 */
function simularCorreccionUnderflow(estructura, hojaId, keysTras, claveEliminada, orden) {
  const d = orden
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const hoja = mapa[hojaId]
  const esRaizUnica = hojaId === estructura.raiz && hoja?.esHoja
  const { pos, dim } = layoutInfo(estructura)

  const baseInfo = {
    claveEliminada,
    keysTras,
    minD: d,
    hojaId,
    hojaPos: pos(hojaId),
  }

  if (esRaizUnica || keysTras.length >= d) {
    const { padreId, separadorFantasma } = buscarSeparadorFantasma(estructura, claveEliminada)
    return {
      tipo: 'fantasma',
      deleteInfo: {
        ...baseInfo,
        tipo: 'fantasma',
        padreId,
        separadorFantasma,
        keysRestantes: keysTras,
        padrePos: pos(padreId),
        padreDim: dim(padreId),
      },
      pasosCorreccion: [],
    }
  }

  const padre = findPadreDeHoja(estructura, hojaId)
  if (!padre) {
    return { tipo: 'ninguno', deleteInfo: { ...baseInfo, tipo: 'ninguno' }, pasosCorreccion: [] }
  }

  const posHijo = padre.hijos.indexOf(hojaId)
  const pasos = []

  pasos.push({
    nodeId: hojaId,
    tipo: 'corregir_entrada',
    mensaje: `📋 corregirUnderflowHoja(${keysTras.length} claves < d=${d}) — evaluando hermanos en orden Comer`,
    highlightedNodes: [hojaId, padre.id],
    underflowNodes: [hojaId],
    deletePhase: FASES_DELETE.UNDERFLOW,
  })

  // 1) Hermano derecho
  if (posHijo + 1 < padre.hijos.length) {
    const hermanoDerId = padre.hijos[posHijo + 1]
    const hermanoDer = mapa[hermanoDerId]
    const keysDer = clavesEnHoja(hermanoDer)

    pasos.push({
      nodeId: hermanoDerId,
      tipo: 'evaluar_der',
      mensaje: `🔎 Revisando hermano DERECHO [${keysDer.join(', ')}]: ¿tiene > d=${d} claves para prestar?`,
      highlightedNodes: [hojaId, hermanoDerId, padre.id],
      highlightedKeys: keysDer,
      hermanoEvaluado: hermanoDerId,
      deletePhase: FASES_DELETE.EVALUAR,
      evaluarHermano: 'derecha',
    })

    if (keysDer.length > d) {
      const clavePrestada = keysDer[0]
      const nuevoSeparador = keysDer[1]
      const keysHojaFinal = ordenarClaves([...keysTras, clavePrestada])
      const keysHermanoFinal = keysDer.slice(1)
      const sepAnterior = padre.separadores[posHijo]

      pasos.push({
        nodeId: hermanoDerId,
        tipo: 'prestamo_decision',
        mensaje: `✓ Hermano derecho tiene ${keysDer.length} > d=${d} → prestarCancionDerecha(): mover clave ${clavePrestada}`,
        highlightedNodes: [hojaId, hermanoDerId, padre.id],
        highlightedKeys: [clavePrestada],
        deletePhase: FASES_DELETE.PRESTAMO,
        animarPrestamo: true,
      })

      pasos.push({
        nodeId: hojaId,
        tipo: 'prestamo_resultado',
        mensaje: `↔ Tras préstamo — hoja: [${keysHojaFinal.join(', ')}] | hermano: [${keysHermanoFinal.join(', ')}]`,
        highlightedNodes: [hojaId, hermanoDerId],
        highlightedKeys: keysHojaFinal,
        underflowNodes: [],
      })

      pasos.push({
        nodeId: padre.id,
        tipo: 'prestamo_separador',
        mensaje: `↑ Actualizar separador padre: ${sepAnterior} → ${nuevoSeparador}`,
        highlightedNodes: [padre.id, hojaId, hermanoDerId],
        activeKey: nuevoSeparador,
        deletePhase: FASES_DELETE.PRESTAMO,
      })

      return {
        tipo: 'prestamo',
        deleteInfo: {
          ...baseInfo,
          tipo: 'prestamo',
          direccion: 'derecha',
          hermanoId: hermanoDerId,
          clavePrestada,
          nuevoSeparador,
          separadorAnterior: sepAnterior,
          padreId: padre.id,
          keysTrasBorrado: keysTras,
          keysHojaFinal,
          keysHermanoFinal,
          keysHermanoAntes: keysDer,
          hojaPos: pos(hojaId),
          hermanoPos: pos(hermanoDerId),
          padrePos: pos(padre.id),
          padreDim: dim(padre.id),
        },
        pasosCorreccion: pasos,
      }
    }

    pasos.push({
      nodeId: hermanoDerId,
      tipo: 'evaluar_der_no',
      mensaje: `✗ Hermano derecho tiene ${keysDer.length} claves (= d=${d}, mínimo) — NO puede prestar`,
      highlightedNodes: [hojaId, hermanoDerId],
      hermanoRechazado: hermanoDerId,
      deletePhase: FASES_DELETE.EVALUAR,
    })
  }

  // 2) Hermano izquierdo
  if (posHijo > 0) {
    const hermanoIzqId = padre.hijos[posHijo - 1]
    const hermanoIzq = mapa[hermanoIzqId]
    const keysIzq = clavesEnHoja(hermanoIzq)

    pasos.push({
      nodeId: hermanoIzqId,
      tipo: 'evaluar_izq',
      mensaje: `🔎 Revisando hermano IZQUIERDO [${keysIzq.join(', ')}]: ¿tiene > d=${d} claves para prestar?`,
      highlightedNodes: [hojaId, hermanoIzqId, padre.id],
      highlightedKeys: keysIzq,
      hermanoEvaluado: hermanoIzqId,
      deletePhase: FASES_DELETE.EVALUAR,
      evaluarHermano: 'izquierda',
    })

    if (keysIzq.length > d) {
      const clavePrestada = keysIzq[keysIzq.length - 1]
      const nuevoSeparador = clavePrestada
      const keysHojaFinal = ordenarClaves([clavePrestada, ...keysTras])
      const keysHermanoFinal = keysIzq.slice(0, -1)
      const sepAnterior = padre.separadores[posHijo - 1]

      pasos.push({
        nodeId: hermanoIzqId,
        tipo: 'prestamo_decision',
        mensaje: `✓ Hermano izquierdo tiene ${keysIzq.length} > d=${d} → prestarCancionIzquierda(): mover clave ${clavePrestada}`,
        highlightedNodes: [hojaId, hermanoIzqId, padre.id],
        highlightedKeys: [clavePrestada],
        deletePhase: FASES_DELETE.PRESTAMO,
        animarPrestamo: true,
      })

      pasos.push({
        nodeId: hojaId,
        tipo: 'prestamo_resultado',
        mensaje: `↔ Tras préstamo — hoja: [${keysHojaFinal.join(', ')}] | hermano: [${keysHermanoFinal.join(', ')}]`,
        highlightedNodes: [hojaId, hermanoIzqId],
        highlightedKeys: keysHojaFinal,
        underflowNodes: [],
      })

      pasos.push({
        nodeId: padre.id,
        tipo: 'prestamo_separador',
        mensaje: `↑ Actualizar separador padre: ${sepAnterior} → ${nuevoSeparador}`,
        highlightedNodes: [padre.id, hojaId, hermanoIzqId],
        activeKey: nuevoSeparador,
        deletePhase: FASES_DELETE.PRESTAMO,
      })

      return {
        tipo: 'prestamo',
        deleteInfo: {
          ...baseInfo,
          tipo: 'prestamo',
          direccion: 'izquierda',
          hermanoId: hermanoIzqId,
          clavePrestada,
          nuevoSeparador,
          separadorAnterior: sepAnterior,
          padreId: padre.id,
          keysTrasBorrado: keysTras,
          keysHojaFinal,
          keysHermanoFinal,
          keysHermanoAntes: keysIzq,
          hojaPos: pos(hojaId),
          hermanoPos: pos(hermanoIzqId),
          padrePos: pos(padre.id),
          padreDim: dim(padre.id),
        },
        pasosCorreccion: pasos,
      }
    }

    pasos.push({
      nodeId: hermanoIzqId,
      tipo: 'evaluar_izq_no',
      mensaje: `✗ Hermano izquierdo tiene ${keysIzq.length} claves (= d=${d}, mínimo) — NO puede prestar`,
      highlightedNodes: [hojaId, hermanoIzqId],
      hermanoRechazado: hermanoIzqId,
      deletePhase: FASES_DELETE.EVALUAR,
    })
  }

  // 3) Fusión — misma prioridad que Go: derecha si existe, si no izquierda
  if (posHijo + 1 < padre.hijos.length) {
    const hermanoDerId = padre.hijos[posHijo + 1]
    const hermanoDer = mapa[hermanoDerId]
    const keysDer = [...clavesEnHoja(hermanoDer)]
    const keysIzq = keysTras
    const keysFusionadas = ordenarClaves([...keysIzq, ...keysDer])
    const separadorEliminado = padre.separadores[posHijo]
    const sepTras = padre.separadores.length - 1
    const esRaizPadre = padre.id === estructura.raiz

    pasos.push({
      nodeId: hojaId,
      tipo: 'fusion_decision',
      mensaje: `⚡ Ningún hermano puede prestar → fusionarHojaConDerecha() (concatenar hoja + hermano derecho)`,
      highlightedNodes: [hojaId, hermanoDerId, padre.id],
      underflowNodes: [hojaId],
      deletePhase: FASES_DELETE.FUSION,
    })

    pasos.push({
      nodeId: hojaId,
      tipo: 'fusion_concatenar',
      mensaje: `🔗 Concatenar: [${keysIzq.join(', ')}] + [${keysDer.join(', ')}] → [${keysFusionadas.join(', ')}]`,
      highlightedNodes: [hojaId, hermanoDerId],
      highlightedKeys: keysFusionadas,
      deletePhase: FASES_DELETE.FUSION,
      animarFusion: true,
      leafChainFrom: hojaId,
      leafChainTo: hermanoDerId,
      hojaAbsorbida: hermanoDerId,
    })

    pasos.push({
      nodeId: hojaId,
      tipo: 'fusion_sequence',
      mensaje: `↪ Re-enlazar sequence set: hoja ${hermanoDerId} desaparece del enlace horizontal (siguienteHoja)`,
      highlightedNodes: [hojaId, hermanoDerId],
      leafChainFrom: hojaId,
      leafChainTo: siguienteEnCadena(estructura, hermanoDerId),
      deletePhase: FASES_DELETE.FUSION,
      atenuarIds: [hermanoDerId],
    })

    pasos.push({
      nodeId: padre.id,
      tipo: 'fusion_separador',
      mensaje: `✂ Eliminar separador ${separadorEliminado} del padre — ya no hay dos hojas que dividir`,
      highlightedNodes: [padre.id, hojaId],
      activeKey: separadorEliminado,
      deletePhase: FASES_DELETE.FUSION,
    })

    if (sepTras < d && !esRaizPadre) {
      pasos.push({
        nodeId: padre.id,
        tipo: 'propagacion_preview',
        mensaje: `⬆ El padre quedará con ${sepTras} separador(es) (< d=${d}) → corregirUnderflowInterno (propagación)`,
        highlightedNodes: [padre.id],
        deletePhase: FASES_DELETE.UNDERFLOW,
      })
    } else if (esRaizPadre && sepTras === 0 && padre.hijos.length - 1 === 1) {
      pasos.push({
        nodeId: padre.id,
        tipo: 'raiz_preview',
        mensaje: `⬇ La raíz quedará con un solo hijo → ajustarRaiz() reducirá la altura del árbol`,
        highlightedNodes: [padre.id],
        deletePhase: FASES_DELETE.RAIZ,
      })
    }

    return {
      tipo: 'fusion',
      deleteInfo: {
        ...baseInfo,
        tipo: 'fusion',
        direccion: 'derecha',
        hermanoId: hermanoDerId,
        hojaFusionadaId: hojaId,
        keysHojaIzqAntes: keysIzq,
        keysHojaDerAntes: keysDer,
        keysFusionadas,
        separadorEliminado,
        padreId: padre.id,
        mergedPos: pos(hojaId),
        mergedDim: dim(hojaId),
        padrePos: pos(padre.id),
        hojaAbsorbida: hermanoDerId,
      },
      pasosCorreccion: pasos,
    }
  }

  // Fusión con izquierda
  const hermanoIzqId = padre.hijos[posHijo - 1]
  const hermanoIzq = mapa[hermanoIzqId]
  const keysIzqAntes = [...clavesEnHoja(hermanoIzq)]
  const keysDerAntes = keysTras
  const keysFusionadas = [...keysIzqAntes, ...keysDerAntes].sort((a, b) => a - b)
  const separadorEliminado = padre.separadores[posHijo - 1]
  const sepTras = padre.separadores.length - 1
  const esRaizPadre = padre.id === estructura.raiz

  pasos.push({
    nodeId: hermanoIzqId,
    tipo: 'fusion_decision',
    mensaje: `⚡ Sin hermano derecho — fusionarHojaConIzquierda() (absorber en hermano izquierdo)`,
    highlightedNodes: [hojaId, hermanoIzqId, padre.id],
    underflowNodes: [hojaId],
    deletePhase: FASES_DELETE.FUSION,
  })

  pasos.push({
    nodeId: hermanoIzqId,
    tipo: 'fusion_concatenar',
    mensaje: `🔗 Concatenar: [${keysIzqAntes.join(', ')}] + [${keysDerAntes.join(', ')}] → [${keysFusionadas.join(', ')}]`,
    highlightedNodes: [hojaId, hermanoIzqId],
    highlightedKeys: keysFusionadas,
    deletePhase: FASES_DELETE.FUSION,
    animarFusion: true,
    leafChainFrom: hermanoIzqId,
    leafChainTo: siguienteEnCadena(estructura, hojaId),
    hojaAbsorbida: hojaId,
  })

  pasos.push({
    nodeId: hermanoIzqId,
    tipo: 'fusion_sequence',
    mensaje: `↪ Re-enlazar sequence set: hoja ${hojaId} desaparece del enlace horizontal`,
    highlightedNodes: [hermanoIzqId, hojaId],
    leafChainFrom: hermanoIzqId,
    leafChainTo: siguienteEnCadena(estructura, hojaId),
    deletePhase: FASES_DELETE.FUSION,
    atenuarIds: [hojaId],
  })

  pasos.push({
    nodeId: padre.id,
    tipo: 'fusion_separador',
    mensaje: `✂ Eliminar separador ${separadorEliminado} del padre`,
    highlightedNodes: [padre.id, hermanoIzqId],
    activeKey: separadorEliminado,
    deletePhase: FASES_DELETE.FUSION,
  })

  if (sepTras < d && !esRaizPadre) {
    pasos.push({
      nodeId: padre.id,
      tipo: 'propagacion_preview',
      mensaje: `⬆ Propagación al padre: corregirUnderflowInterno`,
      highlightedNodes: [padre.id],
      deletePhase: FASES_DELETE.UNDERFLOW,
    })
  } else if (esRaizPadre && sepTras === 0 && padre.hijos.length - 1 === 1) {
    pasos.push({
      nodeId: padre.id,
      tipo: 'raiz_preview',
      mensaje: `⬇ Posible ajustarRaiz() — altura −1`,
      highlightedNodes: [padre.id],
      deletePhase: FASES_DELETE.RAIZ,
    })
  }

  return {
    tipo: 'fusion',
    deleteInfo: {
      ...baseInfo,
      tipo: 'fusion',
      direccion: 'izquierda',
      hermanoId: hermanoIzqId,
      hojaFusionadaId: hermanoIzqId,
      keysHojaIzqAntes: keysIzqAntes,
      keysHojaDerAntes: keysDerAntes,
      keysFusionadas,
      separadorEliminado,
      padreId: padre.id,
      mergedPos: pos(hermanoIzqId),
      mergedDim: dim(hermanoIzqId),
      padrePos: pos(padre.id),
      hojaAbsorbida: hojaId,
    },
    pasosCorreccion: pasos,
  }
}

/** Pasos de navegación + borrado virtual + corrección simulada (ANTES del API). */
export function calcularSecuenciaEliminacion(estructura, clave, orden, idCancion = null) {
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const loc = localizarEntradaEnHojas(estructura, clave, idCancion)
  const pasos = [...loc.nav.slice(0, -1)]
  const minD = orden

  if (!loc.encontrado) {
    const hoja = mapa[loc.hojaId]
    if (hoja?.esHoja && loc.claves?.length) {
      pasos.push({
        nodeId: loc.hojaId,
        tipo: 'hoja_entrada',
        mensaje: `📄 Hoja destino — clave ${loc.claveNorm} no encontrada`,
      })
    }
    pasos.push({
      nodeId: loc.hojaId,
      tipo: 'no_encontrado',
      mensaje: `❌ Clave ${loc.claveNorm}${idCancion != null ? ` (ID ${idCancion})` : ''} no encontrada en el árbol`,
      activeKey: loc.claveNorm,
    })
    return { pasos, encontrado: false, deleteInfo: null }
  }

  const { hojaId, pos, claves, claveNorm, startIdx, cadenaIdx } = loc
  const hoja = mapa[hojaId]
  const hojaNavId = loc.nav[loc.nav.length - 1]?.nodeId
  const cadena = estructura.cadenaHojas ?? []

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

  pasos.push({
    nodeId: hojaId,
    tipo: 'hoja_entrada',
    mensaje: `📄 Hoja destino — borrado físico SOLO aquí${idCancion != null ? ` (registro ID ${idCancion})` : ''}`,
  })

  for (let i = 0; i < pos; i++) {
    const idEtq =
      idCancion != null && hoja?.idsRegistro?.[i] != null ? ` [ID ${hoja.idsRegistro[i]}]` : ''
    pasos.push({
      nodeId: hojaId,
      tipo: 'hoja_comparar',
      mensaje: `${claves[i]}${idEtq} — escaneo hasta el registro correcto (r=(k,a))`,
      activeKey: claves[i],
      activeSlot: i,
    })
  }

  const claveObjetivo = claves[pos]
  pasos.push({
    nodeId: hojaId,
    tipo: 'objetivo',
    mensaje: `🎯 Clave ${claveNorm}${idCancion != null ? ` / ID ${idCancion}` : ''} localizada — preparando borrado físico`,
    activeKey: claveObjetivo,
    activeSlot: pos,
  })

  const keysTras = [...claves]
  keysTras.splice(pos, 1)

  pasos.push({
    nodeId: hojaId,
    tipo: 'eliminar_fisico',
    mensaje: `🗑 Borrado físico: [${claves.join(', ')}] → [${keysTras.join(', ') || '∅'}]`,
    activeKey: claveObjetivo,
    deletingKeys: [claveObjetivo],
    posEliminar: pos,
    borradoVirtual: true,
  })

  const sim = simularCorreccionUnderflow(estructura, hojaId, keysTras, claveObjetivo, orden)
  sim.deleteInfo.posEliminar = pos
  sim.deleteInfo.idCancion = idCancion

  if (sim.tipo === 'fantasma') {
    pasos.push({
      nodeId: hojaId,
      tipo: 'sin_underflow',
      mensaje: `✓ Sin underflow (${keysTras.length} ≥ d=${minD}) — no hace falta corregirUnderflowHoja`,
      highlightedNodes: [hojaId],
      borradoVirtual: true,
    })

    if (sim.deleteInfo.separadorFantasma != null) {
      pasos.push({
        nodeId: sim.deleteInfo.padreId,
        tipo: 'fantasma_explicacion',
        mensaje: `👻 El separador ${sim.deleteInfo.separadorFantasma} permanece en el índice (fantasma) — las rutas siguen siendo válidas`,
        highlightedNodes: [sim.deleteInfo.padreId, hojaId],
        activeKey: sim.deleteInfo.separadorFantasma,
        deletePhase: FASES_DELETE.FANTASMA,
        borradoVirtual: true,
      })
    }
  } else if (sim.tipo !== 'ninguno') {
    pasos.push({
      nodeId: hojaId,
      tipo: 'underflow_detectado',
      mensaje: `⚠ Underflow: ${keysTras.length} clave(s) < d=${minD} — invocando corregirUnderflowHoja()`,
      highlightedNodes: [hojaId],
      underflowNodes: [hojaId],
      deletePhase: FASES_DELETE.UNDERFLOW,
      borradoVirtual: true,
    })
    pasos.push(...sim.pasosCorreccion.map((p) => ({ ...p, borradoVirtual: true })))
  }

  pasos.push({
    nodeId: hojaId,
    tipo: 'aplicar',
    mensaje: `⚙️ Aplicando borrado y corrección en el árbol B+…`,
    highlightedNodes: [hojaId],
  })

  return {
    pasos,
    encontrado: true,
    deleteInfo: sim.deleteInfo,
    hojaId,
    keysTras,
    tipoCorreccion: sim.tipo,
  }
}

/** Pasos finales breves tras refrescar el árbol (solo confirmación, sin re-explicar). */
export function calcularPasosConfirmacion(info, estructura) {
  if (!info || !estructura) return []

  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const pasos = []

  if (info.tipo === 'fantasma') {
    const hId = findHojaPorClaves(estructura, info.keysTras) ?? info.hojaId
    const keys = clavesEnHoja(mapa[hId]) ?? info.keysTras
    pasos.push({
      nodeId: hId,
      tipo: 'confirmacion',
      mensaje: `✅ Resultado: hoja [${keys.join(', ')}]${info.separadorFantasma ? ` — separador ${info.separadorFantasma} fantasma` : ''}`,
      highlightedNodes: [hId],
      highlightedKeys: keys,
    })
    return pasos
  }

  if (info.tipo === 'prestamo') {
    const hId = info.hojaId
    const keys = clavesEnHoja(mapa[hId]) ?? info.keysHojaFinal
    pasos.push({
      nodeId: hId,
      tipo: 'confirmacion',
      mensaje: `✅ Underflow resuelto por REDISTRIBUCIÓN — hoja final [${keys.join(', ')}]`,
      highlightedNodes: [hId, info.hermanoId].filter(Boolean),
      highlightedKeys: keys,
    })
    return pasos
  }

  if (info.tipo === 'fusion') {
    const mergedId =
      findHojaPorClaves(estructura, info.keysFusionadas) ?? info.hojaFusionadaId
    const keys = clavesEnHoja(mapa[mergedId]) ?? info.keysFusionadas
    const hojas = estructura.cadenaHojas?.length ?? 0
    pasos.push({
      nodeId: mergedId,
      tipo: 'confirmacion',
      mensaje: `✅ Underflow resuelto por FUSIÓN — ${hojas} hoja(s), claves [${keys.join(', ')}]`,
      highlightedNodes: [mergedId, info.padreId].filter(Boolean),
      highlightedKeys: keys,
    })

    const prof = estimarProfundidad(estructura)
    if (info.profAntes && prof < info.profAntes) {
      pasos.push({
        nodeId: estructura.raiz,
        tipo: 'confirmacion_raiz',
        mensaje: `⬇ Altura reducida: ${info.profAntes} → ${prof} (ajustarRaiz aplicado)`,
        highlightedNodes: [estructura.raiz],
        deletePhase: FASES_DELETE.RAIZ,
      })
    }
    return pasos
  }

  return pasos
}

// ─── Diff antes/después (verificación tras API) ─────────────────────────────

function findMergeDetails(anterior, nueva, clave, hojaOrigenId) {
  const mapaAnt = Object.fromEntries(anterior.nodos.map((n) => [n.id, n]))
  const cadena = anterior.cadenaHojas || []
  const idx = cadena.indexOf(hojaOrigenId)
  if (idx < 0) return null

  const candidatos = []
  if (idx > 0) candidatos.push([cadena[idx - 1], cadena[idx], 'derecha'])
  if (idx < cadena.length - 1) candidatos.push([cadena[idx], cadena[idx + 1], 'derecha'])
  if (idx > 0) candidatos.push([cadena[idx - 1], cadena[idx], 'izquierda'])

  const seen = new Set()
  for (const [a, b] of candidatos.map(([x, y]) => [x, y])) {
    const key = `${a}-${b}`
    if (seen.has(key)) continue
    seen.add(key)

    const ha = mapaAnt[a]
    const hb = mapaAnt[b]
    if (!ha || !hb) continue

    const keysIzq = a < b ? clavesEnHoja(ha) : clavesEnHoja(hb)
    const keysDer = a < b ? clavesEnHoja(hb) : clavesEnHoja(ha)
    const merged = ordenarClaves(
      quitarUnaClave([...clavesEnHoja(ha), ...clavesEnHoja(hb)], clave),
    )
    const mergedId = findHojaPorClaves(nueva, merged)
    if (mergedId == null) continue

    const padreAnt = findPadreDeHoja(anterior, a) || findPadreDeHoja(anterior, b)
    let separadorEliminado = null
    if (padreAnt) {
      const hi = padreAnt.hijos?.indexOf(a)
      const hj = padreAnt.hijos?.indexOf(b)
      if (hi >= 0 && hj >= 0) {
        const sepIdx = Math.min(hi, hj)
        separadorEliminado = padreAnt.separadores?.[sepIdx] ?? null
      }
    }

    return {
      mergedId,
      keysIzqAntes: quitarUnaClave([...keysIzq], clave),
      keysHojaDerAntes: quitarUnaClave([...keysDer], clave),
      keysFusionadas: merged,
      separadorEliminado,
      padreAntesId: padreAnt?.id ?? null,
      hojaIzqId: a,
      hojaDerId: b,
    }
  }
  return null
}

function detectarPrestamo(anterior, nueva, clave, hojaAntesId, keysTras) {
  const mapaNue = Object.fromEntries(nueva.nodos.map((n) => [n.id, n]))
  const mapaAnt = Object.fromEntries(anterior.nodos.map((n) => [n.id, n]))

  for (const hojaId of nueva.cadenaHojas || []) {
    const h = mapaNue[hojaId]
    if (!h?.esHoja) continue
    const hKeys = clavesEnHoja(h)
    const contieneTras =
      keysTras.length <= hKeys.length &&
      keysTras.every((k) => hKeys.some((hk) => compararClave(hk, k) === 0))
    if (!contieneTras || hKeys.length <= keysTras.length) continue

    const prestada = hKeys.find(
      (k) => !keysTras.some((t) => compararClave(t, k) === 0) && compararClave(k, clave) !== 0,
    )
    if (prestada == null) continue

    for (const hermanoId of nueva.cadenaHojas || []) {
      if (hermanoId === hojaId) continue
      const hermano = mapaNue[hermanoId]
      const keysHermano = clavesEnHoja(hermano)
      if (!keysHermano.some((k) => compararClave(k, prestada) === 0)) continue

      const idxNueH = nueva.cadenaHojas.indexOf(hojaId)
      const idxNueHr = nueva.cadenaHojas.indexOf(hermanoId)
      if (Math.abs(idxNueH - idxNueHr) !== 1) continue

      let padreId = null
      let nuevoSeparador = null
      for (const n of nueva.nodos) {
        if (n.esHoja || !n.hijos) continue
        const hi = n.hijos.indexOf(hojaId)
        const hr = n.hijos.indexOf(hermanoId)
        if (hi >= 0 && hr >= 0) {
          padreId = n.id
          const sepIdx = Math.min(hi, hr)
          nuevoSeparador = n.separadores?.[sepIdx] ?? clavesEnHoja(hermano)[0]
          break
        }
      }

      const direccion = idxNueHr > idxNueH ? 'derecha' : 'izquierda'
      const hermanoAntId =
        direccion === 'derecha'
          ? anterior.cadenaHojas[anterior.cadenaHojas.indexOf(hojaAntesId) + 1]
          : anterior.cadenaHojas[anterior.cadenaHojas.indexOf(hojaAntesId) - 1]
      const hermanoAnt = mapaAnt[hermanoAntId]

      return {
        hojaId,
        hermanoId,
        direccion,
        clavePrestada: prestada,
        nuevoSeparador,
        padreId,
        keysTrasBorrado: keysTras,
        keysHojaFinal: clavesEnHoja(h),
        keysHermanoFinal: clavesEnHoja(hermano),
        keysHermanoAntes: clavesEnHoja(hermanoAnt),
      }
    }
  }
  return null
}

export function analizarEliminacion(anterior, nueva, claveEliminada, orden, idCancion = null) {
  if (!anterior?.nodos?.length || !nueva?.nodos?.length) {
    return { tipo: 'ninguno' }
  }

  const minD = orden
  const { hojaId: hojaAntesId, hoja: hojaAntes, pos } = hojaOrigen(anterior, claveEliminada, idCancion)
  const clavesAntes = clavesEnHoja(hojaAntes)
  if (pos < 0 || pos >= clavesAntes.length) {
    return { tipo: 'no_encontrado' }
  }

  const keysTras = [...clavesAntes]
  keysTras.splice(pos, 1)
  ordenarClaves(keysTras)
  const esRaizUnica = hojaAntesId === anterior.raiz && hojaAntes.esHoja
  const underflow = !esRaizUnica && keysTras.length < minD

  const hojasAntes = anterior.cadenaHojas?.length ?? 0
  const hojasDespues = nueva.cadenaHojas?.length ?? 0
  const profAntes = estimarProfundidad(anterior)
  const profDespues = estimarProfundidad(nueva)

  const base = {
    claveEliminada,
    keysTras,
    minD,
    hojasAntes,
    hojasDespues,
    profAntes,
    profDespues,
    hojaId: hojaAntesId,
  }

  if (!underflow) {
    const { padreId, separadorFantasma } = buscarSeparadorFantasma(anterior, claveEliminada)
    return {
      ...base,
      tipo: 'fantasma',
      padreId,
      separadorFantasma,
      keysRestantes: keysTras,
    }
  }

  if (hojasDespues < hojasAntes) {
    const merge = findMergeDetails(anterior, nueva, claveEliminada, hojaAntesId)
    const alturaReducida = profDespues < profAntes
    return {
      ...base,
      tipo: 'fusion',
      keysFusionadas: merge?.keysFusionadas ?? keysTras,
      hojaFusionadaId: merge?.mergedId,
      alturaReducida,
    }
  }

  const prestamo = detectarPrestamo(anterior, nueva, claveEliminada, hojaAntesId, keysTras)
  if (prestamo) {
    return { ...base, tipo: 'prestamo', ...prestamo }
  }

  return { ...base, tipo: 'underflow' }
}

/** @deprecated Usar calcularSecuenciaEliminacion + calcularPasosConfirmacion */
export function calcularPasosPreEliminacion(estructura, indice, orden) {
  const { pasos, encontrado, hojaId, keysTras } = calcularSecuenciaEliminacion(estructura, indice, orden)
  const prePasos = pasos.filter((p) => p.tipo !== 'aplicar' && !p.borradoVirtual?.toString())
  return { pasos: prePasos, encontrado, keysTras, hojaId }
}

/** @deprecated Usar calcularPasosConfirmacion */
export function calcularPasosPostEliminacion(info, estructura) {
  return calcularPasosConfirmacion(info, estructura)
}
