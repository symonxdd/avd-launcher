import { createRouter, createWebHashHistory } from "vue-router";
import Home from "./views/Home.vue";
import Settings from "./views/Settings.vue";

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/settings', name: 'Settings', component: Settings }
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes
});
