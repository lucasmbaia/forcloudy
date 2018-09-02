package monit

import (
  "testing"
  "context"
)

func TestRunMonit(t *testing.T) {
  Run(context.Background(), 2)
}
