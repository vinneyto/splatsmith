package standalone

import (
	"context"
	"log"

	"github.com/vinneyto/ariadne/api/internal/core"
)

type LogNotifier struct{}

func NewLogNotifier() *LogNotifier { return &LogNotifier{} }

func (n *LogNotifier) NotifyScanCompleted(_ context.Context, userEmail string, scan core.Scan) error {
	log.Printf("scan completed notification: email=%s scan_id=%s status=%s", userEmail, scan.ScanID, scan.Status)
	return nil
}
