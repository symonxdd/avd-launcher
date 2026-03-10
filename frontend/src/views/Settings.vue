<template>
  <div :class="styles.settingsContainer">
    <h2 :class="styles.pageTitle">Settings</h2>
    <div :class="styles.settingsGrid">
      <!-- --- Interface Settings --- -->
      <div :class="styles.settingsGroup">
        <h5 :class="styles.labelHeading">
          <v-icon name="hi-view-grid" :class="styles.headingIcon" />
          Interface
        </h5>

        <div :class="styles.subGroupItems">
          <div v-if="isDarkMode" :class="styles.toggleRow">
            <label for="trueBlackToggle" :class="styles.toggleLabel">Enable True Black</label>
            <label :class="styles.switch">
              <input id="trueBlackToggle" type="checkbox" :checked="themeStore.trueBlack"
                @change="themeStore.toggleTrueBlack" />
              <span :class="styles.slider"></span>
            </label>
          </div>
          <div v-else :class="styles.emptyStateText">
            Switch to Dark Mode to enable True Black.
          </div>
        </div>
      </div>

      <!-- --- System Settings --- -->
      <div :class="styles.settingsGroup">
        <h5 :class="styles.labelHeading">
          <v-icon name="hi-adjustments" :class="styles.headingIcon" />
          System
        </h5>

        <div :class="styles.subGroupItems">
          <div v-if="androidSdkEnv"
            :class="[styles.sdkStatusRow, { [styles.sdkFound]: !!androidSdkEnv.path, [styles.sdkMissing]: !androidSdkEnv.path }]">
            <div :class="styles.statusInfo">
              <div :class="styles.statusLabelRow">
                <span :class="styles.statusLabel">Android SDK Status</span>
                <div v-if="androidSdkEnv.source" :class="styles.infoTooltipTrigger">
                  <v-icon name="hi-information-circle" :class="styles.infoIcon" />
                  <div :class="styles.infoTooltip">
                    <template v-if="androidSdkEnv.source === 'user selected path'">
                      Detected via
                      <span :class="styles.configLink" @click="OpenConfigFolder">
                        user selected path
                      </span>
                    </template>
                    <template v-else>
                      Detected via {{ androidSdkEnv.source }}
                    </template>
                  </div>
                </div>
              </div>
              <span :class="styles.statusPath">{{ androidSdkEnv.path || 'Not found' }}</span>
            </div>
            <div :class="styles.statusIconBox">
              <span v-if="androidSdkEnv.path">✔</span>
              <span v-else>✖</span>
            </div>
          </div>
        </div>
      </div>

      <!-- --- Update Available Notice --- -->
      <div v-if="remoteVersion && isRemoteVersionNewer(remoteVersion, appVersion)" :class="styles.settingsGroup">
        <div>
          <h5 :class="styles.labelHeading">
            <v-icon name="hi-gift" :class="styles.headingIcon" />
            New update ready
          </h5>
          <span :class="styles.versionSubtext" style="margin-left: calc(1.1rem + 14px);">v{{ remoteVersion }}</span>
        </div>

        <div :class="styles.subGroupItems">
          <div :class="styles.updateActionState">
            <button :class="styles.actionButton" @click="handleOpenGithubRelease">
              <v-icon name="fa-github" />
              View on GitHub
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- --- Footer Section --- -->
    <div :class="styles.appFooter">
      <div :class="styles.footerMeta">
        v{{ appVersion }} ({{ environment }}{{ remoteVersion ? (isUpToDate ? ', latest' : ', update available') : '' }})
        <span :class="styles.dotSeparator">•</span>
        Powered by Vue, Go & Wails
      </div>

      <div :class="styles.signature">
        <div :class="styles.tooltipTrigger">
          <span :class="styles.dictionaryLink">Sedulously</span>
          <div :class="styles.customTooltip">
            <span :class="styles.tooltipType">(adverb)</span> — in a way that shows great care and persistent effort.
            <br />
            <small style="opacity: 0.7; margin-top: 4px; display: block;">Synonyms: diligently, persistently</small>
          </div>
        </div> engineered with <span :class="styles.heart">❤️‍🩹</span> by Symon
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { GetAndroidSdkEnv, OpenConfigFolder } from '../../wailsjs/go/app/App'
import { GetLatestVersion } from '../../wailsjs/go/services/UpdateService'
import { BrowserOpenURL } from '../../wailsjs/runtime'
import { useThemeStore } from '../stores/themeStore'
import styles from './Settings.module.css'

const themeStore = useThemeStore()
const androidSdkEnv = ref(null)
// const appVersion = typeof __APP_VERSION__ !== "undefined" ? __APP_VERSION__ : "1.0.0"
const appVersion = "0.0.67"
const environment = import.meta.env.MODE === 'development' ? 'dev' : 'release'

const remoteVersion = ref(null)
const remoteReleaseUrl = ref(null)
const isCheckingForUpdate = ref(true)

const isDarkMode = computed(() => {
  if (themeStore.theme === 'dark') return true
  if (themeStore.theme === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return false
})

const isRemoteVersionNewer = (remote, local) => {
  const parse = (v) => v.split(".").map(Number);
  const [r1, r2, r3] = parse(remote);
  const [l1, l2, l3] = parse(local);
  if (r1 > l1) return true;
  if (r1 === l1 && r2 > l2) return true;
  if (r1 === l1 && r2 === l2 && r3 > l3) return true;
  return false;
}

const isUpToDate = computed(() => {
  if (!remoteVersion.value) return true;
  return !isRemoteVersionNewer(remoteVersion.value, appVersion);
})

const fetchAndroidSdkEnv = async () => {
  try {
    // TEMPORARILY DISABLED
    androidSdkEnv.value = await GetAndroidSdkEnv()
  } catch (error) {
    console.error('Error while running GetAndroidSdkEnv():', error)
  }
}

const checkForUpdate = async () => {
  isCheckingForUpdate.value = true
  try {
    const release = await GetLatestVersion()
    if (!release?.tag_name) return

    const cleanTag = release.tag_name.startsWith("v")
      ? release.tag_name.slice(1)
      : release.tag_name

    remoteVersion.value = cleanTag
    remoteReleaseUrl.value = release.html_url
  } catch (err) {
    console.warn("Version check failed:", err)
  } finally {
    isCheckingForUpdate.value = false
  }
}

const handleOpenGithubRelease = () => {
  try {
    BrowserOpenURL(remoteReleaseUrl.value)
  } catch (err) {
    console.error('Failed to open GitHub link in default browser:', err)
  }
}

onMounted(async () => {
  await fetchAndroidSdkEnv()
  await checkForUpdate()
})
</script>
