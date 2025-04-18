package repository

import (
	"context"
	"monitoring-service/entity"

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
	FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Transcript, error)
}

func NewTranscriptRepository(db *gorm.DB) TranscriptRepository {
	return &transcriptRepository{
		db:             db,
		baseRepository: NewBaseRepository(db),
	}
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

	err = tx.Debug().Model(&entity.Transcript{}).Where("id = ?", id).Updates(&transcript).Error
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

	err := tx.Debug().Model(&entity.Transcript{}).Where("id = ?", id).First(&transcript).Error
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

	err = tx.Debug().Model(&entity.Transcript{}).Where("id = ?", id).Delete(&entity.Transcript{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *transcriptRepository) FindByRegistrationID(ctx context.Context, registrationID string, tx *gorm.DB) ([]entity.Transcript, error) {
	var transcripts []entity.Transcript

	if tx == nil {
		tx = r.db.WithContext(ctx)
	}

	err := tx.Debug().Model(&entity.Transcript{}).Where("registration_id = ?", registrationID).Find(&transcripts).Error
	if err != nil {
		return nil, err
	}

	return transcripts, nil
}
