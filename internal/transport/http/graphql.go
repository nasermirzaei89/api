package http

import (
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/handler"
	"github.com/graphql-go/relay"
	"github.com/nasermirzaei89/api/internal/services/post"
	"github.com/nasermirzaei89/api/internal/services/user"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func (h *handler) handleGraphQL(pretty, graphiQL, playground bool) http.Handler {
	schema := h.newSchema()

	return gqlhandler.New(&gqlhandler.Config{
		Schema:     &schema,
		Pretty:     pretty,
		GraphiQL:   graphiQL,
		Playground: playground,
	})
}

func (h *handler) newSchema() graphql.Schema {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: graphql.Fields{},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})

	query.AddFieldConfig("health", &graphql.Field{
		Type: graphql.NewNonNull(graphql.Boolean),
		Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
			return true, nil
		},
	})

	var typeUser, typePost *graphql.Object

	nodeDefinitions := relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)

			switch resolvedID.Type {
			case "User":
				return h.userSvc.GetUserByUUID(ctx, resolvedID.ID)
			case "Post":
				return h.postSvc.GetPostByUUID(ctx, resolvedID.ID)
			default:
				return nil, errors.New("unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *user.Entity:
				return typeUser
			case *post.Entity:
				return typePost
			default:
				return nil
			}
		},
	})

	query.AddFieldConfig("node", nodeDefinitions.NodeField)

	typeUser = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("User", func(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
				switch obj := obj.(type) {
				case *user.Entity:
					return obj.UUID, nil
				}
				return "", errors.New("object is not a user")
			}),
			"username": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Source.(*user.Entity).Username, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	typePost = graphql.NewObject(graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Post", func(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
				switch obj := obj.(type) {
				case *post.Entity:
					return obj.UUID, nil
				}
				return "", errors.New("object is not a post")
			}),
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
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	postConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "Post",
		NodeType: typePost,
	})

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

	mutation.AddFieldConfig("updatePostByUUID",
		&graphql.Field{
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
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
					Type: graphql.NewNonNull(graphql.String),
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
					Type: graphql.NewNonNull(graphql.String),
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
			Type: postConnectionDefinition.ConnectionType,
			Args: relay.ConnectionArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userID := p.Context.Value(contextKeyUserUUID)
				if userID == nil {
					return nil, errors.New("unauthorized request")
				}

				args := relay.NewConnectionArguments(p.Args)

				res, err := h.postSvc.ListPosts(p.Context)
				if err != nil {
					return nil, err
				}

				data := make([]interface{}, len(res))
				for i := range res {
					data[i] = res[i]
				}

				return relay.ConnectionFromArray(data, args), nil
			},
		},
	)

	query.AddFieldConfig("listPublishedPosts",
		&graphql.Field{
			Type: postConnectionDefinition.ConnectionType,
			Args: relay.ConnectionArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				args := relay.NewConnectionArguments(p.Args)

				res, err := h.postSvc.ListPublishedPosts(p.Context)
				if err != nil {
					return nil, err
				}

				data := make([]interface{}, len(res))
				for i := range res {
					data[i] = res[i]
				}

				return relay.ConnectionFromArray(data, args), nil
			},
		},
	)

	schemaConfig := graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(errors.Wrap(errors.WithStack(err), "error on new schema"))
	}

	return schema
}
