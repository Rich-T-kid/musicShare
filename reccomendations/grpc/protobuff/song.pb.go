// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.12.4
// source: song.proto

package pb

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SongRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // uuid Associated will user in MongoDB
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SongRequest) Reset() {
	*x = SongRequest{}
	mi := &file_song_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SongRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SongRequest) ProtoMessage() {}

func (x *SongRequest) ProtoReflect() protoreflect.Message {
	mi := &file_song_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SongRequest.ProtoReflect.Descriptor instead.
func (*SongRequest) Descriptor() ([]byte, []int) {
	return file_song_proto_rawDescGZIP(), []int{0}
}

func (x *SongRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type SongResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Songs         []*SongBody            `protobuf:"bytes,1,rep,name=songs,proto3" json:"songs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SongResponse) Reset() {
	*x = SongResponse{}
	mi := &file_song_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SongResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SongResponse) ProtoMessage() {}

func (x *SongResponse) ProtoReflect() protoreflect.Message {
	mi := &file_song_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SongResponse.ProtoReflect.Descriptor instead.
func (*SongResponse) Descriptor() ([]byte, []int) {
	return file_song_proto_rawDescGZIP(), []int{1}
}

func (x *SongResponse) GetSongs() []*SongBody {
	if x != nil {
		return x.Songs
	}
	return nil
}

type SongBody struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Artist        string                 `protobuf:"bytes,2,opt,name=artist,proto3" json:"artist,omitempty"`
	SongUri       string                 `protobuf:"bytes,3,opt,name=song_uri,json=songUri,proto3" json:"song_uri,omitempty"` // for now this will only be spotify but just in case we add others later
	Rank          uint32                 `protobuf:"varint,4,opt,name=rank,proto3" json:"rank,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SongBody) Reset() {
	*x = SongBody{}
	mi := &file_song_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SongBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SongBody) ProtoMessage() {}

func (x *SongBody) ProtoReflect() protoreflect.Message {
	mi := &file_song_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SongBody.ProtoReflect.Descriptor instead.
func (*SongBody) Descriptor() ([]byte, []int) {
	return file_song_proto_rawDescGZIP(), []int{2}
}

func (x *SongBody) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SongBody) GetArtist() string {
	if x != nil {
		return x.Artist
	}
	return ""
}

func (x *SongBody) GetSongUri() string {
	if x != nil {
		return x.SongUri
	}
	return ""
}

func (x *SongBody) GetRank() uint32 {
	if x != nil {
		return x.Rank
	}
	return 0
}

var File_song_proto protoreflect.FileDescriptor

var file_song_proto_rawDesc = string([]byte{
	0x0a, 0x0a, 0x73, 0x6f, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x73, 0x6f,
	0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x26, 0x0a, 0x0b, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x42, 0x0a, 0x0c, 0x53, 0x6f, 0x6e, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x6f, 0x6e, 0x67,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73, 0x6f, 0x6e, 0x67, 0x52, 0x65,
	0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x73, 0x6f, 0x6e,
	0x67, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x05, 0x73, 0x6f, 0x6e, 0x67, 0x73, 0x22, 0x65, 0x0a, 0x08,
	0x73, 0x6f, 0x6e, 0x67, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x72, 0x74, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x72,
	0x74, 0x69, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x6f, 0x6e, 0x67, 0x5f, 0x75, 0x72, 0x69,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x6f, 0x6e, 0x67, 0x55, 0x72, 0x69, 0x12,
	0x12, 0x0a, 0x04, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x72,
	0x61, 0x6e, 0x6b, 0x32, 0x5d, 0x0a, 0x0b, 0x53, 0x6f, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x4e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x53, 0x6f, 0x6e, 0x67, 0x12, 0x1f, 0x2e,
	0x73, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20,
	0x2e, 0x73, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x64, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x6f, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x03, 0x5a, 0x01, 0x2e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_song_proto_rawDescOnce sync.Once
	file_song_proto_rawDescData []byte
)

func file_song_proto_rawDescGZIP() []byte {
	file_song_proto_rawDescOnce.Do(func() {
		file_song_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_song_proto_rawDesc), len(file_song_proto_rawDesc)))
	})
	return file_song_proto_rawDescData
}

var file_song_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_song_proto_goTypes = []any{
	(*SongRequest)(nil),  // 0: songRecommendation.SongRequest
	(*SongResponse)(nil), // 1: songRecommendation.SongResponse
	(*SongBody)(nil),     // 2: songRecommendation.songBody
}
var file_song_proto_depIdxs = []int32{
	2, // 0: songRecommendation.SongResponse.songs:type_name -> songRecommendation.songBody
	0, // 1: songRecommendation.SongService.GetSong:input_type -> songRecommendation.SongRequest
	1, // 2: songRecommendation.SongService.GetSong:output_type -> songRecommendation.SongResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_song_proto_init() }
func file_song_proto_init() {
	if File_song_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_song_proto_rawDesc), len(file_song_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_song_proto_goTypes,
		DependencyIndexes: file_song_proto_depIdxs,
		MessageInfos:      file_song_proto_msgTypes,
	}.Build()
	File_song_proto = out.File
	file_song_proto_goTypes = nil
	file_song_proto_depIdxs = nil
}
