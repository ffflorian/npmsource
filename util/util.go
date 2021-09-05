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

package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HasQuery(context *gin.Context, query string) bool {
	var queryValue, hasQuery = context.GetQuery(query)
	return hasQuery && queryValue != "false"
}

func Redirect(context *gin.Context, url string) {
	context.Redirect(302, url)
	context.Abort()
}

func ReturnJSON(context *gin.Context, data interface{}) {
	context.IndentedJSON(http.StatusOK, data)
	context.Abort()
}
