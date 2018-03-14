// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package accessibility

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility(in *jlexer.Lexer, out *GetPartialAXTreeReturns) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nodes":
			if in.IsNull() {
				in.Skip()
				out.Nodes = nil
			} else {
				in.Delim('[')
				if out.Nodes == nil {
					if !in.IsDelim(']') {
						out.Nodes = make([]*AXNode, 0, 8)
					} else {
						out.Nodes = []*AXNode{}
					}
				} else {
					out.Nodes = (out.Nodes)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *AXNode
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(AXNode)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Nodes = append(out.Nodes, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility(out *jwriter.Writer, in GetPartialAXTreeReturns) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Nodes) != 0 {
		const prefix string = ",\"nodes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v2, v3 := range in.Nodes {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetPartialAXTreeReturns) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetPartialAXTreeReturns) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetPartialAXTreeReturns) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetPartialAXTreeReturns) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility1(in *jlexer.Lexer, out *GetPartialAXTreeParams) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nodeId":
			(out.NodeID).UnmarshalEasyJSON(in)
		case "fetchRelatives":
			out.FetchRelatives = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility1(out *jwriter.Writer, in GetPartialAXTreeParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nodeId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.NodeID))
	}
	if in.FetchRelatives {
		const prefix string = ",\"fetchRelatives\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.FetchRelatives))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetPartialAXTreeParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetPartialAXTreeParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetPartialAXTreeParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetPartialAXTreeParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility1(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility2(in *jlexer.Lexer, out *AXValueSource) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			(out.Type).UnmarshalEasyJSON(in)
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(AXValue)
				}
				(*out.Value).UnmarshalEasyJSON(in)
			}
		case "attribute":
			out.Attribute = string(in.String())
		case "attributeValue":
			if in.IsNull() {
				in.Skip()
				out.AttributeValue = nil
			} else {
				if out.AttributeValue == nil {
					out.AttributeValue = new(AXValue)
				}
				(*out.AttributeValue).UnmarshalEasyJSON(in)
			}
		case "superseded":
			out.Superseded = bool(in.Bool())
		case "nativeSource":
			(out.NativeSource).UnmarshalEasyJSON(in)
		case "nativeSourceValue":
			if in.IsNull() {
				in.Skip()
				out.NativeSourceValue = nil
			} else {
				if out.NativeSourceValue == nil {
					out.NativeSourceValue = new(AXValue)
				}
				(*out.NativeSourceValue).UnmarshalEasyJSON(in)
			}
		case "invalid":
			out.Invalid = bool(in.Bool())
		case "invalidReason":
			out.InvalidReason = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility2(out *jwriter.Writer, in AXValueSource) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Type).MarshalEasyJSON(out)
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Value).MarshalEasyJSON(out)
	}
	if in.Attribute != "" {
		const prefix string = ",\"attribute\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Attribute))
	}
	if in.AttributeValue != nil {
		const prefix string = ",\"attributeValue\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.AttributeValue).MarshalEasyJSON(out)
	}
	if in.Superseded {
		const prefix string = ",\"superseded\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Superseded))
	}
	if in.NativeSource != "" {
		const prefix string = ",\"nativeSource\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.NativeSource).MarshalEasyJSON(out)
	}
	if in.NativeSourceValue != nil {
		const prefix string = ",\"nativeSourceValue\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.NativeSourceValue).MarshalEasyJSON(out)
	}
	if in.Invalid {
		const prefix string = ",\"invalid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Invalid))
	}
	if in.InvalidReason != "" {
		const prefix string = ",\"invalidReason\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.InvalidReason))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AXValueSource) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AXValueSource) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AXValueSource) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AXValueSource) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility2(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility3(in *jlexer.Lexer, out *AXValue) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			(out.Type).UnmarshalEasyJSON(in)
		case "value":
			(out.Value).UnmarshalEasyJSON(in)
		case "relatedNodes":
			if in.IsNull() {
				in.Skip()
				out.RelatedNodes = nil
			} else {
				in.Delim('[')
				if out.RelatedNodes == nil {
					if !in.IsDelim(']') {
						out.RelatedNodes = make([]*AXRelatedNode, 0, 8)
					} else {
						out.RelatedNodes = []*AXRelatedNode{}
					}
				} else {
					out.RelatedNodes = (out.RelatedNodes)[:0]
				}
				for !in.IsDelim(']') {
					var v4 *AXRelatedNode
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						if v4 == nil {
							v4 = new(AXRelatedNode)
						}
						(*v4).UnmarshalEasyJSON(in)
					}
					out.RelatedNodes = append(out.RelatedNodes, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "sources":
			if in.IsNull() {
				in.Skip()
				out.Sources = nil
			} else {
				in.Delim('[')
				if out.Sources == nil {
					if !in.IsDelim(']') {
						out.Sources = make([]*AXValueSource, 0, 8)
					} else {
						out.Sources = []*AXValueSource{}
					}
				} else {
					out.Sources = (out.Sources)[:0]
				}
				for !in.IsDelim(']') {
					var v5 *AXValueSource
					if in.IsNull() {
						in.Skip()
						v5 = nil
					} else {
						if v5 == nil {
							v5 = new(AXValueSource)
						}
						(*v5).UnmarshalEasyJSON(in)
					}
					out.Sources = append(out.Sources, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility3(out *jwriter.Writer, in AXValue) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Type).MarshalEasyJSON(out)
	}
	if (in.Value).IsDefined() {
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Value).MarshalEasyJSON(out)
	}
	if len(in.RelatedNodes) != 0 {
		const prefix string = ",\"relatedNodes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v6, v7 := range in.RelatedNodes {
				if v6 > 0 {
					out.RawByte(',')
				}
				if v7 == nil {
					out.RawString("null")
				} else {
					(*v7).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	if len(in.Sources) != 0 {
		const prefix string = ",\"sources\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v8, v9 := range in.Sources {
				if v8 > 0 {
					out.RawByte(',')
				}
				if v9 == nil {
					out.RawString("null")
				} else {
					(*v9).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AXValue) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AXValue) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AXValue) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AXValue) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility3(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility4(in *jlexer.Lexer, out *AXRelatedNode) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "backendDOMNodeId":
			(out.BackendDOMNodeID).UnmarshalEasyJSON(in)
		case "idref":
			out.Idref = string(in.String())
		case "text":
			out.Text = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility4(out *jwriter.Writer, in AXRelatedNode) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"backendDOMNodeId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.BackendDOMNodeID))
	}
	if in.Idref != "" {
		const prefix string = ",\"idref\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Idref))
	}
	if in.Text != "" {
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AXRelatedNode) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AXRelatedNode) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AXRelatedNode) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AXRelatedNode) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility4(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility5(in *jlexer.Lexer, out *AXProperty) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			(out.Name).UnmarshalEasyJSON(in)
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(AXValue)
				}
				(*out.Value).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility5(out *jwriter.Writer, in AXProperty) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Name).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Value == nil {
			out.RawString("null")
		} else {
			(*in.Value).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AXProperty) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AXProperty) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AXProperty) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AXProperty) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility5(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility6(in *jlexer.Lexer, out *AXNode) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "nodeId":
			out.NodeID = AXNodeID(in.String())
		case "ignored":
			out.Ignored = bool(in.Bool())
		case "ignoredReasons":
			if in.IsNull() {
				in.Skip()
				out.IgnoredReasons = nil
			} else {
				in.Delim('[')
				if out.IgnoredReasons == nil {
					if !in.IsDelim(']') {
						out.IgnoredReasons = make([]*AXProperty, 0, 8)
					} else {
						out.IgnoredReasons = []*AXProperty{}
					}
				} else {
					out.IgnoredReasons = (out.IgnoredReasons)[:0]
				}
				for !in.IsDelim(']') {
					var v10 *AXProperty
					if in.IsNull() {
						in.Skip()
						v10 = nil
					} else {
						if v10 == nil {
							v10 = new(AXProperty)
						}
						(*v10).UnmarshalEasyJSON(in)
					}
					out.IgnoredReasons = append(out.IgnoredReasons, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "role":
			if in.IsNull() {
				in.Skip()
				out.Role = nil
			} else {
				if out.Role == nil {
					out.Role = new(AXValue)
				}
				(*out.Role).UnmarshalEasyJSON(in)
			}
		case "name":
			if in.IsNull() {
				in.Skip()
				out.Name = nil
			} else {
				if out.Name == nil {
					out.Name = new(AXValue)
				}
				(*out.Name).UnmarshalEasyJSON(in)
			}
		case "description":
			if in.IsNull() {
				in.Skip()
				out.Description = nil
			} else {
				if out.Description == nil {
					out.Description = new(AXValue)
				}
				(*out.Description).UnmarshalEasyJSON(in)
			}
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(AXValue)
				}
				(*out.Value).UnmarshalEasyJSON(in)
			}
		case "properties":
			if in.IsNull() {
				in.Skip()
				out.Properties = nil
			} else {
				in.Delim('[')
				if out.Properties == nil {
					if !in.IsDelim(']') {
						out.Properties = make([]*AXProperty, 0, 8)
					} else {
						out.Properties = []*AXProperty{}
					}
				} else {
					out.Properties = (out.Properties)[:0]
				}
				for !in.IsDelim(']') {
					var v11 *AXProperty
					if in.IsNull() {
						in.Skip()
						v11 = nil
					} else {
						if v11 == nil {
							v11 = new(AXProperty)
						}
						(*v11).UnmarshalEasyJSON(in)
					}
					out.Properties = append(out.Properties, v11)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "childIds":
			if in.IsNull() {
				in.Skip()
				out.ChildIds = nil
			} else {
				in.Delim('[')
				if out.ChildIds == nil {
					if !in.IsDelim(']') {
						out.ChildIds = make([]AXNodeID, 0, 4)
					} else {
						out.ChildIds = []AXNodeID{}
					}
				} else {
					out.ChildIds = (out.ChildIds)[:0]
				}
				for !in.IsDelim(']') {
					var v12 AXNodeID
					v12 = AXNodeID(in.String())
					out.ChildIds = append(out.ChildIds, v12)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "backendDOMNodeId":
			(out.BackendDOMNodeID).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility6(out *jwriter.Writer, in AXNode) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"nodeId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.NodeID))
	}
	{
		const prefix string = ",\"ignored\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Ignored))
	}
	if len(in.IgnoredReasons) != 0 {
		const prefix string = ",\"ignoredReasons\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v13, v14 := range in.IgnoredReasons {
				if v13 > 0 {
					out.RawByte(',')
				}
				if v14 == nil {
					out.RawString("null")
				} else {
					(*v14).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	if in.Role != nil {
		const prefix string = ",\"role\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Role).MarshalEasyJSON(out)
	}
	if in.Name != nil {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Name).MarshalEasyJSON(out)
	}
	if in.Description != nil {
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Description).MarshalEasyJSON(out)
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Value).MarshalEasyJSON(out)
	}
	if len(in.Properties) != 0 {
		const prefix string = ",\"properties\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v15, v16 := range in.Properties {
				if v15 > 0 {
					out.RawByte(',')
				}
				if v16 == nil {
					out.RawString("null")
				} else {
					(*v16).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	if len(in.ChildIds) != 0 {
		const prefix string = ",\"childIds\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v17, v18 := range in.ChildIds {
				if v17 > 0 {
					out.RawByte(',')
				}
				out.String(string(v18))
			}
			out.RawByte(']')
		}
	}
	if in.BackendDOMNodeID != 0 {
		const prefix string = ",\"backendDOMNodeId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.BackendDOMNodeID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v AXNode) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AXNode) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoAccessibility6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *AXNode) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AXNode) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoAccessibility6(l, v)
}
