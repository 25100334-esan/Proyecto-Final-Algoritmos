<script setup>
import { ref } from 'vue'
import AcademicMode from './views/AcademicMode.vue'
import BusinessMode from './views/BusinessMode.vue'
import SpotifyLogo from './components/icons/SpotifyLogo.vue'

const vistaActiva = ref('academic')

const tabs = [
  { id: 'academic', label: 'Modo Académico', icon: '🎓' },
  { id: 'business', label: 'Modo Negocio', icon: '🎵' },
]
</script>

<template>
  <div id="app-root">
    <header class="app-header">
      <div class="brand">
        <SpotifyLogo class="brand-icon" :size="28" color="#1db954" />
        <div>
          <h1>Spotify B+ Tree</h1>
          <p>Visualización interactiva del árbol B+ — Proyecto Final Algoritmos</p>
        </div>
      </div>

      <nav class="tab-nav">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn"
          :class="{ 'tab-btn--active': vistaActiva === tab.id }"
          @click="vistaActiva = tab.id"
        >
          {{ tab.icon }} {{ tab.label }}
        </button>
      </nav>
    </header>

    <main class="app-main">
      <AcademicMode v-show="vistaActiva === 'academic'" />
      <BusinessMode v-show="vistaActiva === 'business'" />
    </main>
  </div>
</template>

<style scoped>
#app-root {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.5rem;
  background: var(--bg-elevated);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
  gap: 1rem;
  flex-wrap: wrap;
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.brand-icon {
  flex-shrink: 0;
}

.brand h1 {
  font-size: 1.1rem;
  font-weight: 700;
}

.brand p {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.tab-nav {
  display: flex;
  gap: 0.5rem;
}

.tab-btn {
  padding: 0.5rem 1rem;
  background: transparent;
  color: var(--text-secondary);
  font-size: 0.85rem;
  border-radius: 500px;
  border: 1px solid transparent;
}

.tab-btn:hover {
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.05);
}

.tab-btn--active {
  background: var(--bg-card);
  color: var(--text-primary);
  border-color: var(--border);
  font-weight: 600;
}

.app-main {
  flex: 1;
  overflow: hidden;
  min-height: 0;
}
</style>
