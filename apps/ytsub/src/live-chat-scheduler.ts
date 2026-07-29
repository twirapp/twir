export interface RetryScheduler {
	schedule(delay: number, callback: () => void): () => void
}

export const defaultRetryScheduler: RetryScheduler = {
	schedule(delay, callback): () => void {
		const timer = setTimeout(callback, delay)
		return () => clearTimeout(timer)
	},
}
