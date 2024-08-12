package service

import (
	"context"
	"fmt"
	"log/slog"
	pb "user_medic/genproto/user"
	"user_medic/pkg/logger"
	"user_medic/storage/postgres"
)

type MedicineService struct {
	pb.UnimplementedUserServiceServer
	user postgres.MedicineUser
	log  *slog.Logger
}

func NewMedicineService(user postgres.MedicineUser) *MedicineService {
	return &MedicineService{
		user: user,
		log:  logger.NewLogger(),
	}
}

func (ms *MedicineService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error){
	resp,err:=ms.user.RegisterUser(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Register serviceda xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse,error) {
	err:=ms.user.RefReshToken(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Refresh token service da xatolik: %v",err))
		return nil,err
	}
	return &pb.RefreshTokenResponse{
		Message: "SIGNING_KEY yangilandi endi login qiling.",
	},nil
}

func (ms *MedicineService) GetUserProfile(ctx context.Context,req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse,error){
	resp,err:=ms.user.GetUserProfile(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("GetUserProfile service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) UpdateUserProfile(ctx context.Context,req *pb.UpdateUserProfileRequest)(*pb.UpdateUserProfileResponse,error){
	resp,err:=ms.user.UpdateUserProfile(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Updat user profile servicwe da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) LogoutUser(ctx context.Context, token *pb.LogoutUserRequest) (*pb.LogoutUserResponse,error){
	resp,err:=ms.user.LogoutUser(ctx,token)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Logout user service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) GetByUserEmail(ctx context.Context,req *pb.LoginUserRequest) (*pb.RegisterUserResponse, error){
	resp,err:=ms.user.GetByUserEmail(req.Email)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("GetBy User Email service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) (*pb.StoreRefreshTokenRes,error){
	err:=ms.user.StoreRefreshToken(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Store Refresh Token service da xatolik: %v",err))
		return nil,err
	}
	return nil,nil
}

func (ms *MedicineService) GetByUserId(ctx context.Context,req *pb.UserId)(*pb.FLResponse,error){
	resp,err:=ms.user.GetByUserId(ctx,req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("Store Refresh Token service da xatolik: %v",err))
		return nil,err
	}
	return resp,nil
}

func (ms *MedicineService) IdCheck(ctx context.Context,req *pb.UserId)(*pb.Response,error){
	resp,err:=ms.user.IdCheck(req)
	if err!=nil{
		ms.log.Error(fmt.Sprintf("yuborilgan id bazada yuq: %v",err))
		return nil,err
	}
	return resp,nil
}