package main

import "sync"

type SafeUrlMap struct {
	visited map[string]bool
	mu      sync.Mutex
}
