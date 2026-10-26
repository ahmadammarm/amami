import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { VueQueryPlugin } from '@tanstack/vue-query';
import router from './router';
// vue-sonner v2 no longer injects its CSS; without it toasts render in the page flow.
// Imported before style.css so the app's theme can override it.
import 'vue-sonner/style.css';
import './style.css';
import App from './App.vue';
import VueVirtualScroller from 'vue-virtual-scroller';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(VueQueryPlugin);
app.use(VueVirtualScroller);

app.mount('#app');
