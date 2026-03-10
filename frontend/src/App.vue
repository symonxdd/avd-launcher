<template>
  <transition name="fade-app">
    <div class="app-container" v-if="mounted">
      <TopBar />
      <div class="app-body">
        <Sidebar />
        <div class="main-content">
          <router-view />
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Sidebar from './components/Sidebar.vue'
import TopBar from './components/TopBar.vue'
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
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--bg-app);
  color: var(--text-primary);
  overflow: hidden;
}

.app-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.main-content {
  flex: 1;
  background-color: var(--bg-surface);
  border-top-left-radius: 16px;
  border-left: 1px solid var(--border-subtle);
  border-top: 1px solid var(--border-subtle);
  padding: 0 2rem 2rem;
  overflow-y: auto;
  position: relative;

  /* Hide scrollbar but keep functionality */
  scrollbar-width: none;
  /* Firefox */
  -ms-overflow-style: none;
  /* IE and Edge */
}

.main-content::-webkit-scrollbar {
  display: none;
  /* Chrome, Safari and Opera */
}

.fade-app-enter-active {
  transition: opacity 1s ease;
}

.fade-app-enter-from {
  opacity: 0;
}
</style>
