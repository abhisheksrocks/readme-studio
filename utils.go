package main

func empty(s string) bool {
	return len(s) == 0
}

func notEmpty(s string) bool {
	return !empty(s)
}
