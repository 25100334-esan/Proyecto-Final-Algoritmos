<script setup>
import { ref } from 'vue'
import { useDebugLog } from '../../composables/useDebugLog.js'

const { logs, limpiar } = useDebugLog()
const colapsado = ref(false)
</script>

<template>
  <footer class="debug-console" :class="{ 'debug-console--collapsed': colapsado }">
    <header class="debug-console__header" @click="colapsado = !colapsado">
      <span>Consola de Rendimiento — B+-Tree en RAM (Comer)</span>
      <div class="header-actions">
        <button class="btn-clear" @click.stop="limpiar">Limpiar</button>
        <span class="toggle">{{ colapsado ? '▲' : '▼' }}</span>
      </div>
    </header>

    <div v-if="!colapsado" class="debug-console__body">
      <div v-for="log in logs" :key="log.id" class="log-line" :class="'log-line--' + log.tipo">
        <span class="log-time">[{{ log.hora }}]</span>
        {{ log.mensaje }}
      </div>
      <div v-if="!logs.length" class="log-empty">
        Realiza una búsqueda o filtra por popularidad para ver métricas O(log n) y Sequence Set…
      </div>
    </div>
  </footer>
</template>

<style scoped>
.debug-console {
  background: #0a0a0a;
  border-top: 1px solid var(--border);
  font-family: 'Consolas', 'Courier New', monospace;
  font-size: 0.78rem;
  flex-shrink: 0;
}

.debug-console--collapsed {
  max-height: 32px;
  overflow: hidden;
}

.debug-console__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.4rem 1rem;
  background: #1a1a1a;
  cursor: pointer;
  user-select: none;
  color: var(--text-secondary);
  font-size: 0.75rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.btn-clear {
  background: transparent;
  color: var(--text-muted);
  font-size: 0.7rem;
  padding: 0.15rem 0.5rem;
  border: 1px solid var(--border);
}

.btn-clear:hover {
  color: var(--text-primary);
}

.debug-console__body {
  max-height: 120px;
  overflow-y: auto;
  padding: 0.5rem 1rem;
}

.log-line {
  padding: 0.15rem 0;
  color: #ccc;
}

.log-line--success { color: #4ade80; }
.log-line--error { color: #f87171; }
.log-line--info { color: #93c5fd; }
.log-line--warn { color: #fbbf24; }

.log-time {
  color: #666;
  margin-right: 0.5rem;
}

.log-empty {
  color: #555;
  font-style: italic;
}
</style>
