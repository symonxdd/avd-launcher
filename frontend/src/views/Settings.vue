<template>
  <div class="settings-container">
    <div class="settings-main">
      <h2 class="page-title">Settings</h2>

      <!-- Auto-scroll logs -->
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-title">Auto-scroll logs</div>
          <div class="setting-description">Automatically scrolls the logs to the last line</div>
        </div>
        <label class="switch">
          <input type="checkbox" v-model="autoScrollLogs" />
          <span class="slider"></span>
        </label>
      </div>

      <!-- Env vars -->
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-title">Environment variables</div>
          <div class="setting-description">These are the environment variables needed by the tool</div>
          <div class="env-variable">
            <strong>Android SDK: </strong>
            <span>{{ envVariables.ANDROID_HOME || 'Not found' }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- About / Credits -->
    <div class="app-info">
      <div class="app-info-content">
        <div class="app-info-meta">
          v<span>{{ appVersion }} {{ environment }}</span>
        </div>
        <div class="app-info-credits">
          Made with 💜 by Symon from BE
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetEnvVariables } from '../../wailsjs/go/app/App'

const autoScrollLogs = ref(true)
const envVariables = ref({})
const appVersion = __APP_VERSION__ || 'v1.0.0'
const environment = import.meta.env.MODE === 'development' ? '(dev)' : '(prod)'

const fetchEnvVariables = async () => {
  try {
    envVariables.value = await GetEnvVariables()
    console.log('Fetched environment variables:', envVariables.value)
  } catch (error) {
    console.error('Failed to fetch environment variables:', error)
  }
}

onMounted(() => {
  fetchEnvVariables()
})
</script>

<style scoped>
.settings-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 20px 0 20px;
  color: #ccc;
}

.settings-main {
  flex: 1;
}

.page-title {
  font-size: 1.2rem;
  margin-bottom: 16px;
  color: #ccc;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 20px 0;
  padding-bottom: 10px;
  border-bottom: 1px solid #444;
}

.setting-info {
  max-width: 70%;
}

.setting-title {
  font-weight: 600;
  margin-bottom: 4px;
  font-size: 1.1rem;
  color: #fff;
}

.setting-description {
  font-size: 0.9rem;
  color: #aaa;
}

.env-variable {
  margin-top: 8px;
  font-size: 0.9rem;
  color: #ccc;
}

.app-info-meta {
  font-size: 0.9rem;
  color: #888;
}

.app-info-credits {
  font-size: 0.9rem;
  color: #888;
  font-style: italic;
}

/* iOS-style switch */
.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 24px;
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
  background-color: #777;
  transition: 0.4s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #4caf50;
}

input:checked + .slider:before {
  transform: translateX(24px);
}
</style>
