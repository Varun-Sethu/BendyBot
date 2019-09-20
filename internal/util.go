package internal

import (
    "os"
    "fmt"
    "io/ioutil"
    "go/build"
)



// Does what it says
func OpenFileFromStore(dest string) []byte {
	userDict, _ := os.OpenFile(GetAbsFile(fmt.Sprintf("data/%s.json", dest)), os.O_RDONLY|os.O_CREATE, 0666)
	defer userDict.Close()
	dictBytes, _ := ioutil.ReadAll(userDict)
	return dictBytes
}



// Returns the absolute path of an internal file store
func GetAbsFile(internalPath string) string {
    gopath := os.Getenv("GOPATH")
    if gopath == "" {
        gopath = build.Default.GOPATH
    }

    return gopath + "/src/bendy-bot/internal/storage/" + internalPath
}