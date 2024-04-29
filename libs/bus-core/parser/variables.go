package parser

const GetBuiltInVariablesSubject = "parser.get_build_in_variables"

type BuiltInVariable struct {
	Name                string
	Example             string
	Description         string
	Visible             bool
	CanBeUsedInRegistry bool
}
