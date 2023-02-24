package plugin

import (
	"log"

	"github.com/keremdokumaci/goql"
	"github.com/keremdokumaci/goql/pkg/migrations"
)

type Query struct {
	Q string `json:"query"`
}

var (
	gq          = goql.New()
	whitelister goql.WhiteLister
	cacher      goql.Cacher
)

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

	cacher, err = gq.UseGQLCacher()
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}
