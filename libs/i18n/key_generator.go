package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

type GenerateKeysOptions struct {
	Locales    LocalesStore
	Package    string
	BaseLocale string
}

// todo: generate supported locales list

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

	data, err := json.Marshal(opts.Locales)
	if err != nil {
		return "", fmt.Errorf("failed to marshal locales: %v", err)
	}

	var localesData map[string]interface{}
	if err := json.Unmarshal(data, &localesData); err != nil {
		return "", fmt.Errorf("failed to unmarshal locales: %v", err)
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

	baseLocaleData, ok := localesData[opts.BaseLocale]
	if !ok {
		return "", fmt.Errorf("base locale %s not found in locales data", opts.BaseLocale)
	}

	if transMap, ok := baseLocaleData.(map[string]interface{}); ok {
		// Start with empty path to exclude locale prefix
		structDecls := generateStructDecls(transMap, structName, "")
		file.Decls = append(file.Decls, structDecls...)
	}

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

	fields := []*ast.Field{
		{
			Type: &ast.Ident{Name: "twiri18n.TranslationKey"},
		},
	}

	for key, value := range data {
		fieldName := strings.Title(key)
		var fieldType ast.Expr
		newPath := key
		if path != "" {
			newPath = path + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			nestedStructName := structName + fieldName
			decls = append(decls, generateStructDecls(v, nestedStructName, newPath)...)
			fieldType = &ast.Ident{Name: nestedStructName}
		case string:
			nestedStructName := structName + fieldName
			decls = append(decls, generateStringFieldStruct(nestedStructName, newPath)...)
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
	decls = append(decls, generateTranslationKeyImpl(structName, path)...)

	return decls
}

func generateStringFieldStruct(structName, path string) []ast.Decl {
	var decls []ast.Decl

	// Struct definition
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
									Type: &ast.Ident{Name: "twiri18n.TranslationKey"},
								},
								{
									Names: []*ast.Ident{{Name: "Value"}},
									Type:  &ast.Ident{Name: "string"},
								},
							},
						},
					},
				},
			},
		},
	)

	// Add TranslationKey implementation
	decls = append(decls, generateTranslationKeyImpl(structName, path)...)

	return decls
}

func generateTranslationKeyImpl(structName, path string) []ast.Decl {
	var decls []ast.Decl

	// IsTranslationKey method
	decls = append(
		decls, &ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "k"}},
						Type:  &ast.Ident{Name: structName},
					},
				},
			},
			Name: &ast.Ident{Name: "IsTranslationKey"},
			Type: &ast.FuncType{Params: &ast.FieldList{}},
			Body: &ast.BlockStmt{},
		},
	)

	// GetPath method
	decls = append(
		decls,
		&ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "k"}},
						Type:  &ast.Ident{Name: structName},
					},
				},
			},
			Name: &ast.Ident{Name: "GetPath"},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.Ident{Name: "string"},
						},
					},
				},
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
	decls = append(
		decls,
		&ast.FuncDecl{
			Recv: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{{Name: "k"}},
						Type:  &ast.Ident{Name: structName},
					},
				},
			},
			Name: &ast.Ident{Name: "GetPathSlice"},
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.ArrayType{
								Elt: &ast.Ident{Name: "string"},
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.BasicLit{
								Kind: token.STRING,
								// []string{"commands", "followage", "description"}
								Value: fmt.Sprintf(
									`[]string{%s}`, func() string {
										parts := strings.Split(path, ".")
										var quotedParts []string
										for _, part := range parts {
											quotedParts = append(quotedParts, fmt.Sprintf(`"%s"`, part))
										}
										return strings.Join(quotedParts, ", ")
									}(),
								),
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
				Value: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, value)},
			},
		)
	}
	return elements
}
