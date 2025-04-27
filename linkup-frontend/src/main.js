import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

import Vue3Toastify, { toast } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'

import './style.css'

const app = createApp(App)

app.config.globalProperties.$axios = axios

app.use(router)

app.use(Vue3Toastify, {
    autoClose: 3000,
    position: toast.POSITION.TOP_RIGHT,
    pauseOnHover: true,
    draggable: true,
})

app.mount('#app')
