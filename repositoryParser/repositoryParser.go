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

type ParseResult struct {
	Status string;
	Url string;
}

const (
  INVALID_PACKAGE_NAME = "INVALID_NAME";
  INVALID_URL = "INVALID_URL";
  NO_URL_FOUND = "NO_URL_FOUND";
  PACKAGE_NOT_FOUND = "PACKAGE_NOT_FOUND";
  SERVER_ERROR = "SERVER_ERROR";
  SUCCESS = "SUCCESS";
  VERSION_NOT_FOUND = "VERSION_NOT_FOUND";
)

var (
	knownSSLHosts = []string{"bitbucket.org", "github.com", "gitlab.com", "sourceforge.net"};
	logger           = simplelogger.New("npmsource/repositoryParser", true, true);
)

  func GetPackageUrl(rawPackageName string, version string) ParseResult {
    var packageInfo: packageJson.FullMetadata;
    var parsedUrl: string | null = null;
    var parsedRepository: string | null = null;

    const validateResult = validatePackageName(rawPackageName);

    if (!validateResult.validForNewPackages) {
      logger.info(`Invalid package name: "${rawPackageName}"`, validateResult);
      return &ParseResult{status: ParseStatus.INVALID_PACKAGE_NAME};
    }

    try {
      packageInfo = await packageJson(rawPackageName, {fullMetadata: true, version});
    } catch (error) {
      if (error instanceof packageJson.VersionNotFoundError) {
        logger.info(`Version "${version}" not found for package "${rawPackageName}".`);
        return &ParseResult{status: ParseStatus.VERSION_NOT_FOUND};
      }

      if (error instanceof packageJson.PackageNotFoundError) {
        logger.info(`Package "${rawPackageName}" not found.`);
        return &ParseResult{status: ParseStatus.PACKAGE_NOT_FOUND};
      }

      logger.error(error);

      return &ParseResult{status: ParseStatus.SERVER_ERROR};
    }

    if (packageInfo.repository) {
      parsedRepository = parseRepositoryEntry(packageInfo.repository);
    }

    if (parsedRepository) {
      logger.info(`Found repository "${parsedRepository}" for package "${rawPackageName}" (version "${version}").`);
      parsedUrl = parsedRepository;
    } else if (typeof packageInfo?.homepage === "string") {
      logger.info(`Found homepage "${packageInfo.homepage}" for package "${rawPackageName}" (version "${version}").`);
      parsedUrl = packageInfo.homepage;
    } else if (typeof packageInfo.url === "string") {
      logger.info(`Found URL "${packageInfo.url}" for package "${rawPackageName}" (version "${version}").`);
      parsedUrl = packageInfo.url;
    }

    if (!parsedUrl) {
      logger.info(`No source URL found in package "${rawPackageName}".`);
      return &ParseResult{status: ParseStatus.NO_URL_FOUND};
    }

    parsedUrl = parsedUrl.toString().trim().toLowerCase();

    const urlIsValid = validateUrl(parsedUrl);

    if (!urlIsValid) {
      logger.info(`Invalid URL "${parsedUrl}" for package "${rawPackageName}".`);
      return &ParseResult{status: ParseStatus.INVALID_URL};
    }

    parsedUrl = tryHTTPS(parsedUrl);

    return &ParseResult{
      status: ParseStatus.SUCCESS,
      url: parsedUrl,
    };
  }

  func ParseRepositoryEntry(repository string | Record<string, string>) string | null {
    if (typeof repository === "string") {
      return cleanRepositoryUrl(repository);
    }

    if (repository.url) {
      return cleanRepositoryUrl(repository.url);
    }

    return null;
  }

  func cleanRepositoryUrl(repo string) string {
    return repo.replace(/\.git$/, "").replace(/^.*:\/\//, "http://");
  }

  func tryHTTPS(url string) string {
    const parsedURL = new URL(url);
    if (knownSSLHosts.includes(parsedURL.hostname)) {
      parsedURL.protocol = "https:";
    }
    return parsedURL.href;
  }
