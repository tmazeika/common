package common

import (
    "os"
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "os/user"
)

type Config map[string]string

type Storage struct {
    CustomDir string
    Config    Config
}

func (s Storage) dir() (string, error) {
    const DefDirName = ".transhift"

    if len(s.CustomDir) == 0 {
        user, err := user.Current()

        if err != nil {
            return "", err
        }

        return getDir(filepath.Join(user.HomeDir, DefDirName))
    } else {
        return getDir(s.CustomDir)
    }
}

func (s Storage) configFile() (*os.File, error) {
    const FileName = "config.json"
    dir, err := s.dir()

    if err != nil {
        return nil, err
    }

    filePath := filepath.Join(dir, FileName)

    if ! fileExists(filePath, false) {
        data, err := json.MarshalIndent(&s.Config, "", "  ")

        if err != nil {
            return nil, err
        }

        err = ioutil.WriteFile(filePath, data, 0644)

        if err != nil {
            return nil, err
        }
    }

    return getFile(filePath)
}

func (s *Storage) LoadConfig() error {
    file, err := s.configFile()

    if err != nil {
        return err
    }

    defer file.Close()

    return json.NewDecoder(file).Decode(&s.Config)
}

func (s Storage) Certificate(certFileName, keyFileName string) (tls.Certificate, error) {
    dir, err := s.dir()

    if err != nil {
        return tls.Certificate{}, err
    }

    certFilePath := filepath.Join(dir, certFileName)
    keyFilePath := filepath.Join(dir, keyFileName)

    if ! fileExists(certFilePath, false) || ! fileExists(keyFilePath, false) {
        fmt.Print("Generating crypto... ")

        keyData, certData, err := createCertificate()

        if err != nil {
            return tls.Certificate{}, err
        }

        err = ioutil.WriteFile(certFilePath, certData, 0600)

        if err != nil {
            return tls.Certificate{}, err
        }

        err = ioutil.WriteFile(keyFilePath, keyData, 0600)

        if err != nil {
            return tls.Certificate{}, err
        }

        fmt.Println("done")
    }

    return tls.LoadX509KeyPair(certFilePath, keyFilePath)
}

func getFile(path string) (*os.File, error) {
    if fileExists(path, false) {
        return os.Open(path)
    }

    return os.Create(path)
}

func getDir(path string) (string, error) {
    if fileExists(path, true) {
        return path, nil
    }

    err := os.MkdirAll(path, 0700)

    if err != nil {
        return "", err
    }

    return path, nil
}

func fileExists(path string, asDir bool) bool {
    info, err := os.Stat(path)

    if err != nil {
        return false
    }

    if asDir {
        return info.IsDir()
    }

    return info.Mode().IsRegular()
}
