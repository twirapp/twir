import { chmodSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs'
import { tmpdir } from 'node:os'
import { join } from 'node:path'

import { config } from '@twir/config'
import Docker from 'dockerode'

const TIMEOUT_MS = 5000
const MEMORY_LIMIT = 128 * 1024 * 1024
const CPU_LIMIT = 1_000_000_000
const PID_LIMIT = 100
const RUNNER_IMAGE = 'oven/bun:canary'

const docker = new Docker({ socketPath: '/var/run/docker.sock' })

function demuxDockerLogs(buf: Buffer): string {
	const parts: string[] = []
	let offset = 0

	while (offset < buf.length) {
		if (offset + 8 > buf.length) break

		const streamType = buf[offset]
		const size = buf.readUInt32BE(offset + 4)
		offset += 8

		if (offset + size > buf.length) break

		if (streamType === 1 || streamType === 2) {
			parts.push(buf.subarray(offset, offset + size).toString('utf-8'))
		}
		offset += size
	}

	return parts.join('').trim()
}

export interface ExecutionResult {
	result: string
	error: string
}

let imagePulled = false

async function ensureImage() {
	if (imagePulled) return
	try {
		await docker.getImage(RUNNER_IMAGE).inspect()
		imagePulled = true
		return
	} catch {}

	console.log(`[executron] Pulling ${RUNNER_IMAGE}...`)
	const stream = await docker.pull(RUNNER_IMAGE)
	await new Promise<void>((resolve, reject) => {
		docker.modem.followProgress(stream, (err: any) => {
			if (err) reject(err)
			else {
				imagePulled = true
				resolve()
			}
		})
	})
}

export async function executeCode(
	code: string,
	channelId: string,
	secrets: Map<string, string>
): Promise<ExecutionResult> {
	await ensureImage()

	const tmpDir = mkdtempSync(join(tmpdir(), 'executron-'))
	let containerId: string | null = null

	try {
		const secretsObj: Record<string, string> = {}
		for (const [key, value] of secrets) {
			secretsObj[key] = value
		}

		const wrapperContent = `
const __secrets = ${JSON.stringify(secretsObj)};
const twir = {
	secrets: { get(name) { return __secrets[name] ?? null; } },
	channel: { id: ${JSON.stringify(channelId)} }
};

async function __executron_execute() {
${code}
}

try {
	const result = await __executron_execute();
	console.log(JSON.stringify({ result: String(result ?? ''), error: '' }));
} catch (e) {
	console.log(JSON.stringify({ result: '', error: e.message || String(e) }));
}
`

		const wrapperPath = join(tmpDir, 'wrapper.mjs')
		writeFileSync(wrapperPath, wrapperContent)
		chmodSync(tmpDir, 0o755)
		chmodSync(wrapperPath, 0o644)

		const networkMode = config.NODE_ENV !== 'production' ? 'default' : 'container:executron-warp'

		console.log(`[executron] tmpDir: ${tmpDir}, network: ${networkMode}`)

		const container = await docker.createContainer({
			Image: RUNNER_IMAGE,
			Cmd: ['bun', '/code/wrapper.mjs'],
			WorkingDir: '/code',
			Tty: false,
			HostConfig: {
				Mounts: [
					{
						Type: 'bind',
						Source: wrapperPath,
						Target: '/code/wrapper.mjs',
						ReadOnly: true,
						BindOptions: { Propagation: 'rprivate' },
					},
				],
				Memory: MEMORY_LIMIT,
				NanoCpus: CPU_LIMIT,
				PidsLimit: PID_LIMIT,
				NetworkMode: networkMode,
				ReadonlyRootfs: true,
				SecurityOpt: ['no-new-privileges'],
				CapDrop: ['ALL'],
				Tmpfs: { '/tmp': 'size=64M' },
			},
		})

		await container.start()
		containerId = container.id
		console.log(`[executron] Container started: ${containerId}`)

		let timedOut = false
		const waitResult = await Promise.race([
			container.wait(),
			new Promise<never>((_, reject) =>
				setTimeout(() => {
					timedOut = true
					container.stop().catch(() => {})
					reject(new Error('Script execution timed out'))
				}, TIMEOUT_MS)
			),
		])

		if (timedOut) {
			return { result: '', error: 'Script execution timed out' }
		}

		console.log(`[executron] Container exited: ${waitResult.StatusCode}`)

		const logs = await container.logs({ stdout: true, stderr: true })
		const logStr = demuxDockerLogs(logs)
		console.log(`[executron] Logs:`, logStr)

		if (waitResult.StatusCode !== 0) {
			return { result: '', error: logStr || `Container exited with code ${waitResult.StatusCode}` }
		}

		try {
			return JSON.parse(logStr)
		} catch {
			return { result: logStr, error: '' }
		}
	} catch (error: any) {
		console.error(`[executron] Error:`, error.message)
		return { result: '', error: error.message || String(error) }
	} finally {
		if (containerId) {
			try {
				console.log(`[executron] Removing container: ${containerId}`)
				await docker.getContainer(containerId).remove({ force: true })
				console.log(`[executron] Container removed: ${containerId}`)
			} catch (e: any) {
				console.error(`[executron] Failed to remove container:`, e.message)
			}
		} else {
			console.log(`[executron] No container to remove`)
		}
		try {
			rmSync(tmpDir, { recursive: true, force: true })
		} catch {}
	}
}
