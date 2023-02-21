package plugin

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Kong/go-pdk"
	"github.com/keremdokumaci/goql"
	"github.com/keremdokumaci/goql/pkg/gql/query"
	_ "github.com/lib/pq"
)

const (
	VERSION  = "0.0.1"
	PRIORITY = 1000 // Whatever ur priority is..
)

type Config struct {
	Whitelister goql.WhiteLister
}

func New() any {
	cfg := &Config{}

	db, err := connectToPostgres()
	if err != nil {
		log.Panic(err.Error())
	}

	gq := goql.New()
	gq.ConfigureDB(goql.POSTGRES, db)

	wl, err := gq.UseWhitelister()
	if err != nil {
		log.Panic(err.Error())
	}
	cfg.Whitelister = wl

	return cfg
}

func connectToPostgres() (*sql.DB, error) {
	var (
		host     string
		user     string
		password string
		port     int
		dbname   string
	)

	if host = os.Getenv("KONG_PG_HOST"); host == "" {
		log.Panic("KONG_PG_HOST is required")
	}

	if user = os.Getenv("KONG_PG_USER"); user == "" {
		log.Panic("KONG_PG_USER is required")
	}

	if password = os.Getenv("KONG_PG_PASSWORD"); password == "" {
		log.Panic("KONG_PG_PASSWORD is required")
	}

	if port, _ = strconv.Atoi(os.Getenv("KONG_PG_PORT")); port == 0 {
		log.Panic("KONG_PG_PORT is required")
	}

	if dbname = os.Getenv("KONG_DATABASE"); dbname == "" {
		log.Panic("KONG_DATABASE is required")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return sql.Open("postgres", psqlInfo)
}

func (conf Config) Access(kong *pdk.PDK) {
	_ = kong.Log.Debug("[GoQLPlugin].[Access] Start")

	reqRawBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	query, err := query.Parse(string(reqRawBody))
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	allowed := conf.Whitelister.OperationAllowed(context.Background(), query.OperationName()) // TODO: context???
	if !allowed {
		kong.Response.Exit(403, "Query Not Allowed", map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	_ = kong.Log.Debug("[GoQLPlugin].[Access] Finish")
}
