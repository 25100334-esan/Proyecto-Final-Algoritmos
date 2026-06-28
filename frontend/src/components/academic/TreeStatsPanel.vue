<script setup>
defineProps({
  stats: { type: Object, default: null },
  ultimaOperacion: { type: String, default: '' },
  tipoVista: { type: String, default: 'indice' },
  limiteVisualizacion: { type: Number, default: 500 },
})
</script>

<template>
  <div class="tree-stats">
    <h3>Estado del B+-Tree en RAM</h3>
    <p class="tree-stats__aviso">
      Canvas desactivado (más de {{ limiteVisualizacion }} canciones). El motor Go sigue activo.
    </p>

    <p v-if="!stats" class="tree-stats__cargando">Calculando estadísticas del árbol…</p>

    <div v-if="stats" class="tree-stats__grid">
      <div class="stat-card">
        <span class="stat-label">Vista / índice</span>
        <strong>{{ stats.etiquetaIndice }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Canciones</span>
        <strong>{{ stats.totalCanciones }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Nodos internos</span>
        <strong>{{ stats.nodosInternos }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Nodos hoja</span>
        <strong>{{ stats.nodosHoja }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Nodos totales</span>
        <strong>{{ stats.nodosTotales }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Altura</span>
        <strong>{{ stats.altura }}</strong>
      </div>
      <div class="stat-card">
        <span class="stat-label">Entradas en hojas</span>
        <strong>{{ stats.entradasHoja }}</strong>
        <span v-if="stats.entradasHoja > stats.totalCanciones" class="stat-note">claves duplicadas</span>
      </div>
    </div>

    <div v-if="ultimaOperacion" class="tree-stats__op">
      <span class="stat-label">Última operación</span>
      <p>{{ ultimaOperacion }}</p>
    </div>
  </div>
</template>

<style scoped>
.tree-stats {
  flex: 1;
  padding: 1.25rem 1.5rem;
  background: var(--bg-panel);
  border-radius: 8px;
  border: 1px solid var(--border);
  min-height: 320px;
}

.tree-stats h3 {
  margin: 0 0 0.5rem;
  font-size: 1rem;
  color: var(--text-primary);
}

.tree-stats__aviso {
  margin: 0 0 1rem;
  font-size: 0.72rem;
  color: #d97706;
  line-height: 1.4;
}

.tree-stats__cargando {
  margin: 0 0 1rem;
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.tree-stats__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 0.6rem;
  margin-bottom: 1rem;
}

.stat-card {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 6px;
  padding: 0.55rem 0.65rem;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.stat-label {
  font-size: 0.65rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--text-secondary);
}

.stat-card strong {
  font-size: 1.15rem;
  color: var(--accent-academic, #818cf8);
}

.stat-note {
  font-size: 0.62rem;
  color: var(--text-secondary);
}

.tree-stats__op {
  padding: 0.65rem 0.75rem;
  background: rgba(99, 102, 241, 0.12);
  border-radius: 6px;
  border-left: 3px solid var(--accent-academic, #818cf8);
}

.tree-stats__op p {
  margin: 0.25rem 0 0;
  font-size: 0.8rem;
  color: var(--text-primary);
  line-height: 1.45;
}
</style>
