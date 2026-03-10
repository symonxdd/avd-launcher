<template>
  <button @click="toggleTheme" @dblclick.stop class="toggleBtn" aria-label="Toggle theme">
    <transition name="sunmoon" mode="out-in">
      <div :key="themeStore.theme" class="iconContainer">
        <v-icon v-if="isDark" name="hi-sun" :scale="1.1" class="icon" />
        <v-icon v-else name="hi-moon" :scale="1.1" class="icon" />
      </div>
    </transition>
  </button>
</template>

<script setup>
import { computed } from 'vue';
import { useThemeStore } from '../stores/themeStore';

const themeStore = useThemeStore();

const isDark = computed(() => {
  if (themeStore.theme === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
  }
  return themeStore.theme === 'dark';
});

const toggleTheme = () => {
  const nextTheme = isDark.value ? 'light' : 'dark';
  themeStore.setTheme(nextTheme);
};
</script>

<style scoped>
.toggleBtn {
  position: relative;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  cursor: default;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  -webkit-app-region: no-drag;
}

.toggleBtn:hover {
  background: var(--bg-card-active);
  color: var(--text-primary);
}

.iconContainer {
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  display: block;
}

/* Transition animations */
.sunmoon-enter-active,
.sunmoon-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.sunmoon-enter-from {
  opacity: 0;
  transform: scale(0.8) rotate(-45deg);
}

.sunmoon-leave-to {
  opacity: 0;
  transform: scale(0.8) rotate(45deg);
}
</style>
