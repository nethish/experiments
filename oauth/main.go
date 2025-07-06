package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
    "os"
)

var (
    clientID     = os.Getenv("OAUTH_CLIENT_ID")
    clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
    redirectURI  = "http://localhost:8080/callback"
)

func main() {
    http.HandleFunc("/", handleHome)
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/callback", handleCallback)

    fmt.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    tmpl.Execute(w, nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user", clientID, url.QueryEscape(redirectURI))
    http.Redirect(w, r, authURL, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Missing code in callback", http.StatusBadRequest)
        return
    }

    // Exchange code for access token
    resp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
        "client_id":     {clientID},
        "client_secret": {clientSecret},
        "code":          {code},
        "redirect_uri":  {redirectURI},
    })
    if err != nil {
        http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    var tokenResp struct {
        AccessToken string `json:"access_token"`
        Scope       string `json:"scope"`
        TokenType   string `json:"token_type"`
    }

    // GitHub returns form-encoded, not JSON!
    buf := make([]byte, 2048)
    n, _ := resp.Body.Read(buf)
    values, _ := url.ParseQuery(string(buf[:n]))
    tokenResp.AccessToken = values.Get("access_token")

    // Use the access token to fetch user info
    req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
    req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)

    client := &http.Client{}
    userResp, err := client.Do(req)
    if err != nil {
        http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer userResp.Body.Close()

    var userData map[string]interface{}
    json.NewDecoder(userResp.Body).Decode(&userData)

    // Show user data
    tmpl := template.Must(template.ParseFiles("templates/welcome.html"))
    tmpl.Execute(w, userData)
}
