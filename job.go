package main

type JobFunction struct {
	function func(string, ...string)
	output   string
	args     []string
}
