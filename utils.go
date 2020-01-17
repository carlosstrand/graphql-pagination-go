package pagination

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// StringSliceContains - Check if a string slice contains a string
func StringSliceContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

// GetSelectedGraphQLQueryFields - retrieve a list of all requested/selected fields from GraphQL Info
func GetSelectedGraphQLQueryFields(p graphql.ResolveParams) []string {
	fieldNames := make([]string, 0)
	fields := p.Info.FieldASTs
	for _, field := range fields {
		selections := field.SelectionSet.Selections
		for _, selection := range selections {
			fieldNames = append(fieldNames, selection.(*ast.Field).Name.Value)
		}
	}
	return fieldNames
}