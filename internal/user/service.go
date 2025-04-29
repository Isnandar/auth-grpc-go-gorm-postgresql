package user

import pbuser "auth/pb/user"

func toProto(u *UserModel) *pbuser.User {
	return &pbuser.User{
		RoleId:   uint32(u.RoleID),
		RoleName: u.Role.Name,
		Name:     u.Name,
		Email:    u.Email,
	}
}
