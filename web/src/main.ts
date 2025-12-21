import './assets/main.css'
import 'flag-icons/css/flag-icons.min.css'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import i18n from './i18n'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(i18n)
app.use(router)

app.mount('#app')
