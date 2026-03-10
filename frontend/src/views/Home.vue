<template>
  <div class="home-container">

    <div class="page-header">
      <div v-if="androidEnvChecked && !androidHomeMissing" class="page-title-container">
        <h2 class="page-title">Installed AVDs</h2>
        <span class="count-badge" v-show="store.avds.length">{{ store.avds.length }}</span>
      </div>

      <div v-if="androidHomeMissing" class="android-home-warning">
        <i class="bi bi-exclamation-triangle-fill warning-icon"></i>
        <div class="warning-text">
          <span class="warning-text-first-line">ANDROID_HOME is not set</span><br />
          This tool requires the Android SDK to be installed. Please set the ANDROID_HOME environment variable to
          the path of your Android SDK installation.
        </div>
        <button v-if="isWindows" class="btn btn-secondary" @click="openEnvVars">
          Open Environment Variables
        </button>
      </div>
    </div>

    <div v-show="store.avds.length" class="avd-grid">
      <div v-for="avd in store.avds" :key="avd.name" class="avd-card"
        :class="{ 'avd-running': avd.state === AvdState.RUNNING }" @mouseenter="avd.hover = true"
        @mouseleave="avd.hover = false">

        <!-- Card header: avatar + info (name) + menu -->
        <div class="avd-card-header">
          <div class="avd-avatar-container">
            <div class="avd-avatar">{{ avd.name.charAt(0).toUpperCase() }}</div>
          </div>
          <div class="avd-name">{{ avd.name }}</div>
          <button class="menu-button" :class="{ 'menu-visible': avd.hover }" @click="toggleMenu(avd, $event)">
            <i class="bi bi-three-dots"></i>
          </button>
        </div>

        <!-- Status indicator: dot + label (Flows under name/avatar) -->
        <div class="avd-status" :class="getStateClass(avd.state)">
          <div class="status-dot-container">
            <span class="status-dot"></span>
            <span class="status-pulse"></span>
          </div>
          <span class="status-label">{{ avd.state }}</span>
        </div>

        <!-- Animated context menu -->
        <transition name="fade-fast">
          <div v-if="menuAvd === avd" class="context-menu"
            :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }">
            <button @click="openEditDialog(avd)">Edit name</button>
            <button @click="openDeleteDialog(avd)">Delete</button>
          </div>
        </transition>

        <!-- Premium Action buttons with labels -->
        <div class="avd-actions">
          <template v-if="avd.state === AvdState.POWERED_OFF">
            <button class="action-btn-premium launch" @click="startAVD(avd, false)">
              <i class="bi bi-play-fill"></i>
              <span>Launch</span>
            </button>
            <button class="action-btn-premium coldboot" @click="startAVD(avd, true)">
              <i class="bi bi-snow"></i>
              <span>Cold Boot</span>
            </button>
          </template>

          <button class="action-btn-premium stop" v-else :disabled="avd.state !== AvdState.RUNNING"
            @click="stopAVD(avd.name)">
            <i class="bi bi-stop-fill"></i>
            <span>Stop</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit AVD name Dialog -->
    <div v-if="showEditDialog" class="edit-overlay" @click.self="closeEditDialog">
      <div class="edit-dialog">
        <button class="edit-close-button" @click="closeEditDialog">
          <i class="bi bi-x-lg"></i>
        </button>
        <h3>Edit AVD name</h3>
        <input v-model="editAvdName" placeholder="AVD Name" />
        <button class="btn btn-primary mt-3" @click="saveEdit">Save</button>
      </div>
    </div>

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toastMessage" class="toast">{{ toastMessage }}</div>
    </transition>
  </div>


</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ListAVDs, StartAVD, StopAVD, ListRunningAVDs, GetAndroidSdkEnv, OpenEnvironmentVariables } from '../../wailsjs/go/app/App'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'

const store = useAvdStore()

const showEditDialog = ref(false)
const editAvd = ref(null)
const editAvdName = ref('')

// Context menu
const menuAvd = ref(null)
const menuPosition = ref({ x: 0, y: 0 })

// Toasts
const toastMessage = ref('')
let toastTimeout = null

const androidHomeMissing = ref(false)
const androidEnvChecked = ref(false)

const isWindows = navigator.userAgent.includes('Windows')

function toggleMenu(avd, event) {
  if (menuAvd.value === avd) {
    menuAvd.value = null
  } else {
    menuAvd.value = avd
    menuPosition.value = { x: event.clientX, y: event.clientY }
  }
}

function openEditDialog(avd) {
  showEditDialog.value = true
  editAvd.value = avd
  editAvdName.value = avd.name
  menuAvd.value = null
}

function closeEditDialog() {
  showEditDialog.value = false
}

function saveEdit() {
  if (editAvd.value) {
    editAvd.value.name = editAvdName.value.trim()
    showToast('Edit saved ✅')
    closeEditDialog()
  }
}

async function openEnvVars() {
  try {
    await OpenEnvironmentVariables()
  } catch (err) {
    showToast('Failed to open environment settings ❌')
    console.error(err)
  }
}

const startAVD = async (avd, coldBoot = false) => {
  store.appendLog(`[AVD Launcher] Launching ${avd.name} (cold boot: ${coldBoot})...\n\n`)

  try {
    store.updateAvdStatus(avd.name, {
      state: AvdState.LOADING,
    })
    await nextTick()
    await StartAVD(avd.name, coldBoot)
  } catch (error) {
    showToast(`Failed to launch ${avd.name} ❌`)
    store.updateAvdStatus(avd.name, {
      state: AvdState.POWERED_OFF,
    })
  }
}

const stopAVD = async (avdName) => {
  try {
    store.updateAvdStatus(avdName, {
      state: AvdState.SHUTTING_DOWN,
    })
    await nextTick()
    await StopAVD(avdName);
  } catch (e) {
    console.error('Failed to stop AVD:', e);
  }
};

function openDeleteDialog(avd) {
  menuAvd.value = null
  if (confirm(`Are you sure you want to kill emulator for "${avd.name}"?`)) {
    stopAVD(avd.name)
    showToast('AVD killed ✅')
  }
}

// Close menu on click outside
function onClickOutside(event) {
  if (!event.target.closest('.menu-button') && !event.target.closest('.context-menu')) {
    menuAvd.value = null
  }
}

onMounted(async () => {
  try {
    const env = await GetAndroidSdkEnv()
    // env = '' // for debug purposes

    if (!env.ANDROID_HOME || env.ANDROID_HOME === '') {
      androidHomeMissing.value = true
      return
    }
  } catch (error) {
    console.log('Error while running GetAndroidSdkEnv():', error);
  }
  finally {
    androidEnvChecked.value = true
  }

  try {
    const avds = await ListAVDs()
    const runningAvds = await ListRunningAVDs()

    avds.forEach(name => {
      const isRunning = runningAvds?.includes(name)

      // If already in store, update its state
      const existing = store.avds.find(a => a.name === name)
      if (existing) {
        store.updateAvdStatus(name, {
          state: isRunning ? AvdState.RUNNING : AvdState.POWERED_OFF
        })
      } else {
        // Otherwise, add it
        store.avds.push({
          name,
          state: isRunning ? AvdState.RUNNING : AvdState.POWERED_OFF,
          hover: false
        })
      }
    })
  } catch (err) {
    showToast(`Error... ${err}`)
    console.log(err);
    store.avds = []
  }

  document.addEventListener('click', onClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)
})

function showToast(message) {
  toastMessage.value = message
  clearTimeout(toastTimeout)
  toastTimeout = setTimeout(() => {
    toastMessage.value = ''
  }, 4000)
}
</script>

<style scoped>
.home-container {
  padding: 20px 0 0 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-title-container {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 48px;
  /* Stable height to prevent flickering */
}

.page-title {
  font-size: 1.55rem;
  margin: 0;
  color: var(--page-title-color);
  line-height: 1;
}

.count-badge {
  background: var(--color-primary);
  color: white;
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: 20px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  height: fit-content;
}

.theme-dark .count-badge {
  color: #111;
}

/* ──── Card Grid ──── */
.avd-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 20px;
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

/* ──── Radical Glass Card ──── */
.avd-card {
  position: relative;
  background: var(--bg-glass);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 24px 18px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  /* Tighter gap to bring status closer */
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
  cursor: default;
}

.avd-card:hover {
  border-color: var(--border-strong);
}

/* ──── Card Header ──── */
.avd-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avd-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  min-width: 0;
}

.avd-avatar-container {
  flex-shrink: 0;
}

.avd-avatar {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  background: rgba(var(--color-primary-rgb), 0.1);
  border: 1px solid var(--border-glass);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
  backdrop-filter: blur(4px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.theme-dark .avd-avatar {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

.avd-name {
  font-weight: 700;
  font-size: 0.95rem;
  color: var(--text-primary);
  letter-spacing: -0.01em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.2;
}

.menu-button {
  background: none;
  border: none;
  color: var(--text-secondary);
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: default;
  opacity: 0;
  flex-shrink: 0;
  /* NO transition here for instant appearance */
}

.menu-button:hover {
  background: none;
  color: var(--text-primary);
  transform: scale(1.15);
}

.menu-button:active {
  background: none;
  color: var(--text-primary);
  transform: scale(0.95);
}

.avd-card:hover .menu-button {
  opacity: 1;
}

.menu-button:hover {
  background: var(--bg-card-active);
  color: var(--text-primary);
}

/* ──── Context Menu ──── */
.context-menu {
  position: fixed;
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--border-glass);
  border-radius: 12px;
  padding: 6px;
  z-index: 1500;
  min-width: 140px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.context-menu button {
  width: 100%;
  padding: 8px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  color: var(--text-primary);
  font-size: 0.85rem;
  font-weight: 500;
  text-align: left;
  cursor: default;
  transition: background 0.2s ease;
}

.context-menu button:hover {
  background: var(--color-primary);
  color: white;
}

/* ──── Status Indicator ──── */
.avd-status {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: 8px;
  margin-top: -4px;
}

.status-dot-container {
  position: relative;
  width: 8px;
  height: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.status-dot {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--color-status-off);
  z-index: 1;
}

.status-pulse {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--color-status-off);
  z-index: 0;
  opacity: 0;
}

.status-label {
  font-size: 0.82rem;
  font-weight: 500;
  color: var(--text-secondary);
  white-space: nowrap;
}

@keyframes pulse-premium {
  0% {
    transform: scale(1);
    opacity: 0.5;
  }

  70%,
  100% {
    transform: scale(2.2);
    opacity: 0;
  }
}

.avd-status.running .status-dot,
.avd-status.running .status-pulse {
  background: var(--color-status-running);
}

.avd-status.running .status-pulse {
  animation: pulse-premium 2.5s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.avd-status.running .status-label {
  color: var(--color-status-running);
}

.avd-status.loading .status-dot,
.avd-status.shuttingDown .status-dot {
  background: var(--color-status-loading);
}

.avd-status.loading .status-pulse,
.avd-status.shuttingDown .status-pulse {
  background: var(--color-status-loading);
  animation: pulse-premium 1.5s infinite;
}

/* ──── Premium Labeled Buttons ──── */
.avd-actions {
  display: flex;
  flex-direction: row;
  gap: 8px;
  margin-top: 4px;
}

.action-btn-premium {
  flex: 1;
  height: 38px;
  padding: 0 16px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--text-secondary) 8%, transparent);
  border: 1px solid color-mix(in srgb, var(--text-secondary) 15%, transparent);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: default;
  transition: all 0.1s ease;
  /* Snappy transition */
  white-space: nowrap;
}

.action-btn-premium:hover:not(:disabled) {
  background: color-mix(in srgb, var(--text-secondary) 15%, transparent);
  border-color: color-mix(in srgb, var(--text-secondary) 25%, transparent);
  color: var(--text-primary);
}

.action-btn-premium:active:not(:disabled) {
  transform: scale(0.95);
}

.action-btn-premium.launch {
  background: var(--color-primary);
  color: white;
  /* Matches light mode count badge */
  border: none;
}

.action-btn-premium.launch:hover {
  background: var(--color-primary);
  filter: brightness(1.1);
  color: white;
}

.theme-dark .action-btn-premium.launch {
  color: #111;
  /* High contrast on amber */
}

.theme-dark .action-btn-premium.launch:hover {
  color: #111;
}

.action-btn-premium.coldboot i {
  font-size: 0.85rem;
  /* Reduced snowflake size as requested */
}


.action-btn-premium.stop {
  background: rgba(220, 38, 38, 0.1);
  color: #dc2626;
  border: 1px solid rgba(220, 38, 38, 0.2);
}

.theme-dark .action-btn-premium.stop {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.action-btn-premium.stop:hover:not(:disabled) {
  background: #dc2626;
  color: white;
  border-color: #dc2626;
}

.action-btn-premium:disabled {
  opacity: 0.4;
}

.action-btn-premium i {
  font-size: 1.1rem;
}

/* ──── Toast & Dialog ──── */
.toast {
  position: fixed;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-glass);
  padding: 8px 16px;
  border-radius: 50px;
  color: var(--text-primary);
  font-weight: 600;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  z-index: 3000;
  font-size: 0.85rem;
}

.edit-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.edit-dialog {
  background: var(--bg-glass);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border-glass);
  padding: 24px;
  border-radius: 20px;
  width: 90%;
  max-width: 320px;
  color: var(--text-primary);
  text-align: center;
}

.edit-dialog h3 {
  font-size: 1.25rem;
  font-weight: 800;
  margin-bottom: 20px;
  letter-spacing: -0.02em;
}

.edit-dialog input {
  width: 100%;
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-glass);
  padding: 10px 14px;
  border-radius: 10px;
  color: var(--text-primary);
  margin-bottom: 16px;
}

.theme-dark .edit-dialog input {
  background: rgba(255, 255, 255, 0.05);
}

/* ──── Animations ──── */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: all 0.2s ease;
}

.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
