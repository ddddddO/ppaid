package internal

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
)

var ConfigFilePath = ""

type Config struct {
	// 実行されるphpコマンドの前に指定するコマンドを設定
	CommandToSpecifyBeforePHPCommand string

	// 直前に実行されたpucoで選択されたテストファイルパスとカバレッジ対象のパスを残す
	// 設定ファイルに微妙な感じだけど...あと、プロジェクト関係なくなのを何とかしたい気もする
	LatestExecutedData struct {
		SelectedTestFilePaths       []string
		SelectedCoverageTargetPaths []string
	}
}

// 初回起動時に設定ファイルがなかったら以下の値の設定ファイルが作成される
// TODO: それが読み込まれて、実行されるから、不要な設定は消してくださいとREADMEに書いておく
func getDefaultConfig() Config {
	return Config{
		CommandToSpecifyBeforePHPCommand: "docker compose exec app",
	}
}

func LoadConfig() (Config, error) {
	var config Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, err
	}
	ConfigFilePath = filepath.Join(homeDir, ".config", "puco.toml")
	if _, err := os.Stat(ConfigFilePath); os.IsNotExist(err) {
		file, err := os.Create(ConfigFilePath)
		if err != nil {
			return config, err
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		config := getDefaultConfig()
		if err := encoder.Encode(config); err != nil {
			return config, err
		}
		return config, nil
	} else if err != nil {
		return config, err
	}

	if _, err := toml.DecodeFile(ConfigFilePath, &config); err != nil {
		return config, err
	}
	return config, nil
}

func StoreConfig(cfg Config) error {
	file, err := os.OpenFile(ConfigFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(cfg)
}

// 直前の実行時に選択されたテストファイルと合致してるか
func (c Config) IsMatchedTestFile(currentTestFilePath string) bool {
	return slices.Contains(c.LatestExecutedData.SelectedTestFilePaths, currentTestFilePath)
}

// 直前の実行時に選択されたカバレッジ取得対象ファイルと合致してるか
func (c Config) IsMatchedCoverageTargetFile(currentCoverageFilePath string) bool {
	return slices.Contains(c.LatestExecutedData.SelectedCoverageTargetPaths, currentCoverageFilePath)
}
