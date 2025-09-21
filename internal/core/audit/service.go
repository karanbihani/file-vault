package audit

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karanbihani/file-vault/internal/db"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

// LogActivity creates a new audit log entry in a background goroutine.
func (s *Service) LogActivity(ctx context.Context, userID int64, action string, details map[string]interface{}) {
	go func() {
		detailsJSON, err := json.Marshal(details)
		if err != nil {
			log.Printf("ERROR: failed to marshal audit log details: %v", err)
			return
		}

		// --- THIS IS THE FIX ---
		// We construct a pgtype.Int8 struct from our int64 userID.
		// Since we know the userID will always be valid here, we set Valid to true.
		err = s.queries.CreateAuditLog(context.Background(), db.CreateAuditLogParams{
			UserID:  pgtype.Int8{Int64: userID, Valid: true},
			Action:  action,
			Details: detailsJSON,
		})
		// --- END OF FIX ---
		
		if err != nil {
			log.Printf("ERROR: failed to create audit log: %v", err)
		}
	}()
}