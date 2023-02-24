# [WIP] goql-plugin

A Kong Plugin For GraphQL **Caching** And **Whitelisting** That Uses [GoQL](https://github.com/keremdokumaci/goql).

## Install & Run

If you have already had **docker**, **docker compose** and **make** then run

``` make run ``` or ``` make run_bg ```.

To stop the running containers

``` make stop ``` or ``` make stop_all ```.

## How Works

Initialization for GoQL is done in [goql.go](./goql-plugin/plugin/goql.go).

There are 2 global variables called ```whitelister``` and ```cacher``` which will be used in ```Access``` and ```Response```phases.

### Caching

To use GraphQL caching, you have to create ```goql.Cacher``` with using ```UseGQLCacher``` function.

🚨 Cache configuration is **mandatory** to use this feature. So you have to call ```ConfigureCache``` function first.

*There are few different type of cache options by GoQL. You can specify it while configuring cache.*

### Whitelisting

To use Whitelisting, you have to create ```goql.Whitelister``` with using ```UseWhitelister``` function.

🚨 Cache & DB configurations are **mandatory** to use this feature. So you have to call ```ConfigureCache``` and ```ConfigureDB``` functions first.

🚨 Whitelister looks for a table called ```whitelists``` under schema ```goql```. So you have to migrate database. **000001** and **000002** [migrations](./goql-plugin/postgres/) are mandatory. Also migration files can be used to add new whitelist row as seen [00003](./goql-plugin/postgres/000003_add_get_countries_query_to_whitelists.up.sql).

*There are few different type of cache and database options by GoQL. You can specify it while configuring cache and database.*

## Implementation

This repository has end-to-end implementation of goql-plugin for Kong Gateway.

You can copy & paste the repository if you have not got a Kong Gateway yet. Otherwise, you can copy [goql-plugin](./goql-plugin/) to your gateway repository directly.

🚨 Don't forget to add build steps inside [Dockerfile](./Dockerfile) and [environment variables](./.env)
