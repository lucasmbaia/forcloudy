package monit

import (
	"context"
	"testing"
)

func TestRunMonit(t *testing.T) {
	Run(context.Background(), Config{
		Running: 2,
		Topic:   "monitorin",
		Key:     "containers",
	})
}
