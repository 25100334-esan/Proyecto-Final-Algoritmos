<script setup>
import { formatearDuracion } from '../../utils/formatters.js'

defineProps({
  canciones: { type: Array, default: () => [] },
  titulo: { type: String, default: 'Playlist' },
})

defineEmits(['play'])
</script>

<template>
  <section class="playlist">
    <header class="playlist__header">
      <div class="playlist__cover">♫</div>
      <div>
        <p class="playlist__type">Playlist</p>
        <h2>{{ titulo }}</h2>
        <p class="playlist__count">{{ canciones.length }} canciones</p>
      </div>
    </header>

    <table class="playlist__table">
      <thead>
        <tr>
          <th>#</th>
          <th>Título</th>
          <th>Artista</th>
          <th>Álbum</th>
          <th>Duración</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(cancion, i) in canciones"
          :key="cancion.Indice"
          class="track-row"
          @dblclick="$emit('play', cancion)"
        >
          <td>{{ i + 1 }}</td>
          <td class="track-name">{{ cancion.TrackName }}</td>
          <td>{{ cancion.Artists }}</td>
          <td class="album-col">{{ cancion.AlbumName }}</td>
          <td>{{ formatearDuracion(cancion.DurationMs) }}</td>
          <td>
            <button class="mini-play" @click="$emit('play', cancion)">▶</button>
          </td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style scoped>
.playlist__header {
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
  margin-bottom: 1.5rem;
}

.playlist__cover {
  width: 180px;
  height: 180px;
  background: linear-gradient(160deg, #535353, #121212);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}

.playlist__type {
  font-size: 0.75rem;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.1em;
}

.playlist__header h2 {
  font-size: 2.5rem;
  font-weight: 900;
  margin: 0.25rem 0;
}

.playlist__count {
  color: var(--text-secondary);
  font-size: 0.85rem;
}

.playlist__table {
  width: 100%;
  border-collapse: collapse;
}

.playlist__table th {
  text-align: left;
  padding: 0.5rem 0.75rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid var(--border);
}

.playlist__table td {
  padding: 0.6rem 0.75rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.track-row {
  transition: background 0.15s;
}

.track-row:hover {
  background: rgba(255, 255, 255, 0.05);
}

.track-name {
  color: var(--text-primary) !important;
  font-weight: 500;
}

.album-col {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mini-play {
  background: transparent;
  color: var(--text-primary);
  font-size: 0.85rem;
  opacity: 0;
  padding: 0.25rem 0.5rem;
}

.track-row:hover .mini-play {
  opacity: 1;
}

.mini-play:hover {
  color: var(--accent);
}
</style>
