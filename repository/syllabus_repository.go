package repository

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"time"

	"gorm.io/gorm"
)

type syllabusRepository struct {
	db             *gorm.DB
	baseRepository BaseRepository
}

type SyllabusRepository interface {
	Index(ctx context.Context, tx *gorm.DB) ([]entity.Syllabus, error)
	Create(ctx context.Context, syllabus entity.Syllabus, tx *gorm.DB) (entity.Syllabus, error)
	Update(ctx context.Context, id string, syllabus entity.Syllabus, tx *gorm.DB) error
	FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Syllabus, error)
	Destroy(ctx context.Context, id string, tx *gorm.DB) error
	FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Syllabus, error)
	FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Syllabus, error)
	FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Syllabus, int64, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Syllabus, error)
}

func NewSyllabusRepository(db *gorm.DB) SyllabusRepository {
	return &syllabusRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
}

func (r *syllabusRepository) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Syllabus, int64, error) {
	if tx == nil {
		tx = r.db
	}

	// Build the query for unique UserNRPs
	query := tx.Debug().
		Model(&entity.Syllabus{}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("deleted_at IS NULL")

	// Apply user_nrp filter if provided
	if userNRPFilter != "" {
		query = query.Where("user_nrp LIKE ?", "%"+userNRPFilter+"%")
	}

	// Get unique UserNRPs for this advisor (ordered by user_nrp)
	var userNRPs []string
	err := query.Distinct("user_nrp").
		Order("user_nrp ASC"). // Ensure consistent ordering
		Pluck("user_nrp", &userNRPs).Error

	if err != nil {
		return nil, 0, err
	}

	// Get total count of unique UserNRPs
	totalCount := int64(len(userNRPs))

	// Apply pagination to the UserNRPs list
	var paginatedUserNRPs []string
	if pagReq != nil {
		startIdx := pagReq.Offset
		endIdx := pagReq.Offset + pagReq.Limit

		// Check bounds
		if startIdx >= len(userNRPs) {
			paginatedUserNRPs = []string{}
		} else if endIdx > len(userNRPs) {
			paginatedUserNRPs = userNRPs[startIdx:]
		} else {
			paginatedUserNRPs = userNRPs[startIdx:endIdx]
		}
	} else {
		paginatedUserNRPs = userNRPs
	}

	// Return empty if no UserNRPs after pagination
	if len(paginatedUserNRPs) == 0 {
		return make(map[string]entity.Syllabus), totalCount, nil
	}

	// Query for syllabuses with the paginated user NRPs
	var syllabuses []entity.Syllabus
	err = tx.Debug().
		Model(&entity.Syllabus{}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("user_nrp IN ?", paginatedUserNRPs).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&syllabuses).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert slice to map with user_nrp as the key, keeping only the newest syllabus for each NRP
	syllabusMap := make(map[string]entity.Syllabus)
	for _, syllabus := range syllabuses {
		// Only add to map if this NRP doesn't exist yet (since we ordered by created_at DESC, the first one is the newest)
		if _, exists := syllabusMap[syllabus.UserNRP]; !exists {
			syllabusMap[syllabus.UserNRP] = syllabus
		}
	}

	return syllabusMap, totalCount, nil
}

func (r *syllabusRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Syllabus, error) {
	var syllabuses []entity.Syllabus

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Syllabus{}).Find(&syllabuses).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}

	return syllabuses, nil
}

func (r *syllabusRepository) Create(ctx context.Context, syllabus entity.Syllabus, tx *gorm.DB) (entity.Syllabus, error) {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return entity.Syllabus{}, err
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

	err = tx.Debug().Model(&entity.Syllabus{}).Create(&syllabus).Error
	if err != nil {
		return entity.Syllabus{}, err
	}

	return syllabus, nil
}

func (r *syllabusRepository) Update(ctx context.Context, id string, syllabus entity.Syllabus, tx *gorm.DB) error {
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

	err = tx.Debug().
		Model(&entity.Syllabus{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Updates(&syllabus).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *syllabusRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Syllabus, error) {
	var syllabus entity.Syllabus

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Syllabus{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&syllabus).Error
	if err != nil {
		return entity.Syllabus{}, err
	}

	return syllabus, nil
}

func (r *syllabusRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
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

	err = tx.Debug().
		Model(&entity.Syllabus{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *syllabusRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Syllabus, error) {
	var syllabus entity.Syllabus

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Syllabus{}).
		Where("registration_id = ?", registrationID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		First(&syllabus).Error
	if err != nil {
		return entity.Syllabus{}, err
	}

	return syllabus, nil
}

func (r *syllabusRepository) FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Syllabus, error) {
	var syllabuses []entity.Syllabus

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Syllabus{}).
		Where("registration_id = ?", registrationID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&syllabuses).Error
	if err != nil {
		return nil, err
	}

	return syllabuses, nil
}

func (r *syllabusRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Syllabus, error) {
	var syllabuses []entity.Syllabus

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Syllabus{}).
		Where("user_nrp = ?", userNRP).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&syllabuses).Error

	if err != nil {
		return nil, err
	}

	// Group syllabuses by registration_id
	syllabusMap := make(map[string][]entity.Syllabus)
	for _, syllabus := range syllabuses {
		syllabusMap[syllabus.RegistrationID] = append(syllabusMap[syllabus.RegistrationID], syllabus)
	}

	return syllabusMap, nil
}
