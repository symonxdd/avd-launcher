<template>
  <div :class="styles.androidHomeWarning">
    <v-icon name="hi-exclamation" :class="styles.warningIcon" :scale="1.5" />
    <div :class="styles.warningContent">
      <div :class="styles.warningText">
        <div :class="styles.warningTextFirstLine">Android SDK not found</div>
        <p>AVD Launcher could not locate your Android SDK in common locations nor via the ANDROID_HOME environment
          variable.</p>
        <p>Use the 'Select SDK Location' button to choose where your Android SDK is installed, or set the
          ANDROID_HOME environment variable to that location.</p>
        <button v-if="isWindows" :class="styles.envLearnMore" @click="$emit('open-env-info')">
          What are environment variables?
        </button>
      </div>
      <div :class="styles.warningActions">
        <button :class="[styles.btn, styles.btnPrimary, styles.warningBtn]" @click="$emit('select-sdk')">
          Select SDK Location
        </button>
        <button v-if="isWindows" :class="[styles.btn, styles.btnSecondary, styles.warningBtn]" @click="$emit('open-env-vars')">
          Open Environment Variables
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import styles from './MissingSdkWarning.module.css'

defineProps({
  isWindows: {
    type: Boolean,
    required: true
  }
})

defineEmits(['select-sdk', 'open-env-vars', 'open-env-info'])
</script>
