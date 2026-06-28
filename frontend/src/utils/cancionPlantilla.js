/** Valores sugeridos para insertar una canción simulada según el próximo id. */
export function plantillaCancionSimulada(proximoId) {
  const n = proximoId ?? 1
  const generos = ['pop', 'rock', 'latin', 'jazz', 'acoustic', 'electronic']
  return {
    TrackID: `sim_${n}`,
    TrackName: `Canción Simulada ${n}`,
    Artists: 'Artista Demo',
    AlbumName: `Álbum Test ${Math.ceil(n / 5)}`,
    Popularity: Math.min(95, 40 + (n % 56)),
    DurationMs: 180000 + (n % 120) * 1000,
    Explicit: n % 7 === 0,
    Danceability: Number((0.5 + (n % 50) / 100).toFixed(2)),
    Energy: Number((0.4 + (n % 60) / 100).toFixed(2)),
    Key: n % 12,
    Loudness: Number((-8 + (n % 10) * 0.3).toFixed(1)),
    Mode: n % 2,
    Speechiness: Number((0.03 + (n % 20) / 200).toFixed(3)),
    Acousticness: Number((0.05 + (n % 30) / 100).toFixed(2)),
    Instrumentalness: Number(((n % 10) / 100).toFixed(3)),
    Liveness: Number((0.08 + (n % 25) / 100).toFixed(2)),
    Valence: Number((0.45 + (n % 40) / 100).toFixed(2)),
    Tempo: Number((90 + (n % 80)).toFixed(1)),
    TimeSignature: [3, 4, 5][n % 3],
    TrackGenre: generos[n % generos.length],
  }
}
