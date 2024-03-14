package main

import (
	"fmt"
	"gptscratchadmin/internal/flagcomment"
	"gptscratchadmin/internal/getcomments"
	"log"
	"strconv"
	"time"
)

func main() {
    fmt.Println("--------------------------------------------------------\nInput format:\nUsername of project-to-scan's owner (Enter)\nProject-to-scan's ID (Enter)\nFlag comments since _ hours ago (-1 for all comments)\nOpenai key (must be a tier 1 account, leave blank for no ai moderation)\n--------------------------------------------------------")
    var username, readid, readhours, openaikey string
    _, err := fmt.Scan(&username, &readid, &readhours, &openaikey)
    if err != nil {
        log.Fatal(err)
    }
    id, _ := strconv.Atoi(readid)
    hours, _ := strconv.Atoi(readhours)

    fmt.Println("--------------------------------------------------------")

    values, err := getcomments.GetComments(username, id, hours)
    if err != nil { log.Fatal(err) }
    for _, val := range values {
        flag, err := flagcomment.FlagComment(val.Content, openaikey)
        if err != nil { log.Fatal(err) }
        if (flag) { fmt.Println(val.Content) }
        time.Sleep(time.Millisecond * 250)
    }
}
