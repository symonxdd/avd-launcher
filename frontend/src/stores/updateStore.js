import { defineStore } from 'pinia'
import { ref } from 'vue'
import { GetLatestVersion } from '../../wailsjs/go/services/UpdateService'

export const useUpdateStore = defineStore('updateStore', () => {
  const remoteVersion = ref(null)
  const remoteReleaseUrl = ref(null)
  const isCheckingForUpdate = ref(false)
  const lastChecked = ref(null)

  const checkForUpdate = async (force = false) => {
    // If already checked and not forcing, skip
    if (lastChecked.value && !force) return

    isCheckingForUpdate.value = true
    try {
      const release = await GetLatestVersion()
      if (!release?.tag_name) return

      const cleanTag = release.tag_name.startsWith("v")
        ? release.tag_name.slice(1)
        : release.tag_name

      remoteVersion.value = cleanTag
      remoteReleaseUrl.value = release.html_url
      lastChecked.value = Date.now()
    } catch (err) {
      console.warn("Version check failed:", err)
    } finally {
      isCheckingForUpdate.value = false
    }
  }

  return {
    remoteVersion,
    remoteReleaseUrl,
    isCheckingForUpdate,
    lastChecked,
    checkForUpdate
  }
})
