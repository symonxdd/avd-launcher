<template>
  <transition name="fade-app">
    <div :class="styles.appContainer" v-if="mounted">
      <TopBar />
      <div :class="styles.appBody">
        <Sidebar />
        <div :class="styles.mainContent">
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
import styles from './App.module.css'

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
