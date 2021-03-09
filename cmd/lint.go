// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/elastic/elastic-package/internal/docs"
	"github.com/elastic/elastic-package/internal/packages"
	"github.com/elastic/package-spec/code/go/pkg/validator"
)

const lintLongDescription = `Use this command to validate the contents of a package using the package specification (see: https://github.com/elastic/package-spec).

The command ensures that the package is aligned with the package spec and the README file is up-to-date with its template (if present).

Context:
  package`

func setupLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Lint the package",
		Long:  lintLongDescription,
		RunE:  lintCommandAction,
	}
	return cmd
}

func lintCommandAction(cmd *cobra.Command, args []string) error {
	cmd.Println("Lint the package")

	packageRootPath, found, err := packages.FindPackageRoot()
	if !found {
		return errors.New("package root not found")
	}
	if err != nil {
		return errors.Wrap(err, "locating package root failed")
	}

	readmeFiles, err := ioutil.ReadDir(filepath.Join(packageRootPath, "_dev", "build", "docs"))
	if err != nil {
		return errors.Wrapf(err, "failed to return a list of directory entries from %s", packageRootPath)
	}

	for _, readme := range readmeFiles {
		filename := readme.Name()
		ok, err := docs.IsReadmeUpToDate(filename)
		if err != nil {
			return errors.Wrapf(err, "can't check if %s file is up-to-date", filename)
		}
		if !ok {
			return fmt.Errorf("%s file is outdated. Rebuild the package with 'elastic-package build'", filename)
		}
	}

	err = validator.ValidateFromPath(packageRootPath)
	if err != nil {
		return errors.Wrap(err, "linting package failed")
	}

	cmd.Println("Done")
	return nil
}
