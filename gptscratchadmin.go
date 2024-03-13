package main

import (
	"fmt"
	"gptscratchadmin/internal/getcomments"
	"log"
)

func main() {
    values, err := getcomments.GetComments("AnimatorExpands", 975151602)
    if err != nil { log.Fatal(err) }
    for _, val := range values {
        fmt.Println(val.Content)
    }
}
