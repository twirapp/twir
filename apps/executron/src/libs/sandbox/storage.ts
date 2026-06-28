export const storagePolyfill = `
// Storage bridge - communicates with host via capability
let __storageCallback = null;
let __storageResolvers = new Map();
let __nextStorageId = 0;

globalThis.__storageCallback = function(id, result) {
	const entry = __storageResolvers.get(id);
	if (entry) {
		__storageResolvers.delete(id);
		if (result.success) {
			entry.resolve(result.data);
		} else {
			entry.reject(new Error(result.error || 'Storage operation failed'));
		}
	}
};

function __storageOp(action, params) {
	return new Promise((resolve, reject) => {
		const id = __nextStorageId++;
		__storageResolvers.set(id, { resolve, reject });
		startStorageOp(id, JSON.stringify({ action, ...params }));
	});
}

function __getDotPath(obj, path) {
	const parts = path.split('.');
	let current = obj;
	for (const part of parts) {
		if (current == null || typeof current !== 'object') return undefined;
		current = current[part];
	}
	return current;
}

function __setDotPath(obj, path, value) {
	const parts = path.split('.');
	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		if (current[parts[i]] == null || typeof current[parts[i]] !== 'object') {
			current[parts[i]] = {};
		}
		current = current[parts[i]];
	}
	current[parts[parts.length - 1]] = value;
	return obj;
}

function __deleteDotPath(obj, path) {
	const parts = path.split('.');
	let current = obj;
	for (let i = 0; i < parts.length - 1; i++) {
		if (current[parts[i]] == null || typeof current[parts[i]] !== 'object') return false;
		current = current[parts[i]];
	}
	delete current[parts[parts.length - 1]];
	return true;
}

function __hasDotPath(obj, path) {
	const parts = path.split('.');
	let current = obj;
	for (const part of parts) {
		if (current == null || typeof current !== 'object' || !(part in current)) return false;
		current = current[part];
	}
	return true;
}

const storage = {
	async get(key) {
		return __storageOp('get', { key });
	},

	async set(key, value) {
		return __storageOp('set', { key, value });
	},

	async delete(key) {
		return __storageOp('delete', { key });
	},

	async has(key) {
		return __storageOp('has', { key });
	},

	async keys() {
		return __storageOp('keys', {});
	},

	async clear() {
		return __storageOp('clear', {});
	},

	async push(key, ...items) {
		const current = await __storageOp('get', { key });
		const arr = Array.isArray(current) ? current : [];
		arr.push(...items);
		await __storageOp('set', { key, value: arr });
		return arr.length;
	},

	async pop(key) {
		const current = await __storageOp('get', { key });
		if (!Array.isArray(current)) {
			throw new Error('Value at key "' + key + '" is not an array');
		}
		const last = current.pop();
		await __storageOp('set', { key, value: current });
		return last;
	},

	async find(key, predicate) {
		const current = await __storageOp('get', { key });
		if (!Array.isArray(current)) {
			throw new Error('Value at key "' + key + '" is not an array');
		}
		return current.find(predicate);
	},

	async filter(key, predicate) {
		const current = await __storageOp('get', { key });
		if (!Array.isArray(current)) {
			throw new Error('Value at key "' + key + '" is not an array');
		}
		return current.filter(predicate);
	},

	async splice(key, start, deleteCount, ...items) {
		const current = await __storageOp('get', { key });
		if (!Array.isArray(current)) {
			throw new Error('Value at key "' + key + '" is not an array');
		}
		const removed = current.splice(start, deleteCount, ...items);
		await __storageOp('set', { key, value: current });
		return removed;
	},

	async getProperty(key, path) {
		const current = await __storageOp('get', { key });
		if (current == null || typeof current !== 'object') return undefined;
		return __getDotPath(current, path);
	},

	async setProperty(key, path, value) {
		const current = await __storageOp('get', { key });
		const obj = (current != null && typeof current === 'object') ? { ...current } : {};
		__setDotPath(obj, path, value);
		await __storageOp('set', { key, value: obj });
	},

	async deleteProperty(key, path) {
		const current = await __storageOp('get', { key });
		if (current == null || typeof current !== 'object') return false;
		const obj = { ...current };
		const result = __deleteDotPath(obj, path);
		await __storageOp('set', { key, value: obj });
		return result;
	},

	async hasProperty(key, path) {
		const current = await __storageOp('get', { key });
		if (current == null || typeof current !== 'object') return false;
		return __hasDotPath(current, path);
	},

	async merge(key, partial) {
		const current = await __storageOp('get', { key });
		const obj = (current != null && typeof current === 'object') ? { ...current } : {};
		Object.assign(obj, partial);
		await __storageOp('set', { key, value: obj });
		return obj;
	},
};

globalThis.storage = storage;
`
