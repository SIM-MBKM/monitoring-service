package repository

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"

	"github.com/google/uuid"
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
	FindByUserID(ctx context.Context, userNRP string, tx *gorm.DB) ([]entity.ReportSchedule, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.ReportSchedule, error)
	FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNrp string) (map[string][]entity.ReportSchedule, int64, error)
}

func NewReportScheduleRepository(db *gorm.DB) ReportScheduleReposiotry {
	return &reportScheduleRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
}

func (r *reportScheduleRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.ReportSchedule, error) {
	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	var reportSchedules []entity.ReportSchedule
	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Where("user_nrp = ?", userNRP).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&reportSchedules).Error

	if err != nil {
		return nil, err
	}

	userReportSchedules := make(map[string][]entity.ReportSchedule)
	for _, schedule := range reportSchedules {
		// get report based on report_schedule_id order by created_at DESC limit 1
		var report entity.Report
		err = tx.Debug().
			Model(&entity.Report{}).
			Where("report_schedule_id = ?", schedule.ID).
			Order("created_at DESC").
			Limit(1).
			Find(&report).Error
		if err != nil {
			return nil, err
		}

		if report.ID != uuid.Nil {
			schedule.Report = []entity.Report{report}
		}

		userReportSchedules[schedule.RegistrationID] = append(userReportSchedules[schedule.RegistrationID], schedule)
	}

	return userReportSchedules, nil
}

// func (r *reportScheduleRepository) FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNrp string) (map[string][]entity.ReportSchedule, int64, error) {
// 	if tx == nil {
// 		tx = r.db
// 	}

// 	// First, get all unique UserNRPs for this advisor (ordered by user_nrp)
// 	var userNRPs []string
// 	query := tx.Debug().
// 		Model(&entity.ReportSchedule{}).
// 		Where("academic_advisor_email = ?", advisorEmail).
// 		Where("deleted_at IS NULL")

// 	if userNrp != "" {
// 		query = query.Where("user_nrp = ?", userNrp)
// 	}

// 	err := query.Distinct("user_nrp").
// 		Order("user_nrp ASC").
// 		Pluck("user_nrp", &userNRPs).Error

// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Get total count of unique UserNRPs
// 	totalCount := int64(len(userNRPs))

// 	// Apply pagination to the UserNRPs list
// 	var paginatedUserNRPs []string
// 	if pagReq != nil {
// 		startIdx := pagReq.Offset
// 		endIdx := pagReq.Offset + pagReq.Limit

// 		// Check bounds
// 		if startIdx >= len(userNRPs) {
// 			paginatedUserNRPs = []string{}
// 		} else if endIdx > len(userNRPs) {
// 			paginatedUserNRPs = userNRPs[startIdx:]
// 		} else {
// 			paginatedUserNRPs = userNRPs[startIdx:endIdx]
// 		}
// 	} else {
// 		paginatedUserNRPs = userNRPs
// 	}

// 	// Return empty if no UserNRPs after pagination
// 	if len(paginatedUserNRPs) == 0 {
// 		return make(map[string][]entity.ReportSchedule), totalCount, nil
// 	}

// 	// Now fetch the report schedules for only the paginated UserNRPs
// 	var allReportSchedules []entity.ReportSchedule
// 	err = tx.Debug().
// 		Model(&entity.ReportSchedule{}).
// 		Where("academic_advisor_email = ?", advisorEmail).
// 		Where("user_nrp IN ?", paginatedUserNRPs).
// 		Where("deleted_at IS NULL").
// 		Order("user_nrp ASC").
// 		Find(&allReportSchedules).Error

// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// Group them by user_nrp
// 	userReportSchedules := make(map[string][]entity.ReportSchedule)

// 	for _, schedule := range allReportSchedules {
// 		// Get the latest report for this schedule
// 		var report entity.Report
// 		err = tx.Debug().
// 			Model(&entity.Report{}).
// 			Where("report_schedule_id = ?", schedule.ID).
// 			Where("deleted_at IS NULL").
// 			Order("created_at DESC").
// 			Limit(1).
// 			First(&report).Error

// 		if err != nil && err != gorm.ErrRecordNotFound {
// 			return nil, 0, err
// 		}

// 		if report.ID != uuid.Nil {
// 			schedule.Report = []entity.Report{report}
// 		}

// 		userReportSchedules[schedule.UserNRP] = append(userReportSchedules[schedule.UserNRP], schedule)
// 	}

// 	return userReportSchedules, totalCount, nil
// }

// Gunakan subquery untuk pagination yang efisien
func (r *reportScheduleRepository) FindByAdvisorEmailAndGroupByUserID(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNrp string) (map[string][]entity.ReportSchedule, int64, error) {
	if tx == nil {
		tx = r.db
	}

	// Count total unique users
	var totalCount int64
	countQuery := tx.Model(&entity.ReportSchedule{}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("deleted_at IS NULL")

	if userNrp != "" {
		countQuery = countQuery.Where("user_nrp = ?", userNrp)
	}

	err := countQuery.Distinct("user_nrp").Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data with joined reports
	var allReportSchedules []entity.ReportSchedule
	query := tx.Debug().
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").
				Order("created_at DESC").
				Limit(1)
		}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("deleted_at IS NULL")

	if userNrp != "" {
		query = query.Where("user_nrp = ?", userNrp)
	}

	// Subquery untuk pagination berdasarkan user_nrp
	subQuery := tx.Model(&entity.ReportSchedule{}).
		Select("DISTINCT user_nrp").
		Where("academic_advisor_email = ?", advisorEmail).
		Where("deleted_at IS NULL")

	if userNrp != "" {
		subQuery = subQuery.Where("user_nrp = ?", userNrp)
	}

	if pagReq != nil {
		subQuery = subQuery.Limit(pagReq.Limit).Offset(pagReq.Offset)
	}

	err = query.Where("user_nrp IN (?)", subQuery).
		Order("user_nrp ASC").
		Find(&allReportSchedules).Error

	if err != nil {
		return nil, 0, err
	}

	// Group by user_nrp (sama seperti sebelumnya tapi tanpa loop query)
	userReportSchedules := make(map[string][]entity.ReportSchedule)
	for _, schedule := range allReportSchedules {
		userReportSchedules[schedule.UserNRP] = append(userReportSchedules[schedule.UserNRP], schedule)
	}

	return userReportSchedules, totalCount, nil
}

func (r *reportScheduleRepository) FindByUserID(ctx context.Context, userNRP string, tx *gorm.DB) ([]entity.ReportSchedule, error) {
	var reportSchedules []entity.ReportSchedule

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.ReportSchedule{}).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(1)
		}).
		Where("user_nrp = ?", userNRP).
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
		userReportSchedules[schedule.UserNRP] = append(userReportSchedules[schedule.UserNRP], schedule)
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
