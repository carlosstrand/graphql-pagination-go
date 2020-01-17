package pagination

import (
	"github.com/graphql-go/graphql"
)

type PaginatedField struct {
	Name              string                      `json:"name"`
	Type              graphql.Output              `json:"type"`
	Args              graphql.FieldConfigArgument `json:"args"`
	DataResolve       graphql.FieldResolveFn      `json:"-"`
	CountResolve      graphql.FieldResolveFn      `json:"-"`
	DeprecationReason string                      `json:"deprecationReason"`
	Description       string                      `json:"description"`
}

type PaginatedResult struct {
	Data interface{} `json:"data"`
	Count int `json:"count"`
}

type PaginatedResolvers struct {
	ResolveData graphql.FieldResolveFn
	ResolveCount graphql.FieldResolveFn
}

func Paginated(f *PaginatedField) *graphql.Field {
	gqlType := graphql.NewObject(graphql.ObjectConfig{
		Name: f.Name,
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(f.Type),
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	return &graphql.Field{
		Name:              f.Name,
		Type:              gqlType,
		Args:              f.Args,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			fields := GetSelectedGraphQLQueryFields(p)
			withData := StringSliceContains(fields, "data")
			withCount := StringSliceContains(fields, "count")
			r := &PaginatedResult{}
			if withData {
				data, err := f.DataResolve(p)
				if err != nil {
					return nil, err
				}
				r.Data = data
			}
			if withCount {
				count, err := f.CountResolve(p)
				if err != nil {
					return nil, err
				}
				r.Count = count.(int)
			}
			return r, nil
		},
		DeprecationReason: f.DeprecationReason,
		Description:       f.Description,
	}
}