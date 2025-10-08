package auth

import (
    "context"
    "encoding/json"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/github"
)

type GitHubUser struct {
    ID        int64  `json:"id"`
    Login     string `json:"login"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    AvatarURL string `json:"avatar_url"`
}

func GetGitHubOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
    return &oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        Scopes:       []string{"user:email", "repo"},
        Endpoint:     github.Endpoint,
    }
}

func GetGitHubUser(ctx context.Context, config *oauth2.Config, code string) (*GitHubUser, string, error) {
    token, err := config.Exchange(ctx, code)
    if err != nil {
        return nil, "", err
    }

    client := config.Client(ctx, token)
    resp, err := client.Get("https://api.github.com/user")
    if err != nil {
        return nil, "", err
    }
    defer resp.Body.Close()

    var user GitHubUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, "", err
    }

    // Get primary email if not public
    if user.Email == "" {
        emailResp, err := client.Get("https://api.github.com/user/emails")
        if err == nil {
            defer emailResp.Body.Close()
            var emails []struct {
                Email   string `json:"email"`
                Primary bool   `json:"primary"`
            }
            if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
                for _, email := range emails {
                    if email.Primary {
                        user.Email = email.Email
                        break
                    }
                }
            }
        }
    }

    return &user, token.AccessToken, nil
}