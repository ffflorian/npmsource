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

package packageJson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ffflorian/go-tools/simplelogger"
)

type PackageJson struct {
	Homepage   string `json:"Homepage"`
	Repository struct {
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"repository"`
	Url string `json:"url"`
}

const npmRegistryURL = "https://registry.npmjs.org/"

var logger = simplelogger.New("npmsource/packageJson", true, true)

func GetPackageJson(packageName string, version string) (*PackageJson, error) {
	var packageJson *PackageJson

	requestBuffer, requestError := request(packageName)
	if requestError != nil {
		return nil, requestError
	}

	unmarshalError := json.Unmarshal(*requestBuffer, &packageJson)
	if unmarshalError != nil {
		return nil, unmarshalError
	}

	logger.Log("Got packageJson", *packageJson)

	return packageJson, nil
}

func request(urlPath string) (*[]byte, error) {
	httpClient := &http.Client{}
	fullURL := fmt.Sprintf("%s/%s", npmRegistryURL, urlPath)

	logger.Logf("Sending GET request to \"%s\" ...", fullURL)

	request, requestError := http.NewRequest("GET", fullURL, nil)
	if requestError != nil {
		return nil, requestError
	}
	request.Header.Add("Accept", "application/json")

	response, responseError := httpClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}

	defer response.Body.Close()

	logger.Logf("Got response status code \"%d\"", response.StatusCode)

	if response.StatusCode != 200 {
		return nil, errors.New("invalid response status code")
	}

	buffer, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		return nil, readError
	}

	return &buffer, nil
}
