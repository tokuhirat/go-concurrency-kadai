package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// テスト用のダミーファイルを作成する関数
func generateDummyFiles(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// 5つのファイルを作成
	for i := 1; i <= 5; i++ {
		filename := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
		// ダミーの単語データ（"word" という単語を i*10 個書き込む）
		content := strings.Repeat("word ", i*10)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func findFiles(ctx context.Context, root string, paths chan<- string) error {
	defer close(paths)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// ディレクトリはスキップ
		if d.IsDir() {
			return nil
		}
		// .txtファイルのみ対象
		if filepath.Ext(path) == ".txt" {
			select {
			case <-ctx.Done():
				// キャンセルされたら探索を打ち切る
				return ctx.Err()
			case paths <- path: // パスを送信
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("探索エラー: %v\n", err)
		return err
	}
	return nil
}

func countFileWords(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	cnt := 0
	for scanner.Scan() {
		cnt++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return cnt, nil
}

func countWords(ctx context.Context, paths <-chan string, numWords chan<- int, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		select {
		case path, ok := <-paths:
			if !ok {
				return nil
			}

			cnt, err := countFileWords(path)
			if err != nil {
				return err
			}
			numWords <- cnt

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func main() {
	const root string = "./"
	err := generateDummyFiles(root)
	if err != nil {
		fmt.Println("テストデータ生成に失敗しました")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	paths := make(chan string)
	numWords := make(chan int)

	go findFiles(ctx, root, paths)

	numWorker := 4
	var wg sync.WaitGroup
	wg.Add(numWorker)
	for range numWorker {
		go countWords(ctx, paths, numWords, &wg)
	}

	go func() {
		wg.Wait()
		close(numWords)
	}()

	totalNumWords := 0
	numFiles := 0
	for count := range numWords {
		numFiles++
		totalNumWords += count
	}
	fmt.Println("処理完了")
	fmt.Printf("処理ファイル数: %d\n", numFiles)
	fmt.Printf("総単語数: %d\n", totalNumWords)
}
