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
{{ $Prefix := .Prefix }}
{{ range $Name, $Query :=  .Files }}
const {{ $Prefix }}{{ $Name }} = "{{ $Name }}"
{{ end }}

func RegisterQueries(store Store) {
	{{ range $Name, $Query :=  .Files }}
	store.Register({{ $Prefix }}{{ $Name }}, `{{ $Query }}`)
	{{ end }}
}