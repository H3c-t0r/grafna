// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v0alpha1

// ReceiverSpecApplyConfiguration represents an declarative configuration of the ReceiverSpec type for use
// with apply.
type ReceiverSpecApplyConfiguration struct {
	Integrations []IntegrationApplyConfiguration `json:"integrations,omitempty"`
	Title        *string                         `json:"title,omitempty"`
}

// ReceiverSpecApplyConfiguration constructs an declarative configuration of the ReceiverSpec type for use with
// apply.
func ReceiverSpec() *ReceiverSpecApplyConfiguration {
	return &ReceiverSpecApplyConfiguration{}
}

// WithIntegrations adds the given value to the Integrations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Integrations field.
func (b *ReceiverSpecApplyConfiguration) WithIntegrations(values ...*IntegrationApplyConfiguration) *ReceiverSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithIntegrations")
		}
		b.Integrations = append(b.Integrations, *values[i])
	}
	return b
}

// WithTitle sets the Title field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Title field is set to the value of the last call.
func (b *ReceiverSpecApplyConfiguration) WithTitle(value string) *ReceiverSpecApplyConfiguration {
	b.Title = &value
	return b
}
