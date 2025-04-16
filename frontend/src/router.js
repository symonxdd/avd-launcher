import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";
import Logs from "./views/Logs.vue";

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/settings', name: 'Settings', component: Settings },
  { path: '/logs', name: 'Logs', component: Logs }
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes
});
