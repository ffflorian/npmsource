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
	"fmt"
	"net/http"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/gin-gonic/gin"
)

type MainRouteResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

const (
	version       = "0.0.1"
	repositoryUrl = "https://github.com/ffflorian/pkgsource"
	unpkgBase     = "https://unpkg.com/browse"
)

var logger = simplelogger.New("npmsource", true, true)

func hasQuery(context *gin.Context, query string) bool {
	return context.Query(query) != "" && context.Query(query) != "false"
}

func getMain(context *gin.Context) {
	logger.Log("Got request for main page")
	logger.Logf("unpkg query %s", context.Query("unpkg"))

	if hasQuery(context, "unpkg") {
		var redirectUrl = fmt.Sprintf("%s/pkgsource@latest", unpkgBase)

		if hasQuery(context, "raw") {
			logger.Logf("Returning raw unpkg info for main page: \"%s\"", redirectUrl)
			context.IndentedJSON(http.StatusOK, &MainRouteResponseBody{
				Code: http.StatusOK,
				Url:  redirectUrl,
			})
			return
		}

		logger.Logf("Redirecting main page to unpkg: \"%s\"", redirectUrl)
		context.Redirect(http.StatusFound, redirectUrl)
		return
	}

	if hasQuery(context, "raw") {
		context.IndentedJSON(http.StatusOK, &MainRouteResponseBody{
			Code: http.StatusOK,
			Url:  repositoryUrl,
		})
		return
	}

	logger.Logf("Redirecting main page to \"%s\"", repositoryUrl)
	context.Redirect(http.StatusFound, repositoryUrl)
}

func main() {
	router := gin.Default()
	router.GET("/", getMain)
	router.StaticFile("/robots.txt", "./resources/robots.txt")

	router.Run("localhost:8080")
}
