import {createApp} from 'vue'

import MonitororApp from '@/App.vue'
import store from '@/store'

// @ts-ignore
const app = createApp(MonitororApp)
app.use(store)
app.mount('#app-root')
