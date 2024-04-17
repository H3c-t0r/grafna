package rest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/apis/example"
)

const kind = "dummy"

var exampleObj = &example.Pod{TypeMeta: metav1.TypeMeta{}, ObjectMeta: metav1.ObjectMeta{Name: "foo", ResourceVersion: "1"}, Spec: example.PodSpec{}, Status: example.PodStatus{}}
var anotherObj = &example.Pod{TypeMeta: metav1.TypeMeta{}, ObjectMeta: metav1.ObjectMeta{Name: "bar", ResourceVersion: "2"}, Spec: example.PodSpec{}, Status: example.PodStatus{}}
var failingObj = &example.Pod{TypeMeta: metav1.TypeMeta{}, ObjectMeta: metav1.ObjectMeta{Name: "object-fail", ResourceVersion: "2"}, Spec: example.PodSpec{}, Status: example.PodStatus{}}

func TestMode1_Create(t *testing.T) {
	type testCase struct {
		name           string
		input          runtime.Object
		setupLegacyFn  func(m *mock.Mock, input runtime.Object)
		setupStorageFn func(m *mock.Mock, input runtime.Object)
		wantErr        bool
	}
	tests :=
		[]testCase{
			{
				name:  "creating an object only in the legacy store",
				input: exampleObj,
				setupLegacyFn: func(m *mock.Mock, input runtime.Object) {
					m.On("Create", context.Background(), input, mock.Anything, mock.Anything).Return(exampleObj, nil)
				},
				setupStorageFn: func(m *mock.Mock, input runtime.Object) {
					m.On("Create", context.Background(), anotherObj, mock.Anything, mock.Anything).Return(anotherObj, nil)
				},
			},
			{
				name:  "error when creating object in the legacy store fails",
				input: failingObj,
				setupLegacyFn: func(m *mock.Mock, input runtime.Object) {
					m.On("Create", context.Background(), failingObj, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
				},
				wantErr: true,
			},
		}

	for _, tt := range tests {
		l := (LegacyStorage)(nil)
		s := (Storage)(nil)
		m := &mock.Mock{}

		ls := legacyStoreMock{m, l}
		us := storageMock{m, s}

		if tt.setupLegacyFn != nil {
			tt.setupLegacyFn(m, tt.input)
		}
		if tt.setupStorageFn != nil {
			tt.setupStorageFn(m, tt.input)
		}

		dw := SelectDualWriter(Mode1, ls, us)

		obj, err := dw.Create(context.Background(), tt.input, func(context.Context, runtime.Object) error { return nil }, &metav1.CreateOptions{})

		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		us.AssertNotCalled(t, "Create", context.Background(), tt.input, func(context.Context, runtime.Object) error { return nil }, &metav1.CreateOptions{})

		assert.Equal(t, obj, exampleObj)
		assert.NotEqual(t, obj, anotherObj)
	}
}

func TestMode1_Get(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		setupLegacyFn  func(m *mock.Mock, name string)
		setupStorageFn func(m *mock.Mock, name string)
		wantErr        bool
	}
	tests :=
		[]testCase{
			{
				name:  "get an object only in the legacy store",
				input: "foo",
				setupLegacyFn: func(m *mock.Mock, name string) {
					m.On("Get", context.Background(), name, mock.Anything).Return(exampleObj, nil)
				},
				setupStorageFn: func(m *mock.Mock, name string) {
					m.On("Get", context.Background(), name, mock.Anything).Return(anotherObj, nil)
				},
			},
			{
				name:  "error when getting an object in the legacy store fails",
				input: "object-fail",
				setupLegacyFn: func(m *mock.Mock, name string) {
					m.On("Get", context.Background(), name, mock.Anything).Return(nil, errors.New("error"))
				},
				wantErr: true,
			},
		}

	for _, tt := range tests {
		l := (LegacyStorage)(nil)
		s := (Storage)(nil)
		m := &mock.Mock{}

		ls := legacyStoreMock{m, l}
		us := storageMock{m, s}

		if tt.setupLegacyFn != nil {
			tt.setupLegacyFn(m, tt.input)
		}
		if tt.setupStorageFn != nil {
			tt.setupStorageFn(m, tt.input)
		}

		dw := SelectDualWriter(Mode1, ls, us)

		obj, err := dw.Get(context.Background(), tt.input, &metav1.GetOptions{})

		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		us.AssertNotCalled(t, "Get", context.Background(), tt.name, &metav1.GetOptions{})

		assert.Equal(t, obj, exampleObj)
		assert.NotEqual(t, obj, anotherObj)
	}
}

func TestMode1_List(t *testing.T) {
	type testCase struct {
		name           string
		setupLegacyFn  func(m *mock.Mock)
		setupStorageFn func(m *mock.Mock)
		wantErr        bool
	}
	tests :=
		[]testCase{
			{
				name: "error when listing an object in the legacy store is not implemented",
				setupLegacyFn: func(m *mock.Mock) {
					m.On("List", context.Background(), mock.Anything).Return(nil, errors.New("error"))
				},
			},
			// TODO: legacy list is missing
		}

	for _, tt := range tests {
		l := (LegacyStorage)(nil)
		s := (Storage)(nil)
		m := &mock.Mock{}

		ls := legacyStoreMock{m, l}
		us := storageMock{m, s}

		if tt.setupLegacyFn != nil {
			tt.setupLegacyFn(m)
		}
		if tt.setupStorageFn != nil {
			tt.setupStorageFn(m)
		}

		dw := SelectDualWriter(Mode1, ls, us)

		_, err := dw.List(context.Background(), &metainternalversion.ListOptions{})

		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		us.AssertNotCalled(t, "List", context.Background(), &metainternalversion.ListOptions{})
	}
}

func TestMode1_Delete(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		setupLegacyFn  func(m *mock.Mock, name string)
		setupStorageFn func(m *mock.Mock, name string)
		wantErr        bool
	}
	tests :=
		[]testCase{
			{
				name:  "deleting an object in the legacy store",
				input: "foo",
				setupLegacyFn: func(m *mock.Mock, name string) {
					m.On("Delete", context.Background(), name, mock.Anything, mock.Anything).Return(exampleObj, false, nil)
				},
			},
			{
				name:  "error when deleting an object in the legacy store",
				input: "object-fail",
				setupLegacyFn: func(m *mock.Mock, name string) {
					m.On("Delete", context.Background(), name, mock.Anything, mock.Anything).Return(nil, false, errors.New("error"))
				},
				wantErr: true,
			},
		}

	for _, tt := range tests {
		l := (LegacyStorage)(nil)
		s := (Storage)(nil)
		m := &mock.Mock{}

		ls := legacyStoreMock{m, l}
		us := storageMock{m, s}

		if tt.setupLegacyFn != nil {
			tt.setupLegacyFn(m, tt.input)
		}
		if tt.setupStorageFn != nil {
			tt.setupStorageFn(m, tt.input)
		}

		dw := SelectDualWriter(Mode1, ls, us)

		obj, _, err := dw.Delete(context.Background(), tt.input, func(ctx context.Context, obj runtime.Object) error { return nil }, &metav1.DeleteOptions{})

		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		us.AssertNotCalled(t, "Delete", context.Background(), tt.input, func(ctx context.Context, obj runtime.Object) error { return nil }, &metav1.DeleteOptions{})
		assert.Equal(t, obj, exampleObj)
		assert.NotEqual(t, obj, anotherObj)
	}
}

func TestMode1_DeleteCollection(t *testing.T) {
	type testCase struct {
		name           string
		input          *metav1.DeleteOptions
		setupLegacyFn  func(m *mock.Mock, input *metav1.DeleteOptions)
		setupStorageFn func(m *mock.Mock, input *metav1.DeleteOptions)
		wantErr        bool
	}
	tests :=
		[]testCase{
			{
				name:  "deleting a collection in the legacy store",
				input: &metav1.DeleteOptions{TypeMeta: metav1.TypeMeta{Kind: "foo"}},
				setupLegacyFn: func(m *mock.Mock, input *metav1.DeleteOptions) {
					m.On("DeleteCollection", context.Background(), mock.Anything, input, mock.Anything).Return(exampleObj, nil)
				},
			},
			{
				name:  "error deleting a collection in the legacy store",
				input: &metav1.DeleteOptions{TypeMeta: metav1.TypeMeta{Kind: "fail"}},
				setupLegacyFn: func(m *mock.Mock, input *metav1.DeleteOptions) {
					m.On("DeleteCollection", context.Background(), mock.Anything, input, mock.Anything).Return(nil, errors.New("error"))
				},
				wantErr: true,
			},
		}

	for _, tt := range tests {
		l := (LegacyStorage)(nil)
		s := (Storage)(nil)
		m := &mock.Mock{}

		ls := legacyStoreMock{m, l}
		us := storageMock{m, s}

		if tt.setupLegacyFn != nil {
			tt.setupLegacyFn(m, tt.input)
		}
		if tt.setupStorageFn != nil {
			tt.setupStorageFn(m, tt.input)
		}

		dw := SelectDualWriter(Mode1, ls, us)

		obj, err := dw.DeleteCollection(context.Background(), func(ctx context.Context, obj runtime.Object) error { return nil }, tt.input, &metainternalversion.ListOptions{})

		if tt.wantErr {
			assert.Error(t, err)
			continue
		}

		us.AssertNotCalled(t, "DeleteCollection", context.Background(), tt.input, func(ctx context.Context, obj runtime.Object) error { return nil }, &metav1.DeleteOptions{})
		assert.Equal(t, obj, exampleObj)
		assert.NotEqual(t, obj, anotherObj)
	}
}
