package plugin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Kong/go-pdk"
	"github.com/keremdokumaci/goql"
	"github.com/keremdokumaci/goql/pkg/gql/query"
	"github.com/keremdokumaci/goql/pkg/migrations"
	_ "github.com/lib/pq"
)

var (
	gq          = goql.New()
	whitelister goql.WhiteLister
)

const (
	VERSION  = "0.0.1"
	PRIORITY = 1000 // Whatever ur priority is..
)

type Config struct{}

type Query struct {
	Q string `json:"query"`
}

func New() any {
	conf := &Config{}
	return conf
}

func initGoql() error {
	if whitelister != nil {
		return nil
	}

	db, err := connectToPostgres()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	gq.ConfigureCache(goql.INMEMORY).
		ConfigureDB(goql.POSTGRES, db)

	err = migrations.MigratePostgres(db)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	wl, err := gq.UseWhitelister()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	whitelister = wl
	return nil
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
		log.Print("KONG_PG_HOST is required")
	}

	if user = os.Getenv("KONG_PG_USER"); user == "" {
		log.Print("KONG_PG_USER is required")
	}

	if password = os.Getenv("KONG_PG_PASSWORD"); password == "" {
		log.Print("KONG_PG_PASSWORD is required")
	}

	if port, _ = strconv.Atoi(os.Getenv("KONG_PG_PORT")); port == 0 {
		log.Print("KONG_PG_PORT is required")
	}

	if dbname = os.Getenv("KONG_DATABASE"); dbname == "" {
		log.Print("KONG_DATABASE is required")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	return sql.Open("postgres", psqlInfo)
}

func (conf Config) Access(kong *pdk.PDK) {
	defer func() {
		// recover from panic if one occured.
		if err := recover(); err != nil { //catch
			e := fmt.Errorf("[GoQLPlugin].[Access] An unexpected exception occured: %v", err)
			kong.Response.Exit(500, e.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		}
	}()

	_ = kong.Log.Debug("[GoQLPlugin].[Access] Start")

	err := initGoql()
	if err != nil {
		log.Print(err.Error())
	}

	reqRawBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	var q Query
	if err := json.Unmarshal(reqRawBody, &q); err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}
	kong.Log.Debug("[GoQLPlugin].[Access] Incoming Query : " + q.Q)

	query, err := query.Parse(q.Q)
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	operationName := query.OperationName()
	kong.Log.Debug("[GoQLPlugin].[Access] Operation Name : " + operationName)
	kong.Log.Debug(fmt.Sprintf("[GoQLPlugin].[Access] whitelisterke : %v", whitelister))
	allowed := whitelister.OperationAllowed(context.Background(), operationName) // TODO: context???
	if !allowed {
		kong.Response.Exit(403, "Query Not Allowed", map[string][]string{"Content-Type": {"application/json"}}) // TODO: response headers,status etc ...
		return
	}

	_ = kong.Log.Debug("[GoQLPlugin].[Access] Finish")
}
