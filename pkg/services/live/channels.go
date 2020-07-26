package live

import (
	"encoding/json"
	"math/rand"
	"time"
)

// ChannelInfo holds metadata about each channel and is retruned on connection.
// Eventually each plugin should control exactly what is in this structure
type ChannelInfo struct {
	Description string
}

type randomWalkeMessage struct {
	Time  int64
	Value float64
	Min   float64
	Max   float64
}

// RunRandomCSV just for an example
func RunRandomCSV(broker *GrafanaLive, channel string, speedMillis int, dropPercent float64) {
	spread := 50.0

	walker := rand.Float64() * 100
	ticker := time.NewTicker(time.Duration(speedMillis) * time.Millisecond)

	line := randomWalkeMessage{}

	for t := range ticker.C {
		if rand.Float64() <= dropPercent {
			continue //
		}
		delta := rand.Float64() - 0.5
		walker += delta

		line.Time = t.UnixNano() / int64(time.Millisecond)
		line.Value = walker
		line.Min = walker - ((rand.Float64() * spread) + 0.01)
		line.Max = walker + ((rand.Float64() * spread) + 0.01)

		bytes, _ := json.Marshal(&line)
		v := broker.Publish(channel, bytes)
		if !v {
			logger.Warn("write", "channel", channel, "line", line, "ok", v)
		}
	}
}
