package parser

import platformentity "github.com/twirapp/twir/libs/entities/platform"

const GetBuiltInVariablesSubject = "parser.get_build_in_variables"

type BuiltInVariable struct {
	Name                string
	Example             string
	Description         string
	Visible             bool
	CanBeUsedInRegistry bool
	Links               []BuiltInVariableLink
	Platforms           []platformentity.Platform
}

type BuiltInVariableLink struct {
	Name string
	Href string
}
