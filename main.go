/*
Copyright © 2021 Florian Imdahl <git@ffflorian.de>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"github.com/ffflorian/npmsource/routes/errorRoute"
	"github.com/ffflorian/npmsource/routes/mainRoute"
	"github.com/ffflorian/npmsource/routes/packagesRoute"
	"github.com/ffflorian/npmsource/util"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

const version = "0.0.1"

// @title npmsource
// @version 1.0
// @description Find the source of an npm package in an instant

// @contact.name Florian Imdahl
// @contact.url https://ffflorian.de

// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.en.html

// @host npmsource.com
// @BasePath /
func main() {
	util.WriteCommitFile()

	router := gin.New()
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/_swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(errorRoute.InternalError))

	router.GET("/:package/*any", packagesRoute.GetPackage)
	router.GET("/", mainRoute.GetMain)

	router.StaticFile("/robots.txt", "./resources/robots.txt")
	router.StaticFile("/commit", "./resources/commit")
	router.NoRoute(errorRoute.NotFound)

	router.Run(":8080")
}
