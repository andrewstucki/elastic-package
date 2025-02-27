// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package install

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/elastic/elastic-package/internal/configuration/locations"
	"github.com/elastic/elastic-package/internal/version"
)

func checkIfLatestVersionInstalled(elasticPackagePath string) (bool, error) {
	versionFile, err := ioutil.ReadFile(filepath.Join(elasticPackagePath, versionFilename))
	if os.IsExist(err) {
		return false, nil // old version, no version file
	}
	if err != nil {
		return false, errors.Wrap(err, "reading version file failed")
	}
	v := string(versionFile)
	return buildVersionFile(version.CommitHash, version.BuildTime) == v, nil
}

func writeVersionFile(elasticPackagePath *locations.LocationManager) error {
	var err error
	err = writeStaticResource(err,
		filepath.Join(elasticPackagePath.RootDir(), versionFilename),
		buildVersionFile(version.CommitHash, version.BuildTime))
	if err != nil {
		return errors.Wrap(err, "writing static resource failed")
	}
	return nil
}

func buildVersionFile(commitHash, buildTime string) string {
	return fmt.Sprintf("%s-%s", commitHash, buildTime)
}
