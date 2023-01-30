// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     GoTypesJenny
//     LatestJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package librarypanel

// LibraryElementDTOMeta defines model for LibraryElementDTOMeta.
type LibraryElementDTOMeta struct {
	ConnectedDashboards int64                     `json:"connectedDashboards"`
	Created             string                    `json:"created"`
	CreatedBy           LibraryElementDTOMetaUser `json:"createdBy"`
	FolderName          string                    `json:"folderName"`
	FolderUid           string                    `json:"folderUid"`
	Updated             string                    `json:"updated"`
	UpdatedBy           LibraryElementDTOMetaUser `json:"updatedBy"`
}

// LibraryElementDTOMetaUser defines model for LibraryElementDTOMetaUser.
type LibraryElementDTOMetaUser struct {
	AvatarUrl string `json:"avatarUrl"`
	Id        int64  `json:"id"`
	Name      string `json:"name"`
}

// LibraryPanel defines model for LibraryPanel.
type LibraryPanel struct {
	// Panel description
	Description *string `json:"description,omitempty"`

	// Folder UID
	FolderUid *string                `json:"folderUid,omitempty"`
	Meta      *LibraryElementDTOMeta `json:"meta,omitempty"`

	// TODO: should be the same panel schema defined in dashboard
	// Typescript: Omit<Panel, 'gridPos' | 'id' | 'libraryPanel'>;
	Model map[string]interface{} `json:"model"`

	// Panel name (also saved in the model)
	Name string `json:"name"`

	// Dashboard version when this was saved (zero if unknown)
	SchemaVersion *int `json:"schemaVersion,omitempty"`

	// The panel type (from inside the model)
	Type string `json:"type"`

	// Library element UID
	Uid string `json:"uid"`

	// panel version, incremented each time the dashboard is updated.
	Version int64 `json:"version"`
}
