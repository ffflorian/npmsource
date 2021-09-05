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
	"github.com/ffflorian/npmsource/routes/mainRoute"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", mainRoute.GetMain)
	router.StaticFile("/robots.txt", "./resources/robots.txt")

	router.Run("localhost:8080")
}
