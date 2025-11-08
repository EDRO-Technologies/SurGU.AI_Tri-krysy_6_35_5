package storage

import "github.com/jackc/pgtype"

type Statistics struct {
	TotalUsers          int64 `db:"total_users"`
	TotalVoiceQuestions int64 `db:"total_voice_questions"`
	TotalTextQuestions  int64 `db:"total_text_questions"`

	Last10Questions pgtype.TextArray `db:"last_10_questions"`
}

type File struct {
	Url       string `db:"url"`
	ShortName string `db:"short_name"`
}
