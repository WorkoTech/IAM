package externals

import (
    "bytes"
    "os"
    "io/ioutil"
    "encoding/json"
    "net/http"

    "github.com/dgrijalva/jwt-go"
    "github.com/rs/zerolog/log"
)

type OfferItem struct {
    Resource string `json:"resource"`
    Limit int `json:"limit"`
}

type Offer struct {
    Items []OfferItem `json:"items"`
}

type Usage struct {
    WorkspaceCount int `json:"workspace_count"`
    StorageSizeCount int `json:"storage_size_count"`
    DocumentCount int `json:"document_count"`
}

type Billing struct {
    Offer Offer `json:"offer"`
    Usage Usage `json:"usage"`
}

func RetrieveWorkspaceUsage(user *jwt.Token, workspaceId int) (*Billing, error) {
    client := &http.Client{}

    payload := map[string]int{"workspaceId": workspaceId}
    jsonPayload, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST","http://" + os.Getenv("BILLING_HOST") + ":" + os.Getenv("BILLING_PORT") + "/usage/workspace",bytes.NewBuffer(jsonPayload))
    if err != nil {
        return nil, err
    }

    log.Info().Msgf("user %v", user)
    req.Header.Add("Authorization", "Bearer " + user.Raw)
    req.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error().Msgf("error : %v", err)
    }

    var target Billing
    err = json.Unmarshal(bodyBytes, &target)
    if err != nil {
        log.Error().Msgf("error : %v", err.Error())
    }
    log.Info().Msgf("body : %v", target)
    return &target, nil
}

func RetrieveUserUsage(user *jwt.Token) (*Billing, error) {
    client := &http.Client{}

    req, err := http.NewRequest("POST","http://" + os.Getenv("BILLING_HOST") + ":" + os.Getenv("BILLING_PORT") + "/usage/user", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", "Bearer " + user.Raw)
    req.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Error().Msgf("error : %v", err)
    }
    log.Info().Msgf("body : %v", string(bodyBytes))

    var target Billing
    err = json.Unmarshal(bodyBytes, &target)
    if err != nil {
        log.Error().Msgf("error : %v", err.Error())
    }
    log.Info().Msgf("body : %v", target)
    return &target, nil
}

