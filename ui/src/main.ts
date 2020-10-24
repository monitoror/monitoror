import {createApp} from 'vue'

import MonitororApp from '@/App.vue'
import store from '@/store'

// Used to avoid custom scrollbar on macOS
const isMacOs = navigator.platform.toLowerCase().includes('mac')
if (isMacOs) {
  document.body.classList.add('macos')
}

// @ts-ignore
const app = createApp(MonitororApp)
app.use(store)
app.mount('#app-root')
