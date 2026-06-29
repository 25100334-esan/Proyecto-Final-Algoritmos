import { calcularPasosNavegacion, compararClave } from './treeLayout.js'
import { clavesVisuales } from './nodoClaves.js'
import { normalizarClaveArbol } from './claveVista.js'

function idsCoinciden(a, b) {
  if (a == null || b == null) return false
  return Number(a) === Number(b)
}

/**
 * Localiza entrada en hoja(s) del árbol exportado.
 * Con idCancion usa idsRegistro (a_i de Comer) en el sequence set.
 */
export function localizarEntradaEnHojas(estructura, clave, idCancion = null) {
  const mapa = Object.fromEntries(estructura.nodos.map((n) => [n.id, n]))
  const claveNorm = normalizarClaveArbol(estructura, clave)
  const nav = calcularPasosNavegacion(estructura, claveNorm)
  const hojaInicialId = nav[nav.length - 1]?.nodeId
  const cadena = estructura.cadenaHojas ?? []
  let startIdx = cadena.indexOf(hojaInicialId)
  if (startIdx < 0) startIdx = 0

  // Secundarios: localizar por ID en hojas (sequence set) — r = (k, a)
  if (idCancion != null) {
    for (let ci = startIdx; ci < cadena.length; ci++) {
      const hojaId = cadena[ci]
      const hoja = mapa[hojaId]
      if (!hoja?.esHoja) continue
      const claves = clavesVisuales(hoja)
      const ids = hoja.idsRegistro ?? []
      for (let i = 0; i < ids.length; i++) {
        if (idsCoinciden(ids[i], idCancion)) {
          return {
            encontrado: true,
            nav,
            hojaId,
            pos: i,
            claves,
            claveNorm,
            startIdx,
            cadenaIdx: ci,
          }
        }
      }
    }
    return {
      encontrado: false,
      nav,
      hojaId: hojaInicialId,
      claves: [],
      claveNorm,
      startIdx,
      cadenaIdx: startIdx,
    }
  }

  for (let ci = startIdx; ci < cadena.length; ci++) {
    const hojaId = cadena[ci]
    const hoja = mapa[hojaId]
    if (!hoja?.esHoja) continue
    const claves = clavesVisuales(hoja)

    for (let i = 0; i < claves.length; i++) {
      const cmp = compararClave(claves[i], claveNorm)
      if (cmp > 0) {
        return { encontrado: false, nav, hojaId, claves, claveNorm, startIdx, cadenaIdx: ci }
      }
      if (cmp === 0) {
        return {
          encontrado: true,
          nav,
          hojaId,
          pos: i,
          claves,
          claveNorm,
          startIdx,
          cadenaIdx: ci,
        }
      }
    }

    const ultima = claves[claves.length - 1]
    if (!ultima || compararClave(ultima, claveNorm) !== 0) break
  }

  return {
    encontrado: false,
    nav,
    hojaId: hojaInicialId,
    claves: [],
    claveNorm,
    startIdx,
    cadenaIdx: startIdx,
  }
}
