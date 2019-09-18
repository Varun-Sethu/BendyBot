package internal

import (
	"bendy-bot/internal/markov"
    "encoding/base64"
    "encoding/gob"
    "bytes"
    "os"
    "fmt"
    "io/ioutil"
    "go/build"
)



// go binary encoder, this is used when storing structs within a file
func ToGOB64(m markov.Markov) string {
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    err := e.Encode(m)
    if err != nil { fmt.Println(`failed gob Encode`, err) }
    return base64.StdEncoding.EncodeToString(b.Bytes())
}



// go binary decoder, this is used when retreiving structs from a file
func FromGOB64(str string) markov.Markov {
    m := markov.Markov{}
    by, err := base64.StdEncoding.DecodeString(str)
    if err != nil { fmt.Println(`failed base64 Decode`, err); }
    b := bytes.Buffer{}
    b.Write(by)
    d := gob.NewDecoder(&b)
    err = d.Decode(&m)
    if err != nil { fmt.Println(`failed gob Decode`, err); }
    return m
}



// Does what it says
func OpenFileFromStore(dest string) []byte {
	userDict, _ := os.OpenFile(GetAbsFile(fmt.Sprintf("data/%s.dict", dest)), os.O_RDONLY|os.O_CREATE, 0666)
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