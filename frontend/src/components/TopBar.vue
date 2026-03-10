<template>
  <div class="topBar" @dblclick="WindowToggleMaximise">
    <div class="titleContainer">
      <ThemeToggle class="topBtn" />
    </div>

    <div class="winControls">
      <button class="winBtn minimize" @click="WindowMinimise" title="Minimize"></button>
      <button class="winBtn maximize" @click="WindowToggleMaximise" title="Maximize"></button>
      <button class="winBtn close" @click="Quit" title="Close"></button>
    </div>
  </div>
</template>

<script setup>
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime';
import ThemeToggle from './ThemeToggle.vue';
</script>

<style scoped>
.topBar {
  height: 50px;
  width: 100%;
  display: flex;
  justify-content: space-between;
  background-color: var(--bg-app);
  align-items: center;
  padding: 0 1rem;
  user-select: none;
  -webkit-app-region: drag;
  --wails-draggable: drag;
}

.titleContainer {
  display: flex;
  align-items: center;
  gap: 8px;
}

.topBtn {
  margin-left: -0.25rem;
}

.winControls {
  display: flex;
  align-items: center;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.winBtn {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 0.5px solid rgba(0, 0, 0, 0.1);
  cursor: default;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease, filter 0.2s ease, box-shadow 0.2s ease;
  box-shadow: inset 0 0 1px rgba(0, 0, 0, 0.15);
}

.minimize {
  background-color: #e0c65f;
}

.maximize {
  background-color: #57b868;
}

.close {
  background-color: #e37c79;
}

.winBtn:hover {
  filter: brightness(1.1);
  box-shadow: inset 0 0 1px rgba(0, 0, 0, 0.25), 0 1px 3px rgba(0, 0, 0, 0.1);
}

.close:hover {
  background-color: #d75f5b;
}
</style>
