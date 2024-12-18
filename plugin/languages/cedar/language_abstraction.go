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

package cedar

import (
	"bytes"
	"strings"

	"github.com/cedar-policy/cedar-go"

	azlangtypes "github.com/permguard/permguard-abs-language/pkg/languages/types"
	azlangobjs "github.com/permguard/permguard-abs-language/pkg/objects"
	azerrors "github.com/permguard/permguard/pkg/core/errors"
	azlang "github.com/permguard/permguard/pkg/core/languages"
)

const (
	// LanguageName specifies the canonical name of the Cedar language.
	LanguageName = "cedar"

	// LanguageCedar represents the unique identifier for the Cedar language.
	LanguageCedar = "cedar"
	// LanguageCedarID represents the unique identifier for the Cedar language.
	LanguageCedarID = uint32(1)

	// LanguageCedarJSON represents the unique identifier for the JSON-based Cedar language.
	LanguageCedarJSON = "cedar-json"
	// LanguageCedarJSONID represents the unique identifier for the JSON-based Cedar language.
	LanguageCedarJSONID = uint32(2)

	// LanguageSyntaxVersion defines the latest syntax version used by the Cedar language.
	LanguageSyntaxVersion = "*"
	// LanguageSyntaxVersionID defines the latest syntax version ID used by the Cedar language.
	LanguageSyntaxVersionID = uint32(0)
	// LanguageSchemaType specifies the schema type for Cedar language.
	LanguageSchemaType = "schema"
	// LanguageSchemaTypeID specifies the schema type ID for Cedar language.
	LanguageSchemaTypeID = uint32(1)
	// LanguagePolicyType specifies the policy type for Cedar language.
	LanguagePolicyType = "policy"
	// LanguagePolicyTypeID specifies the policy type ID for Cedar language.
	LanguagePolicyTypeID = uint32(2)

	// LanguageFileExtension specifies the standard file extension for Cedar language files.
	LanguageFileExtension = ".cedar"
	// SchemaFileName defines the default filename for the schema definition associated with Cedar.
	SchemaFileName = "schema.json"
)

// CedarLanguageAbstraction is the abstraction for the cedar language.
type CedarLanguageAbstraction struct {
	objMng *azlangobjs.ObjectManager
}

// NewCedarLanguageAbstraction creates a new CedarLanguageAbstraction.
func NewCedarLanguageAbstraction() (*CedarLanguageAbstraction, error) {
	objMng, err := azlangobjs.NewObjectManager()
	if err != nil {
		return nil, err
	}

	return &CedarLanguageAbstraction{
		objMng: objMng,
	}, nil
}

// GetLanguageSpecification returns the specification for the language.
func (abs *CedarLanguageAbstraction) GetLanguageSpecification() azlang.LanguageSpecification {
	return &CedarLanguageSpecification{
		language:                      LanguageName,
		frontendLanguage:              LanguageCedar,
		frontendLanguageID:            LanguageCedarID,
		backendLanguage:               LanguageCedarJSON,
		backendLanguageID:             LanguageCedarJSONID,
		supportedPolicyFileExtensions: []string{LanguageFileExtension},
		supportedSchemaFileNames:      []string{SchemaFileName},
	}
}

// CreateCommitObject creates a commit object.
func (abs *CedarLanguageAbstraction) CreateCommitObject(commit *azlangobjs.Commit) (*azlangobjs.Object, error) {
	return abs.objMng.CreateCommitObject(commit)
}

// ConvertObjectToCommit converts an object to a commit.
func (abs *CedarLanguageAbstraction) ConvertObjectToCommit(obj *azlangobjs.Object) (*azlangobjs.Commit, error) {
	objInfo, err := abs.objMng.GetObjectInfo(obj)
	if err != nil {
		return nil, err
	}

	value, ok := objInfo.GetInstance().(*azlangobjs.Commit)
	if !ok {
		return nil, azerrors.WrapSystemError(azerrors.ErrLanguageFile, "cedar: invalid object type")
	}
	return value, nil
}

// CreateTreeObject creates a tree object.
func (abs *CedarLanguageAbstraction) CreateTreeObject(tree *azlangobjs.Tree) (*azlangobjs.Object, error) {
	return abs.objMng.CreateTreeObject(tree)
}

// ConvertObjectToTree converts an object to a tree.
func (abs *CedarLanguageAbstraction) ConvertObjectToTree(obj *azlangobjs.Object) (*azlangobjs.Tree, error) {
	objInfo, err := abs.objMng.GetObjectInfo(obj)
	if err != nil {
		return nil, err
	}

	value, ok := objInfo.GetInstance().(*azlangobjs.Tree)
	if !ok {
		return nil, azerrors.WrapSystemError(azerrors.ErrLanguageFile, "cedar: invalid object type")
	}
	return value, nil
}

// CreatePolicyBlobObjects creates multi sections policy blob objects.
func (abs *CedarLanguageAbstraction) CreatePolicyBlobObjects(filePath string, data []byte) (*azlangobjs.MultiSectionsObject, error) {
	policySet, err := cedar.NewPolicySetFromBytes(filePath, data)
	if err != nil {
		multiSecObj, err2 := azlangobjs.NewMultiSectionsObject(filePath, 0, nil)
		if err2 != nil {
			return nil, err2
		}
		multiSecObj.AddSectionObjectWithError(0, err)
		return multiSecObj, nil
	}

	policiesMap := policySet.Map()
	multiSecObj, err := azlangobjs.NewMultiSectionsObject(filePath, len(policiesMap), nil)
	if err != nil {
		return nil, err
	}

	const (
		codeType         = azlangtypes.ClassTypePolicy
		lang             = LanguageCedarJSON
		langID           = LanguageCedarJSONID
		langVersion      = LanguageSyntaxVersion
		langVersionID    = LanguageSyntaxVersionID
		langPolicyType   = LanguagePolicyType
		langPolicyTypeID = LanguagePolicyTypeID
	)

	i := -1
	for policyName, policy := range policiesMap {
		i++
		name := string(policyName)
		codeID := name

		header, err := azlangobjs.NewObjectHeader(true, langID, langVersionID, codeID, langPolicyTypeID)
		if err != nil {
			multiSecObj.AddSectionObjectWithError(i, err)
			continue
		}

		policyJson, err := policy.MarshalJSON()
		if err != nil {
			multiSecObj.AddSectionObjectWithError(i, err)
			continue
		}

		obj, err := abs.objMng.CreateBlobObject(header, policyJson)
		if err != nil {
			multiSecObj.AddSectionObjectWithError(i, err)
			continue
		}

		objInfo, err := abs.objMng.GetObjectInfo(obj)
		if err != nil {
			return nil, err
		}

		multiSecObj.AddSectionObjectWithParams(obj, objInfo.GetType(), name, codeID, codeType, lang, langVersion, langPolicyType, i)
	}

	return multiSecObj, nil
}

// ReadPolicyBlobContentBytes reads the policy blob object content bytes.
func (abs *CedarLanguageAbstraction) ReadPolicyBlobContentBytes(obj *azlangobjs.Object) (string, []byte, error) {
	objInfo, err := abs.objMng.GetObjectInfo(obj)
	if err != nil {
		return "", nil, err
	}
	objHeader := objInfo.GetHeader()
	if !objHeader.IsNativeLanguage() {
		return "", nil, azerrors.WrapSystemError(azerrors.ErrLanguageFile, "cedar: invalid object type")
	}
	return "", objInfo.GetObject().GetContent(), nil
}

// CreateMultiPolicyContentBytes creates a multi policy content bytes.
func (abs *CedarLanguageAbstraction) CreateMultiPolicyContentBytes(blocks [][]byte) ([]byte, string, error) {
	var sb strings.Builder
	for i, block := range blocks {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.Write(block)
	}
	return []byte(sb.String()), LanguageFileExtension, nil
}

// CreateSchemaBlobObjects creates multi sections schema blob objects.
func (abs *CedarLanguageAbstraction) CreateSchemaBlobObjects(path string, data []byte) (*azlangobjs.MultiSectionsObject, error) {
	codeID := azlangtypes.ClassTypeSchema
	schemaType := azlangtypes.ClassTypeSchema

	const (
		lang             = LanguageCedarJSON
		langID           = LanguageCedarJSONID
		langVersion      = LanguageSyntaxVersion
		langVersionID    = LanguageSyntaxVersionID
		langPolicyType   = LanguagePolicyType
		langPolicyTypeID = LanguagePolicyTypeID
	)

	multiSecObj, err := azlangobjs.NewMultiSectionsObject(path, 1, nil)
	header, err := azlangobjs.NewObjectHeader(true, langID, langVersionID, codeID, langPolicyTypeID)
	if err != nil {
		multiSecObj.AddSectionObjectWithError(0, err)
		return multiSecObj, nil
	}

	obj, err := abs.objMng.CreateBlobObject(header, data)
	if err != nil {
		multiSecObj.AddSectionObjectWithError(0, err)
		return multiSecObj, nil
	}

	objInfo, err := abs.objMng.GetObjectInfo(obj)
	if err != nil {
		return nil, err
	}

	multiSecObj.AddSectionObjectWithParams(obj, objInfo.GetType(), codeID, codeID, schemaType, lang, langPolicyType, langVersion, 0)
	return multiSecObj, nil
}

// ReadSchemaBlobContentBytes reads the schema blob object content bytes.
func (abs *CedarLanguageAbstraction) ReadSchemaBlobContentBytes(obj *azlangobjs.Object) (string, []byte, error) {
	return "", nil, nil
}

// CreateSchemaContentBytes creates a schema content bytes.
func (abs *CedarLanguageAbstraction) CreateSchemaContentBytes(blocks [][]byte) ([]byte, string, error) {
	return bytes.Join(blocks, nil), LanguageFileExtension, nil
}
