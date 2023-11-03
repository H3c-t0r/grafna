package rest

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	_ rest.Storage              = (*DualWriter)(nil)
	_ rest.Scoper               = (*DualWriter)(nil)
	_ rest.TableConvertor       = (*DualWriter)(nil)
	_ rest.CreaterUpdater       = (*DualWriter)(nil)
	_ rest.CollectionDeleter    = (*DualWriter)(nil)
	_ rest.GracefulDeleter      = (*DualWriter)(nil)
	_ rest.SingularNameProvider = (*DualWriter)(nil)
)

// Storage is a storage implementation that satisfies the same interfaces as genericregistry.Store.
type Storage interface {
	rest.Storage
	rest.StandardStorage
	rest.Scoper
	rest.TableConvertor
	rest.SingularNameProvider
	rest.Getter
}

// LegacyStorage is a storage implementation that writes to the Grafana SQL database.
type LegacyStorage interface {
	rest.Storage
	rest.Scoper
	rest.SingularNameProvider
	rest.TableConvertor
}

// DualWriter is a storage implementation that writes first to LegacyStorage and then to Storage.
// If writing to LegacyStorage fails, the write to Storage is skipped and the error is returned.
// Storage is used for all read operations.
//
// The LegacyStorage implementation must implement the following interfaces:
// - rest.Storage
// - rest.TableConvertor
// - rest.Scoper
// - rest.SingularNameProvider
//
// These interfaces are optional, but they all should be implemented to fully support dual writes:
// - rest.Creater
// - rest.Updater
// - rest.GracefulDeleter
// - rest.CollectionDeleter
type DualWriter struct {
	Storage
	legacy LegacyStorage
}

// NewDualWriter returns a new DualWriter.
func NewDualWriter(legacy LegacyStorage, storage Storage) *DualWriter {
	return &DualWriter{
		Storage: storage,
		legacy:  legacy,
	}
}

// Create overrides the default behavior of the Storage and writes to both the LegacyStorage and Storage.
func (d *DualWriter) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
	if legacy, ok := d.legacy.(rest.Creater); ok {
		created, err := legacy.Create(ctx, obj, createValidation, options)
		if err != nil {
			return nil, err
		}
		obj = created // write the updated version
		rsp, err := d.Storage.Create(ctx, obj, createValidation, options)
		if err != nil {
			// TODO
			fmt.Printf("TODO: delete %v\n", created)
		}
		return rsp, err
	}

	return d.Storage.Create(ctx, obj, createValidation, options)
}

// Update overrides the default behavior of the Storage and writes to both the LegacyStorage and Storage.
func (d *DualWriter) Update(ctx context.Context, name string, objInfo rest.UpdatedObjectInfo, createValidation rest.ValidateObjectFunc, updateValidation rest.ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error) {
	if legacy, ok := d.legacy.(rest.Updater); ok {
		// Will resource version checking work????
		old, err := d.Get(ctx, name, &metav1.GetOptions{})
		if err != nil {
			return nil, false, err
		}
		accessor, err := meta.Accessor(old)
		if err != nil {
			return nil, false, err
		}
		theRV := accessor.GetResourceVersion()
		theUID := accessor.GetUID()

		// Changes applied within standard storage
		updated, err := objInfo.UpdatedObject(ctx, old)
		if err != nil {
			return nil, false, err
		}

		fmt.Printf("AA:UPDATED: %+v\n", updated)
		accessor, err = meta.Accessor(updated)
		if err != nil {
			return nil, false, err
		}
		accessor.SetResourceVersion("") // remove it so it is not a constraint
		obj, created, err := legacy.Update(ctx, name, &updateWrapper{
			upstream: objInfo,
			updated:  updated, // returned as the object that will be updated
		}, createValidation, updateValidation, forceAllowCreate, options)
		if err != nil {
			return obj, created, err
		}

		fmt.Printf("AFTER UPDATE (OBJ): %+v\n", obj)

		accessor, err = meta.Accessor(obj)
		if err != nil {
			return nil, false, err
		}
		accessor.SetResourceVersion(theRV) // the original RV
		accessor.SetUID(theUID)
		objInfo = &updateWrapper{
			upstream: objInfo,
			updated:  obj, // returned as the object that will be updated
		}
	}

	return d.Storage.Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}

// Delete overrides the default behavior of the Storage and delete from both the LegacyStorage and Storage.
func (d *DualWriter) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	if legacy, ok := d.legacy.(rest.GracefulDeleter); ok {
		obj, async, err := legacy.Delete(ctx, name, deleteValidation, options)
		if err != nil {
			return obj, async, err
		}
	}

	return d.Storage.Delete(ctx, name, deleteValidation, options)
}

// DeleteCollection overrides the default behavior of the Storage and delete from both the LegacyStorage and Storage.
func (d *DualWriter) DeleteCollection(ctx context.Context, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error) {
	if legacy, ok := d.legacy.(rest.CollectionDeleter); ok {
		_, err := legacy.DeleteCollection(ctx, deleteValidation, options, listOptions)
		if err != nil {
			return nil, err
		}
	}

	return d.Storage.DeleteCollection(ctx, deleteValidation, options, listOptions)
}

type updateWrapper struct {
	upstream rest.UpdatedObjectInfo
	updated  runtime.Object
}

// Returns preconditions built from the updated object, if applicable.
// May return nil, or a preconditions object containing nil fields,
// if no preconditions can be determined from the updated object.
func (u *updateWrapper) Preconditions() *metav1.Preconditions {
	return u.upstream.Preconditions()
}

// UpdatedObject returns the updated object, given a context and old object.
// The only time an empty oldObj should be passed in is if a "create on update" is occurring (there is no oldObj).
func (u *updateWrapper) UpdatedObject(ctx context.Context, oldObj runtime.Object) (newObj runtime.Object, err error) {
	return u.updated, nil
}
