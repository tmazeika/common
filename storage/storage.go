package storage

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
    AppDir string
    config Config
}

// New creates a Storage struct with the given application directory.
//
// If `appDir` is empty:
// - ~/.transhift will be created if it doesn't exist
// - files will be stored in ~/.transhift
//
// Otherwise:
// - the directory at `appDir` will be created if it doesn't exist, along with
//   its parents
// - files will be stored in the directory at `appDir`
func New(appDir string) (*Storage, error) {
    s := new(Storage)
    return s, s.createAppDir(appDir)
}

func (s *Storage) createAppDir(appDir string) error {
    const (
        Perm = 0700
        DefName = ".transhift"
    )

    if len(s.AppDir) == 0 {
        user, err := user.Current()

        if err != nil {
            return err
        }

        s.AppDir = filepath.Join(user.HomeDir, DefName)
    }

    if dirExists(s.AppDir) {
        return nil
    }

    return os.MkdirAll(s.AppDir, Perm)
}

func (s Storage) configFile() (*os.File, error) {
    const FileName = "config.json"
    dir, err := s.dir()

    if err != nil {
        return nil, err
    }

    filePath := filepath.Join(dir, FileName)

    if ! fileExists(filePath) {
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

    if ! fileExists(certFilePath) || ! fileExists(keyFilePath) {
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
    if fileExists(path) {
        return os.Open(path)
    }

    return os.Create(path)
}

func exists(path string, dir bool) bool {
    info, err := os.Stat(path)

    if err != nil {
        return false
    }

    if dir {
        return info.IsDir()
    }

    return info.Mode().IsRegular()
}

func fileExists(path string) bool {
    return exists(path, false)
}

func dirExists(path string) bool {
    return exists(path, true)
}
