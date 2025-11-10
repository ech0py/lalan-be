
type TermsAndConditionsRepository interface {
}

// menyimpan ke konesksi database
type termsAndConditionsRepository struct {
	db *sqlx.DB
}

func NewTermAndConditions(db *sqlx.DB) TermsAndConditionsRepository{
	retun termsAndConditionsRepository{db:db}
}

func (r *termsAndConditionsRepository) CreateTermAndConditions(tac *model.TermsAndConditionsModel) error {
	// query insert to term and conditions
}

