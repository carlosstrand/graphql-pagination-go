package pagination

import (
	"errors"
	"github.com/graphql-go/graphql"
	"log"
	"testing"
)

func setupGraphQL(t * testing.T, fields graphql.Fields) graphql.Schema {
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return schema
}

func assertData(t * testing.T, r *graphql.Result, paginatedName string, expectedData []string) {
	f := r.Data.(map[string]interface{})[paginatedName]
	data := f.(map[string]interface{})["data"].([]interface{})
	if len(expectedData) != len(data) {
		t.Fatalf("data failed: result data length is not equal to expected result data")
	}
	for idx, item := range data {
		if expectedData[idx] != item {
			t.Errorf("data failed: %s not equal %s", item, expectedData[idx])
		}
	}
}

func assertCount(t * testing.T, r *graphql.Result, paginatedName string, expectedCount int) {
	f := r.Data.(map[string]interface{})[paginatedName]
	count := f.(map[string]interface{})["count"].(int)
	if expectedCount != count{
		t.Fatalf("count failed: %d not equal a %d", count, expectedCount)
	}
}

func assertCountIsNil(t * testing.T, r *graphql.Result, paginatedName string) {
	f := r.Data.(map[string]interface{})[paginatedName]
	value := f.(map[string]interface{})["count"]
	if value != nil {
		t.Fatalf("count failed: %d should be nil", value)
	}
}

func assertDataIsNil(t * testing.T, r *graphql.Result, paginatedName string) {
	f := r.Data.(map[string]interface{})[paginatedName]
	value := f.(map[string]interface{})["data"]
	if value != nil {
		t.Fatalf("count failed: %d should be nil", value)
	}
}

func assertError(t * testing.T, r * graphql.Result, errorStr string) {
	found := false
	for _, e := range r.Errors {
		if e.Message == errorStr {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected error \"%s\" dit not throw", errorStr)
	}
}

func createPaginatedLanguagesSchema(t *testing.T, DataResolve, CountResolve interface{}) graphql.Schema {
	f := &PaginatedField{
		Name:              "Languages",
		Type:              graphql.String,
		Args:              nil,
		DataResolve: func(p graphql.ResolveParams, page Page) (i interface{}, e error) {
			return []string{"Go", "Javascript", "Ruby"}, nil
		},
		CountResolve: func(p graphql.ResolveParams, page Page) (i interface{}, e error) {
			return 3, nil
		},
	}
	if DataResolve != nil {
		f.DataResolve = DataResolve.(PaginatedResolverFn)
	}
	if CountResolve != nil {
		f.CountResolve = CountResolve.(PaginatedResolverFn)
	}
	fields := graphql.Fields{
		"languages": Paginated(f),
	}
	return setupGraphQL(t, fields)
}

func TestPaginated(t *testing.T) {
	schema := createPaginatedLanguagesSchema(t, nil, nil)
	query := `
		{
			languages {
				data
				count
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	assertData(t, r, "languages", []string{
		"Go",
		"Javascript",
		"Ruby",
	})
	assertCount(t, r, "languages", 3)
}

func TestPaginatedRequestOnlyData(t *testing.T) {
	schema := createPaginatedLanguagesSchema(t, nil, nil)
	query := `
		{
			languages {
				data
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	assertData(t, r, "languages", []string{
		"Go",
		"Javascript",
		"Ruby",
	})
	assertCountIsNil(t, r, "languages")
}

func TestPaginatedRequestOnlyCount(t *testing.T) {
	schema := createPaginatedLanguagesSchema(t, nil, nil)
	query := `
		{
			languages {
				count
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	assertCount(t, r, "languages", 3)
	assertDataIsNil(t, r, "languages")
}

func TestPaginatedRequestDataError(t *testing.T) {
	var dataResolve PaginatedResolverFn = func(p graphql.ResolveParams, page Page) (i interface{}, e error) {
		return nil, errors.New("data error")
	}
	schema := createPaginatedLanguagesSchema(t, dataResolve, nil)
	query := `
		{
			languages {
				data
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	assertError(t, r, "data error")
}

func TestPaginatedRequestCountError(t *testing.T) {
	var countResolve PaginatedResolverFn = func(p graphql.ResolveParams, page Page) (i interface{}, e error) {
		return nil, errors.New("count error")
	}
	schema := createPaginatedLanguagesSchema(t, nil, countResolve)
	query := `
		{
			languages {
				count
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	assertError(t, r, "count error")
}

