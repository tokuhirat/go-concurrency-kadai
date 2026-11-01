package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

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

// URLを受け取り、そのページに含まれるリンク（文字列のスライス）を返します
func fetch(url string) []string {
	// 処理時間をランダムにシミュレーション（0〜1.3秒）
	// ※タイムアウト(1秒)の挙動を確認するため
	time.Sleep(time.Duration(rand.Intn(1300)) * time.Millisecond)

	// 擬似的なリンクデータ
	links := map[string][]string{
		"https://example.com":   {"https://example.com/a", "https://example.com/b"},
		"https://example.com/a": {"https://example.com/c"},
		"https://example.com/b": {"https://example.com/a"}, // 重複リンク
		"https://example.com/c": {},
	}

	return links[url]
}

func Crawl(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("取得開始:", url)
	ch := make(chan []string, 1)

	go func() {
		ch <- fetch(url)
	}()

	select {
	case links := <-ch:
		for _, link := range links {
			wg.Add(1)
			go Crawl(link, wg)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("タイムアウト:", url)
	}
}
