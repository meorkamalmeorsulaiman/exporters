package main

import (
    "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"crypto/tls"
)

// A Response struct to map the Entire Response
type Response struct {
    Component []Component `json:"components"`
}

// A Component Struct to map every pokemon to.
type Component struct {
    EntryNo string `json:"id"`
    Name string `json:"name"`
    Status string `json:status`
}

// // A struct to map our Pokemon's Species which includes it's name
// type PokemonSpecies struct {
//     Name string `json:"name"`
// }

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    response, err := http.Get("https://www.cloudflarestatus.com/api/v2/components.json")

    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }

    var responseObject Response
    json.Unmarshal(responseData, &responseObject)

    fmt.Println(len(responseObject.Component))

    for  i := 0; i < len(responseObject.Component); i++ {
        fmt.Print(responseObject.Component[i].Name, " ")
        fmt.Println(responseObject.Component[i].Status)
    }

}