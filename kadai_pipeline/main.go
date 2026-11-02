package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
