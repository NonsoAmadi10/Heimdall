package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/db"
	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"gorm.io/gorm"
)

const (
	alertTypeLowBTCPeers   = "low_btc_peers"
	alertTypeSyncStalled   = "sync_stalled"
	alertTypeBandwidthAnom = "bandwidth_spike"
)

func EvaluateAlerts(database *gorm.DB, metric *utils.ConnectionMetrics) error {
	minBTCPeers := float64(utils.GetEnvInt("ALERT_MIN_BTC_PEERS", 3))
	if err := upsertAlert(
		database,
		alertTypeLowBTCPeers,
		"warning",
		"BTC peer count is below configured threshold",
		float64(metric.NumBTCPeers),
		minBTCPeers,
		float64(metric.NumBTCPeers) < minBTCPeers,
	); err != nil {
		return err
	}

	if err := upsertAlert(
		database,
		alertTypeSyncStalled,
		"critical",
		"Node is not synced to chain",
		boolToFloat(metric.SyncedToChain),
		1,
		!metric.SyncedToChain,
	); err != nil {
		return err
	}

	if err := evaluateBandwidthSpike(database, metric); err != nil {
		return err
	}

	return nil
}

func evaluateBandwidthSpike(database *gorm.DB, metric *utils.ConnectionMetrics) error {
	lookbackSamples := utils.GetEnvInt("ALERT_LOOKBACK_SAMPLES", 10)
	spikeMultiplier := utils.GetEnvFloat("ALERT_BANDWIDTH_SPIKE_MULTIPLIER", 2.5)

	var recent []utils.ConnectionMetrics
	if err := database.
		Where("id != ?", metric.ID).
		Order("timestamp DESC").
		Limit(lookbackSamples).
		Find(&recent).Error; err != nil {
		return fmt.Errorf("failed to load recent metrics for alerting: %w", err)
	}

	if len(recent) < 3 {
		return nil
	}

	var total float64
	for _, m := range recent {
		total += float64(m.BtcdBandwidthIn + m.BtcdBandwidthOut)
	}

	avgBandwidth := total / float64(len(recent))
	currentBandwidth := float64(metric.BtcdBandwidthIn + metric.BtcdBandwidthOut)
	threshold := avgBandwidth * spikeMultiplier

	return upsertAlert(
		database,
		alertTypeBandwidthAnom,
		"warning",
		"BTC total bandwidth has spiked above recent baseline",
		currentBandwidth,
		threshold,
		currentBandwidth > threshold,
	)
}

func upsertAlert(database *gorm.DB, alertType, severity, message string, value, threshold float64, active bool) error {
	var openAlert utils.Alert
	err := database.
		Where("type = ? AND status IN (?, ?)", alertType, utils.AlertStatusOpen, utils.AlertStatusAcknowledged).
		Order("id DESC").
		First(&openAlert).Error

	now := time.Now().UTC()
	if active {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAlert := &utils.Alert{
				Type:        alertType,
				Severity:    severity,
				Status:      utils.AlertStatusOpen,
				Message:     message,
				MetricValue: value,
				Threshold:   threshold,
				TriggeredAt: now,
			}
			return database.Create(newAlert).Error
		}
		if err != nil {
			return err
		}

		openAlert.Severity = severity
		openAlert.Message = message
		openAlert.MetricValue = value
		openAlert.Threshold = threshold
		return database.Save(&openAlert).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	openAlert.Status = utils.AlertStatusResolved
	openAlert.ResolvedAt = &now
	return database.Save(&openAlert).Error
}

func FetchAlerts(status string) ([]utils.Alert, error) {
	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	query := database.Order("triggered_at DESC")
	if status != "" {
		if status != utils.AlertStatusOpen && status != utils.AlertStatusAcknowledged && status != utils.AlertStatusResolved {
			return nil, fmt.Errorf("invalid alert status filter")
		}
		query = query.Where("status = ?", status)
	}

	var alerts []utils.Alert
	if err := query.Find(&alerts).Error; err != nil {
		return nil, err
	}

	return alerts, nil
}

func AcknowledgeAlert(id uint) (*utils.Alert, error) {
	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	var alert utils.Alert
	if err := database.First(&alert, id).Error; err != nil {
		return nil, err
	}

	if alert.Status == utils.AlertStatusResolved {
		return nil, fmt.Errorf("resolved alert cannot be acknowledged")
	}

	now := time.Now().UTC()
	alert.Status = utils.AlertStatusAcknowledged
	alert.AckedAt = &now
	if err := database.Save(&alert).Error; err != nil {
		return nil, err
	}

	return &alert, nil
}

func ResolveAlert(id uint) (*utils.Alert, error) {
	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	var alert utils.Alert
	if err := database.First(&alert, id).Error; err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	alert.Status = utils.AlertStatusResolved
	alert.ResolvedAt = &now
	if err := database.Save(&alert).Error; err != nil {
		return nil, err
	}

	return &alert, nil
}

func boolToFloat(value bool) float64 {
	if value {
		return 1
	}
	return 0
}
