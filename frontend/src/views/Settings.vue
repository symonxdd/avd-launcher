<template>
  <div class="settings-container">
    <div class="settings-main">
      <h2 class="page-title">Settings</h2>
      <div class="setting-item">
        <div class="setting-info">
          <div class="setting-title">Environment variables</div>
          <div class="setting-description">Environment variables needed by the tool</div>
          <div class="env-variable">
            <ul class="env-list">
              <li>
                <strong>ANDROID_HOME (Android SDK): </strong>
                <span>{{ androidSdkEnv.ANDROID_HOME || 'Not found' }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="app-info">
      <div class="app-info-content">
        <div class="app-info-credits">
          Powered by
          <a href="https://wails.io/" target="_blank" rel="noopener" class="wails-link">Wails</a>,
          <a href="https://go.dev/" target="_blank" rel="noopener" class="go-link">Go</a> and
          <a href="https://vuejs.org/" target="_blank" rel="noopener" class="vue-link">Vue</a>
        </div>
        <div class="app-info-bottom-row">
          <div class="app-info-credits">
            Made with ❤️ by Symon from Belgium
          </div>
          <div class="app-info-meta">
            v<span>{{ appVersion }} {{ environment }}</span>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { GetAndroidSdkEnv } from '../../wailsjs/go/app/App'

const autoScrollLogs = ref(true)
const androidSdkEnv = ref({})
const appVersion = __APP_VERSION__ || 'v1.0.0'
const environment = import.meta.env.MODE === 'development' ? '(dev)' : '(prod)'

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
.settings-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 20px 0 0;
  color: #ccc;
}

.settings-main {
  flex: 1;
}

.page-title {
  font-size: 1.3rem;
  margin-bottom: 16px;
  color: #ccc;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 15px;
  padding-bottom: 15px;
}

.setting-item:not(:last-child) {
  border-bottom: 1px solid #3b3b3b;
}

.setting-title {
  font-weight: 600;
  font-size: 1.05rem;
}

.setting-description {
  font-size: 0.85rem;
  color: #aaa;
}

.env-variable {
  margin-top: 8px;
  font-size: 0.85rem;
  color: #ccc;
}

.app-info-meta {
  font-size: 0.9rem;
  color: #888;
}

.app-info-credits {
  font-size: 0.9rem;
  color: #888;
}

.app-info-bottom-row {
  display: flex;
  margin-top: 3px;
}

.app-info-bottom-row .app-info-credits {
  flex: 1;
}

.wails-link {
  color: #DF0000;
  text-decoration: none;
  display: inline-block;
}

.wails-link:visited {
  color: #DF0000;
}

.wails-link:hover {
  text-decoration: underline;
}

.go-link {
  color: #00ADD8;
  text-decoration: none;
  display: inline-block;
}

.go-link:visited {
  color: #00ADD8;
}

.go-link:hover {
  text-decoration: underline;
}

.vue-link {
  color: #42b883;
  text-decoration: none;
  display: inline-block;
}

.vue-link:visited {
  color: #42b883;
}

.vue-link:hover {
  text-decoration: underline;
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
  background-color: #DF0000;
}

input:checked + .slider:before {
  transform: translateX(24px);
}

.env-list {
  list-style-type: disc;
  padding-left: 20px;
  color: #ccc;
}

.env-list li {
  margin-bottom: 6px;
}
</style>
