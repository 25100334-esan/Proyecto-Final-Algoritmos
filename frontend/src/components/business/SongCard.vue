<script setup>
import { formatearDuracion } from '../../utils/formatters.js'

defineProps({
  cancion: { type: Object, required: true },
  playing: { type: Boolean, default: false },
})

defineEmits(['play'])
</script>

<template>
  <article class="song-card">
    <div class="song-card__cover">
      <div class="cover-art">
        <span class="cover-icon">♪</span>
      </div>
      <button
        class="play-btn"
        :class="{ 'play-btn--active': playing }"
        @click="$emit('play', cancion)"
      >
        {{ playing ? '⏸' : '▶' }}
      </button>
    </div>
    <div class="song-card__info">
      <h3>{{ cancion.TrackName }}</h3>
      <p class="artist">{{ cancion.Artists }}</p>
      <p class="album">{{ cancion.AlbumName }}</p>
      <div class="meta">
        <span>ID: {{ cancion.Indice }}</span>
        <span>{{ formatearDuracion(cancion.DurationMs) }}</span>
        <span v-if="cancion.TrackGenre">{{ cancion.TrackGenre }}</span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.song-card {
  display: flex;
  gap: 1.25rem;
  background: var(--bg-card);
  border-radius: 8px;
  padding: 1.25rem;
  transition: background 0.2s;
}

.song-card:hover {
  background: #333;
}

.song-card__cover {
  position: relative;
  flex-shrink: 0;
}

.cover-art {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, #1db954 0%, #191414 100%);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-icon {
  font-size: 3rem;
  color:                rgba(255, 255, 255, 0.3);
}

.play-btn {
  position: absolute;
  bottom: -8px;
  right: -8px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: var(--accent);
  color: #000;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.2s, transform 0.2s;
}

.song-card:hover .play-btn,
.play-btn--active {
  opacity: 1;
  transform: translateY(0);
}

.play-btn:hover {
  background: var(--accent-hover);
  transform: scale(1.05);
}

.song-card__info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 0.25rem;
}

.song-card__info h3 {
  font-size: 1.5rem;
  font-weight: 700;
}

.artist {
  color: var(--text-secondary);
  font-size: 1rem;
}

.album {
  color: var(--text-muted);
  font-size: 0.85rem;
}

.meta {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: var(--text-muted);
}
</style>
