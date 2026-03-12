<template>
  <button @click="layoutStore.toggleLayoutMode" @dblclick.stop :class="[styles.toggleBtn, className]"
    aria-label="Toggle layout mode">
    <div :class="styles.iconContainer" :style="iconStyle">
      <v-icon name="fa-wind" :scale="1.0" :class="styles.icon" />
    </div>
  </button>
</template>

<script setup>
import { computed } from 'vue';
import { useLayoutStore } from '../stores/layoutStore';
import styles from './LayoutToggle.module.css';

defineProps({
  className: {
    type: String,
    default: ''
  }
});

const layoutStore = useLayoutStore();

const iconStyle = computed(() => ({
  transform: layoutStore.layoutMode === 'sidebar' ? 'rotate(0deg)' : 'rotate(90deg)',
  transition: 'transform 0.3s ease-in-out'
}));
</script>
