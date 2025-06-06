package utils

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="refresh" content="0; url={{.AppURL}}" />
  <script>
    window.location.href = "{{.AppURL}}";
    setTimeout(function () {
      window.location.href = "{{.FallbackURL}}";
    }, 3000);
  </script>
  <title>Redirecting...</title>
</head>
<body>
  <p>Redirecting... If nothing happens, <a href="{{.AppURL}}">click here</a>.</p>
</body>
</html>
`

type RedirectData struct {
	AppURL      string
	FallbackURL string
}

func ResetRedirectHandler(c *gin.Context) {
	requestId := c.Query("request_id")
	if requestId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing request_id"})
		return
	}

	appURL := fmt.Sprintf("myhome://auth/reset-password?request_id=%s", requestId)
	playStoreURL := "https://play.google.com/store/apps/details?id=com.morg.home" // change this to your app's package name

	tmpl, err := template.New("redirect").Parse(htmlTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(c.Writer, RedirectData{AppURL: appURL, FallbackURL: playStoreURL}); err != nil {
		c.String(http.StatusInternalServerError, "Template execute error")
	}
}
