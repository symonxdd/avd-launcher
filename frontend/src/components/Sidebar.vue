<template>
  <aside :class="styles.sidebar">
    <div :class="styles.topGroup">
      <div :class="styles.logoContainer">
        <img src="../assets/images/appicon-no-bg.png" alt="App logo" :class="styles.logoImage" />
        <div :class="styles.logoTextContainer">
          <span :class="styles.logoText">AVD Launcher</span>
          <span v-if="isDev" :class="styles.devIndicator">(dev)</span>
        </div>
      </div>
      <div :class="styles.sidebarTop">
        <router-link to="/" :class="[styles.sidebarItem, { [styles.active]: $route.path === '/' }]">
          <v-icon name="hi-device-mobile" :class="styles.icon" />
          <span>AVD</span>
        </router-link>

        <router-link to="/logs" :class="[styles.sidebarItem, { [styles.active]: $route.path === '/logs' }]">
          <v-icon name="hi-document-text" :class="styles.icon" />
          <span>Logs</span>
        </router-link>
      </div>
    </div>

    <div :class="styles.bottomGroup">
      <div :class="styles.sidebarBottom">
        <router-link to="/settings" :class="[styles.sidebarItem, { [styles.active]: $route.path === '/settings' }]">
          <v-icon name="hi-cog" :class="styles.icon" />
          <span>Settings</span>
        </router-link>
      </div>
    </div>
  </aside>
</template>

<script setup>
import styles from './Sidebar.module.css';
// Icons are registered globally in main.js
const isDev = import.meta.env.MODE === 'development';
</script>
