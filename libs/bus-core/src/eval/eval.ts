export interface EvalRequest {
	expression: string
}

export interface EvalResponse {
	result: string
}

export const evalSubject = 'eval.evaluate'
