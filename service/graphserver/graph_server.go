package graphserver

import (
	"database/sql"
	"fmt"
	"log"

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

	log.Println("Start GraphQL Server")
	return nil
}

func (gs *GraphServer) Stop() (err error) {
	defer gs.DB.Close()

	log.Println("Stop GraphQL Server")
	return nil
}

func (gs *GraphServer) connectDB() {

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=graphql",
		gs.Environment.Host, gs.Environment.Port, gs.Environment.Username, gs.Environment.Password, gs.Environment.Database)
	db, err := sql.Open(platform.Database, dbinfo)
	if err != nil {
		log.Fatalln(err)
	}
	gs.DB = db
	log.Println("Success to connect DB")
}

func (gs *GraphServer) setAthorObject() {
	authorObject := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Author",
			Description: "An Author",
			Fields: graphql.Fields{
				"id": &graphql.Field{
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
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The name of the author.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if author, ok := p.Source.(*Author); ok {
							return author.Name, nil
						}

						return nil, nil
					},
				},
				"email": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
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
				"id": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The identifier of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*Author); ok {
							return news.AuthorID, nil
						}

						return nil, nil
					},
				},
				"name": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The name of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*Author); ok {
							return news.Name, nil
						}

						return nil, nil
					},
				},
				"email": &graphql.Field{
					Type:        graphql.NewNonNull(graphql.Int),
					Description: "The email of the news.",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						if news, ok := p.Source.(*Author); ok {
							return news.Email, nil
						}

						return nil, nil
					},
				},
			},
		},
	)

	gs.GQLObject.NewsObject = newsObject
}
