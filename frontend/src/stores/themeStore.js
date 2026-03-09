import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const theme = ref('system') // 'light', 'dark', or 'system'
  const trueBlack = ref(false)
  let hasSetupSystemListener = false

  const applyTheme = () => {
    const root = document.documentElement
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    const themeToApply =
      theme.value === 'system' ? (systemPrefersDark ? 'dark' : 'light') : theme.value

    root.classList.remove('theme-light', 'theme-dark')
    root.classList.add(`theme-${themeToApply}`)

    applyTrueBlack(themeToApply)

    if (!hasSetupSystemListener) {
      window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
        if (theme.value === 'system') {
          applyTheme()
        }
      })
      hasSetupSystemListener = true
    }
  }

  const applyTrueBlack = (resolvedTheme = null) => {
    const root = document.documentElement
    const currentTheme =
      resolvedTheme ??
      (theme.value === 'system'
        ? window.matchMedia('(prefers-color-scheme: dark)').matches
          ? 'dark'
          : 'light'
        : theme.value)
    if (trueBlack.value && currentTheme === 'dark') {
      root.classList.add('true-black')
    } else {
      root.classList.remove('true-black')
    }
  }

  const initTheme = () => {
    const savedTheme = localStorage.getItem('app-theme')
    const savedTrueBlack = localStorage.getItem('app-true-black')

    theme.value = savedTheme || 'system'
    trueBlack.value = savedTrueBlack === 'true'

    applyTheme()
  }

  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('app-theme', newTheme)
    applyTheme()
  }

  const toggleTrueBlack = () => {
    trueBlack.value = !trueBlack.value
    localStorage.setItem('app-true-black', trueBlack.value)
    applyTrueBlack()
  }

  return {
    theme,
    trueBlack,
    initTheme,
    setTheme,
    toggleTrueBlack,
    applyTheme,
  }
})
