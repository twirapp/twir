import { describe, expect, test } from 'bun:test'

import { executeCode } from './executor'

describe('executeCode Lodash global', () => {
	test('supports static and chain APIs through _', async () => {
		const execution = await executeCode(
			`
				const chunks = _.chunk([1, 2, 3, 4, 5], 2);
				const chained = _.chain([1, 2, 3, 4])
					.filter((value) => value % 2 === 0)
					.map((value) => value * 10)
					.value();

				return JSON.stringify({ chunks, chained, version: _.VERSION });
			`,
			'test-channel',
			new Map()
		)

		expect(execution.error).toBe('')
		expect(JSON.parse(execution.result)).toEqual({
			chunks: [[1, 2], [3, 4], [5]],
			chained: [20, 40],
			version: '4.17.21',
		})
	})
})
