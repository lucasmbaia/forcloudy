package docker

import (
	"fmt"
	"testing"
)

func TestListAllContainers(t *testing.T) {
	fmt.Println(ListAllContainers("lucas"))
}
