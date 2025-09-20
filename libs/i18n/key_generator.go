package i18n

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"regexp"
	"strings"
)

type GenerateKeysOptions struct {
	Locales    LocalesStore
	Package    string
	BaseLocale string
}

// todo: generate supported locales list

var varRegex = regexp.MustCompile(`\{(\w+)\}`)

// snakeToCamel converts snake_case or kebab-case to CamelCase.
func snakeToCamel(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.Title(s)
	result := strings.ReplaceAll(s, " ", "")

	// Handle cases where the result starts with a number (invalid Go identifier)
	if len(result) > 0 && result[0] >= '0' && result[0] <= '9' {
		// Convert numbers to words for the first character
		switch result[0] {
		case '0':
			result = "Zero" + result[1:]
		case '1':
			result = "One" + result[1:]
		case '2':
			result = "Two" + result[1:]
		case '3':
			result = "Three" + result[1:]
		case '4':
			result = "Four" + result[1:]
		case '5':
			result = "Five" + result[1:]
		case '6':
			result = "Six" + result[1:]
		case '7':
			result = "Seven" + result[1:]
		case '8':
			result = "Eight" + result[1:]
		case '9':
			result = "Nine" + result[1:]
		}
	}

	return result
}

func GenerateKeysFileContent(opts GenerateKeysOptions) (string, error) {
	if opts.Package == "" {
		return "", fmt.Errorf("package name is required")
	}
	if len(opts.Locales) == 0 {
		return "", fmt.Errorf("locales are required")
	}
	if opts.BaseLocale == "" {
		return "", fmt.Errorf("base locale is required")
	}

	// Load raw nested structure for key generation
	rawStore, err := LoadRawStore("./apps/parser/locales")
	if err != nil {
		return "", fmt.Errorf("failed to load raw store: %v", err)
	}

	fset := token.NewFileSet()
	file := &ast.File{
		Name: &ast.Ident{Name: opts.Package},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: &ast.Ident{Name: "twiri18n"},
						Path: &ast.BasicLit{Kind: token.STRING, Value: `"github.com/twirapp/twir/libs/i18n"`},
					},
				},
			},
		},
	}

	structName := "Keys"
	baseLocaleData, ok := rawStore[opts.BaseLocale]
	if !ok {
		return "", fmt.Errorf("base locale %s not found in raw store", opts.BaseLocale)
	}

	// Convert raw store to the format expected by generateStructDecls
	nestedData := make(map[string]interface{})
	for categoryKey, categoryData := range baseLocaleData {
		for fileKey, fileData := range categoryData {
			// Create nested structure: categoryKey.fileKey.data
			if nestedData[categoryKey] == nil {
				nestedData[categoryKey] = make(map[string]interface{})
			}
			nestedData[categoryKey].(map[string]interface{})[fileKey] = fileData
		}
	}

	structDecls := generateStructDecls(nestedData, structName, "")
	file.Decls = append(file.Decls, structDecls...)

	file.Decls = append(
		file.Decls,
		&ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{Name: "Translations"}},
					Values: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.Ident{Name: structName},
						},
					},
				},
			},
		},
		&ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{{Name: "Store"}},
					Type: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "twiri18n"},
						Sel: &ast.Ident{Name: "LocalesStore"},
					},
					Values: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "twiri18n"},
								Sel: &ast.Ident{Name: "LocalesStore"},
							},
							Elts: buildLocaleStoreLiteral(opts.Locales),
						},
					},
				},
			},
		},
	)

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, file); err != nil {
		return "", fmt.Errorf("failed to print AST: %v", err)
	}
	return buf.String(), nil
}

func generateStructDecls(data map[string]interface{}, structName, path string) []ast.Decl {
	var decls []ast.Decl
	var fields []*ast.Field

	// Process each key-value pair in the data
	for key, value := range data {
		fieldName := snakeToCamel(key)
		var fieldType ast.Expr
		newPath := key
		if path != "" {
			newPath = path + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// This is a nested object, create a nested struct
			nestedStructName := structName + fieldName
			decls = append(decls, generateStructDecls(v, nestedStructName, newPath)...)
			fieldType = &ast.Ident{Name: nestedStructName}
		case string:
			// This is a leaf node (translation string), create a translation key struct
			nestedStructName := structName + fieldName
			decls = append(decls, generateStringFieldStruct(nestedStructName, newPath, v)...)
			fieldType = &ast.Ident{Name: nestedStructName}
		default:
			fieldType = &ast.Ident{Name: "interface{}"}
		}

		fields = append(
			fields, &ast.Field{
				Names: []*ast.Ident{{Name: fieldName}},
				Type:  fieldType,
			},
		)
	}

	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: structName},
				Type: &ast.StructType{
					Fields: &ast.FieldList{List: fields},
				},
			},
		},
	}
	decls = append(decls, structDecl)

	return decls
}

func generateStringFieldStruct(structName, path, value string) []ast.Decl {
	matches := varRegex.FindAllStringSubmatch(value, -1)
	return generateVarStringFieldStruct(structName, path, matches)
}

func generateVarStringFieldStruct(structName, path string, matches [][]string) []ast.Decl {
	var decls []ast.Decl
	varsStructName := structName + "Vars"

	// 1. Create the Vars struct (e.g., KeysCommandsVipsRemovedVars)
	var varFields []*ast.Field
	if len(matches) > 0 {
		for _, match := range matches {
			varFields = append(
				varFields, &ast.Field{
					Names: []*ast.Ident{{Name: snakeToCamel(match[1])}},
					Type:  &ast.Ident{Name: "any"},
				},
			)
		}
	}

	decls = append(
		decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: varsStructName},
					Type: &ast.StructType{
						Fields: &ast.FieldList{List: varFields},
					},
				},
			},
		},
	)

	// 2. Create the key struct with a `Vars` field
	decls = append(
		decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: structName},
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{{Name: "Vars"}},
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "twiri18n"},
										Sel: &ast.Ident{Name: "Vars"},
									},
								},
							},
						},
					},
				},
			},
		},
	)

	// 3. Add method implementations to satisfy the TranslationKey interface
	decls = append(decls, generateTranslationKeyImpl(structName, path)...)
	decls = append(decls, generateGetVarsImpl(structName))
	decls = append(decls, generateSetVarsImpl(structName, varsStructName, matches))

	return decls
}

func generateSetVarsImpl(structName, varsStructName string, matches [][]string) *ast.FuncDecl {
	// Build the body of the function: k.Vars = twiri18n.Vars{ "key": vars.Key, ... }
	var varElements []ast.Expr
	for _, match := range matches {
		varElements = append(
			varElements, &ast.KeyValueExpr{
				Key: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, match[1]), // e.g., "userName"
				},
				Value: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "vars"},
					Sel: &ast.Ident{Name: snakeToCamel(match[1])}, // e.g., vars.UserName
				},
			},
		)
	}

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{Names: []*ast.Ident{{Name: "k"}}, Type: &ast.Ident{Name: structName}},
			},
		},
		Name: &ast.Ident{Name: "SetVars"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Names: []*ast.Ident{{Name: "vars"}}, Type: &ast.Ident{Name: varsStructName}},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.IndexExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "twiri18n"},
								Sel: &ast.Ident{Name: "TranslationKey"},
							},
							Index: &ast.Ident{Name: varsStructName},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   &ast.Ident{Name: "k"},
							Sel: &ast.Ident{Name: "Vars"},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "twiri18n"},
								Sel: &ast.Ident{Name: "Vars"},
							},
							Elts: varElements,
						},
					},
				},
				&ast.ReturnStmt{Results: []ast.Expr{&ast.Ident{Name: "k"}}},
			},
		},
	}
}

func generateGetVarsImpl(structName string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{Names: []*ast.Ident{{Name: "k"}}, Type: &ast.Ident{Name: structName}},
			},
		},
		Name: &ast.Ident{Name: "GetVars"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "twiri18n"},
							Sel: &ast.Ident{Name: "Vars"},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.SelectorExpr{
							X:   &ast.Ident{Name: "k"},
							Sel: &ast.Ident{Name: "Vars"},
						},
					},
				},
			},
		},
	}
}

func generateTranslationKeyImpl(structName, path string) []ast.Decl {
	var decls []ast.Decl
	// IsTranslationKey method
	decls = append(
		decls, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{{Names: []*ast.Ident{{Name: "k"}}, Type: &ast.Ident{Name: structName}}},
			},
			Name: &ast.Ident{Name: "IsTranslationKey"},
			Type: &ast.FuncType{Params: &ast.FieldList{}},
			Body: &ast.BlockStmt{},
		},
	)

	// GetPath method
	decls = append(
		decls, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{{Names: []*ast.Ident{{Name: "k"}}, Type: &ast.Ident{Name: structName}}},
			},
			Name: &ast.Ident{Name: "GetPath"},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "string"}}}},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s"`, path),
							},
						},
					},
				},
			},
		},
	)

	// GetPathSlice method
	var pathElements []ast.Expr
	if path != "" {
		for _, part := range strings.Split(path, ".") {
			pathElements = append(
				pathElements, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, part)},
			)
		}
	}

	decls = append(
		decls, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{{Names: []*ast.Ident{{Name: "k"}}, Type: &ast.Ident{Name: structName}}},
			},
			Name: &ast.Ident{Name: "GetPathSlice"},
			Type: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}}}}},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}},
								Elts: pathElements,
							},
						},
					},
				},
			},
		},
	)
	return decls
}

func buildLocaleStoreLiteral(locales LocalesStore) []ast.Expr {
	var elements []ast.Expr
	for locale, commands := range locales {
		commandElements := buildMapLiteral(commands)
		elements = append(
			elements, &ast.KeyValueExpr{
				Key: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, locale)},
				Value: &ast.CompositeLit{
					Type: &ast.MapType{
						Key: &ast.Ident{Name: "string"},
						Value: &ast.MapType{
							Key: &ast.Ident{Name: "string"},
							Value: &ast.MapType{
								Key:   &ast.Ident{Name: "string"},
								Value: &ast.Ident{Name: "string"},
							},
						},
					},
					Elts: commandElements,
				},
			},
		)
	}
	return elements
}

func buildMapLiteral(data map[string]map[string]map[string]string) []ast.Expr {
	var elements []ast.Expr
	for key, subMap := range data {
		subElements := buildSubMapLiteral(subMap)
		elements = append(
			elements, &ast.KeyValueExpr{
				Key: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, key)},
				Value: &ast.CompositeLit{
					Type: &ast.MapType{
						Key: &ast.Ident{Name: "string"},
						Value: &ast.MapType{
							Key:   &ast.Ident{Name: "string"},
							Value: &ast.Ident{Name: "string"},
						},
					},
					Elts: subElements,
				},
			},
		)
	}
	return elements
}

func buildSubMapLiteral(data map[string]map[string]string) []ast.Expr {
	var elements []ast.Expr
	for key, subMap := range data {
		leafElements := buildLeafMapLiteral(subMap)
		elements = append(
			elements, &ast.KeyValueExpr{
				Key: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, key)},
				Value: &ast.CompositeLit{
					Type: &ast.MapType{
						Key:   &ast.Ident{Name: "string"},
						Value: &ast.Ident{Name: "string"},
					},
					Elts: leafElements,
				},
			},
		)
	}
	return elements
}

func buildLeafMapLiteral(data map[string]string) []ast.Expr {
	var elements []ast.Expr
	for key, value := range data {
		elements = append(
			elements, &ast.KeyValueExpr{
				Key:   &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, key)},
				Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("`%s`", value)},
			},
		)
	}
	return elements
}
