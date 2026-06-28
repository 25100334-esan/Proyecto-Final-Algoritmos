<script setup>
import SpotifyLogo from '../icons/SpotifyLogo.vue'

defineProps({
  cancion: { type: Object, required: true },
  playing: { type: Boolean, default: false },
})

defineEmits(['play'])
</script>

<template>
  <article class="grid-card">
    <div class="grid-card__cover">
      <div class="cover-art">
        <SpotifyLogo :size="36" color="rgba(255,255,255,0.35)" />
      </div>
      <button
        class="play-btn"
        :class="{ 'play-btn--active': playing }"
        aria-label="Reproducir"
        @click="$emit('play', cancion)"
      >
        {{ playing ? '⏸' : '▶' }}
      </button>
    </div>
    <div class="grid-card__body">
      <h3 class="track-name" :title="cancion.TrackName">{{ cancion.TrackName }}</h3>
      <p class="artist" :title="cancion.Artists">{{ cancion.Artists }}</p>
      <div class="badges">
        <span v-if="cancion.Tempo != null" class="badge badge--bpm">
          BPM: {{ Number(cancion.Tempo).toFixed(1) }}
        </span>
        <span v-if="cancion.Popularity != null" class="badge badge--pop">
          Popularidad: {{ cancion.Popularity }}
        </span>
      </div>
    </div>
  </article>
</template>

<style scoped>
.grid-card {
  background: var(--bg-card);
  border-radius: 8px;
  padding: 1rem;
  transition: background 0.2s, transform 0.15s;
  cursor: default;
}

.grid-card:hover {
  background: #333;
}

.grid-card__cover {
  position: relative;
  margin-bottom: 0.85rem;
}

.cover-art {
  aspect-ratio: 1;
  width: 100%;
  background: linear-gradient(135deg, #1db954 0%, #191414 100%);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.35);
}

.play-btn {
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--accent);
  color: #000;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.45);
  opacity: 0;
  transform: translateY(6px);
  transition: opacity 0.2s, transform 0.2s, background 0.15s;
}

.grid-card:hover .play-btn,
.play-btn--active {
  opacity: 1;
  transform: translateY(0);
}

.play-btn:hover {
  background: var(--accent-hover);
  transform: scale(1.06);
}

.track-name {
  font-size: 0.95rem;
  font-weight: 700;
  line-height: 1.3;
  margin-bottom: 0.35rem;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.artist {
  font-size: 0.8rem;
  color: var(--text-secondary);
  margin-bottom: 0.6rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.35rem;
}

.badge {
  font-size: 0.65rem;
  font-weight: 600;
  padding: 0.2rem 0.45rem;
  border-radius: 4px;
  letter-spacing: 0.02em;
}

.badge--bpm {
  background: rgba(99, 102, 241, 0.25);
  color: #a5b4fc;
}

.badge--pop {
  background: rgba(29, 185, 84, 0.2);
  color: #86efac;
}
</style>
