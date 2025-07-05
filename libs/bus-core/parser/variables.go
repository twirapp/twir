package parser

const GetBuiltInVariablesSubject = "parser.get_build_in_variables"

type BuiltInVariable struct {
	Name                string
	Example             string
	Description         string
	Visible             bool
	CanBeUsedInRegistry bool
	Links               []BuiltInVariableLink
}

type BuiltInVariableLink struct {
	Name string
	Href string
}
