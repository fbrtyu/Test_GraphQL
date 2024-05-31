package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"ozon-test/graph"
	switchdb "ozon-test/switchDB"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	dbuser, err := os.LookupEnv("DBUSER")
	if err {
		fmt.Println(dbuser)
	}
	dbpassword, err := os.LookupEnv("DBPASSWORD")
	if err {
		fmt.Println(dbuser)
	}
	dbname, err := os.LookupEnv("DBNAME")
	if err {
		fmt.Println(dbuser)
	}

	// Считывание вариантов хранения данных
	var store string
	fmt.Print("Storage place (A - in memory / B - in database): ")
	fmt.Scan(&store)
	fmt.Println(store)

	// Вызов функции для смены способа хранения данных
	switchdb.SwitchDB(store)

	if store == "B" {
		connStr := "user=" + dbuser + " password=" + dbpassword + " dbname=" + dbname + " sslmode=disable"
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
