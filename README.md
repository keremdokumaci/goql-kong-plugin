# goql-plugin

A Kong Plugin For GraphQL **Caching** And **Whitelisting** That Uses [GoQL](https://github.com/keremdokumaci/goql).

## Install & Run

If you have already had **docker**, **docker compose** and **make** then run

``` make run ``` or ``` make run_bg ```.

To stop the running containers

``` make stop ``` or ``` make stop_all ```.

## Implementation

This repository has end-to-end implementation of goql-plugin for Kong Gateway.

You can copy & paste the repository if you have not got a Kong Gateway yet. Otherwise, you can copy [goql-plugin](./goql-plugin/) to your gateway repository directly.

🚨 Don't forget to add build steps inside [Dockerfile](./Dockerfile) and [environment variables](./.env)
