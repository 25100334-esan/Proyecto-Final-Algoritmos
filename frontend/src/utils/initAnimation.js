/** Máximo de canciones para demo de inicialización paso a paso. */
export const LIMITE_DEMO_INICIAL = 50

export function idDeCancion(cancion) {
  return cancion?.Indice ?? cancion?.indice ?? cancion?.id ?? 0
}

export function nombreDeCancion(cancion) {
  return cancion?.TrackName ?? cancion?.trackName ?? ''
}

export function artistaDeCancion(cancion) {
  return cancion?.Artists ?? cancion?.artists ?? ''
}

/** Pasos antes de llamar al API (conexión y consulta planificada). */
export function calcularPasosPreLecturaBD(limite) {
  return [
    {
      nodeId: 0,
      tipo: 'init_bd_archivo',
      mensaje: `💾 Conectar a dataset.db (SQLite) — archivo en disco con tabla tracks`,
    },
    {
      nodeId: 0,
      tipo: 'init_bd_sql',
      mensaje:
        `📜 Consulta SQL: SELECT id, track_name, artists, … FROM tracks ORDER BY id LIMIT ${limite}`,
    },
    {
      nodeId: 0,
      tipo: 'init_bd_api',
      mensaje: `🌐 POST /api/inicializar-demo → árbol B+ vacío en RAM + lectura de hasta ${limite} fila(s)`,
    },
  ]
}

export function calcularPasosIntroInicializacion(orden, total) {
  return [
    {
      nodeId: 0,
      tipo: 'init_vacio',
      mensaje: `🌱 Árbol B+ vacío creado (d=${orden}) — la raíz es una hoja sin claves [∅]`,
    },
    {
      nodeId: 0,
      tipo: 'init_bd_resultado',
      mensaje: `📥 ${total} registro(s) leídos de SQLite y cargados en memoria (aún NO están en el árbol)`,
    },
    {
      nodeId: 0,
      tipo: 'init_bd_nota',
      mensaje: `ℹ️ Este dataset usa id desde 0: con límite 5 se insertan claves 0, 1, 2, 3 y 4`,
    },
    {
      nodeId: 0,
      tipo: 'init_plan',
      mensaje: `📋 Demo: insertar cada registro en el B+ con la misma lógica que «Insertar (paso a paso)»`,
    },
  ]
}

/** Paso por fila leída de la BD, antes de la animación de inserción. */
export function calcularPasoRegistroBD(numero, total, cancion) {
  const id = idDeCancion(cancion)
  const nombre = nombreDeCancion(cancion)
  const artista = artistaDeCancion(cancion)
  const detalle = nombre
    ? `track_name='${nombre}'${artista ? `, artists='${artista}'` : ''}`
    : 'datos de la fila en memoria'
  return {
    nodeId: 0,
    tipo: 'init_bd_fila',
    mensaje: `🗂 Fila ${numero}/${total} desde BD: id=${id} (${detalle}) → siguiente paso: insertar en B+`,
    activeKey: id,
  }
}

export function calcularPasoEncabezadoCancion(numero, total, cancion) {
  const id = idDeCancion(cancion)
  const nombre = nombreDeCancion(cancion)
  return {
    nodeId: 0,
    tipo: 'init_cancion',
    mensaje: `── Inserción ${numero}/${total}: clave id=${id}${nombre ? ` — «${nombre}»` : ''}`,
    activeKey: id,
  }
}

export function calcularPasosFinInicializacion(total, ms) {
  return [
    {
      nodeId: 0,
      tipo: 'confirmacion',
      mensaje: `✅ Árbol construido: ${total} canciones del dataset insertadas en B+ (${ms.toFixed(0)} ms)`,
    },
  ]
}
