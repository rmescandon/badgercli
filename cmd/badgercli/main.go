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

package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/rmescandon/badgercli"
)

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Parse the command line arguments and execute the command
	parser := flags.NewParser(&badgercli.Command{}, flags.HelpFlag)
	_, err := parser.Parse()

	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp || e.Type == flags.ErrCommandRequired {
				parser.WriteHelp(os.Stdout)
				return nil
			}
		}
		fmt.Println(err)
	}

	return err
}
