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

package ref

// headReferenceConfig represents the configuration for the head.
type headReferenceConfig struct {
	Ref string `toml:"ref"`
}

// headConfig represents the configuration for the head.
type headConfig struct {
	Reference headReferenceConfig `toml:"reference"`
}

// refObjectsConfig represents the configuration for the objects.
type refObjectsConfig struct {
	UpstreamRef string `toml:"upstreamref"`
	LedgerID    string `toml:"ledgerid"`
	Commit      string `toml:"commit"`
}

// refConfig represents the configuration for the ref.
type refConfig struct {
	Objects refObjectsConfig `toml:"objects"`
}
