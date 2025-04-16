<template>
  <div class="logs-container">
    <h2 class="page-title">Logs</h2>
    <textarea class="log-output" readonly :value="store.logs" />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { useAvdStore } from '../stores/avdStore'

const store = useAvdStore()

// Define the handlers so they can be passed to EventsOff
const handleLog = (line) => {
  store.appendLog(line)
}

const handleBooted = (avdName) => {
  store.appendLog(line)
}

const handleShutdown = (avdName) => {
  store.appendLog(line)
}

onMounted(() => {
  EventsOn('avd-log', handleLog)
  EventsOn('avd-booted', handleBooted)
  EventsOn('avd-shutdown', handleShutdown)
})

onBeforeUnmount(() => {
  EventsOff('avd-log', handleLog)
  EventsOff('avd-booted', handleBooted)
  EventsOff('avd-shutdown', handleShutdown)
})
</script>

<style scoped>
.logs-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 0 0 20px;
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
  outline: none;
}
</style>
