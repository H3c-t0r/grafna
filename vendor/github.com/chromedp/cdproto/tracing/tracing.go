// Package tracing provides the Chrome Debugging Protocol
// commands, types, and events for the Tracing domain.
//
// Generated by the chromedp-gen command.
package tracing

// Code generated by chromedp-gen. DO NOT EDIT.

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
)

// EndParams stop trace events collection.
type EndParams struct{}

// End stop trace events collection.
func End() *EndParams {
	return &EndParams{}
}

// Do executes Tracing.end against the provided context.
func (p *EndParams) Do(ctxt context.Context, h cdp.Executor) (err error) {
	return h.Execute(ctxt, CommandEnd, nil, nil)
}

// GetCategoriesParams gets supported tracing categories.
type GetCategoriesParams struct{}

// GetCategories gets supported tracing categories.
func GetCategories() *GetCategoriesParams {
	return &GetCategoriesParams{}
}

// GetCategoriesReturns return values.
type GetCategoriesReturns struct {
	Categories []string `json:"categories,omitempty"` // A list of supported tracing categories.
}

// Do executes Tracing.getCategories against the provided context.
//
// returns:
//   categories - A list of supported tracing categories.
func (p *GetCategoriesParams) Do(ctxt context.Context, h cdp.Executor) (categories []string, err error) {
	// execute
	var res GetCategoriesReturns
	err = h.Execute(ctxt, CommandGetCategories, nil, &res)
	if err != nil {
		return nil, err
	}

	return res.Categories, nil
}

// RecordClockSyncMarkerParams record a clock sync marker in the trace.
type RecordClockSyncMarkerParams struct {
	SyncID string `json:"syncId"` // The ID of this clock sync marker
}

// RecordClockSyncMarker record a clock sync marker in the trace.
//
// parameters:
//   syncID - The ID of this clock sync marker
func RecordClockSyncMarker(syncID string) *RecordClockSyncMarkerParams {
	return &RecordClockSyncMarkerParams{
		SyncID: syncID,
	}
}

// Do executes Tracing.recordClockSyncMarker against the provided context.
func (p *RecordClockSyncMarkerParams) Do(ctxt context.Context, h cdp.Executor) (err error) {
	return h.Execute(ctxt, CommandRecordClockSyncMarker, p, nil)
}

// RequestMemoryDumpParams request a global memory dump.
type RequestMemoryDumpParams struct{}

// RequestMemoryDump request a global memory dump.
func RequestMemoryDump() *RequestMemoryDumpParams {
	return &RequestMemoryDumpParams{}
}

// RequestMemoryDumpReturns return values.
type RequestMemoryDumpReturns struct {
	DumpGUID string `json:"dumpGuid,omitempty"` // GUID of the resulting global memory dump.
	Success  bool   `json:"success,omitempty"`  // True iff the global memory dump succeeded.
}

// Do executes Tracing.requestMemoryDump against the provided context.
//
// returns:
//   dumpGUID - GUID of the resulting global memory dump.
//   success - True iff the global memory dump succeeded.
func (p *RequestMemoryDumpParams) Do(ctxt context.Context, h cdp.Executor) (dumpGUID string, success bool, err error) {
	// execute
	var res RequestMemoryDumpReturns
	err = h.Execute(ctxt, CommandRequestMemoryDump, nil, &res)
	if err != nil {
		return "", false, err
	}

	return res.DumpGUID, res.Success, nil
}

// StartParams start trace events collection.
type StartParams struct {
	BufferUsageReportingInterval float64           `json:"bufferUsageReportingInterval,omitempty"` // If set, the agent will issue bufferUsage events at this interval, specified in milliseconds
	TransferMode                 TransferMode      `json:"transferMode,omitempty"`                 // Whether to report trace events as series of dataCollected events or to save trace to a stream (defaults to ReportEvents).
	StreamCompression            StreamCompression `json:"streamCompression,omitempty"`            // Compression format to use. This only applies when using ReturnAsStream transfer mode (defaults to none)
	TraceConfig                  *TraceConfig      `json:"traceConfig,omitempty"`
}

// Start start trace events collection.
//
// parameters:
func Start() *StartParams {
	return &StartParams{}
}

// WithBufferUsageReportingInterval if set, the agent will issue bufferUsage
// events at this interval, specified in milliseconds.
func (p StartParams) WithBufferUsageReportingInterval(bufferUsageReportingInterval float64) *StartParams {
	p.BufferUsageReportingInterval = bufferUsageReportingInterval
	return &p
}

// WithTransferMode whether to report trace events as series of dataCollected
// events or to save trace to a stream (defaults to ReportEvents).
func (p StartParams) WithTransferMode(transferMode TransferMode) *StartParams {
	p.TransferMode = transferMode
	return &p
}

// WithStreamCompression compression format to use. This only applies when
// using ReturnAsStream transfer mode (defaults to none).
func (p StartParams) WithStreamCompression(streamCompression StreamCompression) *StartParams {
	p.StreamCompression = streamCompression
	return &p
}

// WithTraceConfig [no description].
func (p StartParams) WithTraceConfig(traceConfig *TraceConfig) *StartParams {
	p.TraceConfig = traceConfig
	return &p
}

// Do executes Tracing.start against the provided context.
func (p *StartParams) Do(ctxt context.Context, h cdp.Executor) (err error) {
	return h.Execute(ctxt, CommandStart, p, nil)
}

// Command names.
const (
	CommandEnd                   = "Tracing.end"
	CommandGetCategories         = "Tracing.getCategories"
	CommandRecordClockSyncMarker = "Tracing.recordClockSyncMarker"
	CommandRequestMemoryDump     = "Tracing.requestMemoryDump"
	CommandStart                 = "Tracing.start"
)
