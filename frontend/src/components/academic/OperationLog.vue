<script setup>
defineProps({
  mensajes: { type: Array, default: () => [] },
})
</script>

<template>
  <div class="operation-log">
    <h4>Registro de operaciones</h4>
    <ul>
      <li v-for="msg in mensajes" :key="msg.id" :class="'log--' + msg.tipo">
        <span class="log-time">{{ msg.hora }}</span>
        {{ msg.texto }}
      </li>
      <li v-if="!mensajes.length" class="log-empty">Sin operaciones aún</li>
    </ul>
  </div>
</template>

<style scoped>
.operation-log {
  background: var(--bg-panel);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  max-height: 140px;
  overflow-y: auto;
}

.operation-log h4 {
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  margin-bottom: 0.5rem;
}

.operation-log ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.operation-log li {
  font-size: 0.8rem;
  font-family: 'Consolas', monospace;
  color: var(--text-secondary);
}

.log-time {
  color: var(--text-muted);
  margin-right: 0.5rem;
}

.log--success { color: #34d399; }
.log--error { color: #f87171; }
.log--info { color: #93c5fd; }
.log--warn { color: var(--warning); }
.log-empty { color: var(--text-muted); font-style: italic; }
</style>
