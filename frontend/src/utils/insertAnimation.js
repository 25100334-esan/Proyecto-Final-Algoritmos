import { calcularPasosNavegacion, calcularLayout, compararClave } from './treeLayout.js'
import { FASES_SPLIT } from './splitAnimation.js'
import { clavesDeNodo, formatearClaves } from './formatKeys.js'
import { normalizarClaveArbol } from './claveVista.js'

export { FASES_SPLIT }

function findPadreDeHoja(estructura, hojaId) {
  for (const n of estructura?.nodos || []) {
    if (!n.esHoja && n.hijos?.includes(hojaId)) return n
  }
  return null
}

function layoutInfo(estructura) {
  const layout = calcularLayout(estructura)
  const pos = (id) => (id != null ? layout.posiciones[id] : null)
  const dim = (id) => (id != null ? layout.dimensiones[id] : null)
  return { pos, dim }
}

export function estructuraConInsercionVirtual(estructura, hojaId, claveInsertada) {
  const clone = JSON.parse(JSON.stringify(estructura))
  const nodo = clone.nodos.find((n) => n.id === hojaId)
  const actuales = clavesDeNodo(nodo)
  if (nodo && !actuales.some((k) => compararClave(k, claveInsertada) === 0)) {
    const mezcladas = [...actuales, claveInsertada].sort((a, b) => compararClave(a, b))
    if (estructura.tipoClave === 'int' || !estructura.tipoClave) {
      nodo.indices = mezcladas.map(Number)
    }
    nodo.claves = mezcladas.map(String)
  }
  return clone
}

function simularSplitHoja(estructura, hojaId, keysConInsert, claveInsertada, orden) {
  const d = orden
  const maxPorNodo = d * 2
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const hoja = mapa[hojaId]
  const { pos, dim } = layoutInfo(estructura)
  const hojaPos = pos(hojaId)
  const padre = findPadreDeHoja(estructura, hojaId)
  const esRaizHoja = hojaId === estructura.raiz && hoja?.esHoja

  const baseInfo = {
    claveInsertada,
    hojaDestinoId: hojaId,
    maxPorNodo,
    minD: d,
    hojaPos,
    hojaDim: dim(hojaId),
    huboSplit: false,
  }

  if (keysConInsert.length <= maxPorNodo) {
    return {
      tipo: 'simple',
      insertInfo: {
        ...baseInfo,
        keysFinales: keysConInsert,
      },
      pasosCorreccion: [],
    }
  }

  const medio = Math.floor((keysConInsert.length + 1) / 2)
  const izqKeys = keysConInsert.slice(0, medio)
  const derKeys = keysConInsert.slice(medio)
  const clavePromovida = derKeys[0]

  const hojaDerPos = hojaPos
    ? { x: hojaPos.x + (dim(hojaId)?.width ?? 120) * 0.55, y: hojaPos.y }
    : null
  const padrePos = esRaizHoja
    ? hojaPos ? { x: hojaPos.x, y: hojaPos.y - 90 } : null
    : pos(padre?.id)

  const pasos = []

  pasos.push({
    nodeId: hojaId,
    tipo: 'overflow_detectado',
    mensaje: `⚠ Overflow: ${keysConInsert.length} claves > 2d=${maxPorNodo} — invocando partirHoja()`,
    highlightedNodes: [hojaId, padre?.id].filter(Boolean),
    highlightedKeys: keysConInsert,
    splitPhase: FASES_SPLIT.SPLITTING,
    overflowPreview: true,
    insercionVirtual: true,
    claveInsertada,
    hojaId,
  })

  pasos.push({
    nodeId: hojaId,
    tipo: 'split_partir',
    mensaje: `✂️ partirHoja(): dividir [${keysConInsert.join(', ')}] → izq [${izqKeys.join(', ')}] | der [${derKeys.join(', ')}]`,
    highlightedNodes: [hojaId],
    highlightedKeys: keysConInsert,
    splitPhase: FASES_SPLIT.SPLITTING,
    insercionVirtual: true,
    claveInsertada,
    hojaId,
    izqKeys,
    derKeys,
  })

  pasos.push({
    nodeId: hojaId,
    tipo: 'split_sequence',
    mensaje: `↪ Re-enlazar sequence set: nueva hoja derecha enlazada (siguienteHoja)`,
    highlightedNodes: [hojaId],
    splitPhase: FASES_SPLIT.SPLITTING,
    insercionVirtual: true,
    claveInsertada,
    hojaId,
    leafChainFrom: hojaId,
    leafChainTo: hojaId,
  })

  pasos.push({
    nodeId: hojaId,
    tipo: 'split_promote',
    mensaje: `⬆ Promoción por COPIA: clave ${clavePromovida} sube al índice (permanece en hoja derecha)`,
    highlightedNodes: [hojaId, padre?.id].filter(Boolean),
    highlightedKeys: [clavePromovida, claveInsertada],
    splitPhase: FASES_SPLIT.PROMOTE,
    animarPromote: true,
    insercionVirtual: true,
    claveInsertada,
    hojaId,
    activeKey: clavePromovida,
  })

  if (esRaizHoja) {
    pasos.push({
      nodeId: estructura.raiz,
      tipo: 'split_nueva_raiz',
      mensaje: `🌱 La raíz era hoja → nueva raíz interna con separador ${clavePromovida} y dos hijos hoja`,
      highlightedNodes: [hojaId],
      activeKey: clavePromovida,
      splitPhase: FASES_SPLIT.PROMOTE,
      insercionVirtual: true,
      claveInsertada,
      hojaId,
    })
  } else if (padre) {
    const sepTras = (padre.separadores?.length ?? 0) + 1
    pasos.push({
      nodeId: padre.id,
      tipo: 'split_padre',
      mensaje: `↑ Insertar separador ${clavePromovida} en padre — divide punteros a hoja izq. y hoja der.`,
      highlightedNodes: [padre.id, hojaId],
      activeKey: clavePromovida,
      splitPhase: FASES_SPLIT.PROMOTE,
      insercionVirtual: true,
      claveInsertada,
      hojaId,
    })

    if (sepTras > maxPorNodo) {
      pasos.push({
        nodeId: padre.id,
        tipo: 'split_propagacion',
        mensaje: `⬆ El padre también hace overflow (${sepTras} > 2d=${maxPorNodo}) → partirNodoInterno (propagación)`,
        highlightedNodes: [padre.id],
        splitPhase: FASES_SPLIT.SPLITTING,
        insercionVirtual: true,
        claveInsertada,
        hojaId,
      })
    }
  }

  return {
    tipo: 'split',
    insertInfo: {
      ...baseInfo,
      huboSplit: true,
      keysConInsert,
      izqKeys,
      derKeys,
      clavePromovida,
      medio,
      padreId: padre?.id ?? null,
      esRaizHoja,
      hojaIzqPos: hojaPos,
      hojaDerPos,
      hojaIzqDim: dim(hojaId),
      hojaDerDim: dim(hojaId),
      padrePos,
      padreDim: dim(padre?.id),
    },
    pasosCorreccion: pasos,
  }
}

/** Secuencia completa de inserción ANTES del API (navegación + insert virtual + split simulado). */
export function calcularSecuenciaInsercion(estructura, claveInsertar, orden) {
  const clave = normalizarClaveArbol(estructura, claveInsertar)
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const nav = calcularPasosNavegacion(estructura, clave)
  const hojaId = nav[nav.length - 1]?.nodeId
  const hoja = mapa[hojaId]
  const pasos = [...nav.slice(0, -1)]

  if (!hoja?.esHoja) {
    return { pasos, hojaId: null, insertInfo: null, tipoCorreccion: 'ninguno' }
  }

  const keysActuales = clavesDeNodo(hoja)
  const keysConInsert = [...keysActuales, clave].sort((a, b) => compararClave(a, b))
  const maxPorNodo = orden * 2

  pasos.push({
    nodeId: hojaId,
    tipo: 'hoja_destino',
    mensaje: `📄 Hoja destino — inserción SOLO en nivel inferior (nunca en el índice directamente)`,
    highlightedNodes: [hojaId],
  })

  pasos.push({
    nodeId: hojaId,
    tipo: 'insertar_virtual',
    mensaje: `➕ Insertar clave ${clave}: [${formatearClaves(keysActuales)}] → [${formatearClaves(keysConInsert)}]`,
    highlightedNodes: [hojaId],
    highlightedKeys: keysConInsert,
    activeKey: clave,
    insercionVirtual: true,
    claveInsertada: clave,
    hojaId,
  })

  const sim = simularSplitHoja(estructura, hojaId, keysConInsert, clave, orden)

  if (sim.tipo === 'simple') {
    pasos.push({
      nodeId: hojaId,
      tipo: 'sin_overflow',
      mensaje: `✓ Sin overflow (${keysConInsert.length} ≤ 2d=${maxPorNodo}) — inserción directa, sin split`,
      highlightedNodes: [hojaId],
      highlightedKeys: keysConInsert,
      insercionVirtual: true,
      claveInsertada: clave,
      hojaId,
    })
  } else {
    pasos.push(...sim.pasosCorreccion)
  }

  pasos.push({
    nodeId: hojaId,
    tipo: 'aplicar',
    mensaje: `⚙️ Aplicando inserción en el árbol B+ en memoria…`,
    highlightedNodes: [hojaId],
  })

  return {
    pasos,
    hojaId,
    insertInfo: sim.insertInfo,
    tipoCorreccion: sim.tipo,
  }
}

/** Confirmación breve tras refrescar el árbol. */
export function calcularPasosConfirmacionInsercion(info, estructura, idAsignado) {
  if (!info || !estructura) return []

  const pasos = []
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))

  if (!info.huboSplit) {
    let hId = info.hojaDestinoId
    for (const id of estructura.cadenaHojas || []) {
      const h = mapa[id]
      if (h?.indices?.includes(idAsignado)) {
        hId = id
        break
      }
    }
    const keys = mapa[hId]?.indices ?? info.keysFinales ?? []
    pasos.push({
      nodeId: hId,
      tipo: 'confirmacion',
      mensaje: `✅ Clave ${idAsignado} insertada sin división — hoja [${keys.join(', ')}]`,
      highlightedNodes: [hId],
      highlightedKeys: keys,
    })
    return pasos
  }

  let hojaIzqId = null
  let hojaDerId = null
  for (const id of estructura.cadenaHojas || []) {
    const h = mapa[id]
    if (!h?.esHoja) continue
    if (info.izqKeys?.every((k) => h.indices.includes(k))) hojaIzqId = id
    if (info.derKeys?.every((k) => h.indices.includes(k))) hojaDerId = id
  }

  pasos.push({
    nodeId: hojaDerId ?? hojaIzqId ?? estructura.raiz,
    tipo: 'confirmacion',
    mensaje: `✅ SPLIT completado — separador ${info.clavePromovida} en índice | hojas [${info.izqKeys?.join(', ')}] [${info.derKeys?.join(', ')}]`,
    highlightedNodes: [hojaIzqId, hojaDerId, info.padreId ?? estructura.raiz].filter(Boolean),
    highlightedKeys: [idAsignado, info.clavePromovida],
    splitPhase: FASES_SPLIT.PROMOTE,
  })

  const hojas = estructura.cadenaHojas?.length ?? 0
  pasos.push({
    nodeId: estructura.raiz,
    tipo: 'confirmacion_estructura',
    mensaje: `🌳 Árbol actualizado — ${hojas} hoja(s), ${estructura.nodos.length} nodo(s)`,
    highlightedNodes: [estructura.raiz],
  })

  return pasos
}
