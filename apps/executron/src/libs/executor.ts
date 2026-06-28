import { Agent, expectComplete } from '@isolated-vm/experimental'
import {
	makeDirectResolver,
	makeLinker,
	makeStaticLoader,
} from '@isolated-vm/experimental/utility/linker'

const TIMEOUT_MS = 5000
const FETCH_TIMEOUT_MS = 10000

export interface ExecutionResult {
	result: string
	error: string
}

interface FetchResponse {
	status: number
	statusText: string
	headers: Record<string, string>
	body: string
}

async function hostFetch(url: string, options?: RequestInit): Promise<FetchResponse> {
	const controller = new AbortController()
	const timeout = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS)

	try {
		const response = await fetch(url, {
			...options,
			signal: controller.signal,
		})

		const headers: Record<string, string> = {}
		response.headers.forEach((value, key) => {
			headers[key] = value
		})

		const body = await response.text()

		return {
			status: response.status,
			statusText: response.statusText,
			headers,
			body,
		}
	} finally {
		clearTimeout(timeout)
	}
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

		// Fetch capability — starts a fetch, delivers result via global callback
		const globalRef = await realm.acquireGlobalObject()

		const fetchCapability = await realm.createCapability(
			() => ({
				startFetch(id: unknown, url: unknown, method: unknown, body: unknown) {
					const opts: RequestInit = {}
					if (method) opts.method = String(method)
					if (body) opts.body = String(body)

					void hostFetch(String(url), opts).then(
						async (result) => {
							await globalRef.set('__fetchResultId', id)
							await globalRef.set('__fetchResultData', result)
							await triggerFetch.run(realm)
						},
						async (err) => {
							await globalRef.set('__fetchResultId', id)
							await globalRef.set('__fetchResultData', {
								status: 0,
								statusText: 'Fetch failed',
								headers: {},
								body: String(err?.message ?? err),
							})
							await triggerFetch.run(realm)
						}
					)
				},
			}),
			{ origin: 'twir:fetch' }
		)

		// Trigger script that calls the global callback with a Response-like object
		const triggerFetch = expectComplete(
			await agent.compileScript(`
				if (globalThis.__fetchCallback) {
					const raw = __fetchResultData;
					const response = {
						status: raw.status,
						statusText: raw.statusText,
						headers: raw.headers,
						ok: raw.status >= 200 && raw.status < 300,
						body: raw.body,
						json() { return JSON.parse(raw.body); },
						text() { return raw.body; },
					};
					globalThis.__fetchCallback(__fetchResultId, response);
				}
			`)
		)

		const wrappedCode = `
			import { done, error } from "twir:done";
			import { startFetch } from "twir:fetch";

			const __secrets = ${secretsJson};

			const twir = {
				secrets: {
					get(name) { return __secrets[name] ?? null; }
				},
				channel: { id: ${JSON.stringify(channelId)} }
			};

			let __nextFetchId = 0;
			const __fetchResolvers = new Map();

			globalThis.__fetchCallback = function(id, result) {
				const entry = __fetchResolvers.get(id);
				if (entry) {
					__fetchResolvers.delete(id);
					entry.resolve(result);
				}
			};

			function fetch(url, options) {
				return new Promise((resolve, reject) => {
					const id = __nextFetchId++;
					__fetchResolvers.set(id, { resolve, reject });
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
