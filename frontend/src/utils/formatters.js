export function formatearDuracion(ms) {
  if (!ms && ms !== 0) return '--:--'
  const totalSeg = Math.floor(ms / 1000)
  const min = Math.floor(totalSeg / 60)
  const seg = totalSeg % 60
  return `${min}:${seg.toString().padStart(2, '0')}`
}

export function truncar(texto, max = 28) {
  if (!texto) return ''
  return texto.length > max ? `${texto.slice(0, max)}…` : texto
}
