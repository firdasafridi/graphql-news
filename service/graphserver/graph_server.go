package graphserver

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/lib/pq"

	platform "github.com/firdasafridi/graphql-news"
)

type GraphServer struct {
	Environment *platform.Environment
	DB          *sql.DB
	GQLHandler  *handler.Handler
	GQLObject   *GQLObject
}

type GQLObject struct {
	AuthorObject *graphql.Object
	NewsObject   *graphql.Object
	RootQuery    *graphql.Object
	RootMutation *graphql.Object
}

func NewGraphServer() (gs *GraphServer, err error) {
	gs = &GraphServer{
		Environment: platform.NewEnvironment(),
	}

	if gs.Environment.Debug {
		log.Println(gs)
	}
	return gs, nil
}

func (gs *GraphServer) Start() (err error) {

	gs.connectDB()
	gs.setGraphServer()

	log.Println("Start GraphQL Server")
	return nil
}

func (gs *GraphServer) Stop() (err error) {
	defer gs.DB.Close()

	log.Println("Stop GraphQL Server")
	return nil
}

func (gs *GraphServer) connectDB() {

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=public",
		gs.Environment.Host, gs.Environment.Port, gs.Environment.Username, gs.Environment.Password, gs.Environment.Database)
	log.Println(dbinfo)
	db, err := sql.Open(platform.Database, dbinfo)
	if err != nil {
		log.Fatalln(err)
	}
	gs.DB = db
	log.Println("Success to connect DB")
}

func (gs *GraphServer) setGraphServer() {
	gs.GQLObject = &GQLObject{}
	gs.setAuthorObject()
	gs.setNewsObject()
	gs.setQueryObject()
	gs.setMutationObject()

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    gs.GQLObject.RootQuery,
		Mutation: gs.GQLObject.RootMutation,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	go http.ListenAndServe(":8081", nil)
	log.Println("Success to start graph server")
}

func (gs *GraphServer) setAuthorObject() {
	authorObject := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Author",
			Description: "An Author",
			Fields: graphql.Fields{
				"author_id": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The identifier of the author.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if author, ok := p.Source.(*Author); ok {
							return author.AuthorID, nil
						}

						return nil, nil
					},
				},
				"name": &graphql.Field{
					Type:        graphql.String,
					Description: "The name of the author.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if author, ok := p.Source.(*Author); ok {
							return author.Name, nil
						}

						return nil, nil
					},
				},
				"email": &graphql.Field{
					Type:        graphql.String,
					Description: "The email of the author.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if author, ok := p.Source.(*Author); ok {
							return author.Email, nil
						}

						return nil, nil
					},
				},
			},
		},
	)
	gs.GQLObject.AuthorObject = authorObject
}

func (gs *GraphServer) setNewsObject() {
	newsObject := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "News",
			Description: "A News",
			Fields: graphql.Fields{
				"news_id": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The identifier of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*News); ok {
							return news.NewsID, nil
						}

						return nil, nil
					},
				},
				"title": &graphql.Field{
					Type:        graphql.String,
					Description: "The title of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*News); ok {
							return news.Title, nil
						}

						return nil, nil
					},
				},
				"body": &graphql.Field{
					Type:        graphql.String,
					Description: "The body of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*News); ok {
							return news.Body, nil
						}

						return nil, nil
					},
				},
				"author": &graphql.Field{
					Type:        gs.GQLObject.AuthorObject,
					Description: "The author of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*News); ok {
							author, err := gs.GetAuthorByID(news.AuthorID)
							if err != nil {
								return nil, err
							}
							return author, nil
						}

						return nil, nil
					},
				},
			},
		},
	)

	gs.GQLObject.NewsObject = newsObject
}

func (gs *GraphServer) setQueryObject() {
	queryObject := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"author": &graphql.Field{
					Type:        gs.GQLObject.AuthorObject,
					Description: "get an author by id.",
					Args: graphql.FieldConfigArgument{
						"author_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						id, ok := params.Args["author_id"].(int)
						if !ok {
							return nil, errors.New("author_id must be set")
						}
						author, err := gs.GetAuthorByID(id)
						log.Println(author.Email)
						return author, err
					},
				},
				"authors": &graphql.Field{
					Type:        graphql.NewList(gs.GQLObject.AuthorObject),
					Description: "get all authors.",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						author, err := gs.GetAllAuthor()
						return author, err
					},
				},
				"news": &graphql.Field{
					Type:        gs.GQLObject.NewsObject,
					Description: "get an news by id.",
					Args: graphql.FieldConfigArgument{
						"news_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						id, ok := params.Args["news_id"].(int)
						if !ok {
							return nil, errors.New("news_id must be set")
						}
						news, err := gs.GetNewsByID(id)
						return news, err
					},
				},
				"all_news": &graphql.Field{
					Type:        graphql.NewList(gs.GQLObject.NewsObject),
					Description: "get all news.",
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						news, err := gs.GetAllNews()
						log.Println(news[0])
						return news, err
					},
				},
			},
		},
	)

	gs.GQLObject.RootQuery = queryObject
}

func (gs *GraphServer) setMutationObject() {
	mutationObect := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"insert_author": &graphql.Field{
					Type:        gs.GQLObject.AuthorObject,
					Description: "Insert author",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						name, _ := params.Args["name"].(string)
						email, _ := params.Args["email"].(string)
						newAuthor := &Author{
							Name:  name,
							Email: email,
						}

						lastInsertID, err := gs.InsertAuthor(newAuthor)
						if err != nil {
							return nil, err
						}
						newAuthor.AuthorID = lastInsertID
						return newAuthor, nil
					},
				},
				"update_author": &graphql.Field{
					Type:        gs.GQLObject.AuthorObject,
					Description: "Update author",
					Args: graphql.FieldConfigArgument{
						"author_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"email": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						author_id, ok := params.Args["author_id"].(int)
						if !ok || author_id == 0 {
							return nil, errors.New("author_id can't be nil")
						}
						name, _ := params.Args["name"].(string)
						email, _ := params.Args["email"].(string)
						newAuthor := &Author{
							AuthorID: author_id,
							Name:     name,
							Email:    email,
						}

						_, err := gs.UpdateAuthor(author_id, newAuthor)
						if err != nil {
							return nil, err
						}
						return newAuthor, nil
					},
				},
				"insert_news": &graphql.Field{
					Type:        gs.GQLObject.NewsObject,
					Description: "Insert news",
					Args: graphql.FieldConfigArgument{
						"title": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"body": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"author_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						title, _ := params.Args["title"].(string)
						body, _ := params.Args["body"].(string)
						author_id, _ := params.Args["author_id"].(int)
						newNews := &News{
							Title:    title,
							Body:     body,
							AuthorID: author_id,
						}

						lastInsertID, err := gs.InsertNews(newNews)
						if err != nil {
							return nil, err
						}
						newNews.NewsID = lastInsertID
						return newNews, nil
					},
				},
				"update_news": &graphql.Field{
					Type:        gs.GQLObject.NewsObject,
					Description: "update news",
					Args: graphql.FieldConfigArgument{
						"news_id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
						"title": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"body": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"author_id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						news_id, ok := params.Args["news_id"].(int)
						if !ok || news_id == 0 {
							return nil, errors.New("news_id can't be nil")
						}
						title, _ := params.Args["title"].(string)
						body, _ := params.Args["body"].(string)
						author_id, _ := params.Args["author_id"].(int)
						newNews := &News{
							Title:    title,
							Body:     body,
							AuthorID: author_id,
						}

						_, err := gs.UpdateNews(news_id, newNews)
						if err != nil {
							return nil, err
						}
						return newNews, nil
					},
				},
			},
		},
	)
	gs.GQLObject.RootMutation = mutationObect
}
