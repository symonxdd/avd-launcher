<template>
  <transition name="fade-app">
    <div class="app" v-if="mounted">
      <Sidebar />
      <div class="main-content">
        <router-view />
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from './components/Sidebar.vue'
import { WindowSetTitle } from '../wailsjs/runtime/runtime'

const mounted = ref(false)
onMounted(() => {
  mounted.value = true

  const title = import.meta.env.MODE === 'development'
    ? 'AVD Launcher (dev)'
    : 'AVD Launcher'
  WindowSetTitle(title)
})
</script>

<style>
.app {
  display: flex;
  height: 100vh;
  background: #131313;
  color: #fff;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.fade-app-enter-active {
  transition: opacity 0.8s ease;
}

.fade-app-enter-from {
  opacity: 0;
}
</style>
