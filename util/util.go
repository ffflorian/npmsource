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
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

var logger = simplelogger.New("npmsource/util", true, true)

func HasQueryParameter(context *gin.Context, query string) bool {
	queryValue, hasQuery := context.GetQuery(query)
	return hasQuery && queryValue != "false"
}

func Redirect(context *gin.Context, location string) {
	context.Redirect(302, location)
	context.Abort()
}

func ReturnRedirectURL(context *gin.Context, redirectURL string) {
	context.IndentedJSON(http.StatusOK, &ResponseBody{
		Code: http.StatusOK,
		URL:  redirectURL,
	})
	context.Abort()
}

func ReturnError(context *gin.Context, code int, message string) {
	context.IndentedJSON(code, &ResponseBody{
		Message: message,
		Code:    code,
	})
	context.Abort()
}

func readGitRefFile() ([]byte, error) {
	readFile := ".git/refs/heads/main"
	data, err := ioutil.ReadFile(readFile)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteCommitFile() {
	commitHash, readFileErr := readGitRefFile()
	if readFileErr != nil {
		logger.Errorf("Could not read commit file: %s", readFileErr.Error())
	}

	commitFile := "./resources/commit"
	file, createFileErr := os.Create(commitFile)

	if createFileErr != nil {
		logger.Errorf("Could not create commit file: %s", createFileErr.Error())
	}

	_, writeFileErr := file.WriteString(strings.TrimSpace(string(commitHash)))

	if writeFileErr != nil {
		logger.Errorf("Could not write commit file: %s", writeFileErr.Error())
	}

	file.Close()
}
