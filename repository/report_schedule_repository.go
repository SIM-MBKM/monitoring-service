package repository

import (
	"context"
	"monitoring-service/entity"

	"gorm.io/gorm"
)

type reportScheduleRepository struct {
	db             *gorm.DB
	baseRepository BaseRepository
}

type ReportScheduleReposiotry interface {
	Index(ctx context.Context, tx *gorm.DB) (map[string][]entity.ReportSchedule, error)
	Create(ctx context.Context, reportSchedule entity.ReportSchedule, tx *gorm.DB) (entity.ReportSchedule, error)
	Update(ctx context.Context, id string, reportSchedule entity.ReportSchedule, tx *gorm.DB) error
	FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.ReportSchedule, error)
	Destroy(ctx context.Context, id string, tx *gorm.DB) error
	FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.ReportSchedule, error)
	FindByUserID(ctx context.Context, userID string, tx *gorm.DB) ([]entity.ReportSchedule, error)
	FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB) (map[string][]entity.ReportSchedule, error)
}

func NewReportScheduleRepository(db *gorm.DB) ReportScheduleReposiotry {
	return &reportScheduleRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
}

func (r *reportScheduleRepository) FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB) (map[string][]entity.ReportSchedule, error) {
	if tx == nil {
		tx = r.db
	}

	// Fetch all report schedules with their reports
	var allReportSchedules []entity.ReportSchedule
	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("deleted_at IS NULL").
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Find(&allReportSchedules).Error

	if err != nil {
		return nil, err
	}

	// Group them by user_id
	userReportSchedules := make(map[string][]entity.ReportSchedule)
	for _, schedule := range allReportSchedules {
		userReportSchedules[schedule.UserID] = append(userReportSchedules[schedule.UserID], schedule)
	}

	return userReportSchedules, nil
}

func (r *reportScheduleRepository) FindByUserID(ctx context.Context, userID string, tx *gorm.DB) ([]entity.ReportSchedule, error) {
	var reportSchedules []entity.ReportSchedule

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Find(&reportSchedules).Error

	if err != nil {
		return nil, err
	}

	return reportSchedules, nil
}

// TODO: Add Pagination
func (r *reportScheduleRepository) Index(ctx context.Context, tx *gorm.DB) (map[string][]entity.ReportSchedule, error) {
	if tx == nil {
		tx = r.db
	}

	var allReportSchedules []entity.ReportSchedule
	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Where("deleted_at IS NULL").
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Find(&allReportSchedules).Error

	if err != nil {
		return nil, err
	}

	// Group them by user_id
	userReportSchedules := make(map[string][]entity.ReportSchedule)
	for _, schedule := range allReportSchedules {
		userReportSchedules[schedule.UserID] = append(userReportSchedules[schedule.UserID], schedule)
	}

	return userReportSchedules, nil
}

func (r *reportScheduleRepository) Create(ctx context.Context, reportSchedule entity.ReportSchedule, tx *gorm.DB) (entity.ReportSchedule, error) {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return entity.ReportSchedule{}, err
	}

	defer func() {
		if err != nil {
			r.baseRepository.RollbackTx(ctx, tx)
		} else {
			_, err = r.baseRepository.CommitTx(ctx, tx)
			if err != nil {
				return
			}
		}
	}()
	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err = tx.Debug().Model(&entity.ReportSchedule{}).Create(&reportSchedule).Error
	if err != nil {
		return reportSchedule, err
	}

	return reportSchedule, nil
}

func (r *reportScheduleRepository) Update(ctx context.Context, id string, reportSchedule entity.ReportSchedule, tx *gorm.DB) error {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			r.baseRepository.RollbackTx(ctx, tx)
		} else {
			_, err = r.baseRepository.CommitTx(ctx, tx)
			if err != nil {
				return
			}
		}
	}()

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err = tx.Debug().Model(&entity.ReportSchedule{}).Where("id = ?", id).Updates(reportSchedule).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *reportScheduleRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.ReportSchedule, error) {
	var reportSchedule entity.ReportSchedule
	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("id = ?", id).
		First(&reportSchedule).Error

	if err != nil {
		return reportSchedule, err
	}

	return reportSchedule, nil
}
func (r *reportScheduleRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			r.baseRepository.RollbackTx(ctx, tx)
		} else {
			_, err = r.baseRepository.CommitTx(ctx, tx)
			if err != nil {
				return
			}
		}
	}()

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err = tx.Debug().Model(&entity.ReportSchedule{}).Where("id = ?", id).Delete(&entity.ReportSchedule{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *reportScheduleRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.ReportSchedule, error) {
	var reportSchedules []entity.ReportSchedule

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.ReportSchedule{}).Where("registration_id = ?", registrationID).Find(&reportSchedules).Error
	if err != nil {
		return nil, err
	}

	return reportSchedules, nil
}
