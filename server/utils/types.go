package utils

import "github.com/jackc/pgx/v5/pgtype"

func ToPgtypeText(s *string) pgtype.Text {
    if s == nil {
        return pgtype.Text{Valid: false}
    }
    return pgtype.Text{String: *s, Valid: true}
}