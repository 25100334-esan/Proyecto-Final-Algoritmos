<script setup>
defineProps({
  mensaje: { type: String, default: '' },
  visible: { type: Boolean, default: false },
  pasoActual: { type: Number, default: 0 },
  pasoTotal: { type: Number, default: 0 },
  mostrarSiguiente: { type: Boolean, default: false },
})

const emit = defineEmits(['siguiente'])
</script>

<template>
  <Transition name="fade">
    <div v-if="visible && mensaje" class="step-banner-wrap">
      <div class="step-banner">
        <span v-if="pasoTotal > 0" class="step-counter">
          Paso {{ pasoActual }} / {{ pasoTotal }}
        </span>
        <span class="step-text">{{ mensaje }}</span>
      </div>
      <button
        v-if="mostrarSiguiente"
        class="step-next"
        type="button"
        @click="emit('siguiente')"
      >
        Siguiente →
      </button>
    </div>
  </Transition>
</template>

<style scoped>
.step-banner-wrap {
  position: absolute;
  top: 12px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  max-width: 92%;
}

.step-banner {
  background: rgba(99, 102, 241, 0.95);
  color: #fff;
  padding: 0.6rem 1.25rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.step-counter {
  font-size: 0.7rem;
  opacity: 0.85;
  font-weight: 600;
  letter-spacing: 0.04em;
}

.step-text {
  line-height: 1.35;
}

.step-next {
  padding: 0.45rem 1.25rem;
  background: #fff;
  color: #4338ca;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.25);
}

.step-next:hover {
  background: #eef2ff;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s, transform 0.25s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-8px);
}
</style>
