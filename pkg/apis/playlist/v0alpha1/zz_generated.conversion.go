//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v0alpha1

// import (
// 	unsafe "unsafe"

// 	playlist "github.com/grafana/grafana/pkg/apis/playlist"
// 	conversion "k8s.io/apimachinery/pkg/conversion"
// 	runtime "k8s.io/apimachinery/pkg/runtime"
// )

// // RegisterConversions adds conversion functions to the given scheme.
// // Public to allow building arbitrary schemes.
// func RegisterConversions(s *runtime.Scheme) error {
// 	if err := s.AddGeneratedConversionFunc((*Item)(nil), (*playlist.Item)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_v0alpha1_Item_To_playlist_Item(a.(*Item), b.(*playlist.Item), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*playlist.Item)(nil), (*Item)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_playlist_Item_To_v0alpha1_Item(a.(*playlist.Item), b.(*Item), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*Playlist)(nil), (*playlist.Playlist)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_v0alpha1_Playlist_To_playlist_Playlist(a.(*Playlist), b.(*playlist.Playlist), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*playlist.Playlist)(nil), (*Playlist)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_playlist_Playlist_To_v0alpha1_Playlist(a.(*playlist.Playlist), b.(*Playlist), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*PlaylistList)(nil), (*playlist.PlaylistList)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_v0alpha1_PlaylistList_To_playlist_PlaylistList(a.(*PlaylistList), b.(*playlist.PlaylistList), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*playlist.PlaylistList)(nil), (*PlaylistList)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_playlist_PlaylistList_To_v0alpha1_PlaylistList(a.(*playlist.PlaylistList), b.(*PlaylistList), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*Spec)(nil), (*playlist.Spec)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_v0alpha1_Spec_To_playlist_Spec(a.(*Spec), b.(*playlist.Spec), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	if err := s.AddGeneratedConversionFunc((*playlist.Spec)(nil), (*Spec)(nil), func(a, b interface{}, scope conversion.Scope) error {
// 		return Convert_playlist_Spec_To_v0alpha1_Spec(a.(*playlist.Spec), b.(*Spec), scope)
// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func autoConvert_v0alpha1_Item_To_playlist_Item(in *Item, out *playlist.Item, s conversion.Scope) error {
// 	out.Type = playlist.ItemType(in.Type)
// 	out.Value = in.Value
// 	return nil
// }

// // Convert_v0alpha1_Item_To_playlist_Item is an autogenerated conversion function.
// func Convert_v0alpha1_Item_To_playlist_Item(in *Item, out *playlist.Item, s conversion.Scope) error {
// 	return autoConvert_v0alpha1_Item_To_playlist_Item(in, out, s)
// }

// func autoConvert_playlist_Item_To_v0alpha1_Item(in *playlist.Item, out *Item, s conversion.Scope) error {
// 	out.Type = ItemType(in.Type)
// 	out.Value = in.Value
// 	return nil
// }

// // Convert_playlist_Item_To_v0alpha1_Item is an autogenerated conversion function.
// func Convert_playlist_Item_To_v0alpha1_Item(in *playlist.Item, out *Item, s conversion.Scope) error {
// 	return autoConvert_playlist_Item_To_v0alpha1_Item(in, out, s)
// }

// func autoConvert_v0alpha1_Playlist_To_playlist_Playlist(in *Playlist, out *playlist.Playlist, s conversion.Scope) error {
// 	out.ObjectMeta = in.ObjectMeta
// 	if err := Convert_v0alpha1_Spec_To_playlist_Spec(&in.Spec, &out.Spec, s); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Convert_v0alpha1_Playlist_To_playlist_Playlist is an autogenerated conversion function.
// func Convert_v0alpha1_Playlist_To_playlist_Playlist(in *Playlist, out *playlist.Playlist, s conversion.Scope) error {
// 	return autoConvert_v0alpha1_Playlist_To_playlist_Playlist(in, out, s)
// }

// func autoConvert_playlist_Playlist_To_v0alpha1_Playlist(in *playlist.Playlist, out *Playlist, s conversion.Scope) error {
// 	out.ObjectMeta = in.ObjectMeta
// 	if err := Convert_playlist_Spec_To_v0alpha1_Spec(&in.Spec, &out.Spec, s); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Convert_playlist_Playlist_To_v0alpha1_Playlist is an autogenerated conversion function.
// func Convert_playlist_Playlist_To_v0alpha1_Playlist(in *playlist.Playlist, out *Playlist, s conversion.Scope) error {
// 	return autoConvert_playlist_Playlist_To_v0alpha1_Playlist(in, out, s)
// }

// func autoConvert_v0alpha1_PlaylistList_To_playlist_PlaylistList(in *PlaylistList, out *playlist.PlaylistList, s conversion.Scope) error {
// 	out.ListMeta = in.ListMeta
// 	out.Items = *(*[]playlist.Playlist)(unsafe.Pointer(&in.Items))
// 	return nil
// }

// // Convert_v0alpha1_PlaylistList_To_playlist_PlaylistList is an autogenerated conversion function.
// func Convert_v0alpha1_PlaylistList_To_playlist_PlaylistList(in *PlaylistList, out *playlist.PlaylistList, s conversion.Scope) error {
// 	return autoConvert_v0alpha1_PlaylistList_To_playlist_PlaylistList(in, out, s)
// }

// func autoConvert_playlist_PlaylistList_To_v0alpha1_PlaylistList(in *playlist.PlaylistList, out *PlaylistList, s conversion.Scope) error {
// 	out.ListMeta = in.ListMeta
// 	out.Items = *(*[]Playlist)(unsafe.Pointer(&in.Items))
// 	return nil
// }

// // Convert_playlist_PlaylistList_To_v0alpha1_PlaylistList is an autogenerated conversion function.
// func Convert_playlist_PlaylistList_To_v0alpha1_PlaylistList(in *playlist.PlaylistList, out *PlaylistList, s conversion.Scope) error {
// 	return autoConvert_playlist_PlaylistList_To_v0alpha1_PlaylistList(in, out, s)
// }

// func autoConvert_v0alpha1_Spec_To_playlist_Spec(in *Spec, out *playlist.Spec, s conversion.Scope) error {
// 	out.Title = in.Title
// 	out.Interval = in.Interval
// 	out.Items = *(*[]playlist.Item)(unsafe.Pointer(&in.Items))
// 	return nil
// }

// // Convert_v0alpha1_Spec_To_playlist_Spec is an autogenerated conversion function.
// func Convert_v0alpha1_Spec_To_playlist_Spec(in *Spec, out *playlist.Spec, s conversion.Scope) error {
// 	return autoConvert_v0alpha1_Spec_To_playlist_Spec(in, out, s)
// }

// func autoConvert_playlist_Spec_To_v0alpha1_Spec(in *playlist.Spec, out *Spec, s conversion.Scope) error {
// 	out.Title = in.Title
// 	out.Interval = in.Interval
// 	out.Items = *(*[]Item)(unsafe.Pointer(&in.Items))
// 	return nil
// }

// // Convert_playlist_Spec_To_v0alpha1_Spec is an autogenerated conversion function.
// func Convert_playlist_Spec_To_v0alpha1_Spec(in *playlist.Spec, out *Spec, s conversion.Scope) error {
// 	return autoConvert_playlist_Spec_To_v0alpha1_Spec(in, out, s)
// }
