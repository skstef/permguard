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

package authn

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	aziclicommon "github.com/permguard/permguard/internal/cli/common"
	azcli "github.com/permguard/permguard/pkg/cli"
	azconfigs "github.com/permguard/permguard/pkg/configs"
)

const (
	// commandNameForIdentitiesList is the command name for identities list.
	commandNameForIdentitiesList = "identities.list"
)

// runECommandForListIdentities runs the command for creating an identity.
func runECommandForListIdentities(deps azcli.CliDependenciesProvider, cmd *cobra.Command, v *viper.Viper) error {
	ctx, printer, err := aziclicommon.CreateContextAndPrinter(deps, cmd, v)
	if err != nil {
		color.Red(aziclicommon.ErrorMessageCliBug)
		return aziclicommon.ErrCommandSilent
	}
	aapTarget := ctx.GetAAPTarget()
	client, err := deps.CreateGrpcAAPClient(aapTarget)
	if err != nil {
		printer.Error(fmt.Errorf("invalid aap target %s", aapTarget))
		return aziclicommon.ErrCommandSilent
	}
	page := v.GetInt32(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonPage))
	pageSize := v.GetInt32(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonPageSize))
	accountID := v.GetInt64(azconfigs.FlagName(commandNameForIdentity, aziclicommon.FlagCommonAccountID))
	identitySourceID := v.GetString(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentitySourceID))
	identityID := v.GetString(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentityID))
	kind := v.GetString(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentityKind))
	name := v.GetString(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonName))
	identities, err := client.FetchIdentitiesBy(page, pageSize, accountID, identitySourceID, identityID, kind, name)
	if err != nil {
		printer.Error(err)
		return aziclicommon.ErrCommandSilent
	}
	output := map[string]any{}
	if ctx.IsTerminalOutput() {
		for _, identity := range identities {
			identityID := identity.IdentityID
			identityName := identity.Name
			output[identityID] = identityName
		}
	} else if ctx.IsJSONOutput() {
		output["identities"] = identities
	}
	printer.Print(output)
	return nil
}

// createCommandForIdentityList creates a command for managing identitylist.
func createCommandForIdentityList(deps azcli.CliDependenciesProvider, v *viper.Viper) *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List identities",
		Long: aziclicommon.BuildCliLongTemplate(`This command lists all identities.

Examples:
  # list all identities and output the result in json format
  permguard authn identities list --account 268786704340 --output json
  # list all identities and apply filter by name
  permguard authn identities list --account 268786704340 --name identity1
  # list all identities and apply filter by identity source id
  permguard authn identities list --account 268786704340 --identityid 1da1d9094501425085859c60429163c2
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runECommandForListIdentities(deps, cmd, v)
		},
	}
	command.Flags().Int32P(aziclicommon.FlagCommonPage, aziclicommon.FlagCommonPageShort, 1, "page number")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonPage), command.Flags().Lookup(aziclicommon.FlagCommonPage))
	command.Flags().Int32P(aziclicommon.FlagCommonPageSize, aziclicommon.FlagCommonPageSizeShort, 1000, "page size")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonPageSize), command.Flags().Lookup(aziclicommon.FlagCommonPageSize))
	command.Flags().String(flagIdentitySourceID, "", "identity source id filter")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentitySourceID), command.Flags().Lookup(flagIdentitySourceID))
	command.Flags().String(flagIdentityKind, "", "identity kind filer")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentityKind), command.Flags().Lookup(flagIdentityKind))
	command.Flags().String(flagIdentityID, "", "identity id filter")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, flagIdentityID), command.Flags().Lookup(flagIdentityID))
	command.Flags().String(aziclicommon.FlagCommonName, "", "identity name filter")
	v.BindPFlag(azconfigs.FlagName(commandNameForIdentitiesList, aziclicommon.FlagCommonName), command.Flags().Lookup(aziclicommon.FlagCommonName))
	return command
}