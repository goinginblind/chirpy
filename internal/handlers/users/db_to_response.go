package users

import "github.com/goinginblind/chirpy/internal/database"

/* jesus cant sqlc return just one type for all the queries that are returning the same
fields from the db. The naming here makes me want to shoot myself in the head BRUH.
TODO: Rethink your (life choices) function names
*/

func convertCreRowToCreateParams(user database.CreateUserRow) createUserParams {
	return createUserParams{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
	}
}

func convertLogRowToCreateParams(user database.ChangeUserLoginInfoRow) createUserParams {
	return createUserParams{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
		CreatedAt:   user.CreatedAt.String(),
		UpdatedAt:   user.UpdatedAt.String(),
	}
}

func convertUserToLoginParams(user database.User, token, refreshToken string) loginUserParams {
	return loginUserParams{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}
}
