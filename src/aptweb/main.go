package main

import (
	"fmt"
	"os"
)

func main() {
	// data := ParseDetail(os.Stdin)
	// for k, v := range data {
	// 	fmt.Printf("[%s] %s\n", k, v)
	// }
	data := ParseInstall(os.Stdin)
	for _, url := range data.Urls {
		fmt.Printf("[%s] %s (%s) %s\n", url.Name, url.Url, url.Size, url.Hash)
	}
	//for _, pkg := range data.Packages[GROUP_SUGGESTED] {
	//	fmt.Printf("- %s\n", pkg)
	//}
}
