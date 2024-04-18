package models

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var Store = sessions.NewCookieStore([]byte(JWTSecretKey))

func SetJWTSessionStorage(w http.ResponseWriter, r *http.Request, token string) error {
	session, err := Store.Get(r, "jwt_session")
	if err != nil {
		return err
	}
	session.Values["jwt_token"] = token

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetJWTFromSessionStorage(r *http.Request) (string, error) {
	session, err := Store.Get(r, "jwt_session")
	if err != nil {
		return "", err
	}
	if token, ok := session.Values["jwt_token"].(string); ok {
		return token, nil
	}
	return "", nil
}

func ClearJWTSessionStorage(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, "jwt_session")
	if err != nil {
		return err
	}
	delete(session.Values, "jwt_token")

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}
