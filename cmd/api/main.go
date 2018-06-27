// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 1.0.0
//     Host: localhost:3000
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"github.com/artistomin/friend4me/cmd/api/config"
	"github.com/artistomin/friend4me/cmd/api/mw"
	"github.com/artistomin/friend4me/cmd/api/server"
	"github.com/artistomin/friend4me/cmd/api/service"
	_ "github.com/artistomin/friend4me/cmd/api/swagger"
	"github.com/artistomin/friend4me/internal/account"
	"github.com/artistomin/friend4me/internal/auth"
	"github.com/artistomin/friend4me/internal/platform/postgres"
	"github.com/artistomin/friend4me/internal/rbac"
	"github.com/artistomin/friend4me/internal/user"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

func main() {

	cfg, err := config.Load()
	checkErr(err)

	e := server.New()

	db, err := pgsql.New(cfg.DB)
	checkErr(err)

	addV1Services(cfg, e, db)

	server.Start(e, cfg.Server)
}

func addV1Services(cfg *config.Configuration, e *echo.Echo, db *pg.DB) {

	// Initalize DB interfaces

	userDB := pgsql.NewUserDB(db, e.Logger)
	accDB := pgsql.NewAccountDB(db, e.Logger)

	// Initalize services

	jwt := mw.NewJWT(cfg.JWT)
	authSvc := auth.New(userDB, jwt)
	service.NewAuth(authSvc, e, jwt.MWFunc())

	e.Static("/swaggerui", "cmd/api/swaggerui")

	rbacSvc := rbac.New(userDB)

	v1Router := e.Group("/v1")

	v1Router.Use(jwt.MWFunc())

	// Workaround for Echo's issue with routing.
	// v1Router should be passed to service normally, and then the group name created there
	uR := v1Router.Group("/users")
	service.NewAccount(account.New(accDB, userDB, rbacSvc), uR)
	service.NewUser(user.New(userDB, rbacSvc, authSvc), uR)
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
