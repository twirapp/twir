import { computed } from 'vue'

/**
 * Composable to detect if the user is on macOS
 * Uses modern userAgentData API with fallback to userAgent
 */
export function useIsMac() {
	return computed<boolean>(() => {
		if (typeof navigator === 'undefined') {
			return false
		}

		// Modern way: userAgentData
		if ('userAgentData' in navigator) {
			const uaData = (navigator as any).userAgentData
			return uaData?.platform?.toLowerCase().includes('mac') ?? false
		}

		// Fallback: check userAgent
		return /Mac|iPhone|iPad|iPod/.test(navigator.userAgent)
	})
}
