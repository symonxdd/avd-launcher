<template>
  <div :class="styles.topBar" @dblclick="WindowToggleMaximise">
    <div :class="styles.leftActions">
      <ThemeToggle :class="styles.topBtn" />
      <LayoutToggle :class="styles.topBtn" />
    </div>

    <!-- Centered Branding (Top Nav mode only) -->
    <div v-if="layoutStore.layoutMode === 'top-nav'" :class="styles.brandingCenter">
      <img src="../assets/images/appicon-no-bg.png" alt="App logo" :class="styles.appIcon" />
      <div :class="styles.logoTextContainer">
        <span :class="styles.appTitle">AVD Launcher</span>
        <span v-if="isDev" :class="styles.devIndicator">(dev)</span>
      </div>
    </div>

    <div
      v-if="layoutStore.layoutMode === 'top-nav'"
      :class="styles.navPillContainer"
      @dblclick.stop
    >
      <nav :class="styles.pill">
        <router-link
          v-for="link in navigationLinks"
          :key="link.href"
          :to="link.href"
          :class="[styles.navLink, isActive(link.href) && styles.active]"
          draggable="false"
        >
          <div v-if="isActive(link.href)" :class="styles.activePill" />
          <span :class="styles.navLinkText">{{ link.label }}</span>
        </router-link>
      </nav>
    </div>

    <div :class="styles.winControls">
      <button :class="[styles.winBtn, styles.minimize]" @click="WindowMinimise" title="Minimize"></button>
      <button :class="[styles.winBtn, styles.maximize]" @click="WindowToggleMaximise" title="Maximize"></button>
      <button :class="[styles.winBtn, styles.close]" @click="Quit" title="Close"></button>
    </div>
  </div>
</template>

<script setup>
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime';
import { useRoute } from 'vue-router';
import { useLayoutStore } from '../stores/layoutStore';
import ThemeToggle from './ThemeToggle.vue';
import LayoutToggle from './LayoutToggle.vue';
import styles from './TopBar.module.css';

const route = useRoute();
const layoutStore = useLayoutStore();
const isDev = import.meta.env.MODE === 'development';

const navigationLinks = [
  { href: '/', label: 'Home' },
  { href: '/logs', label: 'Logs' },
  { href: '/settings', label: 'Settings' },
];

const isActive = (path) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};
</script>
