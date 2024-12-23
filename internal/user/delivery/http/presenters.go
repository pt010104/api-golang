package http

import (
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/mongo"
)

type signupReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type forgetPasswordReq struct {
	Email string `json:"email" binding:"required"`
}
type signinReq struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	SessionID string `json:"session_id"`
}
type verifyRequestReq struct {
	Email string `json:"email" binding:"required"`
}

func (r signupReq) toInput() user.CreateUserInput {
	return user.CreateUserInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r signinReq) toInput() user.SignInType {
	return user.SignInType{
		Email:     r.Email,
		Password:  r.Password,
		SessionID: r.SessionID,
	}
}

type verifyUserReq struct {
	UserID string
	Token  string
}

func (r verifyUserReq) validate() error {
	if r.UserID == "" {
		return errWrongHeader
	}

	if r.Token == "" {
		return errWrongQuery
	}

	return nil
}

type resetPasswordReq struct {
	UserID      string
	NewPassword string `json:"new_password" binding:"required"`
	Token       string
}

func (r resetPasswordReq) toInput() user.ResetPasswordInput {
	return user.ResetPasswordInput{
		UserId:  r.UserID,
		NewPass: r.NewPassword,
		Token:   r.Token,
	}
}
func (r verifyUserReq) toInput() user.VerifyUserInput {
	return user.VerifyUserInput{
		UserId: r.UserID,
		Token:  r.Token}
}

type signUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h handler) newSignUpResponse(u models.User) signUpResponse {
	return signUpResponse{
		ID:    u.ID.Hex(),
		Email: u.Email,
		Name:  u.Name,
	}
}

type address_obj struct {
	ID       string `json:"id"`
	Street   string `json:"street"`
	District string `json:"district"`
	City     string `json:"city"`
	Province string `json:"province"`
	Phone    string `json:"phone"`
	Default  bool   `json:"default"`
}

type detailResp struct {
	ID      string        `json:"id"`
	Email   string        `json:"email"`
	Name    string        `json:"name"`
	Avatar  *avatar_obj   `json:"avatar,omitempty"`
	Address []address_obj `json:"addressess,omitempty"`
}

func (h handler) newDetailResp(u user.DetailUserOutput) detailResp {
	var avatar *avatar_obj
	if u.Avatar_URL != "" {
		avatar = &avatar_obj{
			MediaID: u.User.MediaID.Hex(),
			URL:     u.Avatar_URL,
		}
	}

	address := make([]address_obj, 0, len(u.User.Address))
	for _, addr := range u.User.Address {
		address = append(address, address_obj{
			ID:       addr.ID.Hex(),
			Street:   addr.Street,
			District: addr.District,
			City:     addr.City,
			Province: addr.Province,
			Phone:    addr.Phone,
			Default:  addr.Default,
		})
	}

	return detailResp{
		ID:      u.User.ID.Hex(),
		Email:   u.User.Email,
		Name:    u.User.Name,
		Avatar:  avatar,
		Address: address,
	}
}

type distributeNewTokenReq struct {
	UserId       string
	SessionID    string
	RefreshToken string
}

func (r distributeNewTokenReq) validate() error {
	if r.UserId == "" || r.SessionID == "" || r.RefreshToken == "" {
		return errWrongHeader
	}

	return nil
}

func (r distributeNewTokenReq) toInput() user.DistributeNewTokenInput {
	return user.DistributeNewTokenInput{
		UserId:       r.UserId,
		SessionID:    r.SessionID,
		RefreshToken: r.RefreshToken,
	}
}

type avatar_obj struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}
type signInResp struct {
	ID        string     `json:"id"`
	SessionID string     `json:"session_id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Token     user.Token `json:"token"`
}

func (h handler) newSignInResp(output user.SignInOutput) signInResp {
	return signInResp{
		ID:        output.User.ID.Hex(),
		Email:     output.User.Email,
		Username:  output.User.Name,
		Token:     output.Token,
		SessionID: output.SessionID,
	}
}

type distributeNewTokenResp struct {
	NewAccessToken  string `json:"new_access_token"`
	NewRefreshToken string `json:"new_refresh_token"`
}

func (h handler) newDistributeNewTokenResp(output user.DistributeNewTokenOutput) distributeNewTokenResp {
	return distributeNewTokenResp{
		NewAccessToken:  output.Token.AccessToken,
		NewRefreshToken: output.Token.RefreshToken,
	}
}

type updateReq struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	MediaID string `json:"media_id"`
}

func (r updateReq) toInput() user.UpdateInput {
	return user.UpdateInput{
		MediaID: r.MediaID,
		Name:    r.Name,
		Email:   r.Email,
	}
}
func (r updateReq) validate() error {
	if r.MediaID != "" && !mongo.IsObjectID(r.MediaID) {
		return errWrongBody
	}

	return nil
}

type UpdateResp struct {
	MediaID string `json:"media_id"`
}

type addressReq struct {
	Street   string `json:"street" binding:"required"`
	District string `json:"district" binding:"required"`
	City     string `json:"city" binding:"required"`
	Province string `json:"province" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Default  bool   `json:"default"`
}

func (r addressReq) toInput() user.AddAddressInput {
	return user.AddAddressInput{
		Street:   r.Street,
		District: r.District,
		City:     r.City,
		Province: r.Province,
		Phone:    r.Phone,
		Default:  r.Default,
	}
}

type updateAddressReq struct {
	ID       string `json:"id" binding:"required"`
	Street   string `json:"street" binding:"required"`
	District string `json:"district" binding:"required"`
	City     string `json:"city" binding:"required"`
	Province string `json:"province" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Default  bool   `json:"default"`
}

func (r updateAddressReq) toInput() user.UpdateAddressInput {
	return user.UpdateAddressInput{
		AddressID: r.ID,
		Street:    r.Street,
		District:  r.District,
		City:      r.City,
		Province:  r.Province,
		Phone:     r.Phone,
		Default:   r.Default,
	}
}

type detailAddressResp struct {
	Addressess []address_obj `json:"addressess"`
}

func (h handler) newDetailAddressResp(output user.DetailAddressOutput) detailAddressResp {
	address := make([]address_obj, 0, len(output.Addressess))
	for _, addr := range output.Addressess {
		address = append(address, address_obj{
			ID:       addr.ID.Hex(),
			Street:   addr.Street,
			District: addr.District,
			City:     addr.City,
			Province: addr.Province,
			Phone:    addr.Phone,
			Default:  addr.Default,
		})
	}
	return detailAddressResp{Addressess: address}
}
