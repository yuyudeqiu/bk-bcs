/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// KvType is the type of kv
type KvType string

const (
	// KvStr is the type for string kv
	KvStr KvType = "string"
	// KvNumber is the type for number kv
	KvNumber KvType = "number"
	// KvText is the type for text kv
	KvText KvType = "text"
	// KvJson is the type for json kv
	KvJson KvType = "json"
	// KvYAML is the type for yaml kv
	KvYAML KvType = "yaml"
	// KvXml is the type for xml kv
	KvXml KvType = "xml"

	// MaxValueLength max value length 1MB
	MaxValueLength = 1 * 1024 * 1024
)

// ValidateValue the kvType and value match
func (k KvType) ValidateValue(value string) error {

	if value == "" {
		return errors.New("kv value is null")
	}

	if len(value) > MaxValueLength {
		return fmt.Errorf("the length of the value must not exceed %d MB", MaxValueLength)
	}

	switch k {
	case KvStr:
		if strings.Contains(value, "\n") {
			return errors.New("newline characters are not allowed in string-type values")
		}
		return nil
	case KvNumber:
		if !isStringConvertibleToNumber(value) {
			return fmt.Errorf("value is not a number")
		}
		return nil
	case KvText:
		return nil
	case KvJson:
		if !json.Valid([]byte(value)) {
			return fmt.Errorf("value is not a json")
		}
		return nil
	case KvYAML:
		var data interface{}
		if err := yaml.Unmarshal([]byte(value), &data); err != nil {
			return fmt.Errorf("value is not a yaml, err: %v", err)
		}
		return nil
	case KvXml:
		var v interface{}
		if err := xml.Unmarshal([]byte(value), &v); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid key-value type")
	}
}

// Validate checks if the KvType enumeration value is valid.
func (k KvType) Validate() error {
	switch k {
	case KvStr:
	case KvNumber:
	case KvText:
	case KvJson:
	case KvYAML:
	default:
		return errors.New("invalid key-value type")
	}
	return nil
}

// isStringConvertibleToNumber checks if the given string can be converted to an integer or a floating-point number.
func isStringConvertibleToNumber(s string) bool {
	_, err := strconv.Atoi(s)
	if err == nil {
		return true
	}

	_, err = strconv.ParseFloat(s, 64)
	return err == nil

}

// UpsertKvOption is used to define options for inserting or updating key-value data.
type UpsertKvOption struct {
	BizID  uint32
	AppID  uint32
	Key    string
	Value  string
	KvType KvType
}

// Validate is used to validate the effectiveness of the UpsertKvOption structure.
func (o *UpsertKvOption) Validate() error {
	if o.BizID <= 0 {
		return errors.New("invalid biz id, should >= 1")
	}

	if o.AppID <= 0 {
		return errors.New("invalid app id, should >= 1")
	}

	if o.Key == "" {
		return errors.New("kv key is required")
	}

	if o.Value == "" {
		return errors.New("kv value is required")
	}

	if err := o.KvType.ValidateValue(o.Value); err != nil {
		return err
	}

	return nil
}

// GetLastKvOpt is used to define options for retrieving the last key-value data.
type GetLastKvOpt struct {
	BizID uint32
	AppID uint32
	Key   string
}

// Validate is used to validate the effectiveness of the GetLastKvOpt structure.
func (o *GetLastKvOpt) Validate() error {
	if o.BizID <= 0 {
		return errors.New("invalid biz id, should >= 1")
	}

	if o.AppID <= 0 {
		return errors.New("invalid app id, should >= 1")
	}

	if o.Key == "" {
		return errors.New("kv key is required")
	}
	return nil
}

// GetKvByVersion is used to define options for retrieving key-value data by version.
type GetKvByVersion struct {
	BizID   uint32
	AppID   uint32
	Key     string
	Version int
}

// Validate is used to validate the effectiveness of the GetKvByVersion structure.
func (o *GetKvByVersion) Validate() error {
	if o.BizID <= 0 {
		return errors.New("invalid biz id, should >= 1")
	}

	if o.AppID <= 0 {
		return errors.New("invalid app id, should >= 1")
	}

	if o.Key == "" {
		return errors.New("kv key is required")
	}

	if o.Version <= 0 {
		return errors.New("invalid version, should >= 1")
	}

	return nil
}

// DeleteKvOpt is used to define options for deleting key-value data.
type DeleteKvOpt struct {
	BizID uint32
	AppID uint32
	Key   string
}

// Validate is used to validate the effectiveness of the DeleteKvOpt structure.
func (o *DeleteKvOpt) Validate() error {
	if o.BizID <= 0 {
		return errors.New("invalid biz id, should >= 1")
	}

	if o.AppID <= 0 {
		return errors.New("invalid app id, should >= 1")
	}

	if o.Key == "" {
		return errors.New("kv key is required")
	}
	return nil
}

// ListKvOption defines options to list kv.
type ListKvOption struct {
	BizID     uint32    `json:"biz_id"`
	AppID     uint32    `json:"app_id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	SearchKey string    `json:"search_key"`
	All       bool      `json:"all"`
	Page      *BasePage `json:"page"`
}

// Validate is used to validate the effectiveness of the ListKvOption structure.
func (opt *ListKvOption) Validate(po *PageOption) error {
	if opt.BizID <= 0 {
		return errors.New("invalid biz id, should >= 1")
	}
	if opt.AppID <= 0 {
		return errors.New("invalid app id, should >= 1")
	}

	if opt.Page == nil {
		return errors.New("page is null")
	}

	if err := opt.Page.Validate(po); err != nil {
		return err
	}

	return nil
}