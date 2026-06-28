import { lookup } from 'node:dns/promises'

import { Agent, expectComplete } from '@isolated-vm/experimental'
import {
	makeDirectResolver,
	makeLinker,
	makeStaticLoader,
} from '@isolated-vm/experimental/utility/linker'
import { config } from '@twir/config'

const TIMEOUT_MS = 5000
const FETCH_TIMEOUT_MS = 10000
const isProd = config.NODE_ENV === 'production'
const PROXY_URL = 'http://warp-proxy:9091'

const PRIVATE_RANGES: [number, number][] = [
	[0x7F000000, 0x7FFFFFFF],
	[0x0A000000, 0x0AFFFFFF],
	[0xAC100000, 0xAC1FFFFF],
	[0xC0A80000, 0xC0A8FFFF],
	[0xA9FE0000, 0xA9FEFFFF],
	[0x00000000, 0x00FFFFFF],
	[0xE0000000, 0xFFFFFFFF],
	[0x64400000, 0x647FFFFF],
]

const BLOCKED_HOSTNAMES = [
	'localhost',
	'metadata.google.internal',
	'metadata.goog',
]

function ipToLong(ip: string): number {
	return ip.split('.').reduce((acc, octet) => (acc << 8) + Number.parseInt(octet, 10), 0) >>> 0
}

function isPrivateIp(ip: string): boolean {
	const long = ipToLong(ip)
	return PRIVATE_RANGES.some((range) => long >= range[0] && long <= range[1])
}

async function validateUrl(rawUrl: string): Promise<void> {
	const parsed = new URL(rawUrl)

	if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
		throw new Error('Only http and https protocols are allowed')
	}

	const hostname = parsed.hostname

	if (BLOCKED_HOSTNAMES.includes(hostname)) {
		throw new Error(`Blocked hostname: ${hostname}`)
	}

	if (hostname.startsWith('twir_')) {
		throw new Error('Requests to internal services are not allowed')
	}

	if (/^\d+\.\d+\.\d+\.\d+$/.test(hostname)) {
		if (isPrivateIp(hostname)) {
			throw new Error(`Blocked private IP: ${hostname}`)
		}
		return
	}

	try {
		const { address } = await lookup(hostname, { family: 4 })
		if (isPrivateIp(address)) {
			throw new Error(`Blocked private IP: ${address}`)
		}
	} catch (err: any) {
		if (err.message?.startsWith('Blocked')) throw err
	}
}

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
	await validateUrl(url)

	const controller = new AbortController()
	const timeout = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS)

	try {
		const response = await fetch(url, {
			...options,
			signal: controller.signal,
			...(isProd ? { proxy: PROXY_URL } : {}),
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
