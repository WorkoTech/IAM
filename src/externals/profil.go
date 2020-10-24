package externals

import (
    "bytes"
    "os"
    "net/http"
    "encoding/json"

    "github.com/rs/zerolog/log"
)

func CreateUserProfil(token string) error {
    client := &http.Client{}

    req, err := http.NewRequest("POST", "http://" + os.Getenv("GAMIFICATION_HOST") + ":" + os.Getenv("GAMIFICATION_PORT") + "/profile", nil)
    if err != nil {
        return err
    }

    log.Info().Msgf("user %v", token)
    req.Header.Add("Authorization", token)

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}


func LogLoginEvent(token string) error {
    client := &http.Client{}

    payload := map[string]string{"actionTitle": "Total login days"}
    jsonPayload, _ := json.Marshal(payload)

    url := "http://" + os.Getenv("GAMIFICATION_HOST") + ":" + os.Getenv("GAMIFICATION_PORT") + "/action/done"
    req, err := http.NewRequest("POST", url , bytes.NewBuffer(jsonPayload))
    if err != nil {
        return err
    }

    log.Info().Msgf("user %v", token)
    req.Header.Add("Authorization", token)
    req.Header.Add("Content-type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    resp.Body.Close()
    return nil
}
