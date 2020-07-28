package pagination

import (
	"github.com/graphql-go/graphql"
)

type Page struct {
	Skip  int
	Limit int
}

type PaginatedField struct {
	Name                string                       `json:"name"`
	Type                graphql.Output               `json:"type"`
	Args                graphql.FieldConfigArgument  `json:"args"`
	DataResolve         PaginatedResolverFn          `json:"-"`
	CountResolve        PaginatedResolverFn          `json:"-"`
	DataAndCountResolve PaginatedDataCountResolverFn `json:"-"`
	DeprecationReason   string                       `json:"deprecationReason"`
	Description         string                       `json:"description"`
}

type PaginatedResult struct {
	Data  interface{} `json:"data"`
	Count int         `json:"count"`
}

type PaginatedResolverFn func(params graphql.ResolveParams, page Page) (interface{}, error)
type PaginatedDataCountResolverFn func(params graphql.ResolveParams, page Page) (interface{}, int, error)

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

	if f.Args == nil {
		f.Args = graphql.FieldConfigArgument{}
	}

	args := f.Args

	args["skip"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		Description:  "Pagination Skip",
		DefaultValue: 0,
	}

	args["limit"] = &graphql.ArgumentConfig{
		Type:         graphql.Int,
		Description:  "Pagination Limit",
		DefaultValue: 10,
	}

	return &graphql.Field{
		Name: f.Name,
		Type: gqlType,
		Args: f.Args,
		Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
			fields := GetSelectedGraphQLQueryFields(p)
			withData := StringSliceContains(fields, "data")
			withCount := StringSliceContains(fields, "count")
			r := &PaginatedResult{}

			page := Page{
				Limit: p.Args["limit"].(int),
				Skip:  p.Args["skip"].(int),
			}

			if f.DataAndCountResolve != nil {
				data, count, err := f.DataAndCountResolve(p, page)
				if err != nil {
					return nil, err
				}
				r.Data = data
				r.Count = count
				return r, nil
			}

			if withData {
				data, err := f.DataResolve(p, page)
				if err != nil {
					return nil, err
				}
				r.Data = data
			}
			if withCount {
				count, err := f.CountResolve(p, page)
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
