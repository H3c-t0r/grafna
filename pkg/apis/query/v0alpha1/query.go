package v0alpha1

import (
	data "github.com/grafana/grafana-plugin-sdk-go/experimental/apis/data/v0alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Generic query request with shared time across all values
// Copied from: https://github.com/grafana/grafana/blob/main/pkg/api/dtos/models.go#L62
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type QueryDataRequest struct {
	metav1.TypeMeta `json:",inline"`

	// The time range used when not included on each query
	data.QueryDataRequest `json:",inline"`
}
