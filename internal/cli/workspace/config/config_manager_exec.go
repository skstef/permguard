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

package config

import (
	"fmt"
	"strings"

	aziclicommon "github.com/permguard/permguard/internal/cli/common"
	azicliwkscommon "github.com/permguard/permguard/internal/cli/workspace/common"
	azerrors "github.com/permguard/permguard/pkg/core/errors"
)

// ExecInitialize initializes the config resources.
func (m *ConfigManager) ExecInitialize(lang string) error {
	config := config{
		Core: coreConfig{
			ClientVersion: m.ctx.GetClientVersion(),
			Language:      strings.ToLower(lang),
		},
		Remotes:      map[string]remoteConfig{},
		Repositories: map[string]repositoryConfig{},
	}
	return m.saveConfig(false, &config)
}

// ExecAddRemote adds a remote.
func (m *ConfigManager) ExecAddRemote(remote string, server string, aap int, pap int, output map[string]any, out aziclicommon.PrinterOutFunc) (map[string]any, error) {
	if output == nil {
		output = map[string]any{}
	}
	remote, err := azicliwkscommon.SanitizeRemote(remote)
	if err != nil {
		return output, err
	}
	if azicliwkscommon.IsReservedKeyword(remote) {
		return output, azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: remote %s is a reserved keyword", remote))
	}
	server = strings.ToLower(server)
	cfg, err := m.readConfig()
	if err != nil {
		return output, err
	}
	for rmt := range cfg.Remotes {
		if remote == rmt {
			return output, azerrors.WrapSystemError(azerrors.ErrCliRecordExists, fmt.Sprintf("cli: remote %s already exists", remote))
		}
	}
	cfgRemote := remoteConfig{
		Server:  server,
		AAPPort: aap,
		PAPPort: pap,
	}
	cfg.Remotes[remote] = cfgRemote
	m.saveConfig(true, cfg)
	out(nil, "", fmt.Sprintf("Remote %s has been added.", aziclicommon.KeywordText(remote)), nil, true)
	output = map[string]any{}
	if m.ctx.IsJSONOutput() {
		remotes := []any{}
		remoteObj := map[string]any{
			"remote": remote,
			"server": cfgRemote.Server,
			"aap":    cfgRemote.AAPPort,
			"pap":    cfgRemote.PAPPort,
		}
		remotes = append(remotes, remoteObj)
		output = out(output, "remotes", remotes, nil, true)
	}
	return output, nil
}

// ExecRemoveRemote removes a remote.
func (m *ConfigManager) ExecRemoveRemote(remote string, output map[string]any, out aziclicommon.PrinterOutFunc) (map[string]any, error) {
	if output == nil {
		output = map[string]any{}
	}
	remote, err := azicliwkscommon.SanitizeRemote(remote)
	if err != nil {
		return output, err
	}
	cfg, err := m.readConfig()
	if err != nil {
		return output, err
	}
	if _, ok := cfg.Remotes[remote]; !ok {
		return output, azerrors.WrapSystemError(azerrors.ErrCliRecordNotFound, fmt.Sprintf("cli: remote %s does not exist", remote))
	}
	for repository := range cfg.Repositories {
		if cfg.Repositories[repository].Remote == remote {
			return output, azerrors.WrapSystemError(azerrors.ErrCliRecordExists, fmt.Sprintf("cli: remote %s is used by repository %s", remote, repository))
		}
	}
	cfgRemote := cfg.Remotes[remote]
	out(nil, "", fmt.Sprintf("Remote %s has been removed.", aziclicommon.KeywordText(remote)), nil, true)
	output = map[string]any{}
	if m.ctx.IsJSONOutput() {
		remotes := []any{}
		remoteObj := map[string]any{
			"remote": remote,
			"server": cfgRemote.Server,
			"aap":    cfgRemote.AAPPort,
			"pap":    cfgRemote.PAPPort,
		}
		remotes = append(remotes, remoteObj)
		output = out(output, "remotes", remotes, nil, true)
	}
	delete(cfg.Remotes, remote)
	m.saveConfig(true, cfg)
	return output, nil
}

// ExecListRemotes lists the remotes.
func (m *ConfigManager) ExecListRemotes(output map[string]any, out aziclicommon.PrinterOutFunc) (map[string]any, error) {
	if output == nil {
		output = map[string]any{}
	}
	cfg, err := m.readConfig()
	if err != nil {
		return output, err
	}
	if m.ctx.IsTerminalOutput() {
		remotes := []string{}
		for cfgRemote := range cfg.Remotes {
			remotes = append(remotes, cfgRemote)
		}
		if len(remotes) == 0 {
			out(nil, "", "Your workspace doesn't have any remote configured.", nil, true)
		} else {
			out(nil, "", "Your workspace configured remotes:\n", nil, true)
			for _, remote := range remotes {
				out(nil, "", fmt.Sprintf("	- %s", aziclicommon.KeywordText(remote)), nil, true)
			}
			out(nil, "", "\n", nil, false)
		}
	} else if m.ctx.IsJSONOutput() {
		remotes := []any{}
		for cfgRemote := range cfg.Remotes {
			remoteObj := map[string]any{
				"remote": cfgRemote,
				"server": cfg.Remotes[cfgRemote].Server,
				"aap":    cfg.Remotes[cfgRemote].AAPPort,
				"pap":    cfg.Remotes[cfgRemote].PAPPort,
			}
			remotes = append(remotes, remoteObj)
		}
		output = out(output, "remotes", remotes, nil, true)
	}
	return output, nil
}

// ExecAddRepo adds a repo.
func (m *ConfigManager) ExecAddRepo(repoURI, ref, remote, repo, repoID string, account int64, output map[string]any, out aziclicommon.PrinterOutFunc) (map[string]any, error) {
	if output == nil {
		output = map[string]any{}
	}
	cfg, err := m.readConfig()
	if err != nil {
		return output, err
	}
	var cfgRepo repositoryConfig
	exists := false
	for repository := range cfg.Repositories {
		if repository == repo && cfg.Repositories[repoURI].Remote == remote {
			cfgRepo = cfg.Repositories[repoURI]
			exists = true
			break
		}
	}
	if !exists {
		for key, repo := range cfg.Repositories {
			repo.IsHead = false
			cfg.Repositories[key] = repo
		}
		cfgRepo = repositoryConfig{
			Ref:       ref,
			Remote:    remote,
			AccountID: account,
			RepoName:  repo,
			RepoID:    repoID,
			IsHead:    true,
		}
		cfg.Repositories[repoURI] = cfgRepo
		m.saveConfig(true, cfg)
	}
	if m.ctx.IsVerboseTerminalOutput() {
		out(nil, "repo", fmt.Sprintf("Ref successfully set to %s.", aziclicommon.KeywordText(cfgRepo.Ref)), nil, true)
	}
	out(nil, "", fmt.Sprintf("Repo %s has been added.", aziclicommon.KeywordText(repo)), nil, true)
	output = map[string]any{}
	if m.ctx.IsJSONOutput() {
		remotes := []any{}
		remoteObj := map[string]any{
			"repo": repo,
		}
		remotes = append(remotes, remoteObj)
		output = out(output, "repos", remotes, nil, true)
	}
	return output, nil
}

// ExecListRepos lists the repos.
func (m *ConfigManager) ExecListRepos(output map[string]any, out aziclicommon.PrinterOutFunc) (map[string]any, error) {
	if output == nil {
		output = map[string]any{}
	}
	cfg, err := m.readConfig()
	if err != nil {
		return output, err
	}
	if m.ctx.IsTerminalOutput() {
		repos := []string{}
		for cfgRepo := range cfg.Repositories {
			cfgRepoTxt := cfgRepo
			isHead := cfg.Repositories[cfgRepo].IsHead
			if isHead {
				cfgRepoTxt = fmt.Sprintf("*%s", cfgRepo)
			}
			repos = append(repos, cfgRepoTxt)
		}
		if len(repos) == 0 {
			out(nil, "", "Your workspace doesn't have any repository configured.", nil, true)
		} else {
			out(nil, "", "Your workspace configured repositories:\n", nil, true)
			for _, repo := range repos {
				out(nil, "", fmt.Sprintf("	- %s", aziclicommon.KeywordText(repo)), nil, true)
			}
			out(nil, "", "\n", nil, false)
		}
	} else if m.ctx.IsJSONOutput() {
		repos := []any{}
		for cfgRepo := range cfg.Repositories {
			isHead := cfg.Repositories[cfgRepo].IsHead
			repoObj := map[string]any{
				"ref":     cfg.Repositories[cfgRepo].Ref,
				"repo":    cfgRepo,
				"is_head": isHead,
			}
			repos = append(repos, repoObj)
		}
		output = out(output, "repos", repos, nil, true)
	}
	return output, nil
}
