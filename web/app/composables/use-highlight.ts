import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.min.css'

export const useHighlight = () => {
	const detectLanguage = (code: string) => {
		const result = hljs.highlightAuto(code)
		return result.language || 'plaintext'
	}

	const highlight = (code: string, lang?: string) => {
		if (lang && hljs.getLanguage(lang)) {
			return hljs.highlight(code, { language: lang }).value
		}
		return hljs.highlightAuto(code).value
	}

	return { detectLanguage, highlight }
}
