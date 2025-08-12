package config

import (
    "os"
    "encoding/json"
)

const configFileName = ".gatorconfig.json"

type Config struct {
    DBURL string `json:"db_url"`
    CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
    // read .gatorconfig.json and return contents as Config struct
    configPath, err := getConfigFilePath() 
    if err != nil {
        return Config{}, err
    }
    data, err := os.ReadFile(configPath)
    if err != nil {
        return Config{}, err
    }
    var cfg Config
    err = json.Unmarshal(data, &cfg) 
    if err != nil {
        return Config{}, err
    }
    return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
    // set the CurrentUserName field of a Config type
    cfg.CurrentUserName = userName
    err := write(*cfg)
    if err != nil {
        return err
    }
    return nil
}

func getConfigFilePath() (string, error) {
    // get the full file path of the .gatorconfig.json
    homePath, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    configPath := homePath + "/" + configFileName
    return configPath, nil
}

func write(cfg Config) error {
    // write cfg to file
    configPath, err := getConfigFilePath()
    if err != nil {
        return err
    }
    b, err := json.Marshal(cfg)
    if err != nil {
        return err
    }
    err = os.WriteFile(configPath, b, 0660)
    if err != nil {
        return err
    }
    return nil
}
