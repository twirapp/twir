import { h } from 'vue'

import type { SVGProps } from '@tabler/icons-vue'
import type { FunctionalComponent } from 'vue'

interface Opts {
	color?: string
}

export function renderIcon(icon: (props: SVGProps) => FunctionalComponent<SVGProps>, opts?: Opts) {
	return () => h(icon, opts, { default: () => h(icon) })
}
