package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jabuta/dreampicai/pkg/database"
	"github.com/jabuta/dreampicai/pkg/db"
	"github.com/jabuta/dreampicai/pkg/validate"
	"github.com/jabuta/dreampicai/types"
	"github.com/jabuta/dreampicai/view/account"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleAccountIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, account.Index())
}

func HandleAccountCreate(w http.ResponseWriter, r *http.Request) error {
	params := account.AccountParams{
		Username: r.FormValue("username"),
	}
	accountErrors := account.AccountErrors{}
	if ok := validate.New(params, validate.Fields{
		"Username": validate.Rules(validate.Min(5), validate.Max(15)), //needs additional validationrule to exculde special characters and spaces
	}).Validate(&accountErrors); !ok {
		fmt.Println("account in bad format")
		fmt.Println(accountErrors)
		w.WriteHeader(http.StatusBadRequest)
		return render(r, w, account.AccountForm(params, accountErrors, false))
	}
	authedUser := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)

	//TODO create intermediary db fuctions to map apps structs to pgx structs
	if len(authedUser.Username) == 0 {
		_, err := db.Conf.DB.CreateUser(r.Context(), database.CreateUserParams{
			UserID: pgtype.UUID{
				Bytes: authedUser.UserID,
				Valid: true,
			},
			Username: params.Username,
			Email:    authedUser.Email,
			CreatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err)
			return err
		}
	} else {
		_, err := db.Conf.DB.UpdateUser(r.Context(), database.UpdateUserParams{
			UserID: pgtype.UUID{
				Bytes: authedUser.UserID,
				Valid: true,
			},
			Username: params.Username,
			Email:    authedUser.Email,
			UpdatedAt: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			handleError(w, r, http.StatusInternalServerError, err)
			return err
		}
	}

	return render(r, w, account.UpdateResponse())
}
