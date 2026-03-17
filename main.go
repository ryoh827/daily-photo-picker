package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var supportedExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".bmp": true, ".webp": true, ".heic": true, ".heif": true,
	".tiff": true, ".tif": true,
}

type State struct {
	Used []string `json:"used"`
}

func baseDir() string {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, "実行ファイルのパスを取得できません:", err)
		os.Exit(1)
	}
	return filepath.Dir(exe)
}

func photosDir() string {
	return filepath.Join(baseDir(), "photos")
}

func stateFile() string {
	return filepath.Join(baseDir(), "used.json")
}

// photos/ 以下の画像ファイルを再帰的に取得
func getAllPhotos() ([]string, error) {
	dir := photosDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}

	var photos []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if supportedExts[ext] {
			rel, _ := filepath.Rel(baseDir(), path)
			photos = append(photos, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(photos)
	return photos, nil
}

func loadState() State {
	data, err := os.ReadFile(stateFile())
	if err != nil {
		return State{}
	}
	var s State
	json.Unmarshal(data, &s)
	return s
}

func saveState(s State) {
	data, _ := json.MarshalIndent(s, "", "  ")
	os.WriteFile(stateFile(), data, 0644)
}

func cmdPick() {
	all, err := getAllPhotos()
	if err != nil {
		fmt.Fprintln(os.Stderr, "写真の取得に失敗しました:", err)
		os.Exit(1)
	}
	if len(all) == 0 {
		fmt.Printf("写真が見つかりません。%s/ にファイルを配置してください。\n", photosDir())
		os.Exit(1)
	}

	state := loadState()

	// フォルダから削除された写真を used から除去
	allSet := make(map[string]bool, len(all))
	for _, p := range all {
		allSet[p] = true
	}
	var cleanUsed []string
	for _, p := range state.Used {
		if allSet[p] {
			cleanUsed = append(cleanUsed, p)
		}
	}
	state.Used = cleanUsed

	// 未使用の写真を抽出
	usedSet := make(map[string]bool, len(state.Used))
	for _, p := range state.Used {
		usedSet[p] = true
	}
	var available []string
	for _, p := range all {
		if !usedSet[p] {
			available = append(available, p)
		}
	}

	if len(available) == 0 {
		fmt.Printf("全 %d 枚の写真を使い切りました。\n", len(all))
		fmt.Println("リセットするには: daily-photo-picker reset")
		os.Exit(0)
	}

	chosen := available[rand.Intn(len(available))]
	state.Used = append(state.Used, chosen)
	saveState(state)

	fullPath := filepath.Join(baseDir(), chosen)

	fmt.Printf("📷 選ばれた写真: %s\n", chosen)
	fmt.Printf("   パス: %s\n", fullPath)
	fmt.Printf("   残り: %d / %d 枚\n", len(available)-1, len(all))
}

func cmdStatus() {
	all, _ := getAllPhotos()
	state := loadState()

	allSet := make(map[string]bool, len(all))
	for _, p := range all {
		allSet[p] = true
	}
	usedCount := 0
	for _, p := range state.Used {
		if allSet[p] {
			usedCount++
		}
	}

	fmt.Printf("写真フォルダ: %s\n", photosDir())
	fmt.Printf("全写真数:     %d 枚\n", len(all))
	fmt.Printf("使用済み:     %d 枚\n", usedCount)
	fmt.Printf("残り:         %d 枚\n", len(all)-usedCount)
}

func cmdReset() {
	saveState(State{})
	fmt.Println("✓ 使用済みリストをリセットしました。")
}

func main() {
	cmd := "pick"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "pick", "":
		cmdPick()
	case "status":
		cmdStatus()
	case "reset":
		cmdReset()
	case "help", "--help", "-h":
		fmt.Println("daily-photo-picker - 写真ランダム選択ツール")
		fmt.Println()
		fmt.Println("使い方:")
		fmt.Println("  daily-photo-picker          未使用の写真からランダムに1枚選ぶ")
		fmt.Println("  daily-photo-picker status   現在の状態を表示する")
		fmt.Println("  daily-photo-picker reset    使用済みリストをリセットする")
		fmt.Println("  daily-photo-picker help     このヘルプを表示する")
		fmt.Println()
		fmt.Printf("写真フォルダ: %s\n", photosDir())
	default:
		fmt.Fprintf(os.Stderr, "不明なコマンド: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "daily-photo-picker help でヘルプを表示")
		os.Exit(1)
	}
}
