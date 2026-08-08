import urql from '@urql/vue'
import { createApp } from 'vue'

import { urqlClientOptions } from '@/plugins/urql.ts'

import MainApp from './app.vue'
import { router } from './plugins/router.js'
import { loadEruda } from './helpers.js'
import './style.css'

const app = createApp(MainApp)

// Keep the overlay alive in OBS: a single faulty message render must not kill app-wide reactivity.
app.config.errorHandler = (err, _instance, info) => {
	console.error('[overlays] unhandled error', err, info)
}

app.use(router).use(urql, urqlClientOptions)

app.mount('#app')

// eruda devtools
loadEruda()

// refresh the page when new version comes
document.body.addEventListener('plugin_web_update_notice', () => {
	window.location.reload()
})
