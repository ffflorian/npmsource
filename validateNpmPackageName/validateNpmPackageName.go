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

package validateNpmPackageName

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ValidationResult struct {
	ValidForNewPackages bool
	ValidForOldPackages bool
	Errors              []string
	Warnings            []string
}

var scopedPackagePattern = regexp.MustCompile(`^(?:@([^/]+?)[/])?([^/]+?)$`)
var blocklist = []string{
	"node_modules",
	"favicon.ico",
}
var builtins = []string{
	"assert",
	"async_hooks",
	"buffer",
	"child_process",
	"cluster",
	"console",
	"constants",
	"crypto",
	"dgram",
	"dns",
	"domain",
	"events",
	"freelist",
	"fs",
	"http",
	"http2",
	"https",
	"inspector",
	"module",
	"net",
	"os",
	"path",
	"perf_hooks",
	"process",
	"punycode",
	"querystring",
	"readline",
	"repl",
	"stream",
	"string_decoder",
	"sys",
	"timers",
	"tls",
	"trace_events",
	"tty",
	"url",
	"util",
	"v8",
	"vm",
	"wasi",
	"worker_threads",
	"zli",
}

func check(name string) ([]string, []string) {
	var warnings = make([]string, 0)
	var errors = make([]string, 0)

	if name == "" {
		errors = append(errors, "name length must be greater than zero")
	}

	if regexp.MustCompile(`^\.`).MatchString(name) {
		errors = append(errors, "name cannot start with a period")
	}

	if regexp.MustCompile(`^_`).MatchString(name) {
		errors = append(errors, "name cannot start with an underscore")
	}

	if strings.TrimSpace(name) != name {
		errors = append(errors, "name cannot contain leading or trailing spaces")
	}

	// No funny business
	for _, blocklistEntry := range blocklist {
		if strings.ToLower(name) == blocklistEntry {
			errors = append(errors, fmt.Sprintf("%s is a blocked name", blocklistEntry))
		}
	}

	// Generate warnings for stuff that used to be allowed

	// core module names like http, events, util, etc
	for _, builtinEntry := range builtins {
		if strings.ToLower(name) == builtinEntry {
			warnings = append(warnings, fmt.Sprintf("%s is a core module name", builtinEntry))
		}
	}

	// really-long-package-names-------------------------------such--length-----many---wow
	// the thisisareallyreallylongpackagenameitshouldpublishdowenowhavealimittothelengthofpackagenames-poch.
	if len(name) > 214 {
		warnings = append(warnings, "name can no longer contain more than 214 characters")
	}

	// mIxeD CaSe nAMEs
	if strings.ToLower(name) != name {
		warnings = append(warnings, "name can no longer contain capital letters")
	}

	if regexp.MustCompile(`[~"!()*]`).MatchString(strings.Split(name, "/")[1]) {
		warnings = append(warnings, "name can no longer contain special characters (\"~'!()*\")")
	}

	if url.QueryEscape(name) != name {
		// Maybe it's a scoped package name, like @user/package
		var nameMatch = scopedPackagePattern.FindStringSubmatch(name)
		if len(nameMatch) != 0 {
			var user = nameMatch[1]
			var pkg = nameMatch[2]
			if url.QueryEscape(user) == user && url.QueryEscape(pkg) == pkg {
				return warnings, errors
			}
		}

		errors = append(errors, "name can only contain URL-friendly characters")
	}

	return warnings, errors
}

func Validate(name string) ValidationResult {
	var warnings, errors = check(name)
	return ValidationResult{
		ValidForNewPackages: len(errors) == 0 && len(warnings) == 0,
		ValidForOldPackages: len(errors) == 0,
		Errors:              errors,
		Warnings:            warnings,
	}
}
