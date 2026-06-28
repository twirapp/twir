import { Agent, expectComplete } from '@isolated-vm/experimental'
import {
	makeDirectResolver,
	makeLinker,
	makeStaticLoader,
} from '@isolated-vm/experimental/utility/linker'

import { hostFetch } from './host-fetch'

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

	try {
		agent = await Agent.create()
		const realm = await agent.createRealm()

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
					resolveResult(String(result ?? ''))
				},
				error(err: unknown) {
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

					void hostFetch(String(url), opts).then(
						async (result) => {
							pendingFetchControllers.delete(numId)
							await globalRef.set('__fetchResultId', id)
							await globalRef.set('__fetchResultData', result)
							await globalRef.set('__fetchResultAborted', false)
							await triggerFetch.run(realm)
						},
						async (err) => {
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
							json() { return JSON.parse(raw.body); },
							text() { return raw.body; },
						};
						globalThis.__fetchCallback(__fetchResultId, response, false);
					}
				}
			`)
		)

		const wrappedCode = `
			import { done, error } from "twir:done";
			import { startFetch, abortFetch } from "twir:fetch";

			const __secrets = ${secretsJson};

			const twir = {
				secrets: {
					get(name) { return __secrets[name] ?? null; }
				},
				channel: { id: ${JSON.stringify(channelId)} }
			};

			class AbortSignal {
				constructor() {
					this.aborted = false;
					this.onabort = null;
					this._listeners = [];
				}

				addEventListener(type, listener) {
					if (type === 'abort') this._listeners.push(listener);
				}

				removeEventListener(type, listener) {
					if (type === 'abort') {
						this._listeners = this._listeners.filter(l => l !== listener);
					}
				}

				_dispatchAbort() {
					this.aborted = true;
					const event = { type: 'abort' };
					if (this.onabort) this.onabort(event);
					for (const listener of this._listeners) listener(event);
				}

				throwIfAborted() {
					if (this.aborted) throw new DOMException('The operation was aborted.', 'AbortError');
				}
			}

			class AbortController {
				constructor() {
					this.signal = new AbortSignal();
				}

				abort() {
					this.signal._dispatchAbort();
				}
			}

			globalThis.AbortController = AbortController;
			globalThis.AbortSignal = AbortSignal;

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

		const module = expectComplete(await agent.compileModule(wrappedCode))

		const linker = makeLinker(
			makeDirectResolver(),
			makeStaticLoader({
				'twir:done': doneCapability,
				'twir:fetch': fetchCapability,
			})
		)

		await module.link(realm, linker)

		const completion = await module.evaluate(realm)

		if (!completion?.complete) {
			return {
				result: '',
				error: String(completion?.error ?? 'Unknown error'),
			}
		}

		const result = await Promise.race([
			resultPromise,
			new Promise<string>((_, reject) =>
				setTimeout(() => reject(new Error('Script execution timed out')), TIMEOUT_MS)
			),
		])

		return { result, error: '' }
	} catch (error: any) {
		return {
			result: '',
			error: error.message || String(error),
		}
	} finally {
		if (agent) {
			await agent[Symbol.asyncDispose]()
		}
	}
}
