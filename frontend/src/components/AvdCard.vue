<template>
  <div :class="[styles.avdCard, { [styles.running]: avd.state === AvdState.RUNNING }]"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false">

    <!-- Card header: avatar + info (name) + menu -->
    <div :class="styles.avdCardHeader">
      <div :class="styles.avdAvatarContainer">
        <div :class="styles.avdAvatar">{{ avd.name.charAt(0).toUpperCase() }}</div>
      </div>
      <div :class="styles.avdNameContainer">
        <div :class="styles.avdName">{{ avd.displayName || avd.name }}</div>
        <div :class="styles.infoTooltipTrigger"
             @mouseenter="showInfoTooltip"
             @mouseleave="hideInfoTooltip"
             ref="tooltipTrigger">
          <v-icon name="hi-information-circle" :class="styles.infoIcon" />
          <Teleport to="body">
            <transition name="fade-tooltip">
              <div v-if="isTooltipVisible" ref="infoTooltip" :class="styles.infoTooltip" :style="tooltipStyle">
                <div :class="styles.tooltipHeader">AVD Details</div>
                <div :class="styles.tooltipGrid">
                  <div :class="styles.tooltipColumn">
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Platform:</span>
                        <span :class="styles.tooltipValue">{{ avd.androidVersion || 'Android' }}</span>
                      </div>
                      <small :class="styles.tooltipExplanation">API {{ avd.apiLevel }} – {{ avd.androidCodename }}</small>
                    </div>
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Architecture:</span>
                        <span :class="styles.tooltipValue">{{ avd.abi || '?' }}</span>
                      </div>
                      <small :class="styles.tooltipExplanation">The emulated CPU architecture.</small>
                    </div>
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Google Play:</span>
                        <span :class="styles.tooltipValue">{{ avd.hasGooglePlay ? 'Yes' : 'No' }}</span>
                      </div>
                        <small :class="styles.tooltipExplanation">Includes the Play Store app and supports Google services (Maps, Push, etc.).</small>
                    </div>
                  </div>

                  <div :class="styles.tooltipColumn">
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Memory:</span>
                        <span :class="styles.tooltipValue">{{ avd.ramSize || '?' }} MB</span>
                      </div>
                      <small :class="styles.tooltipExplanation">Memory allocated to this device.</small>
                    </div>
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Resolution:</span>
                        <span :class="styles.tooltipValue">{{ avd.resolution || '?' }}</span>
                      </div>
                      <small :class="styles.tooltipExplanation">Emulated screen pixel dimensions.</small>
                    </div>
                    <div :class="styles.tooltipRow">
                      <div :class="styles.tooltipValueGroup">
                        <span :class="styles.tooltipLabel">Storage:</span>
                        <span :class="styles.tooltipValue">{{ avd.diskUsage || 'Calculating...' }}</span>
                      </div>
                      <small :class="styles.tooltipExplanation">Total size on disk.</small>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
          </Teleport>
        </div>
      </div>
      <button :class="[styles.menuButton, 'avd-menu-btn', { [styles.menuVisible]: isHovered }]" @click="$emit('toggle-menu', avd, $event)">
        <v-icon name="hi-dots-horizontal" />
      </button>
    </div>

    <!-- Status indicator: dot + label -->
    <div :class="[styles.avdStatus, styles[getStateClass(avd.state)]]">
      <div :class="styles.statusDotContainer">
        <span :class="styles.statusDot"></span>
        <span :class="styles.statusPulse"></span>
      </div>
      <span :class="styles.statusLabel">{{ avd.state }}</span>
    </div>

    <!-- Premium Action buttons with labels -->
    <div :class="styles.avdActions">
      <template v-if="avd.state === AvdState.POWERED_OFF">
        <button :class="[styles.actionBtnPremium, styles.launch]" @click="$emit('launch', avd, false)">
          <v-icon name="hi-play" :scale="1.1" />
          <span>Launch</span>
        </button>
        <button :class="[styles.actionBtnPremium, styles.coldboot]" @click="$emit('launch', avd, true)">
          <v-icon name="fa-snowflake" :scale="0.85" />
          <span>Cold Boot</span>
        </button>
      </template>
      <button :class="[styles.actionBtnPremium, styles.stop]" v-else :disabled="avd.state !== AvdState.RUNNING"
        @click="$emit('stop', avd.name)">
        <v-icon name="hi-stop" :scale="1.1" />
        <span>Stop</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { AvdState } from '../enums/avdState'
import { getStateClass } from '../utils/helper'
import styles from './AvdCard.module.css'

defineProps({
  avd: {
    type: Object,
    required: true
  }
})

defineEmits(['toggle-menu', 'launch', 'stop'])

const isHovered = ref(false)

const tooltipTrigger = ref(null)
const infoTooltip = ref(null)
const isTooltipVisible = ref(false)
const tooltipStyle = ref({})

async function showInfoTooltip() {
  if (!tooltipTrigger.value) return
  
  const rect = tooltipTrigger.value.getBoundingClientRect()
  
  // Initially position it near the trigger to safely measure
  tooltipStyle.value = {
    top: `${rect.bottom + 12}px`,
    left: `${rect.left}px`,
    '--arrow-left': '50%'
  }
  
  isTooltipVisible.value = true
  
  // Wait for the DOM element to be inserted
  await nextTick()
  if (!isTooltipVisible.value || !infoTooltip.value) return
  
  // Use offsetWidth to get the true width instead of hardcoding
  const tooltipWidth = infoTooltip.value.offsetWidth
  const center = rect.left + rect.width / 2
  
  let left = center - (tooltipWidth / 2)
  let arrowLeft = '50%'
  
  const cardEl = tooltipTrigger.value.closest(`.${styles.avdCard}`)
  if (cardEl) {
    const cardRect = cardEl.getBoundingClientRect()
    const overhang = 14
    // Snap to the card's boundaries with a slight overhang
    if (left < cardRect.left - overhang) {
      left = cardRect.left - overhang
      arrowLeft = `${center - left}px`
    } else if (left + tooltipWidth > cardRect.right + overhang) {
      left = cardRect.right - tooltipWidth + overhang
      arrowLeft = `${center - left}px`
    }
  }

  // Fallback to prevent window edge clipping
  const padding = 20
  if (left < padding) {
    left = padding
    arrowLeft = `${center - left}px`
  } else if (left + tooltipWidth > window.innerWidth - padding) {
    left = window.innerWidth - tooltipWidth - padding
    arrowLeft = `${center - left}px`
  }
  
  tooltipStyle.value = {
    top: `${rect.bottom + 12}px`,
    left: `${left}px`,
    '--arrow-left': arrowLeft
  }
}

function hideInfoTooltip() {
  isTooltipVisible.value = false
}
</script>
