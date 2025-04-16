<template>
  <div class="home-container">

    <div class="page-title-container">
      <h2 class="page-title">Installed AVDs</h2>
      <span class="count-badge" v-if="store.avds.length">{{ store.avds.length }}</span>
    </div>

    <div v-show="store.avds.length" class="avd-grid">
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
          <div class="avd-action-button">
            <template v-if="avd.state === AvdState.POWERED_OFF">
              <button class="icon-button" @click="startAVD(avd, false)">
                <i class="bi bi-play-fill icon-start" title="Start AVD"></i>
              </button>
              <button class="icon-button icon-button-coldboot" @click="startAVD(avd, true)">
                <i class="bi bi-arrow-repeat icon-coldboot" title="Cold Boot"></i>
              </button>
            </template>

            <button class="icon-button" v-else :disabled="avd.state !== AvdState.RUNNING" @click="stopAVD(avd)">
              <i class="bi bi-stop-fill icon-stop" title="Stop AVD"></i>
            </button>
          </div>

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
import { ListAVDs, StartAVD, StopAVD, ListRunningAVDs } from '../../wailsjs/go/app/App'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'

const store = useAvdStore()

const loading = ref(true)
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
  store.appendLog(`[AVD Launcher] Launching ${avd.name} (cold boot: ${coldBoot})...\n\n`)

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

// Close menu on click outside
function onClickOutside(event) {
  if (!event.target.closest('.menu-button') && !event.target.closest('.context-menu')) {
    menuAvd.value = null
  }
}

onMounted(async () => {
  // if (store.avds.length > 0) {
  //   loading.value = false
  //   return
  // }
  loading.value = true

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
  padding: 20px 0px 0 20px;
  /* position: relative; */
}

.page-title {
  font-size: 1.3rem;
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
  background-color: #131313;
  border-radius: 8px;
  padding: 16px;
  width: 215px;
  border: 1px solid #363636;
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
  /* display: flex; */
  /* flex-direction: column; */
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

.page-title-container {
  position: relative;
  display: inline-block;
}

.count-badge {
  position: absolute;
  top: -12px;
  right: -5px;
  background-color: #DF0000;
  color: white;
  font-size: 12px;
  padding: 2px 5px;
  border-radius: 50%;
  line-height: 1;
  text-align: center;
}

.icon-button {
  background: none;
  border: none;
  padding: 6px;
  font-size: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  will-change: filter;
}

.icon-button i {
  transition: filter 0.15s ease;
}

.icon-button:hover i {
  filter: brightness(1.3);
}

.icon-button:disabled {
  opacity: 0.5;
}

.icon-start {
  color: #28a745;
}

.icon-stop {
  color: #dc3545;
}

.avd-action-button {
  display: flex;
  margin-top: 10px;
}

.icon-button-coldboot i {
  font-size: 1.15rem;
  /* Slightly smaller to match the others */
}

.icon-coldboot {
  color: #ffc107;
}

.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 0.3s ease;
}

.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
}
</style>
