package main

import (
	"fmt"
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
