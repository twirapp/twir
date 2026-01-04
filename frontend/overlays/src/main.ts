import urql from '@urql/vue'
import { createApp } from 'vue'

import { urqlClientOptions } from '@/plugins/urql.ts'

import MainApp from './app.vue'
import './style.css'
import { router } from './plugins/router.js'

const app = createApp(MainApp)

app.use(router).use(urql, urqlClientOptions)

app.mount('#app')

// refresh the page when new version comes
document.body.addEventListener('plugin_web_update_notice', () => {
	window.location.reload()
})
