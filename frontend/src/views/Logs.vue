<template>
  <div :class="styles.logsContainer" @keydown.stop>
    <div :class="styles.pageHeader">
      <div :class="styles.pageTitleContainer">
        <h2 :class="styles.pageTitle">Logs</h2>
      </div>
    </div>

    <div v-if="searchActive" :class="styles.logSearchBar">
      <input ref="searchInput" v-model="searchQuery" placeholder="Search logs..." :class="styles.searchInput"
        @keydown.esc.prevent="deactivateSearch" />
      <div :class="styles.searchControls">
        <button @click="prevMatch" title="Previous match">
          <v-icon name="hi-arrow-up" />
        </button>
        <button @click="nextMatch" title="Next match">
          <v-icon name="hi-arrow-down" />
        </button>
        <span :class="styles.matchCounter">{{ matchCounter }}</span>
        <button :class="styles.closeBtn" @click="deactivateSearch">×</button>
      </div>
    </div>
    <div :class="styles.logOutputWrapper">
      <div :class="styles.logOutput" ref="logContainer" v-html="highlightedLogs"></div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed, watch, nextTick } from 'vue'
import { useAvdStore } from '../stores/avdStore'
import styles from './Logs.module.css'

const store = useAvdStore()

const searchActive = ref(false)
const searchQuery = ref('')
const searchInput = ref(null)
const logContainer = ref(null)
const currentMatchIndex = ref(0)
const matches = ref([])
const highlightedLogs = ref('')

// Toggle search mode
const activateSearch = () => {
  searchActive.value = true
  nextTick(() => {
    searchInput.value?.focus()
  })
}
const deactivateSearch = () => {
  searchActive.value = false
  searchQuery.value = ''
  currentMatchIndex.value = 0
}

const matchCounter = computed(() => {
  if (!searchQuery.value || matches.value.length === 0) return 'No results'
  return `${currentMatchIndex.value + 1} of ${matches.value.length}`
})

// Ctrl+F handler
const onKeydown = (e) => {
  if (e.ctrlKey && e.key === 'f') {
    e.preventDefault()

    if (searchActive.value) {
      // If the search bar is already active, select all text in the input field
      searchInput.value.select() // This will always select the text inside the input
    } else {
      // Otherwise, activate the search and focus on the input field
      activateSearch()
    }
  } else if (e.key === 'Escape' && searchActive.value) {
    e.preventDefault()
    deactivateSearch()
  }
}


function scrollToMatch(index) {
  const container = logContainer.value
  if (!container) return

  const elements = container.querySelectorAll(`.${styles.highlight}`)
  if (elements[index]) {
    elements[index].scrollIntoView({ behavior: 'instant' })
  }
}

function nextMatch() {
  if (matches.value.length === 0) return
  currentMatchIndex.value = (currentMatchIndex.value + 1) % matches.value.length
  // updateHighlighting()
  scrollToMatch(currentMatchIndex.value)
}

function prevMatch() {
  if (matches.value.length === 0) return
  currentMatchIndex.value = (currentMatchIndex.value - 1 + matches.value.length) % matches.value.length
  // updateHighlighting()
  scrollToMatch(currentMatchIndex.value)
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)

  // Scroll to bottom when the component is mounted
  nextTick(() => {
    const container = logContainer.value
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  })
})
onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

watch(
  [() => store.logs, searchQuery],
  ([logs, query]) => {
    matches.value = []
    currentMatchIndex.value = 0

    if (!query) {
      highlightedLogs.value = logs
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
      return
    }

    const safeQuery = query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    const regex = new RegExp(safeQuery, 'gi')

    let i = 0
    highlightedLogs.value = logs.replace(regex, match => {
      const className = i === currentMatchIndex.value ? `${styles.highlight} ${styles.activeHighlight}` : styles.highlight
      matches.value.push(i++)
      return `<span class="${className}">${match}</span>`
    })

    console.log('[Search] Matches found:', matches.value.length)
  },
  { immediate: true }
)

watch(() => store.logs, () => {
  const shouldAutoScroll = isNearBottom() && !searchActive.value

  nextTick(() => {
    const container = logContainer.value
    if (container && shouldAutoScroll) {
      container.scrollTop = container.scrollHeight
    }
  })
})

watch(currentMatchIndex, () => {
  if (matches.value.length === 0) return
  updateHighlighting()
  scrollToMatch(currentMatchIndex.value)
})

function isNearBottom(threshold = 100) {
  const container = logContainer.value
  if (!container) return false

  return container.scrollHeight - container.scrollTop - container.clientHeight < threshold
}

// Modify your updateHighlighting method:
function updateHighlighting() {
  const logs = store.logs
  const query = searchQuery.value

  if (!query) {
    highlightedLogs.value = logs
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
    return
  }

  const safeQuery = query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const regex = new RegExp(safeQuery, 'gi')

  let matchIndex = 0
  matches.value = []
  highlightedLogs.value = logs.replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(regex, match => {
    const isActive = matchIndex === currentMatchIndex.value
    const className = isActive ? `${styles.highlight} ${styles.activeHighlight}` : styles.highlight
    matches.value.push(matchIndex)
    matchIndex++
    return `<span class="${className}">${match}</span>`
  })
}
</script>
