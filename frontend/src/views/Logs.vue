<template>
  <div class="logs-container" @keydown.stop>
    <h2 class="page-title">Logs</h2>

    <div v-if="searchActive" class="log-search-bar">
      <input ref="searchInput" v-model="searchQuery" placeholder="Search logs..." class="search-input"
        @keydown.esc.prevent="deactivateSearch" />
      <div class="search-controls">
        <button @click="prevMatch" title="Previous match">
          <i class="bi bi-arrow-up"></i>
        </button>
        <button @click="nextMatch" title="Next match">
          <i class="bi bi-arrow-down"></i>
        </button>
        <span class="match-counter">{{ matchCounter }}</span>
        <button class="close-btn" @click="deactivateSearch">×</button>
      </div>
    </div>
    <div class="log-output" ref="logContainer" v-html="highlightedLogs"></div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, computed, watch, nextTick } from 'vue'
import { useAvdStore } from '../stores/avdStore'

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

  const elements = container.querySelectorAll('.highlight')
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
      const className = i === currentMatchIndex.value ? 'highlight active-highlight' : 'highlight'
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
    const className = isActive ? 'highlight active-highlight' : 'highlight'
    matches.value.push(matchIndex)
    matchIndex++
    return `<span class="${className}">${match}</span>`
  })
}
</script>

<style scoped>
.logs-container {
  position: relative;
  /* allows the search bar to position absolutely inside */
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 0 0 0;
  color: #ccc;
}

.page-title {
  font-size: 1.3rem;
  margin-bottom: 16px;
  color: #ccc;
}

.log-output {
  flex: 1;
  width: 100%;
  height: 100%;
  resize: none;
  background-color: transparent;
  color: #ccc;
  font-family: monospace;
  font-size: 0.9rem;
  padding: 12px;
  border: 1px solid #363636;
  border-radius: 6px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  /* ← this is the fix */
  outline: none;
}

.log-search-bar {
  position: absolute;
  top: -10px;
  right: 25px;
  width: 260px;
  z-index: 10;
  background-color: #2b2b2b;
  border: 1px solid #444;
  border-radius: 6px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.35);
  padding: 6px 10px;
}

.search-input {
  width: 100%;
  font-size: 0.9rem;
  padding: 5px 8px;
  border: none;
  border-radius: 4px;
  background-color: #1e1e1e;
  color: #fff;
  outline: none;
}

::v-deep(.highlight) {
  background-color: #ffd54f;
  color: #000;
  border-radius: 4px;
}

::v-deep(.active-highlight) {
  background-color: #FF9632;
  color: #000;
  border-radius: 4px;
}

.search-controls {
  display: flex;
  gap: 5px;
  align-items: center;
  margin-top: 6px;
}

.search-controls button {
  background: #3a3a3a;
  border: none;
  color: #fff;
  padding: 4px 8px;
  font-size: 0.8rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.2s;
}

.search-controls button:hover {
  background: #555;
}

.close-btn {
  margin-left: auto;
  font-size: 1rem;
  background: transparent;
  color: #aaa;
  padding: 2px 6px;
}

.match-counter {
  font-size: 0.85rem;
  color: #aaa;
  margin-left: 10px;
}
</style>
