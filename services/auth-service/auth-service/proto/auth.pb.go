package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type RegisterRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Email    string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type LoginRequest struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type User struct {
	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username  string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email     string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
}

func (x *User) ProtoReflect() protoreflect.Message {
	return nil
}

type AuthResponse struct {
	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	User  *User  `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	return nil
}

type UserInfoRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *UserInfoRequest) Reset() {
	*x = UserInfoRequest{}
}

func (x *UserInfoRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type UserInfoResponse struct {
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *UserInfoResponse) Reset() {
	*x = UserInfoResponse{}
}

func (x *UserInfoResponse) ProtoReflect() protoreflect.Message {
	return nil
}

type UpdateProfileRequest struct {
	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *UpdateProfileRequest) Reset() {
	*x = UpdateProfileRequest{}
}

func (x *UpdateProfileRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type UpdateResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	User    *User  `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *UpdateResponse) Reset() {
	*x = UpdateResponse{}
}

func (x *UpdateResponse) ProtoReflect() protoreflect.Message {
	return nil
}

type ChangePasswordRequest struct {
	UserId          string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	CurrentPassword string `protobuf:"bytes,2,opt,name=current_password,json=currentPassword,proto3" json:"current_password,omitempty"`
	NewPassword     string `protobuf:"bytes,3,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
}

func (x *ChangePasswordRequest) Reset() {
	*x = ChangePasswordRequest{}
}

func (x *ChangePasswordRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type ChangePasswordResponse struct {
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ChangePasswordResponse) Reset() {
	*x = ChangePasswordResponse{}
}

func (x *ChangePasswordResponse) ProtoReflect() protoreflect.Message {
	return nil
}
