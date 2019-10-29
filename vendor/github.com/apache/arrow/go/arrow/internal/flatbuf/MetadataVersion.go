// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuf

type MetadataVersion = int16
const (
	/// 0.1.0
	MetadataVersionV1 MetadataVersion = 0
	/// 0.2.0
	MetadataVersionV2 MetadataVersion = 1
	/// 0.3.0 -> 0.7.1
	MetadataVersionV3 MetadataVersion = 2
	/// >= 0.8.0
	MetadataVersionV4 MetadataVersion = 3
)

var EnumNamesMetadataVersion = map[MetadataVersion]string{
	MetadataVersionV1:"V1",
	MetadataVersionV2:"V2",
	MetadataVersionV3:"V3",
	MetadataVersionV4:"V4",
}

