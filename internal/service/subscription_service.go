package service

import (
	"errors"

	"github.com/sprrwnnn/subscription-tracker/internal/models"
	"github.com/sprrwnnn/subscription-tracker/internal/repository"

	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
	log  *logrus.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, log *logrus.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		log:  log,
	}
}

func (s *SubscriptionService) Create(req *models.CreateSubscriptionRequest) (*models.Subscription, error) {
	sub, err := req.ToSubscription()
	if err != nil {
		s.log.WithError(err).Warn("Invalid subscription data")
		return nil, err
	}

	if err := s.repo.Create(sub); err != nil {
		s.log.WithError(err).Error("Failed to create subscription")
		return nil, errors.New("database error")
	}

	return sub, nil
}

func (s *SubscriptionService) GetByID(id string) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		s.log.WithError(err).Error("Failed to get subscription")
		return nil, errors.New("database error")
	}
	if sub == nil {
		return nil, errors.New("subscription not found")
	}
	return sub, nil
}

func (s *SubscriptionService) Update(id string, req *models.UpdateSubscriptionRequest) error {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("database error")
	}
	if sub == nil {
		return errors.New("subscription not found")
	}

	updates := make(map[string]interface{})

	if req.ServiceName != nil {
		updates["service_name"] = *req.ServiceName
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.StartDate != nil {
		startDate, err := models.ParseMonthYear(*req.StartDate)
		if err != nil {
			return errors.New("invalid start_date format")
		}
		updates["start_date"] = startDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			updates["end_date"] = nil
		} else {
			endDate, err := models.ParseMonthYear(*req.EndDate)
			if err != nil {
				return errors.New("invalid end_date format")
			}
			updates["end_date"] = endDate
		}
	}

	return s.repo.Update(id, updates)
}

func (s *SubscriptionService) Delete(id string) error {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("database error")
	}
	if sub == nil {
		return errors.New("subscription not found")
	}
	return s.repo.Delete(id)
}

func (s *SubscriptionService) List(userID *string, page, pageSize int) ([]models.Subscription, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.List(userID, pageSize, offset)
}

func (s *SubscriptionService) CalculateCost(req *models.CalculateCostRequest) (int, error) {
	startDate, err := models.ParseMonthYear(req.StartDate)
	if err != nil {
		return 0, errors.New("invalid start_date format")
	}

	endDate, err := models.ParseMonthYear(req.EndDate)
	if err != nil {
		return 0, errors.New("invalid end_date format")
	}

	if endDate.Before(startDate) {
		return 0, errors.New("end_date must be after start_date")
	}

	// Set end date to last day of month
	endDate = endDate.AddDate(0, 1, -1)

	return s.repo.CalculateCost(startDate, endDate, req.UserID, req.ServiceName)
}
