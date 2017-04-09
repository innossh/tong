package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenBrowser opens the site of specified URL in web browser
func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		fmt.Printf("Go to the following link in your browser\n%v\n", url)
	}
}
