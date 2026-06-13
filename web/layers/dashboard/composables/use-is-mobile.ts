import { breakpointsTailwind, createGlobalState, useBreakpoints } from '@vueuse/core'
import { computed } from 'vue'

export const useIsMobile = createGlobalState(() => {
	const breakPoints = useBreakpoints(breakpointsTailwind)
	const isDesktop = breakPoints.greaterOrEqual('md')
	const isMobile = computed(() => !isDesktop.value)

	return {
		isMobile,
		isDesktop,
	}
})
