package flagcomment

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type Choice struct {
    Index int
    Message Message
    Logprobs string
    Finish_reason string
}

type Usage struct {
    Prompt_tokens int
    completion_tokens int
    total_tokens int
}

type Response struct {
    Id string
    Object string
    Created int
    Model string
    Choices []Choice
    Usage Usage
    system_fingerprint string
}

type Message struct {
    Role string `json:"role"`
    Content string `json:"content"`
}
type OpenAiRequest struct {
    Model string `json:"model"`
    Messages []Message `json:"messages"`
}

func FlagComment(comment string, openaikey string) (bool, error) {
    messages := []Message{{Role: "system", Content: "You are an admin for the comments section of a kid friendly coding website where users make animations and games from art drawn themselves and colorful blocks that represent coding languages. I will send a comment to you. If that comments contain a link, a link that has been misspelled in order to conceal that it is a link, homophobic, transphobic, racist, anti-furry, or misogynistic content, trys to promote another game or piece of art, has sexual themes, references a tv show with sexual themes, refrences violence, encourages or partakes in online dating, uses excessively random characters, or asks for other users to follow the commenter, respond with the single word: \"true\". Otherwise respond with the single word: \"false\"."}}
    if strings.Contains(comment, "http") {
        return true, nil
    }
    if openaikey != "" {
        messages = append(messages, Message{Role: "user", Content: comment})
        requestBody := OpenAiRequest{Model: "gpt-3.5-turbo", Messages: messages}

        jsonData, err := json.Marshal(requestBody)
        if err != nil {
            return false, err
        }

        url := "https://api.openai.com/v1/chat/completions"
        req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer " + openaikey)
        if err != nil {
            return false, err
        }

        client := http.Client{Timeout: 10 * time.Second}
        res, err := client.Do(req)
        if err != nil {
            return false, err
        }
        defer res.Body.Close()
        bodyStruct := &Response{}
        body, err := io.ReadAll(res.Body)
        err = json.Unmarshal(body, bodyStruct)
        if err != nil { return false, err }

        if (bodyStruct.Choices[0].Message.Content == "true") { return true, nil }
    }

    return false, nil
}
