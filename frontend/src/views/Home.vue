<template>
  <div :class="styles.homeContainer">

    <div :class="styles.pageHeader">
      <div v-if="!sdkMissing" :class="styles.pageTitleContainer">
        <h2 :class="styles.pageTitle">Installed AVDs</h2>
        <span :class="styles.countBadge" v-show="store.avds.length">{{ store.avds.length }}</span>
      </div>

      <div v-if="sdkMissing" :class="styles.androidHomeWarning">
        <v-icon name="hi-exclamation" :class="styles.warningIcon" :scale="1.5" />
        <div :class="styles.warningContent">
          <div :class="styles.warningText">
            <div :class="styles.warningTextFirstLine">Android SDK not found</div>
            <p>AVD Launcher could not locate your Android SDK in common locations nor via the ANDROID_HOME environment
              variable.</p>
            <p>Use the 'Select SDK Location' button to choose where your Android SDK is installed, or set the
              ANDROID_HOME environment variable to that location.</p>
            <button v-if="isWindows" :class="styles.envLearnMore" @click="openEnvInfo">
              What are environment variables?
            </button>
          </div>
          <div :class="styles.warningActions">
            <button :class="[styles.btn, styles.btnPrimary, styles.warningBtn]" @click="selectSdkPath">
              Select SDK Location
            </button>
            <button v-if="isWindows" :class="[styles.btn, styles.btnSecondary, styles.warningBtn]" @click="openEnvVars">
              Open Environment Variables
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-show="store.avds.length && !sdkMissing" :class="styles.avdGrid">
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
            <v-icon name="hi-dots-horizontal" />
          </button>
        </div>

        <!-- Status indicator: dot + label (Flows under name/avatar) -->
        <div :class="[styles.avdStatus, styles[getStateClass(avd.state)]]">
          <div :class="styles.statusDotContainer">
            <span :class="styles.statusDot"></span>
            <span :class="styles.statusPulse"></span>
          </div>
          <span :class="styles.statusLabel">{{ avd.state }}</span>
          <div v-if="avd.diskUsage" :class="styles.diskUsage">
            <v-icon name="hi-database" :scale="0.7" />
            <span>{{ avd.diskUsage }}</span>
          </div>
        </div>

        <!-- Animated context menu -->
        <Teleport to="body">
          <transition name="fade-fast">
            <div v-if="menuAvd === avd" :class="styles.contextMenu"
              :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }">
              <button @click="openFolder(avd.path)">
                <v-icon name="hi-folder-open" :scale="0.85" />
                <span>Show on Disk</span>
              </button>
              <button @click="openEditDialog(avd)">
                <v-icon name="hi-pencil" :scale="0.85" />
                <span>Rename</span>
              </button>
              <button @click="openDeleteDialog(avd)" :class="styles.deleteItem">
                <v-icon name="hi-trash" :scale="0.85" />
                <span>Delete</span>
              </button>
            </div>
          </transition>
        </Teleport>

        <!-- Premium Action buttons with labels -->
        <div :class="styles.avdActions">
          <template v-if="avd.state === AvdState.POWERED_OFF">
            <button :class="[styles.actionBtnPremium, styles.launch]" @click="startAVD(avd, false)">
              <v-icon name="hi-play" :scale="1.1" />
              <span>Launch</span>
            </button>
            <button :class="[styles.actionBtnPremium, styles.coldboot]" @click="startAVD(avd, true)">
              <v-icon name="fa-snowflake" :scale="0.85" />
              <span>Cold Boot</span>
            </button>
          </template>
          <button :class="[styles.actionBtnPremium, styles.stop]" v-else :disabled="avd.state !== AvdState.RUNNING"
            @click="stopAVD(avd.name)">
            <v-icon name="hi-stop" :scale="1.1" />
            <span>Stop</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit AVD name Modal -->
    <div v-if="showEditDialog" :class="[styles.modalOverlay, { [styles.closing]: isEditClosing }]"
      @animationend="handleEditAnimationEnd">
      <div :class="styles.modal" @click.stop>
        <h3>Rename AVD</h3>
        <p>Enter a new name for <strong>{{ editAvd?.name }}</strong>:</p>
        <input v-model="editAvdName" id="edit-avd-name" name="edit-avd-name" placeholder="New AVD Name"
          @keyup.enter="saveEdit" spellcheck="false" autocorrect="off" autocapitalize="off" />
        <div :class="styles.allowedCharsInfo">
          <v-icon name="hi-sparkles" :scale="0.8" />
          <span>The name can contain uppercase or lowercase letters, numbers, periods, underscores, parentheses, dashes,
            and spaces.</span>
        </div>
        <transition name="fade-slide">
          <div v-if="isCaseOnlyChange" :class="styles.renameHint">
            <v-icon name="hi-information-circle" :scale="0.8" />
            <span>To change only capitalization, first rename to something else, then back to the desired name with the
              correct case.</span>
          </div>
        </transition>
        <div :class="styles.modalActions">
          <button :class="[styles.btn, styles.btnSecondary]" @click="closeEditDialog"
            :disabled="isRenaming">Cancel</button>
          <button :class="[styles.btn, styles.btnPrimary]" @click="saveEdit" :disabled="isRenameDisabled">
            <template v-if="isRenaming">
              <v-icon name="hi-refresh" animation="spin" :scale="0.8" :class="styles.loadingIcon" />
              <span>Renaming...</span>
            </template>
            <span v-else>Confirm</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete AVD Modal -->
    <div v-if="showDeleteDialog" :class="[styles.modalOverlay, { [styles.closing]: isDeleteClosing }]"
      @animationend="handleDeleteAnimationEnd">
      <div :class="styles.modal" @click.stop>
        <h3>Delete AVD</h3>
        <p>Are you sure you want to permanently delete <strong>{{ deleteAvdTarget?.name }} from your disk</strong>?</p>
        <p :class="styles.warningTextSmall">This action cannot be undone.</p>
        <div :class="styles.modalActions">
          <button :class="[styles.btn, styles.btnSecondary]" @click="closeDeleteDialog">Cancel</button>
          <button :class="[styles.btn, styles.btnDanger]" @click="confirmDeleteDialog">Delete</button>
        </div>
      </div>
    </div>

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
import { ref, onMounted, onBeforeUnmount, nextTick, computed } from 'vue'
import { ListAVDs, StartAVD, StopAVD, ListRunningAVDs, GetAndroidSdkEnv, OpenEnvironmentVariables, RenameAVD, DeleteAVD, SelectAndSaveSdkPath, GetAvdInfo, OpenAvdFolder } from '../../wailsjs/go/app/App'
import EnvInfoModal from '../components/EnvInfoModal.vue'
import { useAvdStore } from '../stores/avdStore'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'
import styles from './Home.module.css'

const store = useAvdStore()

const showEditDialog = ref(false)
const isEditClosing = ref(false)
const editAvd = ref(null)
const editAvdName = ref('')
const isRenaming = ref(false)

const isRenameDisabled = computed(() => {
  if (!editAvd.value || isRenaming.value) return true;
  const oldName = editAvd.value.name;
  const newName = editAvdName.value.trim();
  return newName === '' || oldName === newName;
})

const isCaseOnlyChange = computed(() => {
  if (!editAvd.value) return false;
  const oldName = editAvd.value.name;
  const newName = editAvdName.value.trim();
  return oldName !== newName && oldName.toLowerCase() === newName.toLowerCase();
})

const showDeleteDialog = ref(false)
const isDeleteClosing = ref(false)
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
  showEditDialog.value = true
  editAvd.value = avd
  editAvdName.value = avd.name
  menuAvd.value = null
}

function closeEditDialog() {
  isEditClosing.value = true
}

function handleEditAnimationEnd(e) {
  // Only trigger when the overlay's own animation ends
  if (e.target !== e.currentTarget) return;

  if (isEditClosing.value) {
    showEditDialog.value = false
    isEditClosing.value = false
    editAvd.value = null
  }
}

async function saveEdit() {
  if (editAvd.value) {
    const oldName = editAvd.value.name;
    const newName = editAvdName.value.trim();
    if (oldName === newName || newName === '') {
      closeEditDialog();
      return;
    }
    try {
      isRenaming.value = true;
      await RenameAVD(oldName, newName);
      // Update local store to reflect the change
      editAvd.value.name = newName;
      showToast('AVD Renamed ✅');
      closeEditDialog();
    } catch (err) {
      showToast('Failed to rename AVD ❌');
      console.error(err);
    } finally {
      isRenaming.value = false;
    }
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
    // env = '' // for debug purposes

    // TEMPORARILY FORCE THE WARNING FOR TESTING:
    // sdkMissing.value = true;
    // localStorage.setItem('avd_sdk_missing', 'true');
    // store.avds = [];
    // return;

    if (!env.path || env.path === '') {
      sdkMissing.value = true
      localStorage.setItem('avd_sdk_missing', 'true')
      store.avds = [] // Clear any stale AVDs
      return
    }

    sdkMissing.value = false
    localStorage.setItem('avd_sdk_missing', 'false')

    const avds = await ListAVDs()
    const runningAvds = await ListRunningAVDs()

    for (const name of avds) {
      const isRunning = runningAvds?.includes(name)

      let info = null
      try {
        info = await GetAvdInfo(name)
      } catch (e) {
        console.error(`Error fetching info for ${name}:`, e)
      }

      // If already in store, update its state
      const existing = store.avds.find(a => a.name === name)
      const update = {
        state: isRunning ? AvdState.RUNNING : AvdState.POWERED_OFF,
        diskUsage: info?.diskUsage,
        path: info?.path
      }

      if (existing) {
        store.updateAvdStatus(name, update)
      } else {
        // Otherwise, add it
        store.avds.push({
          name,
          ...update,
          hover: false
        })
      }
    }
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

function closeDeleteDialog() {
  isDeleteClosing.value = true
}

function handleDeleteAnimationEnd(e) {
  // Only trigger when the overlay's own animation ends
  if (e.target !== e.currentTarget) return;

  if (isDeleteClosing.value) {
    showDeleteDialog.value = false
    isDeleteClosing.value = false
    deleteAvdTarget.value = null
  }
}

async function confirmDeleteDialog() {
  const avd = deleteAvdTarget.value
  if (!avd) return

  try {
    // Just in case it's currently running, stop it first.
    if (avd.state === AvdState.RUNNING) {
      await stopAVD(avd.name);
    }
    await DeleteAVD(avd.name);

    // Remove from store
    store.avds = store.avds.filter(a => a.name !== avd.name);
    showToast('AVD Deleted ✅');
  } catch (err) {
    showToast('Failed to delete AVD ❌');
    console.error(err);
  } finally {
    closeDeleteDialog()
  }
}

// Close menu on click outside
function onClickOutside(event) {
  if (!event.target.closest(`.${styles.menuButton}`) && !event.target.closest(`.${styles.contextMenu}`)) {
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
