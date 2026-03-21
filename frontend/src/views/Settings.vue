<template>
  <div :class="styles.settingsContainer">
    <div :class="styles.pageHeader">
      <div :class="styles.pageTitleContainer">
        <h2 :class="styles.pageTitle">Settings</h2>
      </div>
    </div>
    <div :class="styles.settingsGrid">
      <!-- --- System Settings --- -->
      <div :class="styles.settingsGroup">
        <div :class="styles.headingRow">
          <h5 :class="styles.labelHeading">
            <v-icon name="hi-adjustments" :class="styles.headingIcon" />
            System
          </h5>
          <button v-if="isWindows" :class="styles.envLearnMore" @click="openEnvInfo">
            Environment variables?
          </button>
        </div>

        <div :class="styles.subGroupItems">
          <div
            :class="[styles.sdkStatusRow, { [styles.sdkFound]: !isCheckingSdk && androidSdkEnv?.path, [styles.sdkMissing]: !isCheckingSdk && !androidSdkEnv?.path }]">
            <div :class="styles.statusInfo">
              <div :class="styles.statusLabelRow">
                <span :class="styles.statusLabel">Android SDK Status</span>
                <div v-if="!isCheckingSdk" :class="styles.infoTooltipTrigger">
                  <v-icon name="hi-information-circle" :class="styles.infoIcon" />
                  <div :class="styles.infoTooltip">
                    <template v-if="androidSdkEnv?.path">
                      <template v-if="androidSdkEnv.source === 'user selected path'">
                        Detected via
                        <span :class="styles.configLink" @click="OpenConfigFolder">
                          user selected path
                        </span>
                      </template>
                      <template v-else>
                        Detected via {{ androidSdkEnv.source }}
                      </template>
                    </template>
                    <template v-else>
                      Android SDK not detected. This is required to manage and launch AVDs.
                    </template>
                  </div>
                </div>
              </div>
              <span :class="[styles.statusPath, { [styles.hidden]: isCheckingSdk }]">
                {{ androidSdkEnv?.path || SDK_MISSING_MSG }}
              </span>
            </div>
            <div :class="[styles.statusIconBox, { [styles.hidden]: isCheckingSdk }]">
              <v-icon v-if="androidSdkEnv?.path" name="hi-check" scale="1.2" />
              <v-icon v-else name="hi-x" scale="1.2" />
            </div>
          </div>

          <div
            :class="[styles.sdkStatusRow, { [styles.accelFound]: !isCheckingAccel && accelInfo?.status === 'available', [styles.accelMissing]: !isCheckingAccel && accelInfo?.status !== 'available' }]">
            <div :class="styles.statusInfo">
              <div :class="styles.statusLabelRow">
                <span :class="styles.statusLabel">VM Acceleration</span>
                <div v-if="!isCheckingAccel" :class="styles.infoTooltipTrigger">
                  <v-icon name="hi-information-circle" :class="styles.infoIcon" />
                  <div :class="styles.infoTooltip">
                    <template v-if="accelInfo?.status === 'available'">
                      {{ accelInfo.hypervisor }} is active.
                    </template>
                    <template v-else>
                      No hypervisor detected. Emulators will run significantly slower.
                    </template>
                    <span :class="styles.externalLink" @click="openAccelDocs">
                      Learn more
                      <v-icon name="hi-external-link" :class="styles.externalLinkIcon" />
                    </span>
                  </div>
                </div>
              </div>
              <span :class="[styles.statusPath, { [styles.hidden]: isCheckingAccel }]">
                {{ accelInfo ? accelInfo.details : 'Checking...' }}
              </span>
            </div>
            <div :class="[styles.statusIconBox, { [styles.hidden]: isCheckingAccel }]">
              <v-icon v-if="accelInfo?.status === 'available'" name="hi-check" scale="1.2" />
              <v-icon v-else name="hi-x" scale="1.2" />
            </div>
          </div>
        </div>
      </div>

      <div :class="styles.settingsRow">
        <!-- --- Interface Settings --- -->
        <div :class="styles.settingsGroup">
          <h5 :class="styles.labelHeading">
            <v-icon name="hi-view-grid" :class="styles.headingIcon" />
            Interface
          </h5>

          <div :class="styles.subGroupItems">
            <div v-if="isDarkMode" :class="styles.toggleRow">
              <label for="trueBlackToggle" :class="styles.toggleLabel">Enable True Black</label>
              <label :class="styles.switch">
                <input id="trueBlackToggle" type="checkbox" :checked="themeStore.trueBlack"
                  @change="themeStore.toggleTrueBlack" />
                <span :class="styles.slider"></span>
              </label>
            </div>
            <div v-else :class="styles.emptyStateText">
              Switch to Dark Mode to enable True Black.
            </div>
          </div>
        </div>

        <!-- --- Update Available Notice --- -->
        <div v-if="updateStore.remoteVersion && isRemoteVersionNewer(updateStore.remoteVersion, appVersion)"
          :class="styles.settingsGroup">
          <div>
            <h5 :class="styles.labelHeading">
              <v-icon name="hi-gift" :class="styles.headingIcon" />
              New update ready
            </h5>
            <span :class="styles.versionSubtext" style="margin-left: calc(1.1rem + 14px);">v{{ updateStore.remoteVersion
            }}</span>
          </div>

          <div :class="styles.subGroupItems">
            <div :class="styles.updateActionState">
              <button :class="styles.actionButton" @click="handleOpenGithubRelease">
                <v-icon name="fa-github" />
                View on GitHub
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- --- Footer Section --- -->
    <div :class="styles.appFooter">
      <div :class="styles.footerMeta">
        v{{ appVersion }} ({{ environment }}{{ versionStatus }})
        <span :class="styles.dotSeparator">•</span>
        Powered by Vue, Go & Wails
      </div>

      <div :class="styles.signature">
        <div :class="styles.tooltipTrigger">
          <span :class="styles.dictionaryLink">Sedulously</span>
          <div :class="styles.customTooltip">
            <span :class="styles.tooltipType">(adverb)</span> — in a way that shows great care and persistent effort.
            <br />
            <small style="opacity: 0.7; margin-top: 4px; display: block;">Synonyms: diligently, persistently</small>
          </div>
        </div> engineered with <span :class="styles.heart">❤️‍🩹</span> by Symon
      </div>
    </div>

    <!-- Env Vars Explanation Modal -->
    <EnvInfoModal :show="showEnvInfoDialog" :is-closing="isEnvInfoClosing" @close="closeEnvInfo"
      @animationend="handleEnvInfoAnimationEnd" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { GetAndroidSdkEnv, OpenConfigFolder, CheckAcceleration } from '../../wailsjs/go/services/SystemService'
import { GetLatestVersion } from '../../wailsjs/go/services/UpdateService'
import { BrowserOpenURL } from '../../wailsjs/runtime'
import { useThemeStore } from '../stores/themeStore'
import { useUpdateStore } from '../stores/updateStore'
import EnvInfoModal from '../components/EnvInfoModal.vue'
import styles from './Settings.module.css'

const themeStore = useThemeStore()
const updateStore = useUpdateStore()
const androidSdkEnv = ref(null)
const SDK_MISSING_MSG = 'No SDK path configured. Please check the AVD tab for more information.'
const isCheckingSdk = ref(true)
const accelInfo = ref(null)
const isCheckingAccel = ref(true)
const appVersion = typeof __APP_VERSION__ !== "undefined" ? __APP_VERSION__ : "1.0.0"
const environment = import.meta.env.MODE === 'development' ? 'dev' : 'release'
const isWindows = navigator.userAgent.includes('Windows')

const showEnvInfoDialog = ref(false)
const isEnvInfoClosing = ref(false)

const isDarkMode = computed(() => {
  if (themeStore.theme === 'dark') return true
  if (themeStore.theme === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  return false
})

const isRemoteVersionNewer = (remote, local) => {
  const parse = (v) => v.split(".").map(Number);
  const [r1, r2, r3] = parse(remote);
  const [l1, l2, l3] = parse(local);
  if (r1 > l1) return true;
  if (r1 === l1 && r2 > l2) return true;
  if (r1 === l1 && r2 === l2 && r3 > l3) return true;
  return false;
}

const isUpToDate = computed(() => {
  if (!updateStore.remoteVersion) return true;
  return !isRemoteVersionNewer(updateStore.remoteVersion, appVersion);
})

const versionStatus = computed(() => {
  if (!updateStore.remoteVersion) return '';
  return isUpToDate.value ? ', latest' : ', update available';
})

const fetchAndroidSdkEnv = async () => {
  try {
    androidSdkEnv.value = await GetAndroidSdkEnv()
  } catch (error) {
    console.error('Error while running GetAndroidSdkEnv():', error)
  }
}

const handleOpenGithubRelease = () => {
  try {
    BrowserOpenURL(updateStore.remoteReleaseUrl)
  } catch (err) {
    console.error('Failed to open GitHub link in default browser:', err)
  }
}

function openAccelDocs() {
  BrowserOpenURL('https://developer.android.com/studio/run/emulator-acceleration')
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

onMounted(async () => {
  try {
    await fetchAndroidSdkEnv()
  } finally {
    isCheckingSdk.value = false
  }

  // Check acceleration in parallel with update check
  CheckAcceleration().then(info => {
    accelInfo.value = info
  }).catch(err => {
    console.error('Error checking acceleration:', err)
    accelInfo.value = { status: 'unavailable', hypervisor: 'Unknown', details: 'Check failed.' }
  }).finally(() => {
    isCheckingAccel.value = false
  })

  await updateStore.checkForUpdate()
})
</script>
