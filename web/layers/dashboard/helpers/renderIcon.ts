import { h, resolveComponent } from 'vue'

interface Opts {
	color?: string
}

export function renderIcon(icon: string, opts?: Opts) {
	return () => h(resolveComponent('Icon'), { name: icon, ...opts })
}
