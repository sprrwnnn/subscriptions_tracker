package repository

import (
	"errors"
	"time"

	"github.com/sprrwnnn/subscription-tracker/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewSubscriptionRepository(db *gorm.DB, log *logrus.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:  db,
		log: log,
	}
}

func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	r.log.WithFields(logrus.Fields{
		"service_name": sub.ServiceName,
		"user_id":      sub.UserID,
		"price":        sub.Price,
	}).Info("Creating new subscription")

	return r.db.Create(sub).Error
}

func (r *SubscriptionRepository) GetByID(id string) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.Where("id = ?", id).First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sub, err
}

func (r *SubscriptionRepository) Update(id string, updates map[string]interface{}) error {
	r.log.WithFields(logrus.Fields{
		"id":      id,
		"updates": updates,
	}).Info("Updating subscription")

	result := r.db.Model(&models.Subscription{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *SubscriptionRepository) Delete(id string) error {
	r.log.WithField("id", id).Info("Deleting subscription")

	result := r.db.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription not found")
	}
	return nil
}

func (r *SubscriptionRepository) List(userID *string, limit, offset int) ([]models.Subscription, int64, error) {
	var subscriptions []models.Subscription
	var total int64

	query := r.db.Model(&models.Subscription{})
	if userID != nil && *userID != "" {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		r.log.WithError(err).Error("Failed to count subscriptions")
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&subscriptions).Error
	if err != nil {
		r.log.WithError(err).Error("Failed to list subscriptions")
		return nil, 0, err
	}

	r.log.WithFields(logrus.Fields{
		"user_id": userID,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	}).Info("Listed subscriptions")

	return subscriptions, total, nil
}

func (r *SubscriptionRepository) CalculateCost(startDate, endDate time.Time, userID, serviceName *string) (int, error) {
	var subscriptions []models.Subscription

	query := r.db.Model(&models.Subscription{}).
		Where("start_date <= ?", endDate).
		Where("end_date IS NULL OR end_date >= ?", startDate)

	if userID != nil && *userID != "" {
		query = query.Where("user_id = ?", *userID)
	}

	if serviceName != nil && *serviceName != "" {
		query = query.Where("service_name = ?", *serviceName)
	}

	if err := query.Find(&subscriptions).Error; err != nil {
		r.log.WithError(err).Error("Failed to fetch subscriptions for cost calculation")
		return 0, err
	}

	total := 0
	for _, sub := range subscriptions {
		subStart := max(sub.StartDate, startDate)
		subEnd := min(getEndDate(&sub), endDate)

		if subEnd.Before(subStart) {
			continue
		}

		months := calculateMonthsBetween(subStart, subEnd)
		total += sub.Price * months

		r.log.WithFields(logrus.Fields{
			"subscription_id": sub.ID,
			"service_name":    sub.ServiceName,
			"months":          months,
			"price":           sub.Price,
			"contribution":    sub.Price * months,
		}).Debug("Calculated contribution for subscription")
	}

	r.log.WithFields(logrus.Fields{
		"start_date":          startDate.Format("2006-01-02"),
		"end_date":            endDate.Format("2006-01-02"),
		"user_id":             userID,
		"service_name":        serviceName,
		"total":               total,
		"subscriptions_count": len(subscriptions),
	}).Info("Calculated total cost")

	return total, nil
}

func getEndDate(sub *models.Subscription) time.Time {
	if sub.EndDate != nil {
		return *sub.EndDate
	}
	// If no end date, assume it's ongoing (use far future date)
	return time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func calculateMonthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	return years*12 + months + 1
}
