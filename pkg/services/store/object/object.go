// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: object.proto

package object

// Will be replaced with something from the SDK
type UserInfo struct {
	Id    int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`      // internal grafana ID
	Login string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"` // string ID?
}

// The canonical object/document data -- this represents the raw bytes and storage level metadata
type RawObject struct {
	// Unique ID
	UID string `protobuf:"bytes,1,opt,name=UID,proto3" json:"UID,omitempty"`
	// Identify the object kind.  This kind will be used to apply a schema to the body and
	// will trigger additional indexing behavior.
	Kind string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	// Time in epoch milliseconds that the object was modified
	Modified int64 `protobuf:"varint,3,opt,name=modified,proto3" json:"modified,omitempty"`
	// Who modified the object
	ModifiedBy *UserInfo `protobuf:"bytes,4,opt,name=modified_by,json=modifiedBy,proto3" json:"modified_by,omitempty"`
	// Content Length
	Size int64 `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`
	// MD5 digest of the body
	ETag string `protobuf:"bytes,6,opt,name=ETag,proto3" json:"ETag,omitempty"`
	// Raw bytes of the storage object.  The kind will determine what is a valid payload
	Body []byte `protobuf:"bytes,7,opt,name=body,proto3" json:"body,omitempty"`
	// The version will change when the object is saved. It is not necessarily sortable
	//
	// NOTE: currently managed by the dashboard+dashboard_version tables
	Version string `protobuf:"bytes,8,opt,name=version,proto3" json:"version,omitempty"`
	// optional "save" or "commit" message
	//
	// NOTE: currently managed by the dashboard_version table, and will be returned from a "history" command
	Comment string `protobuf:"bytes,9,opt,name=comment,proto3" json:"comment,omitempty"`
	// Location (path/repo/etc) that defines the canonocal form
	//
	// NOTE: currently managed by the dashboard_provisioning table
	SyncSrc string `protobuf:"bytes,10,opt,name=sync_src,json=syncSrc,proto3" json:"sync_src,omitempty"`
	// Time in epoch milliseconds that the object was last synced with an external system (provisioning/git)
	//
	// NOTE: currently managed by the dashboard_provisioning table
	SyncTime int64 `protobuf:"varint,11,opt,name=sync_time,json=syncTime,proto3" json:"sync_time,omitempty"`
}

// Searchable fields extracted from the object
type ObjectErrorInfo struct {
	Code    int64  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"` // TODO... registry somewhere...  should be limited to most severe issues
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details string `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

type ExternalReference struct {
	// datasource, panel
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// prometheus / heatmap
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Unique ID for this object
	UID string `protobuf:"bytes,3,opt,name=UID,proto3" json:"UID,omitempty"`
}
