package http

import (
	"github.com/graphql-go/graphql"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (h *handler) handleGraphQL() http.HandlerFunc {
	schema := h.newSchema()

	type Request struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if r.Method == http.MethodPost {
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				respond(w, r, badRequest("invalid request body"))
				return
			}
		}

		if r.Method == http.MethodGet {
			req.Query = r.URL.Query().Get("query")
			req.OperationName = r.URL.Query().Get("operationName")

			if variables := r.URL.Query().Get("variables"); variables != "" {
				err := json.Unmarshal([]byte(variables), &req.Variables)
				if err != nil {
					respond(w, r, badRequest("invalid graphql variables"))
					return
				}
			}
		}

		res := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  req.Query,
			RootObject:     nil,
			VariableValues: req.Variables,
			OperationName:  req.OperationName,
			Context:        r.Context(),
		})

		if res.HasErrors() {
			// TODO: manage error type
			respond(w, r, badRequest("error in request", setExtension("errors", res.Errors)))
			return
		}

		respond(w, r, res)
	}
}

func (h *handler) newSchema() graphql.Schema {
	types := make([]graphql.Type, 0)

	query := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})

	typeUser := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*user.Entity).UUID, nil
				},
			},
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*user.Entity).Username, nil
				},
			},
		},
	})

	types = append(types, typeUser)

	typePost := graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*post.Entity).UUID, nil
				},
			},
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*post.Entity).Title, nil
				},
			},
			"slug": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*post.Entity).Slug, nil
				},
			},
			"contentMarkdown": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*post.Entity).ContentMarkdown, nil
				},
			},
			"contentHTML": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*post.Entity).ContentHTML, nil
				},
			},
			"publishedAt": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					publishedAt := p.Source.(*post.Entity).PublishedAt
					if publishedAt != nil {
						return publishedAt.Format(time.RFC3339), nil
					}

					return nil, nil
				},
			},
		},
	})

	types = append(types, typeUser)

	typeLogInRequest := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "LogInRequest",
		Fields: graphql.InputObjectConfigFieldMap{
			"username": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"password": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	types = append(types, typeLogInRequest)

	typeLogInResponse := graphql.NewObject(graphql.ObjectConfig{
		Name: "LogInResponse",
		Fields: graphql.Fields{
			"accessToken": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*user.LogInResponse).AccessToken, nil
				},
			},
			"user": &graphql.Field{
				Type: graphql.NewNonNull(typeUser),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return h.userSvc.GetUserByUUID(p.Context, p.Source.(*user.LogInResponse).UserUUID)
				},
			},
		},
	})

	types = append(types, typeLogInResponse)

	query.AddFieldConfig("me",
		&graphql.Field{
			Type: graphql.NewNonNull(typeUser),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				return h.userSvc.GetUserByUUID(p.Context, userID.(string))
			},
		},
	)

	mutation.AddFieldConfig("logIn",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"request": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(typeLogInRequest),
				},
			},
			Type: graphql.NewNonNull(typeLogInResponse),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				req := p.Args["request"].(map[string]interface{})

				return h.userSvc.LogIn(p.Context, user.LogInRequest{
					Username: req["username"].(string),
					Password: req["password"].(string),
				})
			},
		},
	)

	typeCreatePostRequest := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreatePostRequest",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"slug": &graphql.InputObjectFieldConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"contentMarkdown": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	types = append(types, typeCreatePostRequest)

	mutation.AddFieldConfig("createPost",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"request": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(typeCreatePostRequest),
				},
			},
			Type: graphql.NewNonNull(typePost),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				req := p.Args["request"].(map[string]interface{})

				return h.postSvc.CreatePost(p.Context, post.CreatePostRequest{
					Title:           req["title"].(string),
					Slug:            req["slug"].(string),
					ContentMarkdown: req["contentMarkdown"].(string),
				})
			},
		},
	)

	typeUpdatePostByUUIDRequest := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdatePostByUUIDRequest",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"slug": &graphql.InputObjectFieldConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"contentMarkdown": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	types = append(types, typeUpdatePostByUUIDRequest)

	mutation.AddFieldConfig("updatePostByUUID",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"request": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(typeUpdatePostByUUIDRequest),
				},
			},
			Type: graphql.NewNonNull(typePost),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				req := p.Args["request"].(map[string]interface{})

				return h.postSvc.UpdatePostByUUID(p.Context, p.Args["uuid"].(string), post.UpdatePostByUUIDRequest{
					Title:           req["title"].(string),
					Slug:            req["slug"].(string),
					ContentMarkdown: req["contentMarkdown"].(string),
				})
			},
		},
	)

	mutation.AddFieldConfig("publishPostByUUID",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Type: graphql.NewNonNull(typePost),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				return h.postSvc.PublishPostByUUID(p.Context, p.Args["uuid"].(string))
			},
		},
	)

	query.AddFieldConfig("getPostByUUID",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Type: graphql.NewNonNull(typePost),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				return h.postSvc.GetPostByUUID(p.Context, p.Args["uuid"].(string))
			},
		},
	)

	query.AddFieldConfig("getPublishedPostBySlug",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"slug": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Type: graphql.NewNonNull(typePost),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return h.postSvc.GetPublishedPostBySlug(p.Context, p.Args["slug"].(string))
			},
		},
	)

	query.AddFieldConfig("listPosts",
		&graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(typePost))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				return h.postSvc.ListPosts(p.Context)
			},
		},
	)

	query.AddFieldConfig("listPublishedPosts",
		&graphql.Field{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(typePost))),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return h.postSvc.ListPublishedPosts(p.Context)
			},
		},
	)

	schemaConfig := graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
		Types:    types,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(errors.Wrap(errors.WithStack(err), "error on new schema"))
	}

	return schema
}
