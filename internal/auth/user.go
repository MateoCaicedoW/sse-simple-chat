package auth

import (
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func SetUserID(session *sessions.Session) string {
	if session.Values["user_id"] != nil {
		return session.Values["user_id"].(string)
	}

	id := uuid.New().String()
	session.Values["user_id"] = id
	return id
}

func GetUserID(session *sessions.Session) string {
	if session.Values["user_id"] == nil {
		return ""
	}

	return session.Values["user_id"].(string)
}
