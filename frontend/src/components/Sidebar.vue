<template>
  <aside class="sidebar">
    <div class="topGroup">
      <div class="logoContainer">
        <img src="../assets/images/appicon-no-bg.png" alt="App logo" class="logoImage" />
        <div class="logoTextContainer">
          <span class="logoText">AVD Launcher</span>
          <span v-if="isDev" class="devIndicator">(dev)</span>
        </div>
      </div>
      <div class="sidebar-top">
        <router-link to="/" class="sidebar-item" :class="{ active: $route.path === '/' }">
          <v-icon name="hi-device-mobile" class="icon" />
          <span>AVD</span>
        </router-link>

        <router-link to="/logs" class="sidebar-item" :class="{ active: $route.path === '/logs' }">
          <v-icon name="hi-document-text" class="icon" />
          <span>Logs</span>
        </router-link>
      </div>
    </div>

    <div class="bottomGroup">
      <div class="sidebar-bottom">
        <router-link to="/settings" class="sidebar-item" :class="{ active: $route.path === '/settings' }">
          <v-icon name="hi-cog" class="icon" />
          <span>Settings</span>
        </router-link>
      </div>
    </div>
  </aside>
</template>

<script setup>
// Icons are registered globally in main.js
const isDev = import.meta.env.MODE === 'development';
</script>

<style scoped>
.sidebar {
  width: 230px;
  background-color: var(--bg-app);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 0.5rem 0;
  height: 100%;
}

.logoContainer {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 1.5rem 1.25rem 2rem 1.25rem;
  transition: transform 0.3s ease-in-out;
  transform: scale(1);
}

.logoContainer:hover {
  animation: pulseGlow 3s ease-in-out infinite;
}

@keyframes pulseGlow {

  0%,
  100% {
    transform: scale(1);
    filter: brightness(1);
  }

  50% {
    transform: scale(1.05);
    filter: brightness(1.2);
  }
}

.logoImage {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.logoText {
  font-size: 1.2rem;
  font-weight: 500;
  color: var(--text-primary);
}

.devIndicator {
  position: absolute;
  top: 100%;
  left: 0;
  font-size: 0.7rem;
  color: var(--text-secondary);
  font-weight: 400;
  line-height: 1;
  margin-top: -4px;
  white-space: nowrap;
}

.logoTextContainer {
  position: relative;
  display: flex;
  align-items: center;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  margin: 4px 8px;
  color: var(--text-secondary);
  text-decoration: none;
  font-weight: 500;
  font-size: 0.95rem;
  border: 1px solid transparent;
  border-radius: 16px;
  cursor: default;
  -webkit-user-drag: none;
  position: relative;
  transition: background 0.2s ease, color 0.2s ease, border 0.2s ease;
}

.sidebar-item > .icon,
.sidebar-item > span {
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1), color 0.2s ease;
}

.sidebar-item:hover > .icon,
.sidebar-item:hover > span {
  transform: translateX(4px);
}

.sidebar-item:hover {
  background: var(--bg-card-active);
  color: var(--text-primary);
  border: 1px solid var(--border-strong);
}

.sidebar-item.active {
  background: var(--bg-card-active);
  color: var(--text-primary);
  border: 1px solid var(--border-strong);
}

.icon {
  font-size: 18px;
  color: var(--text-secondary);
  width: 18px;
  height: 18px;
}

.sidebar-item:hover .icon,
.sidebar-item.active .icon {
  color: var(--text-primary);
}

.topGroup {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.bottomGroup {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
