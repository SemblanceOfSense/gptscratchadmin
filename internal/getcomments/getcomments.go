package getcomments

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type author struct {
    Id int
    Username string
    Scratchteam bool
    Image string
}

type ScratchComment []struct {
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

func GetComments(username string, projectId int) (ScratchComment, error) {
    req, err := http.NewRequest("GET", "https://api.scratch.mit.edu/users/" + username + "/projects/" + strconv.Itoa(projectId) + "/comments", nil)
    if err != nil { return ScratchComment{}, err }

    resp, err := http.DefaultClient.Do(req)
    if err != nil { return ScratchComment{}, err }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)

    commentsStruct := &ScratchComment{}

    json.Unmarshal(body, commentsStruct)

    return *commentsStruct, nil
}
