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

package repositoryParser

import (
	"fmt"
	URL "net/url"
	"regexp"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/ffflorian/npmsource/packageJson"
	"github.com/ffflorian/npmsource/validateNpmPackageName"
)

type ParseResult struct {
	Status string
	Url    string
}

const (
	INVALID_PACKAGE_NAME = "INVALID_NAME"
	INVALID_URL          = "INVALID_URL"
	NO_URL_FOUND         = "NO_URL_FOUND"
	PACKAGE_NOT_FOUND    = "PACKAGE_NOT_FOUND"
	SERVER_ERROR         = "SERVER_ERROR"
	SUCCESS              = "SUCCESS"
	VERSION_NOT_FOUND    = "VERSION_NOT_FOUND"
)

var (
	knownSSLHosts = []string{"bitbucket.org", "github.com", "gitlab.com", "sourceforge.net"}
	logger        = simplelogger.New("npmsource/repositoryParser", true, true)
)

func cleanUrl(url string) (*string, error) {
	url = regexp.MustCompile(`\.git$`).ReplaceAllString(url, "")

	var parsedURL, urlParseError = URL.Parse(url)

	if urlParseError != nil {
		return nil, urlParseError
	}

	var protocol = "http:"

	for _, knownSSLHost := range knownSSLHosts {
		if knownSSLHost == parsedURL.Hostname() {
			protocol = "https:"
		}
	}

	var cleanURL = fmt.Sprintf("%s//%s/%s", protocol, parsedURL.Hostname(), parsedURL.Path)
	return &cleanURL, nil
}

func GetPackageUrl(rawPackageName string, version string) ParseResult {
	validateResult := validateNpmPackageName.Validate(rawPackageName)

	if !validateResult.ValidForNewPackages {
		logger.Logf("Invalid package name: \"%s\" %s", rawPackageName, validateResult)
		return ParseResult{Status: INVALID_PACKAGE_NAME}
	}

	packageInfo, packageInfoError = packageJson.GetPackageJson(rawPackageName, version)
}
