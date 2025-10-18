package internal

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// 実行されるphpコマンドの前に指定するコマンドを設定
	CommandToSpecifyBeforePHPCommand string
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
	filePath := filepath.Join(homeDir, ".config", "ppaid.toml")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
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

	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return config, err
	}
	return config, nil
}
