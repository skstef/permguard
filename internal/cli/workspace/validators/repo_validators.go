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

package validators

import (
	"fmt"
	"strconv"
	"strings"

	azerrors "github.com/permguard/permguard/pkg/extensions/errors"
	azvalidators "github.com/permguard/permguard/pkg/authz/validators"
)

// SanitizeRepo sanitizes the repo name.
func SanitizeRepo(repo string) (string, int64, string, error) {
	repo = strings.ToLower(repo)
	items := strings.Split(repo, "/")
	if len(items) < 3 {
		return "", 0, "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid repo %s", repo))
	}
	remoteName, err := SanitizeRemote(items[0])
	if err != nil {
		return "", 0, "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid remote %s", remoteName))
	}
	accountIDStr := items[1]
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		return "", 0, "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid account id %s", accountIDStr))
	}
	err = azvalidators.ValidateAccountID("repo", accountID)
	if err != nil {
		return "", 0, "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid account id %s", accountIDStr))
	}
	repoName := items[2]
	err = azvalidators.ValidateName("repo", repoName)
	if err != nil {
		return "", 0, "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid repo name %s", repoName))
	}
	return remoteName, accountID, repoName, nil
}