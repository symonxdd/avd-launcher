<template>
  <div class="settings-container">
    <div class="settings-main">
      <h2 class="page-title">Settings</h2>

      <!-- Theming Setting -->
      <div class="setting-group theme-setting-group">
        <h5 class="label-heading">
          Theme
        </h5>
        <small class="label-subtext">Choose app appearance</small>
        <div class="btn-group">
          <button class="btn" @click="setTheme('light')"
            :class="{ active: themeStore.theme === 'light' }">Light</button>
          <button class="btn" @click="setTheme('dark')" :class="{ active: themeStore.theme === 'dark' }">Dark</button>
          <button class="btn" @click="setTheme('system')"
            :class="{ active: themeStore.theme === 'system' }">System</button>
        </div>

        <!-- ✅ Place the True Black toggle here -->
        <div v-if="isDarkMode" class="true-black-toggle">
          <label class="switch">
            <input type="checkbox" :checked="themeStore.trueBlack" @change="themeStore.toggleTrueBlack" />
            <span class="slider"></span>
          </label>
          <span class="toggle-label">Enable True Black</span>
        </div>
      </div>
    </div>

    <!-- App Info Section -->
    <div class="app-info-section">
      <div v-if="androidSdkEnv" class="sdk-status"
        :class="{ 'sdk-found': !!androidSdkEnv.ANDROID_HOME, 'sdk-missing': !androidSdkEnv.ANDROID_HOME }"
        :data-tooltip="sdkTooltipText">
        <div class="status-icon">
          <span v-if="androidSdkEnv.ANDROID_HOME">✔</span>
          <span v-else>✖</span>
        </div>
        <div class="status-text">
          <strong>Android SDK: </strong>
          <span>
            {{ androidSdkEnv.ANDROID_HOME || 'Android SDK not installed/found' }}
          </span>
        </div>
      </div>

    </div>

    <!-- Footer Section -->
    <div class="app-footer">
      <div class="footer-section footer-left app-info-meta">
        v{{ appVersion }} {{ environment }}
      </div>
      <div class="footer-section footer-center app-info-credits">
        Powered by
        <a href="https://wails.io/" target="_blank" rel="noopener" class="wails-link">Wails</a>,
        <a href="https://go.dev/" target="_blank" rel="noopener" class="go-link">Go</a> and
        <a href="https://vuejs.org/" target="_blank" rel="noopener" class="vue-link">Vue</a>
      </div>
      <div class="footer-section footer-right app-info-credits">
        Made with ❤️ by Symon from Belgium
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { GetAndroidSdkEnv } from '../../wailsjs/go/app/App'
import { useThemeStore } from '../stores/themeStore'

const themeStore = useThemeStore()
const androidSdkEnv = ref(null)
const appVersion = __APP_VERSION__ || 'v1.0.0'
const environment = import.meta.env.MODE === 'development' ? '(dev)' : '(prod)'

const sdkTooltipText = computed(() => {
  if (!androidSdkEnv.value) return ''
  return androidSdkEnv.value.ANDROID_HOME
    ? 'SDK found through ANDROID_HOME env variable'
    : 'This tool requires the Android SDK to be installed and the ANDROID_HOME env variable to be set'
})

const fetchAndroidSdkEnv = async () => {
  try {
    androidSdkEnv.value = await GetAndroidSdkEnv()
  } catch (error) {
    console.error('Error while running GetAndroidSdkEnv():', error)
  }
}

const isDarkMode = computed(() => {
  if (themeStore.theme === 'dark') return true
  if (themeStore.theme === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return false
})

const setTheme = (mode) => {
  themeStore.setTheme(mode)
}

onMounted(async () => {
  await fetchAndroidSdkEnv()
})
</script>

<style scoped>
/* --- Container Layout --- */
.settings-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 0px 0 0;
  color: var(--text-color);
}

/* --- Page Title --- */
.page-title {
  font-size: 1.55rem;
  margin-bottom: 16px;
  color: var(--page-title-color);
}

/* --- Environment Variable Block --- */
.setting-item {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;
  padding-bottom: 15px;
}

.setting-title {
  font-weight: 600;
  font-size: 1.05rem;
}

.setting-description {
  font-size: 0.85rem;
  color: var(--text-color);
}

.env-variable {
  margin-top: 8px;
  font-size: 0.85rem;
  color: var(--text-color);
}

.env-list {
  list-style-type: disc;
  padding-left: 20px;
  color: var(--text-color);
}

.env-list li {
  margin-bottom: 6px;
}

/* --- Theming Section --- */
.setting-group {
  background-color: var(--background-secondary);
  /* border: 1px solid var(--border-color); */
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
  width: fit-content;
}

.label-heading {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--text-color);
}

.label-subtext {
  display: block;
  font-size: 0.8rem;
  color: var(--secondary-text-color);
  margin-bottom: 10px;
}

.btn-group {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.btn {
  background-color: transparent;
  color: var(--text-color);
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  font-weight: 400;
  font-family: Nunito, sans-serif;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  transition: background-color 0.2s ease-in-out, color 0.2s ease-in-out;
}

.btn:hover {
  background-color: var(--btn-hover-bg);
}

.btn.active {
  border-color: #8e44ad;
  background-color: rgba(142, 68, 173, 0.2);
  color: #8e44ad;
}

/* --- Footer Layout --- */
.app-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  padding-top: 15px;
  font-size: 0.75rem;
  background-color: var(--background-secondary);
  color: var(--secondary-text-color);
}

.footer-section {
  padding: 5px 10px;
}

.footer-left {
  order: 2;
}

.footer-center {
  order: 1;
  text-align: center;
}

.footer-right {
  order: 3;
  text-align: right;
}

.app-info-section {
  /* margin-top: 20px;
  padding: 16px;
  background-color: var(--background-secondary);
  border-radius: 12px;
  width: fit-content;
  font-size: 0.9rem;
  color: var(--text-color); */

  margin-top: auto;
  /* Push it to the bottom of the flex container */
  margin-bottom: 20px;
  /* Spacing above the footer */
  padding: 16px;
  background-color: var(--background-secondary);
  border-radius: 12px;
  width: fit-content;
  font-size: 0.9rem;
  color: var(--text-color);
}

/* Tooltip via data attribute */
.sdk-status::after {
  content: attr(data-tooltip);
  position: absolute;
  bottom: calc(100% + 6px);
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--tooltip-bg, #333);
  color: var(--tooltip-text, #fff);
  font-weight: 500;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  white-space: nowrap;
  z-index: 1000;
  opacity: 0;
  pointer-events: none;
  transition: none;
}

/* Show tooltip only on hover */
.sdk-status:hover::after {
  opacity: 1;
}

.sdk-status {
  position: relative;
  cursor: default;

  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 999px;
  transition: background-color 0.3s ease-in-out;
}

.sdk-found {
  background-color: rgba(46, 204, 113, 0.15);
  /* light green bg */
  color: #2ecc71;
  /* green text */
}

.sdk-missing {
  background-color: rgba(231, 76, 60, 0.15);
  /* light red bg */
  color: #e74c3c;
  /* red text */
}

.status-icon {
  width: 1.1rem;
  height: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  cursor: default;
  position: relative;
}

.status-text {
  font-family: 'Consolas', monospace;
}

.status-text span {
  font-weight: 400;
  color: inherit;
}

/* Light/Dark tooltip support */
:root {
  --tooltip-bg: #333;
  --tooltip-text: #fff;
}

body[data-theme='light'] {
  --tooltip-bg: #eee;
  --tooltip-text: #111;
}

.true-black-toggle {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 16px;
  font-size: 0.9rem;
  color: var(--text-color);
}

.toggle-label {
  font-weight: 500;
}

/* Material 3-like switch */
.switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--border-color);
  transition: 0.3s;
  border-radius: 34px;
}

.slider:before {
  position: absolute;
  content: '';
  height: 16px;
  width: 16px;
  left: 3px;
  bottom: 3px;
  background-color: var(--bg-color);
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #8e44ad;
}

input:checked + .slider:before {
  transform: translateX(18px);
}


/* --- Link Colors --- */
.wails-link {
  color: #df0000;
  text-decoration: none;
}

.go-link {
  color: #00add8;
  text-decoration: none;
}

.vue-link {
  color: #42b883;
  text-decoration: none;
}

.wails-link:hover,
.go-link:hover,
.vue-link:hover {
  text-decoration: underline;
}
</style>
