// Copyright 2024 Nitro Agility S.r.l.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package servers

import (
	"flag"

	"github.com/spf13/viper"

	azcopier "github.com/permguard/permguard-common/pkg/extensions/copier"
	azvalidators "github.com/permguard/permguard-common/pkg/extensions/validators"
	azservices "github.com/permguard/permguard/pkg/agents/services"
	azstorage "github.com/permguard/permguard/pkg/agents/storage"
	azoptions "github.com/permguard/permguard/pkg/cli/options"
	azerrors "github.com/permguard/permguard/pkg/core/errors"
)

const (
	flagPrefixServer  = "server"
	flagSuffixAppData = "appdata"
)

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	host                 azservices.HostKind
	debug                bool
	logLevel             string
	appData              string
	centralStorageEngine azstorage.StorageKind
	storages             []azstorage.StorageKind
	storagesFactories    map[azstorage.StorageKind]azstorage.StorageFactoryProvider
	services             []azservices.ServiceKind
	servicesFactories    map[azservices.ServiceKind]azservices.ServiceFactoryProvider
}

// newServerConfig creates a new server factory configuration.
func newServerConfig(host azservices.HostKind, centralStorageEngine azstorage.StorageKind,
	storages []azstorage.StorageKind, storagesFactories map[azstorage.StorageKind]azstorage.StorageFactoryProvider,
	services []azservices.ServiceKind, servicesFactories map[azservices.ServiceKind]azservices.ServiceFactoryProvider,
) (*ServerConfig, error) {
	return &ServerConfig{
		host:                 host,
		centralStorageEngine: centralStorageEngine,
		storages:             azcopier.CopySlice(storages),
		storagesFactories:    azcopier.CopyMap(storagesFactories),
		services:             azcopier.CopySlice(services),
		servicesFactories:    azcopier.CopyMap(servicesFactories),
	}, nil
}

// GetHost returns the host kind.
func (c *ServerConfig) GetHost() azservices.HostKind {
	return c.host
}

// GetCentralStorageEngine returns the central storage engine.
func (c *ServerConfig) GetCentralStorageEngine() azstorage.StorageKind {
	return c.centralStorageEngine
}

// GetStorages returns service kinds.
func (c *ServerConfig) GetStorages() []azstorage.StorageKind {
	return azcopier.CopySlice(c.storages)
}

// GetStoragesFactories returns factories.
func (c *ServerConfig) GetStoragesFactories() map[azstorage.StorageKind]azstorage.StorageFactoryProvider {
	return azcopier.CopyMap(c.storagesFactories)
}

// GetServices returns service kinds.
func (c *ServerConfig) GetServices() []azservices.ServiceKind {
	return azcopier.CopySlice(c.services)
}

// GetServicesFactories returns factories.
func (c *ServerConfig) GetServicesFactories() map[azservices.ServiceKind]azservices.ServiceFactoryProvider {
	return azcopier.CopyMap(c.servicesFactories)
}

// GetAppData returns the zone data.
func (c *ServerConfig) GetAppData() string {
	return c.appData
}

// AddFlags adds flags.
func (c *ServerConfig) AddFlags(flagSet *flag.FlagSet) error {
	err := azoptions.AddFlagsForCommon(flagSet)
	if err != nil {
		return err
	}
	flagSet.String(azoptions.FlagName(flagPrefixServer, flagSuffixAppData), "./", "directory to be used as zone data")
	for _, fcty := range c.storagesFactories {
		config, _ := fcty.GetFactoryConfig()
		err = config.AddFlags(flagSet)
		if err != nil {
			return err
		}
	}
	for _, fcty := range c.servicesFactories {
		config, _ := fcty.GetFactoryConfig()
		err = config.AddFlags(flagSet)
		if err != nil {
			return err
		}
	}
	return nil
}

// InitFromViper initializes the configuration from viper.
func (c *ServerConfig) InitFromViper(v *viper.Viper) error {
	debug, logLevel, err := azoptions.InitFromViperForCommon(v)
	if err != nil {
		return err
	}
	c.debug = debug
	c.logLevel = logLevel
	c.appData = v.GetString(azoptions.FlagName(flagPrefixServer, flagSuffixAppData))
	if !azvalidators.IsValidPath(c.appData) {
		return azerrors.WrapSystemErrorWithMessage(azerrors.ErrCliArguments, "invalid zone data directory")
	}
	for _, fcty := range c.storagesFactories {
		config, err := fcty.GetFactoryConfig()
		if err != nil {
			return err
		}
		err = config.InitFromViper(v)
		if err != nil {
			return err
		}
	}
	for _, fcty := range c.servicesFactories {
		config, err := fcty.GetFactoryConfig()
		if err != nil {
			return err
		}
		err = config.InitFromViper(v)
		if err != nil {
			return err
		}
	}
	return nil
}
