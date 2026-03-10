<template>
  <div :class="styles.homeContainer">

    <div :class="styles.pageHeader">
      <div v-if="androidEnvChecked && !androidHomeMissing" :class="styles.pageTitleContainer">
        <h2 :class="styles.pageTitle">Installed AVDs</h2>
        <span :class="styles.countBadge" v-show="store.avds.length">{{ store.avds.length }}</span>
      </div>

      <div v-if="androidHomeMissing" :class="styles.androidHomeWarning">
        <i class="bi bi-exclamation-triangle-fill" :class="styles.warningIcon"></i>
        <div :class="styles.warningText">
          <span :class="styles.warningTextFirstLine">ANDROID_HOME is not set</span><br />
          This tool requires the Android SDK to be installed. Please set the ANDROID_HOME environment variable to
          the path of your Android SDK installation.
        </div>
        <button v-if="isWindows" :class="[styles.btn, styles.btnSecondary]" @click="openEnvVars">
          Open Environment Variables
        </button>
      </div>
    </div>

    <div v-show="store.avds.length" :class="styles.avdGrid">
      <div v-for="avd in store.avds" :key="avd.name"
        :class="[styles.avdCard, { [styles.running]: avd.state === AvdState.RUNNING }]" @mouseenter="avd.hover = true"
        @mouseleave="avd.hover = false">

        <!-- Card header: avatar + info (name) + menu -->
        <div :class="styles.avdCardHeader">
          <div :class="styles.avdAvatarContainer">
            <div :class="styles.avdAvatar">{{ avd.name.charAt(0).toUpperCase() }}</div>
          </div>
          <div :class="styles.avdName">{{ avd.name }}</div>
          <button :class="[styles.menuButton, { [styles.menuVisible]: avd.hover }]" @click="toggleMenu(avd, $event)">
            <i class="bi bi-three-dots"></i>
          </button>
        </div>

        <!-- Status indicator: dot + label (Flows under name/avatar) -->
        <div :class="[styles.avdStatus, styles[getStateClass(avd.state)]]">
          <div :class="styles.statusDotContainer">
            <span :class="styles.statusDot"></span>
            <span :class="styles.statusPulse"></span>
          </div>
          <span :class="styles.statusLabel">{{ avd.state }}</span>
        </div>

        <!-- Animated context menu -->
        <transition name="fade-fast">
          <div v-if="menuAvd === avd" :class="styles.contextMenu"
            :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }">
            <button @click="openEditDialog(avd)">Edit name</button>
            <button @click="openDeleteDialog(avd)">Delete</button>
          </div>
        </transition>

        <!-- Premium Action buttons with labels -->
        <div :class="styles.avdActions">
          <template v-if="avd.state === AvdState.POWERED_OFF">
            <button :class="[styles.actionBtnPremium, styles.launch]" @click="startAVD(avd, false)">
              <i class="bi bi-play-fill"></i>
              <span>Launch</span>
            </button>
            <button :class="[styles.actionBtnPremium, styles.coldboot]" @click="startAVD(avd, true)">
              <i class="bi bi-snow"></i>
              <span>Cold Boot</span>
            </button>
          </template>

          <button :class="[styles.actionBtnPremium, styles.stop]" v-else :disabled="avd.state !== AvdState.RUNNING"
            @click="stopAVD(avd.name)">
            <i class="bi bi-stop-fill"></i>
            <span>Stop</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit AVD name Dialog -->
    <div v-if="showEditDialog" :class="styles.editOverlay" @click.self="closeEditDialog">
      <div :class="styles.editDialog">
        <button :class="styles.editCloseButton" @click="closeEditDialog">
          <i class="bi bi-x-lg"></i>
        </button>
        <h3>Edit AVD name</h3>
        <input v-model="editAvdName" placeholder="AVD Name" />
        <button :class="[styles.btn, styles.btnPrimary, 'mt-3']" @click="saveEdit">Save</button>
      </div>
    </div>

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toastMessage" :class="styles.toast">{{ toastMessage }}</div>
    </transition>
  </div>


</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ListAVDs, StartAVD, StopAVD, ListRunningAVDs, GetAndroidSdkEnv, OpenEnvironmentVariables } from '../../wailsjs/go/app/App'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'
import styles from './Home.module.css'

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
  if (!event.target.closest(`.${styles.menuButton}`) && !event.target.closest(`.${styles.contextMenu}`)) {
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
