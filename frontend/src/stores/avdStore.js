import { defineStore } from 'pinia'
import { reactive, ref } from 'vue'

export const useAvdStore = defineStore('avd', () => {
  const avds = ref([])
  const logs = ref('')

  function updateAvdStatus(avdName, status) {
    const idx = avds.value.findIndex(a => a.name === avdName)
    if (idx !== -1) {
      const updated = { ...avds.value[idx], ...status }
      avds.value.splice(idx, 1, updated) // replace with new object to trigger reactivity
    }
  }

  function appendLog(message) {
    logs.value += message + '\n'
  }

  function resetLogs() {
    logs.value = ''
  }

  return {
    avds,
    logs,
    updateAvdStatus,
    appendLog,
    resetLogs
  }
})
