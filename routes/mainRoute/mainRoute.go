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

package mainRoute

import (
	"fmt"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/ffflorian/npmsource/util"
	"github.com/gin-gonic/gin"
)

type MainRouteResponseBody struct {
	Code int    `json:"code"`
	URL  string `json:"url"`
}

const (
	repositoryURL = "https://github.com/ffflorian/npmsource"
	unpkgBase     = "https://unpkg.com/browse"
)

var logger = simplelogger.New("npmsource/routes/main", true, true)

func GetMain(context *gin.Context) {
	logger.Log("Got request for main page")

	if util.HasQueryParameter(context, "unpkg") {
		redirectURL := fmt.Sprintf("%s/npmsource@latest", unpkgBase)

		if util.HasQueryParameter(context, "raw") {
			logger.Logf("Returning raw unpkg info for main page: \"%s\"", redirectURL)
			util.ReturnRedirectURL(context, redirectURL)
			return
		}

		logger.Logf("Redirecting main page to unpkg: \"%s\"", redirectURL)
		util.Redirect(context, redirectURL)
		return
	}

	if util.HasQueryParameter(context, "raw") {
		util.ReturnRedirectURL(context, repositoryURL)
		return
	}

	logger.Logf("Redirecting main page to \"%s\"", repositoryURL)
	util.Redirect(context, repositoryURL)
}
