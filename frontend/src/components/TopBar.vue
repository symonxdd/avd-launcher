<template>
  <div :class="styles.topBar" @dblclick="WindowToggleMaximise">
    <div :class="styles.titleContainer">
      <ThemeToggle :class="styles.topBtn" />
    </div>

    <div :class="styles.winControls">
      <button :class="[styles.winBtn, styles.minimize]" @click="WindowMinimise" title="Minimize"></button>
      <button :class="[styles.winBtn, styles.maximize]" @click="WindowToggleMaximise" title="Maximize"></button>
      <button :class="[styles.winBtn, styles.close]" @click="Quit" title="Close"></button>
    </div>
  </div>
</template>

<script setup>
import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime';
import ThemeToggle from './ThemeToggle.vue';
import styles from './TopBar.module.css';
</script>
