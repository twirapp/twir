const pendingTimers = new Map<number, ReturnType<typeof setTimeout>>()
const pendingIntervals = new Map<number, ReturnType<typeof setInterval>>()

export function handleStartTimer(
	id: number,
	delay: number,
	onFire: (id: number) => void,
): void {
	const timer = setTimeout(() => {
		pendingTimers.delete(id)
		onFire(id)
	}, delay)
	pendingTimers.set(id, timer)
}

export function handleClearTimer(id: number): void {
	const timer = pendingTimers.get(id)
	if (timer) {
		clearTimeout(timer)
		pendingTimers.delete(id)
	}
	const interval = pendingIntervals.get(id)
	if (interval) {
		clearInterval(interval)
		pendingIntervals.delete(id)
	}
}

export function handleStartInterval(
	id: number,
	delay: number,
	onFire: (id: number) => void,
): void {
	const interval = setInterval(() => {
		onFire(id)
	}, delay)
	pendingIntervals.set(id, interval)
}

export function clearAllTimers(): void {
	for (const timer of pendingTimers.values()) clearTimeout(timer)
	for (const interval of pendingIntervals.values()) clearInterval(interval)
	pendingTimers.clear()
	pendingIntervals.clear()
}
