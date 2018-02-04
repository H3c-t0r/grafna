// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package domdebugger

import (
	json "encoding/json"
	runtime "github.com/chromedp/cdproto/runtime"
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

func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger(in *jlexer.Lexer, out *SetXHRBreakpointParams) {
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
		case "url":
			out.URL = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger(out *jwriter.Writer, in SetXHRBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetXHRBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetXHRBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetXHRBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetXHRBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger1(in *jlexer.Lexer, out *SetInstrumentationBreakpointParams) {
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
		case "eventName":
			out.EventName = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger1(out *jwriter.Writer, in SetInstrumentationBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"eventName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.EventName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetInstrumentationBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetInstrumentationBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetInstrumentationBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetInstrumentationBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger1(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger2(in *jlexer.Lexer, out *SetEventListenerBreakpointParams) {
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
		case "eventName":
			out.EventName = string(in.String())
		case "targetName":
			out.TargetName = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger2(out *jwriter.Writer, in SetEventListenerBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"eventName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.EventName))
	}
	if in.TargetName != "" {
		const prefix string = ",\"targetName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TargetName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetEventListenerBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetEventListenerBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetEventListenerBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetEventListenerBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger2(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger3(in *jlexer.Lexer, out *SetDOMBreakpointParams) {
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
		case "type":
			(out.Type).UnmarshalEasyJSON(in)
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger3(out *jwriter.Writer, in SetDOMBreakpointParams) {
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SetDOMBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SetDOMBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SetDOMBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SetDOMBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger3(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger4(in *jlexer.Lexer, out *RemoveXHRBreakpointParams) {
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
		case "url":
			out.URL = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger4(out *jwriter.Writer, in RemoveXHRBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"url\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.URL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoveXHRBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoveXHRBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoveXHRBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoveXHRBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger4(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger5(in *jlexer.Lexer, out *RemoveInstrumentationBreakpointParams) {
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
		case "eventName":
			out.EventName = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger5(out *jwriter.Writer, in RemoveInstrumentationBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"eventName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.EventName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoveInstrumentationBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoveInstrumentationBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoveInstrumentationBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoveInstrumentationBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger5(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger6(in *jlexer.Lexer, out *RemoveEventListenerBreakpointParams) {
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
		case "eventName":
			out.EventName = string(in.String())
		case "targetName":
			out.TargetName = string(in.String())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger6(out *jwriter.Writer, in RemoveEventListenerBreakpointParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"eventName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.EventName))
	}
	if in.TargetName != "" {
		const prefix string = ",\"targetName\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TargetName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoveEventListenerBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoveEventListenerBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoveEventListenerBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoveEventListenerBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger6(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger7(in *jlexer.Lexer, out *RemoveDOMBreakpointParams) {
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
		case "type":
			(out.Type).UnmarshalEasyJSON(in)
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger7(out *jwriter.Writer, in RemoveDOMBreakpointParams) {
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
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RemoveDOMBreakpointParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RemoveDOMBreakpointParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RemoveDOMBreakpointParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RemoveDOMBreakpointParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger7(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger8(in *jlexer.Lexer, out *GetEventListenersReturns) {
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
		case "listeners":
			if in.IsNull() {
				in.Skip()
				out.Listeners = nil
			} else {
				in.Delim('[')
				if out.Listeners == nil {
					if !in.IsDelim(']') {
						out.Listeners = make([]*EventListener, 0, 8)
					} else {
						out.Listeners = []*EventListener{}
					}
				} else {
					out.Listeners = (out.Listeners)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *EventListener
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(EventListener)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Listeners = append(out.Listeners, v1)
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger8(out *jwriter.Writer, in GetEventListenersReturns) {
	out.RawByte('{')
	first := true
	_ = first
	if len(in.Listeners) != 0 {
		const prefix string = ",\"listeners\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v2, v3 := range in.Listeners {
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
func (v GetEventListenersReturns) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetEventListenersReturns) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetEventListenersReturns) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetEventListenersReturns) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger8(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger9(in *jlexer.Lexer, out *GetEventListenersParams) {
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
		case "objectId":
			out.ObjectID = runtime.RemoteObjectID(in.String())
		case "depth":
			out.Depth = int64(in.Int64())
		case "pierce":
			out.Pierce = bool(in.Bool())
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger9(out *jwriter.Writer, in GetEventListenersParams) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"objectId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ObjectID))
	}
	if in.Depth != 0 {
		const prefix string = ",\"depth\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Depth))
	}
	if in.Pierce {
		const prefix string = ",\"pierce\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Pierce))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetEventListenersParams) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetEventListenersParams) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetEventListenersParams) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetEventListenersParams) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger9(l, v)
}
func easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger10(in *jlexer.Lexer, out *EventListener) {
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
			out.Type = string(in.String())
		case "useCapture":
			out.UseCapture = bool(in.Bool())
		case "passive":
			out.Passive = bool(in.Bool())
		case "once":
			out.Once = bool(in.Bool())
		case "scriptId":
			out.ScriptID = runtime.ScriptID(in.String())
		case "lineNumber":
			out.LineNumber = int64(in.Int64())
		case "columnNumber":
			out.ColumnNumber = int64(in.Int64())
		case "handler":
			if in.IsNull() {
				in.Skip()
				out.Handler = nil
			} else {
				if out.Handler == nil {
					out.Handler = new(runtime.RemoteObject)
				}
				(*out.Handler).UnmarshalEasyJSON(in)
			}
		case "originalHandler":
			if in.IsNull() {
				in.Skip()
				out.OriginalHandler = nil
			} else {
				if out.OriginalHandler == nil {
					out.OriginalHandler = new(runtime.RemoteObject)
				}
				(*out.OriginalHandler).UnmarshalEasyJSON(in)
			}
		case "backendNodeId":
			(out.BackendNodeID).UnmarshalEasyJSON(in)
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
func easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger10(out *jwriter.Writer, in EventListener) {
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
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"useCapture\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.UseCapture))
	}
	{
		const prefix string = ",\"passive\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Passive))
	}
	{
		const prefix string = ",\"once\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Once))
	}
	{
		const prefix string = ",\"scriptId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ScriptID))
	}
	{
		const prefix string = ",\"lineNumber\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.LineNumber))
	}
	{
		const prefix string = ",\"columnNumber\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.ColumnNumber))
	}
	if in.Handler != nil {
		const prefix string = ",\"handler\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Handler).MarshalEasyJSON(out)
	}
	if in.OriginalHandler != nil {
		const prefix string = ",\"originalHandler\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.OriginalHandler).MarshalEasyJSON(out)
	}
	if in.BackendNodeID != 0 {
		const prefix string = ",\"backendNodeId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.BackendNodeID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventListener) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventListener) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC5a4559bEncodeGithubComChromedpCdprotoDomdebugger10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventListener) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventListener) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC5a4559bDecodeGithubComChromedpCdprotoDomdebugger10(l, v)
}
