package main

import (
	"fmt"
	"forcloudy/monitoring/utils"
)

func main() {
	for {
		fmt.Println(utils.CpuUsageContainerUnix("4c168fc6095a69f2f76d8da3dc29ec5d122db563f1532df82a334daf99d1e9d3"))
	}

	fmt.Println("vim-go")
}
