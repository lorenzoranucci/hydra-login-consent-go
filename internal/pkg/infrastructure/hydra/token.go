package hydra

import (
	"fmt"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/pkg/domain"
)

type IDToken struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccessToken struct {
	Mercure Mercure `json:"mercure"`
	User    User    `json:"user"`
}

type Mercure struct {
	Subscribe []string `json:"subscribe"`
}

type User struct {
	Email string   `json:"email"`
	Id    int      `json:"id"`
	Roles []string `json:"roles"`
	UUID  string   `json:"uuid"`
}

func getIDTokenFromUser(user domain.User) IDToken {
	return IDToken{
		Name:  fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email: user.Email,
	}
}

func getAccessTokenFromUser(user domain.User) AccessToken {
	return AccessToken{
		Mercure: Mercure{
			Subscribe: []string{fmt.Sprintf("users/%s", user.UUID.String())},
		},
		User: User{
			Email: user.Email,
			Id:    user.ID,
			Roles: user.Roles,
			UUID:  user.UUID.String(),
		},
	}
}
