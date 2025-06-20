package users

import "github.com/goinginblind/chirpy/internal/database"

func dbUserRowToCreateParams(user database.CreateUserRow) createUserParams {
	return createUserParams{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

func dbUserToLoginParams(user database.User, token, refreshToken string) loginUserParams {
	return loginUserParams{
		ID:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Token:        token,
		RefreshToken: refreshToken,
	}
}
