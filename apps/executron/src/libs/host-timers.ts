const MIN_INTERVAL_DELAY = 100
const MAX_TIMERS = 20
const MAX_FIRES = 100

export interface TimerContext {
	activeTimers: number
	totalFires: number
	timerIds: Set<number>
	intervalIds: Set<number>
	timerRefs: Map<number, ReturnType<typeof setTimeout>>
	intervalRefs: Map<number, ReturnType<typeof setInterval>>
}

export function createTimerContext(): TimerContext {
	return {
		activeTimers: 0,
		totalFires: 0,
		timerIds: new Set(),
		intervalIds: new Set(),
		timerRefs: new Map(),
		intervalRefs: new Map(),
	}
}

export function handleStartTimer(
	ctx: TimerContext,
	id: number,
	delay: number,
	onFire: (id: number) => void,
): void {
	if (ctx.activeTimers >= MAX_TIMERS) {
		throw new Error(`Timer limit exceeded: max ${MAX_TIMERS} active timers`)
	}
	if (delay < MIN_INTERVAL_DELAY) {
		delay = MIN_INTERVAL_DELAY
	}

	ctx.activeTimers++
	ctx.timerIds.add(id)

	const timer = setTimeout(() => {
		ctx.timerIds.delete(id)
		ctx.timerRefs.delete(id)
		ctx.activeTimers--
		ctx.totalFires++
		if (ctx.totalFires > MAX_FIRES) {
			return
		}
		onFire(id)
	}, delay)
	ctx.timerRefs.set(id, timer)
}

export function handleClearTimer(ctx: TimerContext, id: number): void {
	if (ctx.timerIds.has(id)) {
		ctx.timerIds.delete(id)
		ctx.activeTimers--
		const ref = ctx.timerRefs.get(id)
		if (ref) {
			clearTimeout(ref)
			ctx.timerRefs.delete(id)
		}
	}
	if (ctx.intervalIds.has(id)) {
		ctx.intervalIds.delete(id)
		ctx.activeTimers--
		const ref = ctx.intervalRefs.get(id)
		if (ref) {
			clearInterval(ref)
			ctx.intervalRefs.delete(id)
		}
	}
}

export function handleStartInterval(
	ctx: TimerContext,
	id: number,
	delay: number,
	onFire: (id: number) => void,
): void {
	if (ctx.activeTimers >= MAX_TIMERS) {
		throw new Error(`Timer limit exceeded: max ${MAX_TIMERS} active timers`)
	}
	if (delay < MIN_INTERVAL_DELAY) {
		delay = MIN_INTERVAL_DELAY
	}

	ctx.activeTimers++
	ctx.intervalIds.add(id)

	const interval = setInterval(() => {
		ctx.totalFires++
		if (ctx.totalFires > MAX_FIRES) {
			clearInterval(interval)
			ctx.intervalIds.delete(id)
			ctx.intervalRefs.delete(id)
			ctx.activeTimers--
			return
		}
		onFire(id)
	}, delay)
	ctx.intervalRefs.set(id, interval)
}

export function clearAllTimers(ctx: TimerContext): void {
	for (const ref of ctx.timerRefs.values()) {
		clearTimeout(ref)
	}
	for (const ref of ctx.intervalRefs.values()) {
		clearInterval(ref)
	}
	ctx.activeTimers = 0
	ctx.totalFires = 0
	ctx.timerIds.clear()
	ctx.intervalIds.clear()
	ctx.timerRefs.clear()
	ctx.intervalRefs.clear()
}
