// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package badgercli

import (
	"encoding/json"
	"fmt"
)

// GetCommand is the command to get an object from a badger database
type GetCommand struct {
	Args struct {
		Path string `positional-arg-name:"path" required:"yes"`
	} `positional-args:"yes"`
	Dir string `long:"dir" short:"d" description:"Directory where database is"`
}

// Execute gets an object from a badger database
func (c *GetCommand) Execute(args []string) error {
	if len(c.Args.Path) == 0 {
		return newEmptyArgument("path")
	}

	if len(c.Dir) == 0 {
		c.Dir = defaultDbPath
	}

	var obj interface{}
	if err := dbExec(c.Dir, get, c.Args.Path, &obj); err != nil {
		return err
	}

	b, err := json.Marshal(&obj)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
