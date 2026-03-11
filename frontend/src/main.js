import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { router } from "./router";
import { useThemeStore } from './stores/themeStore'
import './styles/theme.css';
import './style.css';
import './styles/flicker.css';
import './styles/native.css';
import './styles/focus.css';

import { addIcons, OhVueIcon } from "oh-vue-icons";
import { HiDeviceMobile, HiDocumentText, HiCog, HiSun, HiMoon, HiAdjustments, HiViewGrid, HiGift, HiPencil, HiTrash, HiPlay, HiStop, HiExclamation, HiDotsHorizontal, HiArrowUp, HiArrowDown, HiInformationCircle, HiCheck, HiX, FaGithub, FaSnowflake } from "oh-vue-icons/icons";

addIcons(HiDeviceMobile, HiDocumentText, HiCog, HiSun, HiMoon, HiAdjustments, HiViewGrid, HiGift, HiPencil, HiTrash, HiPlay, HiStop, HiExclamation, HiDotsHorizontal, HiArrowUp, HiArrowDown, HiInformationCircle, HiCheck, HiX, FaGithub, FaSnowflake);

const app = createApp(App)
app.use(createPinia())
app.use(router)

app.component("v-icon", OhVueIcon);

app.mount('#app')

// Apply theme after mount
const themeStore = useThemeStore();
themeStore.initTheme();