import antfu from '@antfu/eslint-config'

export default antfu({
	typescript: true,
	astro: true,
	yaml: true,
	toml: false,
	jsonc: false,
	markdown: false,
	rules: {
		'curly': 'off',
		'no-unused-vars': 'off',
		'no-var': 'error',
		'unused-imports/no-unused-imports': 'error',
		'style/no-tabs': 'off',
		'style/indent': ['error', 'tab']
	},
	vue: {
		overrides: {
			'vue/component-name-in-template-casing': ['error', 'kebab-case', {
				registeredComponentsOnly: false
			}],
			'vue/block-order': ['error', {
				order: [['script', 'template'], 'style']
			}],
			'vue/multi-word-component-names': [
				'off'
			],
			'vue/max-attributes-per-line': 'off',
			'vue/static-class-names-order': 'off',
			'vue/attribute-hyphenation': 'off',
			'vue/html-self-closing': 'off',
			'vue/html-indent': ['error', 'tab'],
			'vue/no-v-text-v-html-on-component': 'off'
		}
	},
	stylistic: {
		overrides: {
			'node/prefer-global': 'off',
			'node/prefer-global/buffer': 'off',
			'antfu/no-import-dist': 'off',
			'no-console': 'off',
			'style/semi': ['error', 'never'],
			'style/comma-dangle': ['error', 'never'],
			'style/arrow-parens': 'off',
			'style/quotes': [
				'error',
				'single',
				{
					allowTemplateLiterals: true
				}
			],
			'style/brace-style': [
				'error'
			],
			'style/comma-spacing': 'off',
			'style/func-call-spacing': 'off',
			'prefer-const': [
				'error',
				{
					destructuring: 'all',
					ignoreReadBeforeAssign: false
				}
			],
			'import/order': [
				'error',
				{
					'groups': [
						'builtin',
						'external',
						[
							'internal'
						],
						[
							'parent',
							'sibling'
						],
						'index'
					],
					'newlines-between': 'always',
					'alphabetize': {
						order: 'asc',
						caseInsensitive: true
					},
					'pathGroups': [
						{
							pattern: 'src/**',
							group: 'internal',
							position: 'after'
						}
					]
				}
			],
			'import/no-cycle': [
				2,
				{
					maxDepth: 1
				}
			],
			'import/newline-after-import': [
				'error',
				{
					count: 1
				}
			],
			'style/object-curly-spacing': [
				2,
				'always'
			]
		}
	}
})
