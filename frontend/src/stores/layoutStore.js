import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
  const layoutMode = ref(localStorage.getItem('app-layout-mode') || 'top-nav')

  const setLayoutMode = (mode) => {
    layoutMode.value = mode
    localStorage.setItem('app-layout-mode', mode)
    // Dispatch storage event to notify other components if needed (e.g. across tabs)
    window.dispatchEvent(new Event("storage"))
  }

  const toggleLayoutMode = () => {
    const nextMode = layoutMode.value === 'sidebar' ? 'top-nav' : 'sidebar'
    setLayoutMode(nextMode)
  }

  return {
    layoutMode,
    setLayoutMode,
    toggleLayoutMode
  }
})
