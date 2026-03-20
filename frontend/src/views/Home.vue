<template>
  <div :class="styles.homeContainer">

    <div :class="styles.pageHeader">
      <div v-if="!sdkMissing" :class="styles.pageTitleContainer">
        <h2 :class="styles.pageTitle">Installed AVDs</h2>
        <span :class="styles.countBadge" v-show="store.avds.length">{{ store.avds.length }}</span>
      </div>

      <MissingSdkWarning v-if="sdkMissing" :is-windows="isWindows" @select-sdk="selectSdkPath"
        @open-env-vars="openEnvVars" @open-env-info="openEnvInfo" />
    </div>

    <div v-show="store.avds.length && !sdkMissing" :class="styles.avdGrid">
      <AvdCard v-for="avd in store.avds" :key="avd.name" :avd="avd" @toggle-menu="toggleMenu" @launch="startAVD"
        @stop="stopAVD" />
    </div>

    <!-- Animated context menu -->
    <Teleport to="body">
      <transition name="fade-fast">
        <div v-if="menuAvd" :class="styles.contextMenu"
          :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }">
          <button @click="openFolder(menuAvd.path)">
            <v-icon name="hi-folder-open" :scale="0.85" />
            <span>Show on Disk</span>
          </button>
          <button @click="openEditDialog(menuAvd)">
            <v-icon name="hi-pencil" :scale="0.85" />
            <span>Rename</span>
          </button>
          <button @click="openDeleteDialog(menuAvd)" :class="styles.deleteItem">
            <v-icon name="hi-trash" :scale="0.85" />
            <span>Delete</span>
          </button>
        </div>
      </transition>
    </Teleport>

    <!-- Modals -->
    <RenameAvdModal :show="showEditDialog" :avd="editAvd" @close="showEditDialog = false" @success="onRenameSuccess"
      @error="onRenameError" />
    <DeleteAvdModal :show="showDeleteDialog" :avd="deleteAvdTarget" @close="showDeleteDialog = false"
      @success="onDeleteSuccess" @error="onDeleteError" />

    <!-- Env Vars Explanation Modal -->
    <EnvInfoModal :show="showEnvInfoDialog" :is-closing="isEnvInfoClosing" @close="closeEnvInfo"
      @animationend="handleEnvInfoAnimationEnd" />

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toastMessage" :class="styles.toast">{{ toastMessage }}</div>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ListAVDs, StartAVD, StopAVD, GetAndroidSdkEnv, OpenEnvironmentVariables, SelectAndSaveSdkPath, GetAvdInfo, GetAvdDiskUsage, OpenAvdFolder } from '../../wailsjs/go/app/App'
import EnvInfoModal from '../components/EnvInfoModal.vue'
import MissingSdkWarning from '../components/MissingSdkWarning.vue'
import RenameAvdModal from '../components/RenameAvdModal.vue'
import DeleteAvdModal from '../components/DeleteAvdModal.vue'
import AvdCard from '../components/AvdCard.vue'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import styles from './Home.module.css'

const store = useAvdStore()

const showEditDialog = ref(false)
const editAvd = ref(null)

const showDeleteDialog = ref(false)
const deleteAvdTarget = ref(null)

const showEnvInfoDialog = ref(false)
const isEnvInfoClosing = ref(false)

// Context menu
const menuAvd = ref(null)
const menuPosition = ref({ x: 0, y: 0 })

// Toasts
const toastMessage = ref('')
let toastTimeout = null

const sdkMissing = ref(localStorage.getItem('avd_sdk_missing') === 'true')
const androidEnvChecked = ref(false)

const isWindows = navigator.userAgent.includes('Windows')

function openFolder(path) {
  if (path) {
    OpenAvdFolder(path)
  }
}

function toggleMenu(avd, event) {
  if (menuAvd.value === avd) {
    menuAvd.value = null
  } else {
    menuAvd.value = avd
    const rect = event.currentTarget.getBoundingClientRect()
    // Open menu slightly below and to the left of the button for better UX
    menuPosition.value = {
      x: rect.left - 110,
      y: rect.bottom + 8
    }
  }
}

function openEditDialog(avd) {
  editAvd.value = avd
  showEditDialog.value = true
  menuAvd.value = null
}

async function onRenameSuccess() {
  await initData();
  showToast('AVD Renamed ✅');
}

function onRenameError(err) {
  showToast('Failed to rename AVD ❌');
  console.error(err);
}

function openEnvVars() {
  OpenEnvironmentVariables().catch(err => {
    showToast('Failed to open environment settings ❌')
    console.error(err)
  })
}

function openEnvInfo() {
  showEnvInfoDialog.value = true
  isEnvInfoClosing.value = false
}

function closeEnvInfo() {
  isEnvInfoClosing.value = true
}

function handleEnvInfoAnimationEnd(e) {
  if (e.target !== e.currentTarget) return;
  if (isEnvInfoClosing.value) {
    showEnvInfoDialog.value = false
    isEnvInfoClosing.value = false
  }
}

async function selectSdkPath() {
  try {
    const path = await SelectAndSaveSdkPath()
    if (path) {
      sdkMissing.value = false
      localStorage.setItem('avd_sdk_missing', 'false')
      await initData()
    }
  } catch (err) {
    showToast(`${err}`)
    console.error(err)
  }
}

async function initData() {
  try {
    const env = await GetAndroidSdkEnv()

    if (!env.path || env.path === '') {
      sdkMissing.value = true
      localStorage.setItem('avd_sdk_missing', 'true')
      store.avds = []
      return
    }

    sdkMissing.value = false
    localStorage.setItem('avd_sdk_missing', 'false')

    // Step 1: Get AVD names
    const avdNames = await ListAVDs()

    // Step 2: Get info + running state
    const infos = await Promise.all(
      avdNames.map(name => GetAvdInfo(name).catch(e => {
        console.error(`Error fetching info for ${name}:`, e)
        return null
      }))
    )

    // Step 2.5: Remove AVDs that are no longer present
    store.avds = store.avds.filter(a => avdNames.includes(a.name))

    // Step 3: Show/Update cards
    for (let i = 0; i < avdNames.length; i++) {
      const name = avdNames[i]
      const info = infos[i]
      const isRunning = info?.running ?? false

      const existingIndex = store.avds.findIndex(a => a.name === name)
      const update = {
        state: isRunning ? AvdState.RUNNING : AvdState.POWERED_OFF,
        path: info?.path
      }

      if (existingIndex !== -1) {
        store.updateAvdStatus(name, {
          ...update,
          displayName: info?.displayName,
          apiLevel: info?.apiLevel,
          androidVersion: info?.androidVersion,
          androidCodename: info?.androidCodename,
          abi: info?.abi,
          ramSize: info?.ramSize,
          resolution: info?.resolution,
          hasGooglePlay: info?.hasGooglePlay
        })
      } else {
        store.avds.push({
          name,
          ...update,
          displayName: info?.displayName,
          apiLevel: info?.apiLevel,
          androidVersion: info?.androidVersion,
          androidCodename: info?.androidCodename,
          abi: info?.abi,
          ramSize: info?.ramSize,
          resolution: info?.resolution,
          hasGooglePlay: info?.hasGooglePlay,
          diskUsage: null
        })
      }
    }

    // Step 4: Fetch disk usage in the background
    Promise.all(
      avdNames.map(name =>
        GetAvdDiskUsage(name)
          .then(usage => {
            store.updateAvdStatus(name, { diskUsage: usage })
          })
          .catch(e => console.error(`Error fetching disk usage for ${name}:`, e))
      )
    )
  } catch (err) {
    showToast(`Error... ${err}`)
    console.log(err);
    store.avds = []
  } finally {
    androidEnvChecked.value = true
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
  deleteAvdTarget.value = avd
  showDeleteDialog.value = true
}

function onDeleteSuccess(avdName) {
  store.avds = store.avds.filter(a => a.name !== avdName);
  showToast('AVD Deleted ✅');
}

function onDeleteError(err) {
  showToast('Failed to delete AVD ❌');
  console.error(err);
}

// Close menu on click outside
function onClickOutside(event) {
  // Check if click is inside the menu itself
  if (event.target.closest(`.${styles.contextMenu}`)) return

  // Check if click is on any card's menu button
  if (!event.target.closest('.avd-menu-btn')) {
    menuAvd.value = null
  }
}

onMounted(async () => {
  await initData()
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
