<template>
  <nav class="sidebar">
    <TitleLogo />

    <ul class="nav-list">
      <li class="nav-item" :class="{ active: $route.path === '/' }" @click="$router.push('/')">
        <i class="bi bi-house-door nav-icon"></i>
        <span class="nav-text">AVDs</span>
      </li>
      <li class="nav-item" :class="{ active: $route.path === '/logs' }" @click="$router.push('/logs')">
        <i class="bi bi-terminal nav-icon"></i>
        <span class="nav-text">Logs</span>
      </li>
    </ul>

    <!-- Settings at the bottom -->
    <div class="nav-bottom">
      <li class="nav-item" :class="{ active: $route.path === '/settings' }" @click="$router.push('/settings')"
        :style="{ marginBottom: '0' }">
        <i class="bi bi-gear nav-icon"></i>
        <span class="nav-text">Settings</span>
      </li>
    </div>
  </nav>
</template>

<script setup>
import TitleLogo from './TitleLogo.vue'
</script>

<style scoped>
.sidebar {
  width: 195px;
  min-width: 195px;
  background-color: transparent;
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 15px;
}

.nav-list {
  flex-grow: 1;
  /* 👈 Important: pushes the nav-bottom to the bottom */
  list-style: none;
  padding: 0;
  margin-top: 30px;
}

.nav-bottom {
  margin-top: auto;
  /* 👈 Pushes this div to the bottom */
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  color: #919191;
  font-size: 0.95rem;
  border-radius: 6px;
  margin-bottom: 10px;
  transition: background 0.2s, color 0.2s;
}

.nav-item:hover {
  background-color: #252525;
  color: #d3d3d3;
}

.nav-item.active:hover {
  background-color: #252525;
  color: #d3d3d3;
}

.nav-item.active {
  background-color: transparent;
  color: #d3d3d3;
  font-weight: bold;
}

.nav-icon {
  font-size: 1rem;
  margin-right: 8px;
}
</style>
