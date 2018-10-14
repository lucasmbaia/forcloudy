package els

import (
  "testing"
  "context"
)

func TestNewClient(t *testing.T) {
  if _, err := NewClient(context.Background(), "http://172.16.95.185:9200"); err != nil {
    t.Fatal(err)
  }
}
