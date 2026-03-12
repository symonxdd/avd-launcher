<template>
  <transition name="fade-app">
    <div :class="styles.appContainer" v-if="mounted">
      <TopBar />
      <div :class="[styles.appBody, layoutStore.layoutMode === 'top-nav' && styles.topNavMode, 'layout-' + layoutStore.layoutMode]">
        <Sidebar v-if="layoutStore.layoutMode === 'sidebar'" />
        <div :class="[styles.mainContent, layoutStore.layoutMode === 'top-nav' && styles.topNavContent]">
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
import { useLayoutStore } from './stores/layoutStore'
import styles from './App.module.css'

const store = useAvdStore()
const layoutStore = useLayoutStore()
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
