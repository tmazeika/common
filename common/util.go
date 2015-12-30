package common

func ArrayContains(array []Type, value Type) bool {
    for _, v := range array {
        if v == value {
            return true
        }
    }

    return false
}
