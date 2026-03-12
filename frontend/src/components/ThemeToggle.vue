<template>
  <button @click="toggleTheme" @dblclick.stop :class="styles.toggleBtn" aria-label="Toggle theme">
    <transition name="sunmoon" mode="out-in">
      <div :key="themeStore.theme" :class="styles.iconContainer">
        <v-icon v-if="isDark" name="hi-moon" :scale="1.1" :class="styles.icon" />
        <v-icon v-else name="hi-sun" :scale="1.1" :class="styles.icon" />
      </div>
    </transition>
  </button>
</template>

<script setup>
import { computed } from 'vue';
import { useThemeStore } from '../stores/themeStore';
import styles from './ThemeToggle.module.css';

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
