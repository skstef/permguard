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
	"strings"

	azerrors "github.com/permguard/permguard/pkg/extensions/errors"
	azvalidators "github.com/permguard/permguard/pkg/extensions/validators"
)

// SanitizeRemote sanitizes the remote name.
func SanitizeRemote(remote string) (string, error) {
	remote = strings.ToLower(remote)
	if !azvalidators.ValidateSimpleName(remote) {
		return "", azerrors.WrapSystemError(azerrors.ErrCliInput, fmt.Sprintf("cli: invalid remote name %s", remote))
	}
	return remote, nil
}