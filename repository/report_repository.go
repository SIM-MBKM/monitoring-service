package repository

import (
	"context"
	"monitoring-service/entity"

	"gorm.io/gorm"
)

type reportRepository struct {
	db             *gorm.DB
	baseRepository BaseRepository
}

type ReportRepository interface {
	Index(ctx context.Context, tx *gorm.DB) ([]entity.Report, error)
	Create(ctx context.Context, report entity.Report, tx *gorm.DB) (entity.Report, error)
	Update(ctx context.Context, report entity.Report, tx *gorm.DB) error
	FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Report, error)
	Destroy(ctx context.Context, id string, tx *gorm.DB) error
	FindByReportScheduleID(ctx context.Context, reportScheduleID string, tx *gorm.DB) ([]entity.Report, error)
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
}
func (r *reportRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Report, error) {
	var reports []entity.Report

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Report{}).Find(&reports).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}

	return reports, nil
}
func (r *reportRepository) Create(ctx context.Context, report entity.Report, tx *gorm.DB) (entity.Report, error) {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return entity.Report{}, err
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

	err = tx.Debug().Model(&entity.Report{}).Create(&report).Error
	if err != nil {
		return entity.Report{}, err
	}

	return report, nil
}
func (r *reportRepository) Update(ctx context.Context, report entity.Report, tx *gorm.DB) error {
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

	err = tx.Debug().Model(&entity.Report{}).Where("id = ?", report.ID).Updates(&report).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *reportRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Report, error) {
	var report entity.Report

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Report{}).Where("id = ?", id).First(&report).Error
	if err != nil {
		return entity.Report{}, err
	}

	return report, nil
}
func (r *reportRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
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

	err = tx.Debug().Model(&entity.Report{}).Where("id = ?", id).Delete(&entity.Report{}).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *reportRepository) FindByReportScheduleID(ctx context.Context, reportScheduleID string, tx *gorm.DB) ([]entity.Report, error) {
	var reports []entity.Report

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Report{}).Where("report_schedule_id = ?", reportScheduleID).Find(&reports).Error
	if err != nil {
		return nil, err
	}

	return reports, nil
}
