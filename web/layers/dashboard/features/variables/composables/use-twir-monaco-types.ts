import { computed } from 'vue'

import { useSecretsApi } from '~~/layers/dashboard/features/channels-secret/composables/use-secrets-api'
import { useStorageApi } from '~~/layers/dashboard/features/storage/composables/use-storage-api'

export function useTwirMonacoTypes() {
	const secretsApi = useSecretsApi()
	const secrets = secretsApi.secrets

	const storageApi = useStorageApi()
	const storageEntries = storageApi.entries

	const twirTypeDefinitions = computed(() => {
		const secretNames = secrets.value.map((s) => s.name)

		const secretsOverload = secretNames.length > 0
			? secretNames.map((name) => `        get(name: '${name}'): string | null;`).join('\n')
			: '        get(name: string): string | null;'

		const storageKeys = storageEntries.value.map((e) => e.key)

		const storageKeyUnion = storageKeys.length > 0
			? storageKeys.map((k) => `'${k}'`).join(' | ')
			: 'string'

		return `
interface TwirChannel {
    /** Current channel ID */
    id: string;
}

interface TwirSecrets {
${secretsOverload}
}

interface TwirStorage {
    /** Get a value by key */
    get<T = any>(key: ${storageKeyUnion}): Promise<T | null>;
    /** Set a value by key */
    set(key: string, value: any): Promise<void>;
    /** Delete a key */
    delete(key: ${storageKeyUnion}): Promise<boolean>;
    /** Check if a key exists */
    has(key: ${storageKeyUnion}): Promise<boolean>;
    /** List all keys */
    keys(): Promise<string[]>;
    /** Delete all keys */
    clear(): Promise<void>;

    // Array helpers (value at key must be an array)
    /** Append items to array at key */
    push(key: string, ...items: any[]): Promise<number>;
    /** Remove and return last element */
    pop<T = any>(key: string): Promise<T>;
    /** Find first matching element */
    find<T = any>(key: string, predicate: (item: T) => boolean): Promise<T | undefined>;
    /** Filter array elements */
    filter<T = any>(key: string, predicate: (item: T) => boolean): Promise<T[]>;
    /** Array splice */
    splice<T = any>(key: string, start: number, deleteCount: number, ...items: any[]): Promise<T[]>;

    // Object helpers (value at key must be an object, path is dot-notation like "address.city")
    /** Get nested value via dot path */
    getProperty<T = any>(key: string, path: string): Promise<T | undefined>;
    /** Set nested value via dot path */
    setProperty(key: string, path: string, value: any): Promise<void>;
    /** Delete nested property via dot path */
    deleteProperty(key: string, path: string): Promise<boolean>;
    /** Check if nested path exists */
    hasProperty(key: string, path: string): Promise<boolean>;
    /** Shallow merge partial object into existing */
    merge<T extends object = Record<string, any>>(key: string, partial: Partial<T>): Promise<T>;
}

interface Twir {
    secrets: TwirSecrets;
    storage: TwirStorage;
    channel: TwirChannel;
}

declare const twir: Twir;

// Sandbox APIs

interface AbortSignal {
    readonly aborted: boolean;
    onabort: ((event: Event) => void) | null;
    addEventListener(type: 'abort', listener: (event: Event) => void): void;
    removeEventListener(type: 'abort', listener: (event: Event) => void): void;
    throwIfAborted(): void;
}

interface AbortController {
    readonly signal: AbortSignal;
    abort(): void;
}

declare const AbortController: {
    prototype: AbortController;
    new(): AbortController;
};

declare const AbortSignal: {
    prototype: AbortSignal;
    new(): AbortSignal;
};

interface Response {
    readonly status: number;
    readonly statusText: string;
    readonly headers: Record<string, string>;
    readonly ok: boolean;
    readonly body: string;
    json(): any;
    text(): string;
}

interface RequestInit {
    method?: string;
    body?: string;
    headers?: Record<string, string>;
    signal?: AbortSignal;
}

declare function fetch(url: string, options?: RequestInit): Promise<Response>;

type TimerId = number;
declare function setTimeout(callback: () => void, ms?: number): TimerId;
declare function clearTimeout(id: TimerId): void;
declare function setInterval(callback: () => void, ms?: number): TimerId;
declare function clearInterval(id: TimerId): void;

declare class URL {
    constructor(url: string, base?: string);
    protocol: string;
    hostname: string;
    port: string;
    pathname: string;
    search: string;
    hash: string;
    searchParams: URLSearchParams;
    readonly origin: string;
    readonly host: string;
    readonly href: string;
    toString(): string;
}

declare class URLSearchParams {
    constructor(init?: string | string[][] | Record<string, string>);
    append(name: string, value: string): void;
    delete(name: string): void;
    get(name: string): string | null;
    getAll(name: string): string[];
    has(name: string): boolean;
    set(name: string, value: string): void;
    sort(): void;
    entries(): IterableIterator<[string, string]>;
    keys(): IterableIterator<string>;
    values(): IterableIterator<string>;
    [Symbol.iterator](): IterableIterator<[string, string]>;
    toString(): string;
}

declare class TextEncoder {
    readonly encoding: string;
    encode(input?: string): Uint8Array;
    encodeInto(src: string, dest: Uint8Array): { read: number; written: number };
}

declare class TextDecoder {
    constructor(encoding?: string);
    readonly encoding: string;
    decode(input?: ArrayBuffer | ArrayBufferView): string;
}

declare function btoa(data: string): string;
declare function atob(data: string): string;

interface Crypto {
    randomUUID(): string;
}

declare const crypto: Crypto;
`
	})

	return { twirTypeDefinitions }
}
