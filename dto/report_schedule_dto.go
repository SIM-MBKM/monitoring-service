package dto

type (
	ReportScheduleRequest struct {
		UserID               string `json:"user_id" validate:"required"`
		UserNRP              string `json:"user_nrp" validate:"required"`
		RegistrationID       string `json:"registration_id" validate:"required"`
		AcademicAdvisorID    string `json:"academic_advisor_id" validate:"required"`
		AcademicAdvisorEmail string `json:"academic_advisor_email" validate:"required"`
		ReportType           string `json:"report_type" validate:"required"`
		Week                 int    `json:"week" validate:"required"`
		StartDate            string `json:"start_date" validate:"required"`
		EndDate              string `json:"end_date" validate:"required"`
	}

	ReportScheduleResponse struct {
		ID                   string          `json:"id"`
		UserID               string          `json:"user_id"`
		UserNRP              string          `json:"user_nrp"`
		RegistrationID       string          `json:"registration_id"`
		ActivityName         string          `json:"activity_name"`
		AcademicAdvisorID    string          `json:"academic_advisor_id"`
		AcademicAdvisorEmail string          `json:"academic_advisor_email"`
		ReportType           string          `json:"report_type"`
		Week                 int             `json:"week"`
		StartDate            string          `json:"start_date"`
		EndDate              string          `json:"end_date"`
		Report               *ReportResponse `json:"report"`
	}

	ReportScheduleByAdvisorResponse struct {
		Reports map[string][]ReportScheduleResponse `json:"reports"`
	}

	ReportScheduleByStudentResponse struct {
		Reports map[string][]ReportScheduleResponse `json:"reports"`
	}

	ReportScheduleAdvisorRequest struct {
		UserNRP string `json:"user_nrp"`
	}
)
