package plugin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/keremdokumaci/goql/pkg/gql/query"
	_ "github.com/lib/pq"
)

const (
	VERSION  = "0.0.1"
	PRIORITY = 1000 // Whatever ur priority is..
)

type Config struct{}

func New() any {
	conf := &Config{}

	err := initGoql()
	if err != nil {
		log.Print(err.Error())
	}

	return conf
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
			kong.Response.Exit(500, e.Error(), map[string][]string{"Content-Type": {"application/json"}})
		}
	}()

	_ = kong.Log.Debug("[GoQLPlugin].[Access] Start")

	reqRawBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	var q Query
	if err := json.Unmarshal(reqRawBody, &q); err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	query, err := query.Parse(q.Q)
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	operationName := query.OperationName()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ctxCancel()

	allowed, err := whitelister.OperationAllowed(ctx, operationName)
	if err != nil {
		kong.Response.Exit(500, err.Error(), map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	if !allowed {
		kong.Response.Exit(403, "Query Not Allowed", map[string][]string{"Content-Type": {"application/json"}})
		return
	}

	res := cacher.GetQueryCache(q.Q)
	if res != nil {
		kong.Response.Exit(200, res.(string), map[string][]string{"Content-Type": {"application/json"}})
	}

	_ = kong.Log.Debug("[GoQLPlugin].[Access] Finish")
}

func (conf Config) Response(kong *pdk.PDK) {
	status, _ := kong.Response.GetStatus()
	if status < 200 || status > 299 {
		return
	}

	reqRawBody, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Info("\n[GoQLPlugin].[Response] error while getting request body\n" + err.Error())
		return
	}

	var q Query
	if err := json.Unmarshal(reqRawBody, &q); err != nil {
		kong.Log.Info("\n[GoQLPlugin].[Response] error while unmarshaling request body\n" + err.Error())
		return
	}

	response, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		kong.Log.Info("\n[GoQLPlugin].[Response] error while getting response body\n" + err.Error())
		return
	}

	err = cacher.CacheQuery(q.Q, response, time.Second*20)
	if err != nil {
		kong.Log.Info("\n[GoQLPlugin].[Response] error while caching query\n" + err.Error())
		return
	}
}
