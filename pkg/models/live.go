package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/components/securejsondata"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// ChannelPublisher writes data into a channel. Note that permissions are not checked.
type ChannelPublisher func(orgID int64, channel string, data []byte) error

// ChannelClientCount will return the number of clients for a channel
type ChannelClientCount func(orgID int64, channel string) (int, error)

// SubscribeEvent contains subscription data.
type SubscribeEvent struct {
	Channel string
	Path    string
}

// SubscribeReply is a reaction to SubscribeEvent.
type SubscribeReply struct {
	Presence  bool
	JoinLeave bool
	Recover   bool
	Data      json.RawMessage
}

// PublishEvent contains publication data.
type PublishEvent struct {
	Channel string
	Path    string
	Data    json.RawMessage
}

// PublishReply is a reaction to PublishEvent.
type PublishReply struct {
	// By default, it's a handler responsibility to publish data
	// into a stream upon OnPublish but returning a data here
	// will make Grafana Live publish data itself (i.e. stream handler
	// just works as permission proxy in this case).
	Data json.RawMessage
	// HistorySize sets a stream history size.
	HistorySize int
	// HistoryTTL is a time that messages will live in stream history.
	HistoryTTL time.Duration
}

// ChannelHandler defines the core channel behavior
type ChannelHandler interface {
	// OnSubscribe is called when a client wants to subscribe to a channel
	OnSubscribe(ctx context.Context, user *SignedInUser, e SubscribeEvent) (SubscribeReply, backend.SubscribeStreamStatus, error)

	// OnPublish is called when a client writes a message to the channel websocket.
	OnPublish(ctx context.Context, user *SignedInUser, e PublishEvent) (PublishReply, backend.PublishStreamStatus, error)
}

// ChannelHandlerFactory should be implemented by all core features.
type ChannelHandlerFactory interface {
	// GetHandlerForPath gets a ChannelHandler for a path.
	// This is called fast and often -- it must be synchronized
	GetHandlerForPath(path string) (ChannelHandler, error)
}

// DashboardActivityChannel is a service to advertise dashboard activity
type DashboardActivityChannel interface {
	// Called when a dashboard is saved -- this includes the error so we can support a
	// gitops workflow that knows if the value was saved to the local database or not
	// in many cases all direct save requests will fail, but the request should be forwarded
	// to any gitops observers
	DashboardSaved(orgID int64, user *UserDisplayDTO, message string, dashboard *Dashboard, err error) error

	// Called when a dashboard is deleted
	DashboardDeleted(orgID int64, user *UserDisplayDTO, uid string) error

	// Experimental! Indicate is GitOps is active.  This really means
	// someone is subscribed to the `grafana/dashboards/gitops` channel
	HasGitOpsObserver(orgID int64) bool
}

type LiveMessage struct {
	Id        int64
	OrgId     int64
	Channel   string
	Data      json.RawMessage
	Published time.Time
}

type SaveLiveMessageQuery struct {
	OrgId   int64
	Channel string
	Data    json.RawMessage
}

type GetLiveMessageQuery struct {
	OrgId   int64
	Channel string
}

type RemoteWriteConfig struct {
	// Enabled to enable remote write for a channel.
	Enabled bool `json:"enabled,omitempty" yaml:"enabled"`
	// Endpoint to send streaming frames to.
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint"`
	// User is a user for remote write request.
	User string `json:"user,omitempty" yaml:"user"`
	// SampleMilliseconds allow setting minimal time before
	// different remote writes for a channel. 0 means no sampling interval.
	SampleMilliseconds int64 `json:"sampleMilliseconds,omitempty" yaml:"sampleMilliseconds"`
}

// LiveChannelRulePlainConfig contains various channel configuration options.
type LiveChannelRulePlainConfig struct {
	RemoteWrite *RemoteWriteConfig `json:"remoteWrite,omitempty"  yaml:"remoteWrite"`
}

var (
	_ driver.Valuer = LiveChannelRulePlainConfig{}
	_ sql.Scanner   = &LiveChannelRulePlainConfig{}
)

func (a LiveChannelRulePlainConfig) Value() (driver.Value, error) {
	d, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(d), nil
}

func (a *LiveChannelRulePlainConfig) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return json.Unmarshal([]byte(v), &a)
	case []uint8:
		return json.Unmarshal(v, &a)
	default:
		return fmt.Errorf("type assertion on scan failed: got %T", value)
	}
}

// LiveChannelRule represents channel rules saved in database.
type LiveChannelRule struct {
	Id      int64                         `json:"id"`
	OrgId   int64                         `json:"orgId"`
	Version int                           `json:"version"`
	Pattern string                        `json:"pattern"`
	Created time.Time                     `json:"-"`
	Updated time.Time                     `json:"-"`
	Config  LiveChannelRulePlainConfig    `json:"config"`
	Secure  securejsondata.SecureJsonData `json:"secure"`
}

// Also acts as api DTO.
type ListLiveChannelRuleCommand struct {
	OrgId int64 `json:"orgId"`
}

// Also acts as api DTO.
type GetLiveChannelRuleCommand struct {
	Id      int64  `json:"id"`
	OrgId   int64  `json:"orgId"`
	Pattern string `json:"pattern"`
}

// Also acts as api DTO.
type CreateLiveChannelRuleCommand struct {
	OrgId   int64                      `json:"orgId" binding:"Required"`
	Pattern string                     `json:"pattern" binding:"Required"`
	Config  LiveChannelRulePlainConfig `json:"config"`
	Secure  map[string]string          `json:"secure"`
}

// Also acts as api DTO.
type UpdateLiveChannelRuleCommand struct {
	Id      int64                      `json:"id" binding:"Required"`
	OrgId   int64                      `json:"orgId" binding:"Required"`
	Version int                        `json:"version"`
	Pattern string                     `json:"pattern"`
	Config  LiveChannelRulePlainConfig `json:"config"`
	Secure  map[string]string          `json:"secure"`
}

// Also acts as api DTO.
type DeleteLiveChannelRuleCommand struct {
	Id      int64  `json:"id"`
	Pattern string `json:"pattern"`
	OrgId   int64  `json:"orgId"`
}

var (
	ErrLiveChannelRuleExists             = errors.New("channel rule with the same pattern already exists")
	ErrLiveChannelRuleUpdatingOldVersion = errors.New("trying to update old version of live channel rule")
	ErrLiveChannelRuleNotFound           = errors.New("channel rule not found")
	ErrLiveChannelRuleIdentifierNotSet   = errors.New("channel rule identifier not set")
)
