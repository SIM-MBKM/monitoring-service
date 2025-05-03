package repository

import (
	"context"
	"monitoring-service/dto"
	"monitoring-service/entity"
	"time"

	"gorm.io/gorm"
)

type transcriptRepository struct {
	db             *gorm.DB
	baseRepository BaseRepository
}

type TranscriptRepository interface {
	Index(ctx context.Context, tx *gorm.DB) ([]entity.Transcript, error)
	Create(ctx context.Context, transcript entity.Transcript, tx *gorm.DB) (entity.Transcript, error)
	Update(ctx context.Context, id string, transcript entity.Transcript, tx *gorm.DB) error
	FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Transcript, error)
	Destroy(ctx context.Context, id string, tx *gorm.DB) error
	FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Transcript, error)
	FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Transcript, error)
	FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Transcript, int64, error)
	FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Transcript, error)
}

func NewTranscriptRepository(db *gorm.DB) TranscriptRepository {
	return &transcriptRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
}

func (r *transcriptRepository) FindByAdvisorEmailAndGroupByUserNRP(ctx context.Context, advisorEmail string, tx *gorm.DB, pagReq *dto.PaginationRequest, userNRPFilter string) (map[string]entity.Transcript, int64, error) {
	if tx == nil {
		tx = r.db
	}

	// Build the query for unique UserNRPs
	query := tx.Debug().
		Model(&entity.Transcript{}).
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
		return make(map[string]entity.Transcript), totalCount, nil
	}

	// Query for transcripts with the paginated user NRPs
	var transcripts []entity.Transcript
	err = tx.Debug().
		Model(&entity.Transcript{}).
		Where("academic_advisor_email = ?", advisorEmail).
		Where("user_nrp IN ?", paginatedUserNRPs).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&transcripts).Error

	if err != nil {
		return nil, 0, err
	}

	// Convert slice to map with user_nrp as the key, keeping only the newest transcript for each NRP
	transcriptMap := make(map[string]entity.Transcript)
	for _, transcript := range transcripts {
		// Only add to map if this NRP doesn't exist yet (since we ordered by created_at DESC, the first one is the newest)
		if _, exists := transcriptMap[transcript.UserNRP]; !exists {
			transcriptMap[transcript.UserNRP] = transcript
		}
	}

	return transcriptMap, totalCount, nil
}

func (r *transcriptRepository) Index(ctx context.Context, tx *gorm.DB) ([]entity.Transcript, error) {
	var transcripts []entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Transcript{}).Find(&transcripts).Where("deleted_at IS NULL").Error
	if err != nil {
		return nil, err
	}

	return transcripts, nil
}

func (r *transcriptRepository) Create(ctx context.Context, transcript entity.Transcript, tx *gorm.DB) (entity.Transcript, error) {
	tx, err := r.baseRepository.BeginTx(ctx)
	if err != nil {
		return entity.Transcript{}, err
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

	err = tx.Debug().Model(&entity.Transcript{}).Create(&transcript).Error
	if err != nil {
		return entity.Transcript{}, err
	}

	return transcript, nil
}

func (r *transcriptRepository) Update(ctx context.Context, id string, transcript entity.Transcript, tx *gorm.DB) error {
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
		Model(&entity.Transcript{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Updates(&transcript).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transcriptRepository) FindByID(ctx context.Context, id string, tx *gorm.DB) (entity.Transcript, error) {
	var transcript entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Transcript{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&transcript).Error
	if err != nil {
		return entity.Transcript{}, err
	}

	return transcript, nil
}

func (r *transcriptRepository) Destroy(ctx context.Context, id string, tx *gorm.DB) error {
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
		Model(&entity.Transcript{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transcriptRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) (entity.Transcript, error) {
	var transcripts entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Transcript{}).
		Where("registration_id = ?", registrationID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		First(&transcripts).Error
	if err != nil {
		return entity.Transcript{}, err
	}

	return transcripts, nil
}

func (r *transcriptRepository) FindAllByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Transcript, error) {
	var transcripts []entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Transcript{}).
		Where("registration_id = ?", registrationID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&transcripts).Error
	if err != nil {
		return nil, err
	}

	return transcripts, nil
}

func (r *transcriptRepository) FindByUserNRPAndGroupByRegistrationID(ctx context.Context, userNRP string, tx *gorm.DB) (map[string][]entity.Transcript, error) {
	var transcripts []entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().
		Model(&entity.Transcript{}).
		Where("user_nrp = ?", userNRP).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&transcripts).Error

	if err != nil {
		return nil, err
	}

	// Group transcripts by registration_id
	transcriptMap := make(map[string][]entity.Transcript)
	for _, transcript := range transcripts {
		transcriptMap[transcript.RegistrationID] = append(transcriptMap[transcript.RegistrationID], transcript)
	}

	return transcriptMap, nil
}
