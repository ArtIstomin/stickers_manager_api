package actions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
)

var botToken = envy.Get("TELEGRAM_BOT_TOKEN", "")

type telegramClaims struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AuthDate  int64  `json:"authDate"`
	jwt.StandardClaims
}

// AuthCallback default implementation.
func AuthCallback(c buffalo.Context) error {
	urlParams := c.Params().(url.Values)

	if invalidAuth := !checkTelegramAuthorization(urlParams); invalidAuth {
		return c.Error(http.StatusUnprocessableEntity, errors.New("Invalid data"))
	}

	id, _ := strconv.Atoi(urlParams.Get("id"))
	authDate, _ := strconv.ParseInt(urlParams.Get("auth_date"), 10, 64)

	claims := telegramClaims{
		id,
		urlParams.Get("first_name"),
		urlParams.Get("last_name"),
		authDate,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(14 * 24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(envy.Get("JWT_SECRET", "")))

	if err != nil {
		return err
	}

	return c.Render(200, r.JSON(map[string]string{"token": ss}))
}

func checkTelegramAuthorization(params url.Values) bool {
	var telegramHash string
	var strs []string

	for k, v := range params {
		if k == "hash" {
			telegramHash = v[0]
			continue
		}
		strs = append(strs, k+"="+v[0])
	}

	sort.Strings(strs)

	sha256hash := sha256.New()
	sha256hash.Write([]byte(botToken))
	mac := hmac.New(sha256.New, sha256hash.Sum(nil))
	mac.Write([]byte(strings.Join(strs, "\n")))
	paramsHash := hex.EncodeToString(mac.Sum(nil))

	return telegramHash == paramsHash
}
