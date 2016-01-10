package protocol

//import (
//    "crypto/sha256"
//    "os"
//    "io"
//)
//
//func CalculateFileHash(file *os.File) ([]byte, error) {
//    hash := sha256.New()
//
//    if _, err := io.Copy(hash, file); err != nil {
//        return nil, err
//    }
//
//    return hash.Sum(nil), nil
//}
//
//func CalculateBytesHash(bytes []byte) ([]byte, error) {
//    hash := sha256.New()
//
//    if _, err := hash.Write(bytes); err != nil {
//        return nil, err
//    }
//
//    return hash.Sum(nil), nil
//}
