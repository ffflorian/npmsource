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

package errorRoute

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorRouteResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func InternalError(context *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		context.IndentedJSON(http.StatusInternalServerError, &ErrorRouteResponseBody{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error: %s", err),
		})
	}
	context.AbortWithStatus(http.StatusInternalServerError)
}

func NotFound(context *gin.Context) {
	context.IndentedJSON(http.StatusNotFound, &ErrorRouteResponseBody{
		Code:    http.StatusNotFound,
		Message: "Not found",
	})
}
