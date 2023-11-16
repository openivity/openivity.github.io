// import './assets/bootstrap.scss'
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

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
