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
import { useThemeStore } from '../stores/themeStore'
import styles from './Settings.module.css'

const themeStore = useThemeStore()
const androidSdkEnv = ref(null)
const appVersion = typeof __APP_VERSION__ !== "undefined" ? __APP_VERSION__ : "1.0.0"
const environment = import.meta.env.MODE === 'development' ? 'dev' : 'release'
const remoteVersion = ref(null) // Not implemented in AVD Launcher yet

const isDarkMode = computed(() => {
  if (themeStore.theme === 'dark') return true
  if (themeStore.theme === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return false
})

const isUpToDate = computed(() => true) // Placeholder for now

const fetchAndroidSdkEnv = async () => {
  try {
    androidSdkEnv.value = await GetAndroidSdkEnv()
  } catch (error) {
    console.error('Error while running GetAndroidSdkEnv():', error)
  }
}

onMounted(async () => {
  await fetchAndroidSdkEnv()
})
</script>
