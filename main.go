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
	"net/http"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/gin-gonic/gin"
)

const version = "0.0.1"

var logger = simplelogger.New("npmsource", false, true)

type MainRouteResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

func getMain(c *gin.Context) {
	logger.Log("Got request for main page")

	response := &MainRouteResponseBody{
		Code: http.StatusOK,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func getVersion(c *gin.Context) {
	logger.Log("Got request for version")

	response := &MainRouteResponseBody{
		Code:    http.StatusOK,
		Version: version,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func main() {
	router := gin.Default()
	router.GET("/", getMain)
	router.GET("/version", getVersion)
	router.StaticFile("/robots.txt", "./resources/robots.txt")

	router.Run("localhost:8080")
}
