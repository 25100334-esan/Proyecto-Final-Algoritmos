import { formatearClaves, clavesDeNodo } from './formatKeys.js'
import { clavesVisuales, separadoresVisuales, indicesNumericos } from './nodoClaves.js'
import { normalizarClaveArbol, clavesUnicasEnEstructura } from './claveVista.js'

const NODE_HEIGHT = 54
const V_GAP = 80
const H_GAP_BASE = 40
const PTR_SLOT_W = 18
const LEAF_CHAIN_Y_OFFSET = 28

export { NODE_HEIGHT, V_GAP, H_GAP_BASE, LEAF_CHAIN_Y_OFFSET }

/** Ancho mínimo de celda según longitud del número. */
function anchoCelda(valor, esPuntero = false) {
  if (esPuntero) return PTR_SLOT_W
  const texto = String(valor)
  return Math.max(48, texto.length * 8.5 + 16)
}

/** Dimensiones y celdas de un nodo para dibujo tipo diagrama académico. */
export function calcularDimensiones(nodo) {
  if (nodo.esHoja) {
    const keys = clavesVisuales(nodo).length ? clavesVisuales(nodo) : [null]
    const celdas = keys.map((k) => ({
      tipo: 'key',
      valor: k,
      ancho: anchoCelda(k ?? '∅'),
    }))
    const width = celdas.reduce((s, c) => s + c.ancho, 0)
    return { width, height: NODE_HEIGHT, celdas }
  }

  const seps = separadoresVisuales(nodo)
  const celdas = []
  for (let i = 0; i <= seps.length; i++) {
    celdas.push({ tipo: 'ptr', indice: i, ancho: PTR_SLOT_W })
    if (i < seps.length) {
      celdas.push({ tipo: 'sep', valor: seps[i], ancho: anchoCelda(seps[i]) })
    }
  }
  const width = celdas.reduce((s, c) => s + c.ancho, 0)
  return { width, height: NODE_HEIGHT, celdas, numHijos: (nodo.hijos || []).length }
}

function anchoSubarbol(id, mapa, dims, hGap) {
  const nodo = mapa[id]
  if (!nodo || nodo.esHoja) return dims[id]?.width ?? 60
  const hijos = nodo.hijos || []
  if (!hijos.length) return dims[id].width
  const suma = hijos.reduce((s, h) => s + anchoSubarbol(h, mapa, dims, hGap), 0)
  return Math.max(dims[id].width, suma + (hijos.length - 1) * hGap)
}

/** Posición del ancla inferior de un puntero (índice de hijo 0..n). */
export function anclaPuntero(nodeId, ptrIndex, posiciones, dimensiones) {
  const pos = posiciones[nodeId]
  const dim = dimensiones[nodeId]
  if (!pos || !dim) return null

  let offset = -dim.width / 2
  let ptrCount = 0
  for (const celda of dim.celdas) {
    if (celda.tipo === 'ptr') {
      if (ptrCount === ptrIndex) {
        return { x: pos.x + offset + celda.ancho / 2, y: pos.y + dim.height / 2 }
      }
      ptrCount++
    }
    offset += celda.ancho
  }
  return { x: pos.x, y: pos.y + dim.height / 2 }
}

/**
 * Layout recursivo: padres centrados sobre hijos, como diagrama de libros.
 */
export function calcularLayout(estructura) {
  if (!estructura?.nodos?.length) {
    return {
      posiciones: {},
      dimensiones: {},
      edges: [],
      leafChain: [],
      width: 800,
      height: 400,
    }
  }

  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const dimensiones = {}
  estructura.nodos.forEach((n) => {
    dimensiones[n.id] = calcularDimensiones(n)
  })

  const maxAnchoNodo = Math.max(...Object.values(dimensiones).map((d) => d.width), 60)
  const hGap = Math.max(H_GAP_BASE, maxAnchoNodo * 0.3)

  const posiciones = {}
  const edges = []

  function layout(id, centerX, y) {
    const nodo = mapa[id]
    posiciones[id] = { x: centerX, y }

    if (nodo.esHoja) return

    const hijos = nodo.hijos || []
    const anchos = hijos.map((h) => anchoSubarbol(h, mapa, dimensiones, hGap))
    const total = anchos.reduce((a, b) => a + b, 0) + (hijos.length - 1) * hGap
    let x = centerX - total / 2

    hijos.forEach((hijoId, i) => {
      const cx = x + anchos[i] / 2
      layout(hijoId, cx, y + NODE_HEIGHT + V_GAP)
      edges.push({ from: id, to: hijoId, ptrIndex: i, key: `${id}-${i}` })
      x += anchos[i] + hGap
    })
  }

  layout(estructura.raiz, 0, 0)

  const xs = Object.values(posiciones).map((p) => p.x)
  const ys = Object.values(posiciones).map((p) => p.y)
  const minX = Math.min(...xs) - 80
  const maxX = Math.max(...xs) + 80
  const maxY = Math.max(...ys) + NODE_HEIGHT + LEAF_CHAIN_Y_OFFSET + 50

  const leafChain = []
  for (let i = 0; i < estructura.cadenaHojas.length - 1; i++) {
    leafChain.push({ from: estructura.cadenaHojas[i], to: estructura.cadenaHojas[i + 1] })
  }

  return {
    posiciones,
    dimensiones,
    edges,
    leafChain,
    width: maxX - minX,
    height: maxY + 20,
    offsetX: (minX + maxX) / 2,
  }
}

/** Ruta simple de ids de nodos visitados. */
export function calcularRutaBusqueda(estructura, indice) {
  return calcularPasosNavegacion(estructura, indice).map((p) => p.nodeId)
}

/**
 * Índice de hijo al bajar (misma regla que indiceHijo en interno.go).
 * Únicas: clave >= separador → derecha. Duplicados: solo clave > separador → derecha.
 */
export function indiceHijoEnEstructura(estructura, clave, separadores) {
  const unicas = clavesUnicasEnEstructura(estructura)
  let i = 0
  while (i < separadores.length) {
    const cmp = compararClave(clave, separadores[i])
    if (unicas) {
      if (cmp < 0) break
      i++
    } else if (cmp > 0) {
      i++
    } else {
      break
    }
  }
  return i
}

/**
 * Pasos detallados con mensajes de comparación para animación académica.
 */
export function calcularPasosNavegacion(estructura, clave) {
  if (!estructura?.nodos?.length) return []

  const claveBusqueda = normalizarClaveArbol(estructura, clave)
  const unicas = clavesUnicasEnEstructura(estructura)

  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const pasos = []
  let actual = estructura.raiz

  pasos.push({
    nodeId: actual,
    tipo: 'inicio',
    mensaje: `🔍 Buscar clave ${claveBusqueda}: comenzamos en la raíz del índice`,
    activeEdge: null,
    activeSlot: null,
    activeKey: null,
  })

  while (actual != null) {
    const nodo = mapa[actual]
    if (!nodo) break

    if (nodo.esHoja) {
      pasos.push({
        nodeId: actual,
        tipo: 'hoja',
        mensaje: `📄 Hoja alcanzada — buscamos ${claveBusqueda} entre [${formatearClaves(clavesVisuales(nodo))}]`,
        activeEdge: null,
        activeSlot: null,
        activeKey: claveBusqueda,
      })
      break
    }

    const seps = separadoresVisuales(nodo)
    const i = indiceHijoEnEstructura(estructura, claveBusqueda, seps)

    for (let j = 0; j < i; j++) {
      const cmp = compararClave(claveBusqueda, seps[j])
      const igual = cmp === 0
      pasos.push({
        nodeId: actual,
        tipo: 'comparar',
        mensaje: igual
          ? unicas
            ? `${claveBusqueda} = ${seps[j]} (separador) → índice único: seguir al hijo derecho`
            : `${claveBusqueda} = ${seps[j]} (separador) → índice secundario: primer bloque de duplicados a la izquierda`
          : `${claveBusqueda} > ${seps[j]} → a la derecha del separador ${seps[j]}`,
        activeEdge: null,
        activeSlot: j,
        activeKey: seps[j],
        separador: seps[j],
      })
    }

    let msg
    if (seps.length === 0) {
      msg = `Sin separadores → descender al único hijo`
    } else if (i === 0) {
      msg = `${claveBusqueda} < ${seps[0]} → descender por puntero izquierdo (hijo 0)`
    } else if (i < seps.length && compararClave(claveBusqueda, seps[i]) === 0 && !unicas) {
      msg = `${claveBusqueda} = ${seps[i]} → descender por puntero ${i} (hoja candidata de duplicados)`
    } else if (i < seps.length) {
      msg = `${claveBusqueda} < ${seps[i]} → descender por puntero ${i} (entre ${seps[i - 1]} y ${seps[i]})`
    } else {
      msg = `${claveBusqueda} > ${seps[seps.length - 1]} → descender por puntero derecho (hijo ${i})`
    }

    pasos.push({
      nodeId: actual,
      tipo: 'decision',
      mensaje: `⬇ ${msg}`,
      activeEdge: `${actual}-${i}`,
      activeSlot: i,
      activeKey: null,
      hijoIndex: i,
    })

    actual = nodo.hijos?.[i] ?? null
  }

  return pasos
}

function compararClave(a, b) {
  const na = Number(a)
  const nb = Number(b)
  if (!Number.isNaN(na) && !Number.isNaN(nb)) return na - nb
  return String(a).localeCompare(String(b))
}

export { compararClave }

export function detectarSplits(estructuraAnterior, estructuraNueva) {
  if (!estructuraAnterior?.nodos?.length || !estructuraNueva?.nodos?.length) {
    return { huboSplit: false, nodosAnimar: [] }
  }

  const hojasAntes = estructuraAnterior.cadenaHojas.length
  const hojasDespues = estructuraNueva.cadenaHojas.length
  const nodosAntes = estructuraAnterior.nodos.length
  const nodosDespues = estructuraNueva.nodos.length
  const huboSplit = hojasDespues > hojasAntes || nodosDespues > nodosAntes

  let nodosAnimar = []
  if (huboSplit) {
    const nuevasHojas = Math.max(hojasDespues - hojasAntes, 1)
    nodosAnimar = estructuraNueva.cadenaHojas.slice(-nuevasHojas)
  }

  return { huboSplit, nodosAnimar }
}

export function estimarProfundidad(estructura) {
  if (!estructura?.nodos?.length) return 0
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))

  function profundidad(id) {
    const nodo = mapa[id]
    if (!nodo || nodo.esHoja) return 1
    return 1 + Math.max(...nodo.hijos.map(profundidad))
  }

  return profundidad(estructura.raiz)
}
