<template>
  <div class="settings-container">
    <h2 class="page-title">Settings</h2>
    <div class="settings-grid">
      <!-- --- Interface Settings --- -->
      <div class="settings-group">
        <h5 class="label-heading">
          <v-icon name="hi-view-grid" class="heading-icon" />
          Interface
        </h5>

        <div class="sub-group-items">
          <div v-if="isDarkMode" class="toggle-row">
            <label for="trueBlackToggle" class="toggle-label">Enable True Black</label>
            <label class="switch">
              <input id="trueBlackToggle" type="checkbox" :checked="themeStore.trueBlack"
                @change="themeStore.toggleTrueBlack" />
              <span class="slider"></span>
            </label>
          </div>
          <div v-else class="empty-state-text">
            Switch to Dark Mode to enable True Black.
          </div>
        </div>
      </div>

      <!-- --- System Settings --- -->
      <div class="settings-group">
        <h5 class="label-heading">
          <v-icon name="hi-adjustments" class="heading-icon" />
          System
        </h5>

        <div class="sub-group-items">
          <div v-if="androidSdkEnv" class="sdk-status-row"
            :class="{ 'sdk-found': !!androidSdkEnv.ANDROID_HOME, 'sdk-missing': !androidSdkEnv.ANDROID_HOME }">
            <div class="status-info">
              <span class="status-label">Android SDK Status</span>
              <span class="status-path">{{ androidSdkEnv.ANDROID_HOME || 'Not found' }}</span>
            </div>
            <div class="status-icon-box">
              <span v-if="androidSdkEnv.ANDROID_HOME">✔</span>
              <span v-else>✖</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- --- Footer Section --- -->
    <div class="app-footer">
      <div class="footer-meta">
        v{{ appVersion }} ({{ environment }}{{ remoteVersion ? (isUpToDate ? ', latest' : ', update available') : '' }})
        <span class="dot-separator">•</span>
        Powered by Vue, Go & Wails
      </div>

      <div class="signature">
        <div class="tooltip-trigger">
          <span class="dictionary-link">Sedulously</span>
          <div class="custom-tooltip">
            <span class="tooltip-type">(adverb)</span> — in a way that shows great care and persistent effort.
            <br />
            <small style="opacity: 0.7; margin-top: 4px; display: block;">Synonyms: diligently, persistently</small>
          </div>
        </div> engineered with <span class="heart">❤️‍🩹</span> by Symon
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

<style scoped>
/* --- Container Layout --- */
.settings-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  color: var(--text-primary);
  padding: 20px 0 0 0;
}

.page-title {
  font-size: 1.55rem;
  margin-bottom: 24px;
  color: var(--page-title-color);
}

/* --- Vertical Layout for Groups --- */
.settings-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
  justify-content: flex-start;
  align-items: stretch;
  padding-bottom: 24px;
}

/* --- Shared Group Styling --- */
.settings-group {
  padding: 24px;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 24px;
  flex: 0 0 auto;
  min-width: 280px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-sizing: border-box;
  transition: border-color 0.2s ease;
}

.label-heading {
  margin: 0;
  display: flex;
  align-items: center;
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.01em;
}

.heading-icon {
  color: var(--text-primary);
  margin-right: 12px;
  font-size: 1.1rem;
  opacity: 0.8;
  flex-shrink: 0;
}

.sub-group-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* --- Toggles --- */
.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-size: 0.9rem;
  color: var(--text-primary);
}

.toggle-label {
  font-weight: 500;
  user-select: none;
  cursor: default;
}

.empty-state-text {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-style: italic;
}

/* --- Switch Control (Material 3 style) --- */
.switch {
  position: relative;
  display: inline-block;
  width: 38px;
  height: 20px;
  flex-shrink: 0;
  cursor: default;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #515151;
  border-radius: 999px;
  transition: background-color 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.2);
}

.slider::before {
  content: "";
  position: absolute;
  height: 16px;
  width: 16px;
  left: 2px;
  top: 2px;
  background-color: #ffffff;
  border-radius: 50%;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.switch input:checked + .slider {
  background: linear-gradient(135deg, #fbbb03 0%, #e21d0a 100%);
  box-shadow: 0 0 8px rgba(226, 29, 10, 0.3);
}

.switch input:checked + .slider::before {
  transform: translateX(18px);
}

/* --- SDK Status Section --- */
.sdk-status-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px;
  border-radius: 16px;
  background-color: var(--bg-surface);
  border: 1px solid var(--border-subtle);
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.status-label {
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--text-primary);
}

.status-path {
  font-size: 0.75rem;
  color: var(--text-secondary);
  font-family: 'Consolas', monospace;
  word-break: break-all;
}

.status-icon-box {
  font-weight: bold;
  font-size: 1rem;
  padding-right: 12px;
}

.sdk-found .status-icon-box {
  color: #2ecc71;
}

.sdk-missing .status-icon-box {
  color: #e74c3c;
}

/* --- Footer Section --- */
.app-footer {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  gap: 16px;
  padding: 32px 0 24px;
}

.footer-meta {
  font-size: 0.72rem;
  color: var(--text-secondary);
  font-weight: 400;
  opacity: 0.5;
  letter-spacing: 0.02em;
  display: flex;
  align-items: center;
  gap: 8px;
}

.dot-separator {
  opacity: 0.6;
}

.signature {
  font-size: 0.78rem;
  color: var(--text-secondary);
  font-weight: 450;
  opacity: 0.9;
}

.heart {
  color: var(--color-primary);
  margin: 0 2px;
}

/* Tooltips */
.tooltip-trigger {
  position: relative;
  display: inline-block;
  cursor: default;
}

.custom-tooltip {
  position: absolute;
  bottom: 125%;
  left: 50%;
  transform: translateX(-50%) translateY(4px);
  width: 260px;
  padding: 10px 14px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 0.75rem;
  line-height: 1.4;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  pointer-events: none;
  opacity: 0;
  transition: all 0.25s ease;
  z-index: 1000;
  backdrop-filter: blur(8px);
}

.tooltip-trigger:hover .custom-tooltip {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.tooltip-type {
  font-style: italic;
  margin-right: 4px;
}

.dictionary-link {
  text-decoration: underline dotted color-mix(in srgb, currentColor 40%, transparent);
  text-underline-offset: 3px;
  transition: color 0.2s ease;
}

.dictionary-link:hover {
  text-decoration-color: currentColor;
}

/* Responsive adjustment */
@media (max-width: 768px) {
  .settings-grid {
    flex-direction: column;
  }
}
</style>
