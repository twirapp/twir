package types

type VariableHandlerParams struct {
	Key    string
	Params *string
}

type VariableHandlerResult struct {
	Result string
}

type VariableHandler func(data VariableHandlerParams) (*VariableHandlerResult, error)
type Variable struct {
	Name    string
	Handler VariableHandler
}