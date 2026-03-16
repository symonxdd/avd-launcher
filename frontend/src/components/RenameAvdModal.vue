<template>
  <div v-if="showInternal" :class="[styles.modalOverlay, { [styles.closing]: isClosing }]"
    @animationend="handleAnimationEnd">
    <div :class="styles.modal" @click.stop>
      <h3>Rename AVD</h3>
      <p>Enter a new name for <strong>{{ avd?.displayName || avd?.name }}</strong>:</p>
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
        <button :class="[styles.btn, styles.btnSecondary]" @click="closeModal"
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
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { RenameAVD } from '../../wailsjs/go/app/App'
import styles from './RenameAvdModal.module.css'

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
const editAvdName = ref('')
const isRenaming = ref(false)

watch(() => props.show, (newVal) => {
  if (newVal) {
    showInternal.value = true
    isClosing.value = false
    editAvdName.value = props.avd ? (props.avd.displayName || props.avd.name) : ''
    isRenaming.value = false
  } else {
    // If parent forces close without animation
    if (!isClosing.value) {
        showInternal.value = false
    }
  }
}, { immediate: true })

const isRenameDisabled = computed(() => {
  if (!props.avd || isRenaming.value) return true;
  const oldDisplayName = props.avd.displayName || props.avd.name;
  const newName = editAvdName.value.trim();
  return newName === '' || oldDisplayName === newName;
})

const isCaseOnlyChange = computed(() => {
  if (!props.avd) return false;
  const oldDisplayName = props.avd.displayName || props.avd.name;
  const newName = editAvdName.value.trim();
  return oldDisplayName !== newName && oldDisplayName.toLowerCase() === newName.toLowerCase();
})

function closeModal() {
  if (isRenaming.value) return;
  isClosing.value = true
}

function handleAnimationEnd(e) {
  // Only trigger when the overlay's own animation ends
  if (e.target !== e.currentTarget) return;

  if (isClosing.value) {
    showInternal.value = false
    isClosing.value = false
    emit('close')
  }
}

async function saveEdit() {
  if (!props.avd || isRenameDisabled.value) return;

  const oldInternalName = props.avd.name;
  const oldDisplayName = props.avd.displayName || oldInternalName;
  const newName = editAvdName.value.trim();
  
  if (oldDisplayName === newName || newName === '') {
    closeModal();
    return;
  }

  try {
    isRenaming.value = true;
    await RenameAVD(oldInternalName, newName);
    // Let the parent handle store updates and toasts
    emit('success', { oldName: oldInternalName, newName });
    closeModal();
  } catch (err) {
    emit('error', err);
  } finally {
    isRenaming.value = false;
  }
}
</script>
