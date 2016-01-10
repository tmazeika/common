package storage

import (
    "os"
    "io/ioutil"
    "path/filepath"
    "encoding/json"
    "os/user"
)

const ConfigMode = 0644

type Config map[string]interface{}

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
func New(appDir string, defConf Config) (*Storage, error) {
    s := Storage{
        config: defConf,
    }

    return &s, s.createAppDir(appDir)
}

func (s *Storage) Config() (Config, error) {
    const Name = "config.json"

    if err := s.loadConfig(Name); err != nil {
        return err
    }

    return s.config, nil
}

func (s *Storage) loadConfig(name string) error {
    if ! s.fileExists(name) {
        err := s.saveConfig(name)

        if err != nil {
            return err
        }
    }

    file, err := s.file(name, ConfigMode)

    if err != nil {
        return err
    }

    defer file.Close()
    return json.NewDecoder(file).Decode(&s.config)
}

func (s *Storage) saveConfig(name string) error {
    data, err := json.MarshalIndent(&s.config, "", "  ")

    if err != nil {
        return err
    }

    return ioutil.WriteFile(s.filePath(name), data, ConfigMode)
}

func (s *Storage) createAppDir(path string) error {
    const (
        Mode = 0700
        DefName = ".transhift"
    )

    if len(path) == 0 {
        user, err := user.Current()

        if err != nil {
            return err
        }

        s.AppDir = filepath.Join(user.HomeDir, DefName)
    } else {
        s.AppDir = path
    }

    if dirExists(s.AppDir) {
        return nil
    }

    return os.MkdirAll(s.AppDir, Mode)
}

func (s Storage) file(name string, mode os.FileMode) (*os.File, error) {
    return file(s.filePath(name), mode)
}

func (s Storage) fileExists(name string) bool {
    return fileExists(s.filePath(name))
}

func (s Storage) filePath(name string) string {
    return filepath.Join(s.AppDir, name)
}

func file(path string, mode os.FileMode) (file *os.File, err error) {
    if fileExists(path) {
        return os.Open(path), nil
    }

    file, err = os.Create(path)

    if err != nil {
        return
    }

    err = file.Chmod(mode)
    return
}

//func (s Storage) Certificate(certFileName, keyFileName string) (tls.Certificate, error) {
//    dir, err := s.dir()
//
//    if err != nil {
//        return tls.Certificate{}, err
//    }
//
//    certFilePath := filepath.Join(dir, certFileName)
//    keyFilePath := filepath.Join(dir, keyFileName)
//
//    if ! fileExists(certFilePath) || ! fileExists(keyFilePath) {
//        fmt.Print("Generating crypto... ")
//
//        keyData, certData, err := createCertificate()
//
//        if err != nil {
//            return tls.Certificate{}, err
//        }
//
//        err = ioutil.WriteFile(certFilePath, certData, 0600)
//
//        if err != nil {
//            return tls.Certificate{}, err
//        }
//
//        err = ioutil.WriteFile(keyFilePath, keyData, 0600)
//
//        if err != nil {
//            return tls.Certificate{}, err
//        }
//
//        fmt.Println("done")
//    }
//
//    return tls.LoadX509KeyPair(certFilePath, keyFilePath)
//}

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
