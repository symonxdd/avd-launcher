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
import { useAvdStore } from './stores/avdStore'

const store = useAvdStore()
const mounted = ref(false)

onMounted(() => {
  store.startAvdWatcher()

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
  color: #e9e9e9;
}

.main-content {
  flex: 1;
  padding: 1rem;
}

.fade-app-enter-active {
  transition: opacity 1s ease;
}

.fade-app-enter-from {
  opacity: 0;
}
</style>
