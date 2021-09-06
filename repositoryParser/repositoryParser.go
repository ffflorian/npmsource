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
	"strings"

	"github.com/ffflorian/go-tools/simplelogger"
	"github.com/ffflorian/npmsource/packageJson"
	"github.com/ffflorian/npmsource/validateNpmPackageName"
)

type ParseResult struct {
	Status string
	URL    string
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

var logger = simplelogger.New("npmsource/repositoryParser", true, true)

func cleanURL(url string) (*string, error) {
	knownSSLHosts := []string{"bitbucket.org", "github.com", "gitlab.com", "sourceforge.net"}
	url = regexp.MustCompile(`\.git$`).ReplaceAllString(url, "")
	protocol := "http"

	parsedURL, urlParseError := URL.Parse(url)

	if urlParseError != nil {
		return nil, urlParseError
	}

	for _, knownSSLHost := range knownSSLHosts {
		if knownSSLHost == parsedURL.Hostname() {
			protocol = "https"
			break
		}
	}

	cleanURL := fmt.Sprintf("%s://%s/%s", protocol, parsedURL.Hostname(), parsedURL.Path)
	return &cleanURL, nil
}

func GetPackageURL(rawPackageName string, version string) ParseResult {
	validateResult := validateNpmPackageName.Validate(rawPackageName)
	foundURL := ""

	if !validateResult.ValidForNewPackages {
		logger.Logf("Invalid package name: \"%s\" %s", rawPackageName, validateResult)
		return ParseResult{Status: INVALID_PACKAGE_NAME}
	}

	packageInfo, packageInfoError := packageJson.GetPackageJson(rawPackageName, version)

	if packageInfoError != nil {
		if strings.Contains(packageInfoError.Error(), "for package") {
			logger.Logf("Version \"%s\" not found for package \"%s\"", version, rawPackageName)
			return ParseResult{Status: VERSION_NOT_FOUND}
		}

		if strings.Contains(packageInfoError.Error(), "could not be found") {
			logger.Logf("Package \"%s\" not found", rawPackageName)
			return ParseResult{Status: PACKAGE_NOT_FOUND}
		}
	}

	if packageInfo.Repository.URL != "" {
		parsedRepository := packageInfo.Repository.URL
		logger.Logf("Found repository \"parsedRepository\" for package \"%s\" (version \"%s\").", rawPackageName, version)
		foundURL = parsedRepository
	} else if packageInfo.Homepage != "" {
		logger.Logf("Found homepage \"%s\" for package \"%s\" (version \"%s\").", packageInfo.Homepage, version)
		foundURL = packageInfo.Homepage
	} else if packageInfo.URL != "" {
		logger.Logf("Found URL \"%s\" for package \"%s\" (version \"%s\").", packageInfo.URL, version)
		foundURL = packageInfo.URL
	}

	if foundURL != "" {
		logger.Logf("No source URL found in package \"%s\".", rawPackageName)
		return ParseResult{Status: NO_URL_FOUND}
	}

	cleanURL, cleanURLError := cleanURL(foundURL)
	if cleanURLError != nil {
		logger.Logf("Invalid URL \"%s\" for package \"%s\".", foundURL, rawPackageName)
		return ParseResult{Status: INVALID_URL}
	}

	return ParseResult{
		Status: SUCCESS,
		URL:    *cleanURL,
	}
}
