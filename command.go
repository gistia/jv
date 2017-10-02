package main

// Quit closes the main view
func Quit(args []string) bool {
	// Close the main view
	CurView().Quit()
	return false
}
