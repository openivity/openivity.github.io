import './assets/fontawesome/css/brands.min.css'
import './assets/fontawesome/css/fontawesome.min.css'
import './assets/fontawesome/css/solid.min.css'
import './assets/main.scss'

import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import router from './router'

// Import all of Bootstrap's JS
import 'bootstrap'

// Enable vselect globally
import vSelect from 'vue-select'

const app = createApp(App)

app.component('v-select', vSelect)
app.use(createPinia())
app.use(router)

app.mount('#app')
