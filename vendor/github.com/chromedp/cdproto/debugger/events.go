package debugger

// Code generated by chromedp-gen. DO NOT EDIT.

import (
	"github.com/chromedp/cdproto/runtime"
	"github.com/mailru/easyjson"
)

// EventBreakpointResolved fired when breakpoint is resolved to an actual
// script and location.
type EventBreakpointResolved struct {
	BreakpointID BreakpointID `json:"breakpointId"` // Breakpoint unique identifier.
	Location     *Location    `json:"location"`     // Actual breakpoint location.
}

// EventPaused fired when the virtual machine stopped on breakpoint or
// exception or any other stop criteria.
type EventPaused struct {
	CallFrames            []*CallFrame          `json:"callFrames"` // Call stack the virtual machine stopped on.
	Reason                PausedReason          `json:"reason"`     // Pause reason.
	Data                  easyjson.RawMessage   `json:"data,omitempty"`
	HitBreakpoints        []string              `json:"hitBreakpoints,omitempty"`        // Hit breakpoints IDs
	AsyncStackTrace       *runtime.StackTrace   `json:"asyncStackTrace,omitempty"`       // Async stack trace, if any.
	AsyncStackTraceID     *runtime.StackTraceID `json:"asyncStackTraceId,omitempty"`     // Async stack trace, if any.
	AsyncCallStackTraceID *runtime.StackTraceID `json:"asyncCallStackTraceId,omitempty"` // Just scheduled async call will have this stack trace as parent stack during async execution. This field is available only after Debugger.stepInto call with breakOnAsynCall flag.
}

// EventResumed fired when the virtual machine resumed execution.
type EventResumed struct{}

// EventScriptFailedToParse fired when virtual machine fails to parse the
// script.
type EventScriptFailedToParse struct {
	ScriptID                runtime.ScriptID           `json:"scriptId"`           // Identifier of the script parsed.
	URL                     string                     `json:"url"`                // URL or name of the script parsed (if any).
	StartLine               int64                      `json:"startLine"`          // Line offset of the script within the resource with given URL (for script tags).
	StartColumn             int64                      `json:"startColumn"`        // Column offset of the script within the resource with given URL.
	EndLine                 int64                      `json:"endLine"`            // Last line of the script.
	EndColumn               int64                      `json:"endColumn"`          // Length of the last line of the script.
	ExecutionContextID      runtime.ExecutionContextID `json:"executionContextId"` // Specifies script creation context.
	Hash                    string                     `json:"hash"`               // Content hash of the script.
	ExecutionContextAuxData easyjson.RawMessage        `json:"executionContextAuxData,omitempty"`
	SourceMapURL            string                     `json:"sourceMapURL,omitempty"` // URL of source map associated with script (if any).
	HasSourceURL            bool                       `json:"hasSourceURL,omitempty"` // True, if this script has sourceURL.
	IsModule                bool                       `json:"isModule,omitempty"`     // True, if this script is ES6 module.
	Length                  int64                      `json:"length,omitempty"`       // This script length.
	StackTrace              *runtime.StackTrace        `json:"stackTrace,omitempty"`   // JavaScript top stack frame of where the script parsed event was triggered if available.
}

// EventScriptParsed fired when virtual machine parses script. This event is
// also fired for all known and uncollected scripts upon enabling debugger.
type EventScriptParsed struct {
	ScriptID                runtime.ScriptID           `json:"scriptId"`           // Identifier of the script parsed.
	URL                     string                     `json:"url"`                // URL or name of the script parsed (if any).
	StartLine               int64                      `json:"startLine"`          // Line offset of the script within the resource with given URL (for script tags).
	StartColumn             int64                      `json:"startColumn"`        // Column offset of the script within the resource with given URL.
	EndLine                 int64                      `json:"endLine"`            // Last line of the script.
	EndColumn               int64                      `json:"endColumn"`          // Length of the last line of the script.
	ExecutionContextID      runtime.ExecutionContextID `json:"executionContextId"` // Specifies script creation context.
	Hash                    string                     `json:"hash"`               // Content hash of the script.
	ExecutionContextAuxData easyjson.RawMessage        `json:"executionContextAuxData,omitempty"`
	IsLiveEdit              bool                       `json:"isLiveEdit,omitempty"`   // True, if this script is generated as a result of the live edit operation.
	SourceMapURL            string                     `json:"sourceMapURL,omitempty"` // URL of source map associated with script (if any).
	HasSourceURL            bool                       `json:"hasSourceURL,omitempty"` // True, if this script has sourceURL.
	IsModule                bool                       `json:"isModule,omitempty"`     // True, if this script is ES6 module.
	Length                  int64                      `json:"length,omitempty"`       // This script length.
	StackTrace              *runtime.StackTrace        `json:"stackTrace,omitempty"`   // JavaScript top stack frame of where the script parsed event was triggered if available.
}
