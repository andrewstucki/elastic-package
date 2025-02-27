// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package stack

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/elastic/elastic-package/internal/compose"
	"github.com/elastic/elastic-package/internal/configuration/locations"
	"github.com/elastic/elastic-package/internal/install"
)

const snapshotDefinitionFile = "snapshot.yml"

func dockerComposeBuild(options Options) error {
	locationManager, err := locations.NewLocationManager()
	if err != nil {
		return errors.Wrap(err, "locating stack directory failed")
	}

	c, err := compose.NewProject(DockerComposeProjectName, filepath.Join(locationManager.StackDir(), snapshotDefinitionFile))
	if err != nil {
		return errors.Wrap(err, "could not create docker compose project")
	}

	opts := compose.CommandOptions{
		Services: withIsReadyServices(withDependentServices(options.Services)),
	}

	if err := c.Build(opts); err != nil {
		return errors.Wrap(err, "running command failed")
	}
	return nil
}

func dockerComposePull(options Options) error {
	locationManager, err := locations.NewLocationManager()
	if err != nil {
		return errors.Wrap(err, "locating stack directory failed")
	}

	c, err := compose.NewProject(DockerComposeProjectName, filepath.Join(locationManager.StackDir(), snapshotDefinitionFile))
	if err != nil {
		return errors.Wrap(err, "could not create docker compose project")
	}

	appConfig, err := install.Configuration()
	if err != nil {
		return errors.Wrap(err, "can't read application configuration")
	}

	opts := compose.CommandOptions{
		Env:      appConfig.StackImageRefs(options.StackVersion).AsEnv(),
		Services: withIsReadyServices(withDependentServices(options.Services)),
	}

	if err := c.Pull(opts); err != nil {
		return errors.Wrap(err, "running command failed")
	}
	return nil
}

func dockerComposeUp(options Options) error {
	locationManager, err := locations.NewLocationManager()
	if err != nil {
		return errors.Wrap(err, "locating stack directory failed")
	}

	c, err := compose.NewProject(DockerComposeProjectName, filepath.Join(locationManager.StackDir(), snapshotDefinitionFile))
	if err != nil {
		return errors.Wrap(err, "could not create docker compose project")
	}

	var args []string
	if options.DaemonMode {
		args = append(args, "-d")
	}

	appConfig, err := install.Configuration()
	if err != nil {
		return errors.Wrap(err, "can't read application configuration")
	}

	opts := compose.CommandOptions{
		Env:       appConfig.StackImageRefs(options.StackVersion).AsEnv(),
		ExtraArgs: args,
		Services:  withIsReadyServices(withDependentServices(options.Services)),
	}

	if err := c.Up(opts); err != nil {
		return errors.Wrap(err, "running command failed")
	}
	return nil
}

func dockerComposeDown() error {
	locationManager, err := locations.NewLocationManager()
	if err != nil {
		return errors.Wrap(err, "locating stack directory failed")
	}

	c, err := compose.NewProject(DockerComposeProjectName, filepath.Join(locationManager.StackDir(), snapshotDefinitionFile))
	if err != nil {
		return errors.Wrap(err, "could not create docker compose project")
	}

	if err := c.Down(compose.CommandOptions{}); err != nil {
		return errors.Wrap(err, "running command failed")
	}
	return nil
}

func dockerComposeLogs(serviceName string) ([]byte, error) {
	locationManager, err := locations.NewLocationManager()
	if err != nil {
		return nil, errors.Wrap(err, "locating stack directory failed")
	}

	c, err := compose.NewProject(DockerComposeProjectName, filepath.Join(locationManager.StackDir(), snapshotDefinitionFile))
	if err != nil {
		return nil, errors.Wrap(err, "could not create docker compose project")
	}

	opts := compose.CommandOptions{
		Services: []string{serviceName},
	}

	out, err := c.Logs(opts)
	if err != nil {
		return nil, errors.Wrap(err, "running command failed")
	}
	return out, nil
}

func withDependentServices(services []string) []string {
	for _, aService := range services {
		if aService == "elastic-agent" {
			return []string{} // elastic-agent service requires to load all other services
		}
	}
	return services
}

func withIsReadyServices(services []string) []string {
	if len(services) == 0 {
		return services // load all defined services
	}

	var allServices []string
	for _, aService := range services {
		allServices = append(allServices, aService, fmt.Sprintf("%s_is_ready", aService))
	}
	return allServices
}
