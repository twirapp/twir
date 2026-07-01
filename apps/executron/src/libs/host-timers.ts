const MIN_INTERVAL_DELAY = 100
const MAX_TIMERS = 20
const MAX_FIRES = 100

export interface TimerContext {
	activeTimers: number
	totalFires: number
	timerIds: Set<number>
	intervalIds: Set<number>
}

export function createTimerContext(): TimerContext {
	return {
		activeTimers: 0,
		totalFires: 0,
		timerIds: new Set(),
		intervalIds: new Set(),
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
		ctx.activeTimers--
		ctx.totalFires++
		if (ctx.totalFires > MAX_FIRES) {
			return
		}
		onFire(id)
	}, delay)
}

export function handleClearTimer(ctx: TimerContext, id: number): void {
	if (ctx.timerIds.has(id)) {
		ctx.timerIds.delete(id)
		ctx.activeTimers--
	}
	if (ctx.intervalIds.has(id)) {
		ctx.intervalIds.delete(id)
		ctx.activeTimers--
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
			ctx.activeTimers--
			return
		}
		onFire(id)
	}, delay)
}

export function clearAllTimers(ctx: TimerContext): void {
	ctx.activeTimers = 0
	ctx.totalFires = 0
	ctx.timerIds.clear()
	ctx.intervalIds.clear()
}
