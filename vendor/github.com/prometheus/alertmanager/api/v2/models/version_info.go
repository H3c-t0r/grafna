// Code generated by go-swagger; DO NOT EDIT.

// Copyright Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// VersionInfo version info
// swagger:model versionInfo
type VersionInfo struct {

	// branch
	// Required: true
	Branch *string `json:"branch"`

	// build date
	// Required: true
	BuildDate *string `json:"buildDate"`

	// build user
	// Required: true
	BuildUser *string `json:"buildUser"`

	// go version
	// Required: true
	GoVersion *string `json:"goVersion"`

	// revision
	// Required: true
	Revision *string `json:"revision"`

	// version
	// Required: true
	Version *string `json:"version"`
}

// Validate validates this version info
func (m *VersionInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBranch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBuildDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBuildUser(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGoVersion(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRevision(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *VersionInfo) validateBranch(formats strfmt.Registry) error {

	if err := validate.Required("branch", "body", m.Branch); err != nil {
		return err
	}

	return nil
}

func (m *VersionInfo) validateBuildDate(formats strfmt.Registry) error {

	if err := validate.Required("buildDate", "body", m.BuildDate); err != nil {
		return err
	}

	return nil
}

func (m *VersionInfo) validateBuildUser(formats strfmt.Registry) error {

	if err := validate.Required("buildUser", "body", m.BuildUser); err != nil {
		return err
	}

	return nil
}

func (m *VersionInfo) validateGoVersion(formats strfmt.Registry) error {

	if err := validate.Required("goVersion", "body", m.GoVersion); err != nil {
		return err
	}

	return nil
}

func (m *VersionInfo) validateRevision(formats strfmt.Registry) error {

	if err := validate.Required("revision", "body", m.Revision); err != nil {
		return err
	}

	return nil
}

func (m *VersionInfo) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *VersionInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VersionInfo) UnmarshalBinary(b []byte) error {
	var res VersionInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
