<template>
  <div v-if="showInternal" :class="[styles.modalOverlay, { [styles.closing]: isClosing }]"
    @animationend="handleAnimationEnd">
    <div :class="styles.modal" @click.stop>
      <h3>Delete AVD</h3>
      <p>Are you sure you want to permanently delete <strong>{{ avd?.name }} from your disk</strong>?</p>
      <p :class="styles.warningTextSmall">This action cannot be undone.</p>
      <div :class="styles.modalActions">
        <button :class="[styles.btn, styles.btnSecondary]" @click="closeModal" :disabled="isDeleting">Cancel</button>
        <button :class="[styles.btn, styles.btnDanger]" @click="confirmDelete" :disabled="isDeleting">
          <template v-if="isDeleting">
            <v-icon name="hi-refresh" animation="spin" :scale="0.8" :class="styles.loadingIcon" />
            <span>Deleting...</span>
          </template>
          <span v-else>Delete</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { DeleteAVD, StopAVD } from '../../wailsjs/go/app/App'
import { AvdState } from '../enums/avdState'
import styles from './DeleteAvdModal.module.css'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  avd: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'success', 'error'])

const showInternal = ref(false)
const isClosing = ref(false)
const isDeleting = ref(false)

watch(() => props.show, (newVal) => {
  if (newVal) {
    showInternal.value = true
    isClosing.value = false
    isDeleting.value = false
  } else {
    if (!isClosing.value) {
      showInternal.value = false
    }
  }
}, { immediate: true })

function closeModal() {
  if (isDeleting.value) return;
  isClosing.value = true
}

function handleAnimationEnd(e) {
  if (e.target !== e.currentTarget) return;

  if (isClosing.value) {
    showInternal.value = false
    isClosing.value = false
    emit('close')
  }
}

async function confirmDelete() {
  if (!props.avd || isDeleting.value) return;

  try {
    isDeleting.value = true
    // Stop it first if it is running
    if (props.avd.state === AvdState.RUNNING) {
      await StopAVD(props.avd.name);
    }
    await DeleteAVD(props.avd.name);
    emit('success', props.avd.name);
    closeModal();
  } catch (err) {
    emit('error', err);
  } finally {
    isDeleting.value = false;
  }
}
</script>
