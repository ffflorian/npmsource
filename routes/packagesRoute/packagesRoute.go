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
	"regexp"
	"strings"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/ffflorian/npmsource/util"
	"github.com/gin-gonic/gin"
)

type PackagesRouteResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

const (
	repositoryUrl = "https://github.com/ffflorian/npmsource"
	unpkgBase     = "https://unpkg.com/browse"
)

var (
	packageNameRegex = regexp.MustCompile(`^\\/((?:@[^@/]+/)?[^@/]+)(?:@([^@/]+))?\\/?$`)
	logger           = simplelogger.New("npmsource/routes/main", true, true)
)

func GetPackage(context *gin.Context) {
      packageName := strings.TrimSpace(context.Param("package"))
      version: string = "undefined"

      logger.Logf("Got request for package \"%s\" (version \"%s\").", packageName, version);

      if (util.HasQueryParameter(context, "unpkg")) {
        const redirectUrl = `${unpkgBase}/${packageName}@${version}/`;

        if (!validateUrl(redirectUrl)) {
          return response
            .status(HTTP_STATUS.BAD_REQUEST)
            .json({code: HTTP_STATUS.BAD_REQUEST, message: `Invalid URL: ${redirectUrl}`});
        }

        if (util.HasQueryParameter(context, "raw")) {
          logger.Logf("Returning raw unpkg info for \"%s\": \"%s\" ...", packageName, redirectUrl);
          return util.ReturnJSON(&PackagesRouteResponseBody{
            code: HTTP_STATUS.OK,
            url: redirectUrl,
          });
        }

        logger.Logf("Redirecting package \"%s\" to unpkg: \"%s\" ...", packageName, redirectUrl);
        return util.Redirect(redirectUrl);
      }

      const parseResult = "" // await RepositoryParser.getPackageUrl(packageName, version);

      var errorCode int;
      var errorMessage string;

      switch (parseResult) {
        case "SUCCESS": {
          var redirectUrl = "parseResult.url";

          if (util.HasQueryParameter(request, "raw")) {
            logger.Logf("Returning raw info for \"%s\": \"%s\" ...", packageName, redirectUrl);
            return util.ReturnJSON(&PackagesRouteResponseBody{
              code: HTTP_STATUS.OK,
              url: redirectUrl,
            });
          }

          logger.Logf("Redirecting package \"%s\" to \"%s\" ...", packageName, redirectUrl);
          return util.Redirect(redirectUrl);
        }

        case ParseStatus.INVALID_PACKAGE_NAME: {
          errorCode = HTTP_STATUS.UNPROCESSABLE_ENTITY;
          errorMessage = "Invalid package name";
          break;
        }

        case ParseStatus.INVALID_URL:
        case ParseStatus.NO_URL_FOUND: {
          errorCode = HTTP_STATUS.NOT_FOUND;
          errorMessage = `No source URL found. Please visit https://www.npmjs.com/package/${packageName}.`;
          break;
        }

        case ParseStatus.PACKAGE_NOT_FOUND: {
          errorCode = HTTP_STATUS.NOT_FOUND;
          errorMessage = "Package not found";
          break;
        }

        case ParseStatus.VERSION_NOT_FOUND: {
          errorCode = HTTP_STATUS.NOT_FOUND;
          errorMessage = "Version not found";
          break;
        }

        case ParseStatus.SERVER_ERROR:
        default: {
          errorCode = HTTP_STATUS.INTERNAL_SERVER_ERROR;
          errorMessage = "Internal server error";
          break;
        }
      }

      return context.IndentedJSON(errorCode, &PackagesRouteResponseBody{code: errorCode, message: errorMessage});
    }
}
