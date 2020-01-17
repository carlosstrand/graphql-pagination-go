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
    "languages": pagination.Paginated(&pagination.PaginatedField{
      Name: "Languages",
      Type: graphql.String,
      Args: nil,
      DataResolve: func(p graphql.ResolveParams, page pagination.Page) (i interface{}, e error) {
          return []string{"Go", "Javascript", "Ruby"}, nil
      },
      CountResolve: func(p graphql.ResolveParams, page pagination.Page) (i interface{}, e error) {
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

## Skip and Limit

This library already has the `limit` and `skip` arguments ready to be used in a query with the database or external service. See the following example:

```go
  var DataResolver = func(p graphql.ResolveParams, page pagination.Page) (i interface{}, e error) {
      users, err := users.FindMany(db.Filter{
        Limit: page.Limit,
        Skip: page.Skip,
      })
      if err != nil {
        return nil, err
      }
      return users, nil
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
