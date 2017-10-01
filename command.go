package main

// Quit closes the main view
func Quit(args []string) {
	// Close the main view
	CurView().Quit(true)
}
