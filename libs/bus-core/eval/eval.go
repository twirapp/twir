package eval

const (
	EvalEvaluateSubject = "eval.evaluate"
)

type EvalRequest struct {
	Expression string `json:"expression"`
}

type EvalResponse struct {
	Result string `json:"result"`
}
