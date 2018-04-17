/*
 *  Copyright (C) 2018 Pierre Marchand <pierre.m@atelier-cartographique.be>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Affero General Public License as published by
 *  the Free Software Foundation, version 3 of the License.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/base64"
)

func base64decode(input string) string {
    content, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        // log.Printf("Error:base64.StdEncoding.DecodeString (%s)", err.Error())
        return input
    }
    return string(content)
}

{{ $Prefix := .Prefix }}
{{ range $Name, $Js :=  .Files }}
var {{ $Prefix }}{{ $Name }} = base64decode(`{{ base64 $Js }}`)
{{ end }}