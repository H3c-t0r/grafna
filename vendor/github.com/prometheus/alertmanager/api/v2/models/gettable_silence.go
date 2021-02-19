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

// GettableSilence gettable silence
// swagger:model gettableSilence
type GettableSilence struct {

	// id
	// Required: true
	ID *string `json:"id"`

	// status
	// Required: true
	Status *SilenceStatus `json:"status"`

	// updated at
	// Required: true
	// Format: date-time
	UpdatedAt *strfmt.DateTime `json:"updatedAt"`

	Silence
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *GettableSilence) UnmarshalJSON(raw []byte) error {
	// AO0
	var dataAO0 struct {
		ID *string `json:"id"`

		Status *SilenceStatus `json:"status"`

		UpdatedAt *strfmt.DateTime `json:"updatedAt"`
	}
	if err := swag.ReadJSON(raw, &dataAO0); err != nil {
		return err
	}

	m.ID = dataAO0.ID

	m.Status = dataAO0.Status

	m.UpdatedAt = dataAO0.UpdatedAt

	// AO1
	var aO1 Silence
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.Silence = aO1

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m GettableSilence) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	var dataAO0 struct {
		ID *string `json:"id"`

		Status *SilenceStatus `json:"status"`

		UpdatedAt *strfmt.DateTime `json:"updatedAt"`
	}

	dataAO0.ID = m.ID

	dataAO0.Status = m.Status

	dataAO0.UpdatedAt = m.UpdatedAt

	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0)
	if errAO0 != nil {
		return nil, errAO0
	}
	_parts = append(_parts, jsonDataAO0)

	aO1, err := swag.WriteJSON(m.Silence)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this gettable silence
func (m *GettableSilence) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	// validation for a type composition with Silence
	if err := m.Silence.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *GettableSilence) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *GettableSilence) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	if m.Status != nil {
		if err := m.Status.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("status")
			}
			return err
		}
	}

	return nil
}

func (m *GettableSilence) validateUpdatedAt(formats strfmt.Registry) error {

	if err := validate.Required("updatedAt", "body", m.UpdatedAt); err != nil {
		return err
	}

	if err := validate.FormatOf("updatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *GettableSilence) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GettableSilence) UnmarshalBinary(b []byte) error {
	var res GettableSilence
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
