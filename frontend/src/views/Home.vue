<template>
  <div class="home-container">
    <h2 class="page-title">
      Installed AVDs
      <transition name="fade" mode="out-in">
        <span v-if="!loading" :key="store.avds.length">
          ({{ store.avds.length }})
        </span>
      </transition>
    </h2>

    <transition v-if="loading" name="fade-overlay">
      <div v-if="loading" class="loading-state">
        <span class="spinner"></span>
        Loading AVDs ...
      </div>
    </transition>

    <div v-else-if="store.avds.length" class="avd-grid">
      <div v-for="avd in store.avds" :key="avd.name" class="avd-card" :class="{ 'avd-running': avd.running }"
        @mouseenter="avd.hover = true" @mouseleave="avd.hover = false">

        <!-- Three dots menu -->
        <button v-if="avd.hover" class="menu-button" @click="toggleMenu(avd, $event)">
          <i class="bi bi-three-dots"></i>
        </button>

        <!-- Animated context menu -->
        <transition name="fade-fast">
          <div v-if="menuAvd === avd" class="context-menu"
            :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }">
            <button @click="openEditDialog(avd)">Edit name</button>
            <button @click="openDeleteDialog(avd)">Delete</button>
          </div>
        </transition>

        <div class="avd-name">{{ avd.name }}</div>

        <div class="avd-status" :class="getStateClass(avd.state)">
          {{ avd.state }}
        </div>

        <div class="avd-buttons">
          <div class="avd-launch-buttons">
            <button class="btn btn-primary" :class="{ 'disabled': avd.state !== AvdState.POWERED_OFF }"
              :disabled="avd.state !== AvdState.POWERED_OFF" @click="startAVD(avd, false)">
              Start
            </button>
            <button class="btn btn-secondary" :class="{ 'disabled': avd.state !== AvdState.POWERED_OFF }"
              :disabled="avd.state !== AvdState.POWERED_OFF" @click="startAVD(avd, true)">
              Cold Boot
            </button>
          </div>

          <div v-if="avd.state === AvdState.RUNNING" class="avd-stop-button">
            <button class="btn btn-close" @click="stopAVD(avd)">
              <i class="bi bi-x-lg"></i>
            </button>
          </div>
        </div>

      </div>
    </div>
    <div v-else class="no-avds">No AVDs found.</div>

    <!-- Logs container -->
    <!-- <div v-if="store.logs.length" class="logs" :class="{ 'logs-expanded': isLogsExpanded }" @click="toggleLogs"
      ref="logsContainer">
      <div class="logs-content">{{ store.logs }}</div>
    </div> -->

    <!-- <div v-if="isLogsExpanded" class="overlay" @click="toggleLogs"></div> -->

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
import { ListAVDs, StartAVD, StopAVD, ListRunningAVDs } from '../../wailsjs/go/app/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'

const store = useAvdStore()

const isLogsExpanded = ref(false)
const loading = ref(true)
const logsContainer = ref(null)
const showEditDialog = ref(false)
const editAvd = ref(null)
const editAvdName = ref('')

// Context menu
const menuAvd = ref(null)
const menuPosition = ref({ x: 0, y: 0 })

// Toasts
const toastMessage = ref('')
let toastTimeout = null

// Log listener
let logListener = null

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

const startAVD = async (avd, coldBoot = false) => {
  store.appendLog(`[AVD Launcher] Launching ${avd.name} (cold boot: ${coldBoot})`)

  try {
    store.updateAvdStatus(avd.name, {
      state: AvdState.LOADING,
      launchMode: coldBoot ? 'cold' : 'normal'
    })
    await nextTick()
    await StartAVD(avd.name, coldBoot)

    // showToast(`AVD ${avd.name} launching... 🚀`)
  } catch (error) {
    showToast(`Failed to launch ${avd.name} ❌`)
    store.updateAvdStatus(avd.name, {
      loading: false,
      running: false,
      launchMode: coldBoot ? 'cold' : 'normal'
    })
  }

  // isLogsExpanded.value = true
  // await nextTick()
  // scrollLogsToBottom()
}

const stopAVD = async (avd) => {
  try {
    store.updateAvdStatus(avd.name, {
      state: AvdState.SHUTTING_DOWN,
    })
    await nextTick()
    await StopAVD(avd.name);
  } catch (e) {
    console.error('Failed to stop AVD:', e);
  }
};

function openDeleteDialog(avd) {
  menuAvd.value = null
  if (confirm(`Are you sure you want to kill emulator for "${avd.name}"?`)) {
    stopAVD(avd)
    showToast('AVD killed ✅')
  }
}

function toggleLogs() {
  isLogsExpanded.value = !isLogsExpanded.value
}

function scrollLogsToBottom() {
  if (logsContainer.value) {
    nextTick(() => {
      logsContainer.value.scrollTop = logsContainer.value.scrollHeight
    })
  }
}

// Close menu on click outside
function onClickOutside(event) {
  if (!event.target.closest('.menu-button') && !event.target.closest('.context-menu')) {
    menuAvd.value = null
  }
}

onMounted(async () => {
  if (store.avds.length > 0) {
    loading.value = false
    return
  }

  try {
    const avds = await ListAVDs()
    const runningAvds = await ListRunningAVDs()

    store.avds = avds.map(name => ({
      name,
      launchMode: null,
      state: runningAvds?.includes(name) ? AvdState.RUNNING : AvdState.POWERED_OFF,
      hover: false
    }))
  } catch (err) {
    showToast(`Error... ${err}`)
    console.log(err);
    store.avds = []
  } finally {
    loading.value = false
  }

  logListener = EventsOn('avd-log', (line) => {
    store.appendLog(line)
    scrollLogsToBottom()
  })

  EventsOn('avd-booted', (name) => {
    store.updateAvdStatus(name, {
      state: AvdState.RUNNING,
      launchMode: null
    })
    showToast(`${name} is running`)
  })

  EventsOn("avd-shutdown", (name) => {
    store.updateAvdStatus(name, {
      state: AvdState.POWERED_OFF,
      launchMode: null
    })
  });

  document.addEventListener('click', onClickOutside)
})

onBeforeUnmount(() => {
  if (logListener) EventsOff('avd-log', logListener)
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
  padding: 20px 10px;
  /* position: relative; */
}

.page-title {
  font-size: 1.2rem;
  margin-bottom: 16px;
  color: #ccc;
}

.avd-grid {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
  align-items: flex-start;
  /* flex-direction: column; */
}

.no-avds {
  color: #999;
  font-style: italic;
}

.avd-card {
  position: relative;
  background-color: #222;
  border-radius: 8px;
  padding: 16px;
  width: 215px;
  border: 1px solid transparent;
  transition: border 0.3s ease;
}

.avd-card:hover {
  border: 1px solid #474747;
}

.avd-running {
  /* border: 1px solid #147029; */
}

.menu-button {
  position: absolute;
  top: 8px;
  right: 8px;
  background: transparent;
  border: none;
  color: #bbb;
  font-size: 1.1rem;
}

.context-menu {
  border-radius: 14px;
  position: fixed;
  background-color: #333333;
  border: 1px solid #555;
  padding: 4px;
  z-index: 1500;
  display: flex;
  flex-direction: column;
  min-width: 120px;
  box-shadow: 0 10px 10px rgba(0, 0, 0, 0.2);
}

.context-menu button {
  border-radius: 9px;
  background: transparent;
  border: none;
  color: #fff;
  text-align: left;
  padding: 8px;
  width: 100%;
}

.context-menu button:hover {
  background: #444;
}

.avd-name {
  font-weight: 600;
  font-size: 0.9rem;
  color: #eee;
}

.avd-status {
  font-size: 0.85rem;
}

.avd-status.running {
  color: #28a745;
}

.avd-status.poweredOff {
  color: gray;
}

.avd-status.loading {
  color: #f39c12;
}

.avd-status.shuttingDown {
  color: #f39c12;
}

.avd-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 20px;
  /* Adjust as needed to reserve space */
}

.avd-launch-buttons {
  display: flex;
  gap: 10px;
}

.avd-stop-button {
  display: flex;
  justify-content: flex-start;
}

.btn {
  font-family: inherit;
  padding: 6px 12px;
  font-size: 0.9rem;
  border-radius: 4px;
  border: 1px solid transparent;
  transition: all 0.2s ease;
}

.btn-primary {
  background-color: #0d6efd;
  color: white;
}

.btn-primary:hover {
  background-color: #005dc0;
}

.btn-secondary {
  background-color: #000000;
  color: white;
  /* background-color: #6c63ff; */
}

.btn-secondary:hover {
  background-color: #222222;
}

.btn-close {
  width: 32px;
  height: 32px;
  background-color: #b12e3b;
  color: white;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  background-color: #862932;
}

.logs {
  position: fixed;
  bottom: 0;
  background: #2c2c2c;
  color: white;
  padding: 6px 16px;
  font-family: monospace;
  height: 30px;
  width: 100%;
  white-space: pre-wrap;
  border-radius: 10px 10px 0 0;
  transition: height 0.3s ease;
  z-index: 1000;
  border: 1px solid #474747;
  border-bottom: none;
  overflow-y: auto;
  font-size: 13px;
}

.logs-content {
  line-height: 1.5;
}

.logs-expanded {
  height: 450px;
  font-size: 13px;
}

.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

/* Edit dialog */
.edit-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
}

.edit-dialog {
  background: #2c2c2c;
  padding: 24px;
  border-radius: 8px;
  position: relative;
  color: white;
  width: 300px;
  text-align: center;
}

.edit-close-button {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #444;
  border: none;
  color: #fff;
  border-radius: 4px;
  padding: 5px 7px;
}

.edit-dialog input {
  width: 100%;
  padding: 8px;
  margin-top: 12px;
  border-radius: 4px;
  border: 1px solid #555;
  background: #1c1c1c;
  color: #c4c4c4;
  margin-bottom: 16px;
}

/* Toast */
.toast {
  position: fixed;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  background: #080808;
  color: white;
  padding: 10px 20px;
  border-radius: 15px;
  z-index: 3000;
  opacity: 0.9;
}

.disabled {
  background-color: #555;
  color: #999;
}

.disabled:hover {
  background-color: #555;
  color: #999;
}

.loading-state {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #181818;
  /* or match your background */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  z-index: 9999;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top-color: #4caf50;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.7s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 0.3s ease;
}

.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
}

.fade-overlay-enter-active,
.fade-overlay-leave-active {
  transition: opacity 2s ease;
}

.fade-overlay-enter-from,
.fade-overlay-leave-to {
  opacity: 0;
}

.fade-overlay-enter-to,
.fade-overlay-leave-from {
  opacity: 1;
}
</style>
