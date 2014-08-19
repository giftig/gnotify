package gcalendar

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"

  "code.google.com/p/goauth2/oauth"
  "github.com/skratchdot/open-golang/open"

  "xantoria.com/config"
)

type CalendarEvents struct {
  Kind, Updated string
  Items []CalendarEvent
}
type CalendarEvent struct {
  Status, Summary string
  Start, End string
}


func authenticate() (transport *oauth.Transport) {
  googleConfig := config.Config.Auth.Google
  code := googleConfig.Account.Code

  // Configure and create the OAuth Transport
  oauthConfig := &oauth.Config{
    ClientId: googleConfig.ClientID,
    ClientSecret: googleConfig.Secret,
    RedirectURL: googleConfig.RedirectURI,
    Scope: googleConfig.Scope,
    AuthURL: googleConfig.AuthEndpoint,
    TokenURL: googleConfig.TokenEndpoint,
    TokenCache: oauth.CacheFile("_oauth_cache.json"),
  }
  transport = &oauth.Transport{Config: oauthConfig}

  token, err := oauthConfig.TokenCache.Token()

  // We don't have a cached token: we'll need to request one
  if err != nil {
    if code == "" {
      log.Print("The account code needs to be set in the config.")
      open.Run(oauthConfig.AuthCodeURL(""))
      return
    }
    if token, err = transport.Exchange(code); err != nil {
      log.Fatal("OAuth token exchange failed", err)
    }
  }

  transport.Token = token
  return
}

func GetCalendar() {
  transport := authenticate()

  r, err := transport.Client().Get(fmt.Sprintf(
    "https://www.googleapis.com/calendar/v3/calendars/%s/events",
    config.Config.Auth.Google.Account.CalendarID,
  ))
  if err != nil {
    log.Fatal("SYNC: Request failed:", err)
  }
  defer r.Body.Close()

  responseText, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Print("SYNC: Error reading response body")
  }

  // FIXME: Don't just print the data; create Notifications from the data
  // Make sure no duplicate Notifications are created though
  var data CalendarEvents
  json.Unmarshal(responseText, &data)
  fmt.Printf("DATA ===== %s", responseText)
  fmt.Println()
}
