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

package packagesRoute

import (
	"fmt"
	"net/http"
	URL "net/url"
	"regexp"
	"strings"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/ffflorian/npmsource/repositoryParser"
	"github.com/ffflorian/npmsource/util"
	"github.com/gin-gonic/gin"
)

type PackagesRouteResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

var (
	packageNameRegex = regexp.MustCompile(`^/((?:@[^@/]+/)?[^@/]+)(?:@([^@/]+))?/?$`)
	logger           = simplelogger.New("npmsource/routes/main", true, true)
	unpkgBase        = "https://unpkg.com/browse"
)

func GetPackage(context *gin.Context) {
	if !packageNameRegex.MatchString(context.Request.URL.Path) {
		logger.Logf("Got request \"%s\", doesn't match", context.Request.URL.Path)
		return
	}

	packageParam := fmt.Sprintf("/%s", strings.TrimSpace(context.Param("package")))

	matches := packageNameRegex.FindStringSubmatch(packageParam)
	rawPackageName := matches[1]
	version := matches[2]

	var errorCode int
	var errorMessage string

	logger.Logf("Got request for package \"%s\" (version \"%s\").", rawPackageName, version)

	if util.HasQueryParameter(context, "unpkg") {
		redirectURL := fmt.Sprintf("%s/%s/%s", unpkgBase, rawPackageName, version)

		_, urlParseError := URL.Parse(redirectURL)

		if urlParseError != nil {
			util.ReturnError(context, http.StatusBadRequest, fmt.Sprintf("Invalid URL: %s", redirectURL))
			return
		}

		if util.HasQueryParameter(context, "raw") {
			logger.Logf("Returning raw unpkg info for \"%s\": \"%s\" ...", rawPackageName, redirectURL)
			util.ReturnRedirectURL(context, redirectURL)
			return
		}

		logger.Logf("Redirecting package \"%s\" to unpkg: \"%s\" ...", rawPackageName, redirectURL)
		util.Redirect(context, redirectURL)
		return
	}

	parseResult := repositoryParser.GetPackageURL(rawPackageName, version)

	switch parseResult.Status {
	case repositoryParser.SUCCESS:
		{
			redirectURL := parseResult.URL

			if util.HasQueryParameter(context, "raw") {
				logger.Logf("Returning raw info for \"%s\": \"%s\" ...", rawPackageName, redirectURL)
				util.ReturnRedirectURL(context, redirectURL)
				return
			}

			logger.Logf("Redirecting package \"%s\" to \"%s\" ...", rawPackageName, redirectURL)
			util.Redirect(context, redirectURL)
			return
		}

	case repositoryParser.INVALID_PACKAGE_NAME:
		{
			errorCode = http.StatusUnprocessableEntity
			errorMessage = "Invalid package name"
			break
		}

	case repositoryParser.INVALID_URL:
	case repositoryParser.NO_URL_FOUND:
		{
			errorCode = http.StatusNotFound
			errorMessage = fmt.Sprintf("No source URL found. Please visit https://www.npmjs.com/package/%s.", rawPackageName)
			break
		}

	case repositoryParser.PACKAGE_NOT_FOUND:
		{
			errorCode = http.StatusNotFound
			errorMessage = "Package not found"
			break
		}

	case repositoryParser.VERSION_NOT_FOUND:
		{
			errorCode = http.StatusNotFound
			errorMessage = "Version not found"
			break
		}

	case repositoryParser.SERVER_ERROR:
	default:
		{
			errorCode = http.StatusInternalServerError
			errorMessage = "Internal server error"
			break
		}
	}

	context.IndentedJSON(errorCode, &PackagesRouteResponseBody{Code: errorCode, Message: errorMessage})
}
