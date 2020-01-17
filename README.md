[![Documentation](https://godoc.org/github.com/carlosstrand/graphql-pagination-go?status.svg)](http://godoc.org/github.com/carlosstrand/graphql-pagination-go)
[![Actions Status](https://github.com/carlosstrand/graphql-pagination-go/workflows/Go/badge.svg)](https://github.com/carlosstrand/graphql-pagination-go/actions)
[![Coverage Status](https://coveralls.io/repos/github/carlosstrand/graphql-pagination-go/badge.svg?branch=master)](https://coveralls.io/github/carlosstrand/graphql-pagination-go?branch=master)

# graphql-pagination-go

This library makes it easy to create paginated fields for graphql-go. We currently have the following features:

- [x] Simple Pagination (data & count)
- [x] Separated data and count resolvers
- [x] Resolvers are executed only for requested fields


Example:

```go
  fields := graphql.Fields{
    "languages": Paginated(&PaginatedField{
      Name: "Languages",
      Type: graphql.String,
      Args: nil,
      DataResolve: func(p graphql.ResolveParams) (i interface{}, e error) {
          return []string{"Go", "Javascript", "Ruby"}, nil
      },
      CountResolve: func(p graphql.ResolveParams) (i interface{}, e error) {
          return 3, nil
      },
    }),
  }
  rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
  schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
```

Now you can query as below:

```gql
  query {
    languages {
      data
      count
    }
  }
```


## Resolve only requested fields

In some datasources or databases like MongoDB, calling a count comes at an additional cost and is not always used. Thus, this library takes care of resolvers only of the requested fields (data and / or count).

```gql
  query {
    languages {
      count
    }
  }
```

As you can see in example above, only the `CountResolve` you be called and the query will not have the cost of calling DataResolve because `data` was not requested.
