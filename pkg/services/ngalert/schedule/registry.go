package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/grafana/grafana/pkg/services/ngalert/models"
)

type alertRuleInfoRegistry struct {
	mu            sync.Mutex
	alertRuleInfo map[models.AlertRuleKey]*alertRuleInfo
}

// getOrCreateInfo gets rule routine information from registry by the key. If it does not exist, it creates a new one.
// Returns a pointer to the rule routine information and a flag that indicates whether it is a new struct or not.
func (r *alertRuleInfoRegistry) getOrCreateInfo(context context.Context, key models.AlertRuleKey) (*alertRuleInfo, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, ok := r.alertRuleInfo[key]
	if !ok {
		info = newAlertRuleInfo(context)
		r.alertRuleInfo[key] = info
	}
	return info, !ok
}

// get returns the channel for the specific alert rule
// if the key does not exist returns an error
func (r *alertRuleInfoRegistry) get(key models.AlertRuleKey) (*alertRuleInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, ok := r.alertRuleInfo[key]
	if !ok {
		return nil, fmt.Errorf("%v key not found", key)
	}
	return info, nil
}

func (r *alertRuleInfoRegistry) exists(key models.AlertRuleKey) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.alertRuleInfo[key]
	return ok
}

// del removes pair that has specific key from alertRuleInfo.
// Returns 2-tuple where the first element is value of the removed pair
// and the second element indicates whether element with the specified key existed.
func (r *alertRuleInfoRegistry) del(key models.AlertRuleKey) (*alertRuleInfo, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	info, ok := r.alertRuleInfo[key]
	if ok {
		delete(r.alertRuleInfo, key)
	}
	return info, ok
}

func (r *alertRuleInfoRegistry) keyMap() map[models.AlertRuleKey]struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	definitionsIDs := make(map[models.AlertRuleKey]struct{}, len(r.alertRuleInfo))
	for k := range r.alertRuleInfo {
		definitionsIDs[k] = struct{}{}
	}
	return definitionsIDs
}

type ruleVersion int64

type alertRuleInfo struct {
	evalCh   chan *evaluation
	updateCh chan ruleVersion
	ctx      context.Context
	stop     context.CancelFunc
}

func newAlertRuleInfo(parent context.Context) *alertRuleInfo {
	ctx, cancel := context.WithCancel(parent)
	return &alertRuleInfo{evalCh: make(chan *evaluation), updateCh: make(chan ruleVersion), ctx: ctx, stop: cancel}
}

// eval signals the rule evaluation routine to perform the evaluation of the rule. Does nothing if the loop is stopped.
// Before sending a message into the channel, it does non-blocking read to make sure that there is no concurrent send operation.
// Returns a tuple where first element is
//   - true when message was sent
//   - false when the send operation is stopped
// the second element contains a dropped message that was sent by a concurrent sender.
func (a *alertRuleInfo) eval(t time.Time, version int64) (bool, *evaluation) {
	// read the channel in unblocking manner to make sure that there is no concurrent send operation.
	var droppedMsg *evaluation
	select {
	case droppedMsg = <-a.evalCh:
	default:
	}

	select {
	case a.evalCh <- &evaluation{
		scheduledAt: t,
		version:     version,
	}:
		return true, droppedMsg
	case <-a.ctx.Done():
		return false, droppedMsg
	}
}

// update signals the rule evaluation routine to update the internal state. Does nothing if the loop is stopped.
func (a *alertRuleInfo) update(version ruleVersion) bool {
	// check if the channel is not empty.
	msg := version
	select {
	case v := <-a.updateCh:
		// if it has a version pick the greatest one.
		if v > msg {
			msg = v
		}
	case <-a.ctx.Done():
		return false
	default:
	}

	select {
	case a.updateCh <- msg:
		return true
	case <-a.ctx.Done():
		return false
	}
}

type evaluation struct {
	scheduledAt time.Time
	version     int64
}

type schedulableAlertRulesRegistry struct {
	rules map[models.AlertRuleKey]*models.SchedulableAlertRule
	mu    sync.Mutex
}

// all returns all rules in the registry.
func (r *schedulableAlertRulesRegistry) all() []*models.SchedulableAlertRule {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]*models.SchedulableAlertRule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result
}

func (r *schedulableAlertRulesRegistry) get(k models.AlertRuleKey) *models.SchedulableAlertRule {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rules[k]
}

// set replaces all rules in the registry.
func (r *schedulableAlertRulesRegistry) set(rules []*models.SchedulableAlertRule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = make(map[models.AlertRuleKey]*models.SchedulableAlertRule)
	for _, rule := range rules {
		r.rules[rule.GetKey()] = rule
	}
}

// update inserts or replaces a rule in the registry.
func (r *schedulableAlertRulesRegistry) update(rule *models.SchedulableAlertRule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.GetKey()] = rule
}

// del removes pair that has specific key from schedulableAlertRulesRegistry.
// Returns 2-tuple where the first element is value of the removed pair
// and the second element indicates whether element with the specified key existed.
func (r *schedulableAlertRulesRegistry) del(k models.AlertRuleKey) (*models.SchedulableAlertRule, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule, ok := r.rules[k]
	if ok {
		delete(r.rules, k)
	}
	return rule, ok
}
