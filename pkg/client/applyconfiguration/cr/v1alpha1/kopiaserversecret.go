/*
Copyright 2024 The Kanister Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// KopiaServerSecretApplyConfiguration represents an declarative configuration of the KopiaServerSecret type for use
// with apply.
type KopiaServerSecretApplyConfiguration struct {
	Username       *string                                 `json:"username,omitempty"`
	Hostname       *string                                 `json:"hostname,omitempty"`
	UserPassphrase *KopiaServerSecretRefApplyConfiguration `json:"userPassphrase,omitempty"`
	TLSCert        *KopiaServerSecretRefApplyConfiguration `json:"tlsCert,omitempty"`
	ConnectOptions map[string]int                          `json:"connectOptions,omitempty"`
}

// KopiaServerSecretApplyConfiguration constructs an declarative configuration of the KopiaServerSecret type for use with
// apply.
func KopiaServerSecret() *KopiaServerSecretApplyConfiguration {
	return &KopiaServerSecretApplyConfiguration{}
}

// WithUsername sets the Username field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Username field is set to the value of the last call.
func (b *KopiaServerSecretApplyConfiguration) WithUsername(value string) *KopiaServerSecretApplyConfiguration {
	b.Username = &value
	return b
}

// WithHostname sets the Hostname field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Hostname field is set to the value of the last call.
func (b *KopiaServerSecretApplyConfiguration) WithHostname(value string) *KopiaServerSecretApplyConfiguration {
	b.Hostname = &value
	return b
}

// WithUserPassphrase sets the UserPassphrase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UserPassphrase field is set to the value of the last call.
func (b *KopiaServerSecretApplyConfiguration) WithUserPassphrase(value *KopiaServerSecretRefApplyConfiguration) *KopiaServerSecretApplyConfiguration {
	b.UserPassphrase = value
	return b
}

// WithTLSCert sets the TLSCert field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLSCert field is set to the value of the last call.
func (b *KopiaServerSecretApplyConfiguration) WithTLSCert(value *KopiaServerSecretRefApplyConfiguration) *KopiaServerSecretApplyConfiguration {
	b.TLSCert = value
	return b
}

// WithConnectOptions puts the entries into the ConnectOptions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the ConnectOptions field,
// overwriting an existing map entries in ConnectOptions field with the same key.
func (b *KopiaServerSecretApplyConfiguration) WithConnectOptions(entries map[string]int) *KopiaServerSecretApplyConfiguration {
	if b.ConnectOptions == nil && len(entries) > 0 {
		b.ConnectOptions = make(map[string]int, len(entries))
	}
	for k, v := range entries {
		b.ConnectOptions[k] = v
	}
	return b
}
