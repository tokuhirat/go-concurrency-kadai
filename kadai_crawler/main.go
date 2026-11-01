package main

import "sync"

type SafeUrlMap struct {
	visited map[string]bool
	mu      sync.Mutex
}

func (m *SafeUrlMap) Visit(url string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.visited[url] {
		return
	}
	m.visited[url] = true
}
