package getcomments

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type author struct {
    Id int
    Username string
    Scratchteam bool
    Image string
}

type ScratchComment struct {
    Id int
    Parent_id int
    Commentee_id int
    Content string
    Datetime_created string
    Datetime_modified string
    Visibility string
    Author author
    Reply_count int
}

func GetComments(username string, projectId int, hours int) ([]ScratchComment, error) {
    allComments := []ScratchComment{}

    out:
    for i := 0; i < 100000000 /* arbitrarily large number */; i++ {
        url := "https://api.scratch.mit.edu/users/" + username + "/projects/" + strconv.Itoa(projectId) + "/comments" + "?offset=" + strconv.Itoa(i * 40) + "&limit=40"
        req, err := http.NewRequest("GET", url, nil)
        if err != nil { return []ScratchComment{}, err }

        resp, err := http.DefaultClient.Do(req)
        if err != nil { return []ScratchComment{}, err }
        defer resp.Body.Close()
        body, err := io.ReadAll(resp.Body)

        commentsStruct := []ScratchComment{}

        err = json.Unmarshal(body, &commentsStruct)
        if err != nil { return []ScratchComment{}, err }
        if (len(commentsStruct) == 0) { break }

        const layout = "2021-05-17T04:00:00.000Z"
        for _, v := range commentsStruct {
            t, _ := time.Parse(time.RFC3339Nano, v.Datetime_created)
            if ((time.Since(t).Hours() > float64(hours) && hours > -1)) { break out }
            allComments = append(allComments, v)
        }
        time.Sleep(time.Millisecond * 250)
    }

    return allComments, nil
}
