import { Agent, expectComplete } from '@isolated-vm/experimental'
import {
	makeDirectResolver,
	makeLinker,
	makeStaticLoader,
} from '@isolated-vm/experimental/utility/linker'

import { hostFetch } from './host-fetch'
import { handleStorageOperation } from './host-storage'
import {
	clearAllTimers,
	handleClearTimer,
	handleStartInterval,
	handleStartTimer,
} from './host-timers'
import { abortControllerPolyfill } from './sandbox/abort-controller'
import { base64Polyfill } from './sandbox/base64'
import { cryptoPolyfill } from './sandbox/crypto'
import { storagePolyfill } from './sandbox/storage'
import { textEncodingPolyfill } from './sandbox/text-encoding'
import { urlPolyfill } from './sandbox/url'

const TIMEOUT_MS = 5000

export interface ExecutionResult {
	result: string
	error: string
}

export async function executeCode(
	code: string,
	channelId: string,
	secrets: Map<string, string>
): Promise<ExecutionResult> {
	let agent: Agent | null = null

	console.log(code)

	try {
		console.log(
			`[executron] Starting execution for channel ${channelId}, code length: ${code.length}`
		)
		const startTime = Date.now()

		agent = await Agent.create()
		console.log(`[executron] Agent created in ${Date.now() - startTime}ms`)

		const realm = await agent.createRealm()
		console.log(`[executron] Realm created in ${Date.now() - startTime}ms`)

		const secretsObj: Record<string, string> = {}
		for (const [key, value] of secrets) {
			secretsObj[key] = value
		}

		const secretsJson = JSON.stringify(secretsObj)

		let resolveResult: (value: string) => void
		let rejectResult: (error: string) => void
		const resultPromise = new Promise<string>((resolve, reject) => {
			resolveResult = resolve
			rejectResult = reject
		})

		const doneCapability = await realm.createCapability(
			() => ({
				done(result: unknown) {
					console.log(
						`[executron] done() called, elapsed: ${Date.now() - startTime}ms, result length: ${String(result ?? '').length}`
					)
					resolveResult(String(result ?? ''))
				},
				error(err: unknown) {
					console.log(
						`[executron] error() called, elapsed: ${Date.now() - startTime}ms, error: ${err}`
					)
					rejectResult(String(err ?? 'Unknown error'))
				},
			}),
			{ origin: 'twir:done' }
		)

		const globalRef = await realm.acquireGlobalObject()

		const pendingFetchControllers = new Map<number, AbortController>()

		const fetchCapability = await realm.createCapability(
			() => ({
				startFetch(id: unknown, url: unknown, method: unknown, body: unknown) {
					const numId = Number(id)
					const controller = new AbortController()
					pendingFetchControllers.set(numId, controller)

					const opts: RequestInit = {}
					if (method) opts.method = String(method)
					if (body) opts.body = String(body)
					opts.signal = controller.signal

					console.log(`[executron] Fetch #${numId}: ${method || 'GET'} ${url}`)
					const fetchStart = Date.now()

					void hostFetch(String(url), opts).then(
						async (result) => {
							console.log(
								`[executron] Fetch #${numId} done: status=${result.status} in ${Date.now() - fetchStart}ms`
							)
							pendingFetchControllers.delete(numId)
							await globalRef.set('__fetchResultId', id)
							await globalRef.set('__fetchResultData', result)
							await globalRef.set('__fetchResultAborted', false)
							await triggerFetch.run(realm)
						},
						async (err) => {
							console.error(`[executron] Fetch #${numId} error:`, err.message)
							pendingFetchControllers.delete(numId)
							const isAbort = controller.signal.aborted
							await globalRef.set('__fetchResultId', id)
							await globalRef.set('__fetchResultData', {
								status: 0,
								statusText: isAbort ? 'Aborted' : 'Fetch failed',
								headers: {},
								body: String(err?.message ?? err),
							})
							await globalRef.set('__fetchResultAborted', isAbort)
							await triggerFetch.run(realm)
						}
					)
				},
				abortFetch(id: unknown) {
					const numId = Number(id)
					const controller = pendingFetchControllers.get(numId)
					if (controller) {
						controller.abort()
						pendingFetchControllers.delete(numId)
					}
				},
			}),
			{ origin: 'twir:fetch' }
		)

		const storageCapability = await realm.createCapability(
			() => ({
				startStorageOp(id: unknown, opJson: unknown) {
					const numId = Number(id)
					const op = JSON.parse(String(opJson))
					console.log(`[executron] Storage op: ${op.action} key=${op.key ?? '-'}`)

					void handleStorageOperation(channelId, op).then(
						async (result) => {
							console.log(`[executron] Storage op done: ${op.action} success=${result.success}`)
							await globalRef.set('__storageResultId', numId)
							await globalRef.set('__storageResultData', result)
							await triggerStorage.run(realm)
						},
						async (err) => {
							console.error(`[executron] Storage op error:`, err.message)
							await globalRef.set('__storageResultId', numId)
							await globalRef.set('__storageResultData', {
								success: false,
								error: err.message || String(err),
							})
							await triggerStorage.run(realm)
						}
					)
				},
			}),
			{ origin: 'twir:storage' }
		)

		const triggerFetch = expectComplete(
			await agent.compileScript(`
				if (globalThis.__fetchCallback) {
					const raw = __fetchResultData;
					if (__fetchResultAborted) {
						globalThis.__fetchCallback(__fetchResultId, null, true);
					} else {
						const response = {
							status: raw.status,
							statusText: raw.statusText,
							headers: raw.headers,
							ok: raw.status >= 200 && raw.status < 300,
							body: raw.body,
							json() { return Promise.resolve().then(() => JSON.parse(raw.body)); },
							text() { return Promise.resolve(raw.body); },
						};
						globalThis.__fetchCallback(__fetchResultId, response, false);
					}
				}
			`)
		)

		const triggerStorage = expectComplete(
			await agent.compileScript(`
				if (globalThis.__storageCallback) {
					globalThis.__storageCallback(__storageResultId, __storageResultData);
				}
			`)
		)

		const timerCapability = await realm.createCapability(
			() => ({
				startTimer(id: unknown, delay: unknown) {
					handleStartTimer(Number(id), Number(delay), (timerId) => {
						void (async () => {
							await globalRef.set('__timerFireId', timerId)
							await triggerTimer.run(realm)
						})()
					})
				},
				clearTimer(id: unknown) {
					handleClearTimer(Number(id))
				},
				startInterval(id: unknown, delay: unknown) {
					handleStartInterval(Number(id), Number(delay), (timerId) => {
						void (async () => {
							await globalRef.set('__timerFireId', timerId)
							await triggerTimer.run(realm)
						})()
					})
				},
			}),
			{ origin: 'twir:timers' }
		)

		const triggerTimer = expectComplete(
			await agent.compileScript(`
				if (globalThis.__timerCallback) {
					globalThis.__timerCallback(__timerFireId);
				}
			`)
		)

		const wrappedCode = `
			import { done, error } from "twir:done";
			import { startFetch, abortFetch } from "twir:fetch";
			import { startTimer, clearTimer, startInterval } from "twir:timers";
			import { startStorageOp } from "twir:storage";

			const __secrets = ${secretsJson};

			const twir = {
				secrets: {
					get(name) { return __secrets[name] ?? null; }
				},
				channel: { id: ${JSON.stringify(channelId)} }
			};

			${urlPolyfill}
			${textEncodingPolyfill}
			${base64Polyfill}
			${cryptoPolyfill}
			${abortControllerPolyfill}
			${storagePolyfill}

			twir.storage = storage;

			const __timerCallbacks = new Map();

			globalThis.__timerCallback = function(id) {
				const fn = __timerCallbacks.get(id);
				if (fn) {
					__timerCallbacks.delete(id);
					try { fn(); } catch(e) {}
				}
			};

			globalThis.setTimeout = function(fn, ms) {
				const id = Math.floor(Math.random() * 2147483647);
				__timerCallbacks.set(id, fn);
				startTimer(id, ms || 0);
				return id;
			};

			globalThis.clearTimeout = function(id) {
				__timerCallbacks.delete(id);
				clearTimer(id);
			};

			globalThis.setInterval = function(fn, ms) {
				const id = Math.floor(Math.random() * 2147483647);
				__timerCallbacks.set(id, fn);
				startInterval(id, ms || 0);
				return id;
			};

			globalThis.clearInterval = function(id) {
				__timerCallbacks.delete(id);
				clearTimer(id);
			};

			let __nextFetchId = 0;
			const __fetchResolvers = new Map();

			globalThis.__fetchCallback = function(id, result, aborted) {
				const entry = __fetchResolvers.get(id);
				if (entry) {
					__fetchResolvers.delete(id);
					if (aborted) {
						entry.reject(new DOMException('The operation was aborted.', 'AbortError'));
					} else {
						entry.resolve(result);
					}
				}
			};

			function fetch(url, options) {
				const signal = options?.signal;
				if (signal?.aborted) {
					return Promise.reject(new DOMException('The operation was aborted.', 'AbortError'));
				}

				return new Promise((resolve, reject) => {
					const id = __nextFetchId++;
					__fetchResolvers.set(id, { resolve, reject });

					if (signal) {
						const onAbort = () => {
							__fetchResolvers.delete(id);
							abortFetch(id);
							reject(new DOMException('The operation was aborted.', 'AbortError'));
						};
						signal.addEventListener('abort', onAbort);
					}

					startFetch(id, url, options?.method || null, options?.body || null);
				});
			}

			globalThis.fetch = fetch;

			void (async () => {
				try {
					const result = await (async () => {
						${code}
					})();
					done(result);
				} catch (e) {
					error(e?.message ?? String(e));
				}
			})();
		`

		console.log(
			`[executron] Capabilities created in ${Date.now() - startTime}ms, compiling module...`
		)

		const module = expectComplete(await agent.compileModule(wrappedCode))
		console.log(`[executron] Module compiled in ${Date.now() - startTime}ms, linking...`)

		const linker = makeLinker(
			makeDirectResolver(),
			makeStaticLoader({
				'twir:done': doneCapability,
				'twir:fetch': fetchCapability,
				'twir:storage': storageCapability,
				'twir:timers': timerCapability,
			})
		)

		await module.link(realm, linker)
		const linkDone = Date.now()
		console.log(
			`[executron] Module linked in ${linkDone - startTime}ms (link step: ${linkDone - startTime}ms), evaluating...`
		)

		const completion = await module.evaluate(realm)
		const evalDone = Date.now()
		console.log(
			`[executron] Module evaluated in ${evalDone - startTime}ms (eval step: ${evalDone - linkDone}ms), complete=${completion?.complete}`
		)

		if (!completion?.complete) {
			console.error(`[executron] Module evaluation failed:`, completion?.error)
			return {
				result: '',
				error: String(completion?.error ?? 'Unknown error'),
			}
		}

		const remaining = TIMEOUT_MS - (Date.now() - startTime)
		console.log(
			`[executron] Waiting for result promise (remaining budget: ${remaining}ms of ${TIMEOUT_MS}ms)...`
		)

		const result = await Promise.race([
			resultPromise,
			new Promise<string>((_, reject) =>
				setTimeout(() => {
					console.error(
						`[executron] TIMEOUT after ${Date.now() - startTime}ms total. done()/error() was never called.`
					)
					reject(new Error('Script execution timed out'))
				}, TIMEOUT_MS)
			),
		])

		console.log(`[executron] Execution completed in ${Date.now() - startTime}ms`)
		return { result, error: '' }
	} catch (error: any) {
		console.error(`[executron] Execution failed:`, error.message)
		return {
			result: '',
			error: error.message || String(error),
		}
	} finally {
		clearAllTimers()
		if (agent) {
			await agent[Symbol.asyncDispose]()
			console.log(`[executron] Agent disposed`)
		}
	}
}
