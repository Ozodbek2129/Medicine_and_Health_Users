package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	pb "user_medic/genproto/user"
	"user_medic/pkg/logger"

	"github.com/google/uuid"
)

type MedicineUser struct {
	db  *sql.DB
	log *slog.Logger
}

func NewMedicineUser(db *sql.DB) *MedicineUser {
	return &MedicineUser{
		db:  db,
		log: logger.NewLogger(),
	}
}

func (m *MedicineUser) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	query := `insert into users(
				id, email, password_hash, first_name,
				last_name, date_of_birth, gender, role, created_at, updated_at
			)values(
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	id := uuid.NewString()
	newtime := time.Now()

	_, err := m.db.ExecContext(ctx, query, id, req.Email, req.Password, req.FirstName,
		req.LastName, req.DateOfBirthday, req.Gender, req.Role, newtime, newtime)

	if err != nil {
		m.log.Error(fmt.Sprintf("Registerda bazdaga qushishda xatolik: %v", err))
		return nil, err
	}

	return &pb.RegisterUserResponse{
		Id:             id,
		Email:          req.Email,
		Password:       req.Password,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		DateOfBirthday: req.DateOfBirthday,
		Gender:         req.Gender,
		Role:           req.Role,
		CreatedAt:      newtime.String(),
		UpdatedAt:      newtime.String(),
	}, nil
}

func (m MedicineUser) GetByUserEmail(email string) (*pb.RegisterUserResponse, error) {
	resp := pb.RegisterUserResponse{}
	query := `select 
				id, email, password_hash, first_name, 
				last_name, date_of_birth, gender, role, created_at, updated_at
			from 
				users
			where
				email=$1 and deleted_at is null`

	err := m.db.QueryRow(query, email).Scan(&resp.Id, &resp.Email, &resp.Password, &resp.FirstName, &resp.LastName, &resp.DateOfBirthday,
		&resp.Gender, &resp.Role, &resp.CreatedAt, &resp.UpdatedAt)

	if err != nil {
		m.log.Error(fmt.Sprintf("Email buyicha xatolik: %v", err))
		return nil, err
	}
	return &resp, nil
}

func (m *MedicineUser) StoreRefreshToken(ctx context.Context, req *pb.StoreRefreshTokenReq) error {
	query := `insert into refresh_token(
				id, user_id, refresh, created_at, updated_at
			)values(
				$1.$2,$3,$4,$5)`

	id := uuid.NewString()
	newtime := time.Now()
	_, err := m.db.ExecContext(ctx, query, id, req.UserId, req.RefreshToken, newtime, newtime)
	if err != nil {
		m.log.Error(fmt.Sprintf("Refresh token ni bazaga qushishda xatolik: %v", err))
		return nil
	}
	return nil
}

func (m *MedicineUser) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	resp := pb.GetUserProfileResponse{}
	query := `select 
				id, email, password_hash, first_name, 
				last_name, date_of_birth, gender, role, created_at, updated_at
			from 
				users
			where
				email=$1 and deleted_at is null`

	err := m.db.QueryRow(query, req.Email).Scan(&resp.Id, &resp.Email, &resp.Password, &resp.FirstName, &resp.LastName, &resp.DateOfBirthday,
		&resp.Gender, &resp.Role, &resp.CreatedAt, &resp.UpdatedAt)

	if err != nil {
		m.log.Error(fmt.Sprintf("Email buyicha xatolik: %v", err))
		return nil, err
	}
	return &resp, nil
}

func (m *MedicineUser) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserProfileResponse, error) {
	query := `UPDATE 
				users
			SET 
				email = $1, password_hash = $2, first_name = $3, 
				last_name = $4, date_of_birth = $5, gender = $6, role = $7, updated_at = $8 
			WHERE 
				id = $9`

	_, err := m.db.ExecContext(ctx, query, req.Email, req.Password, req.FirstName, req.LastName,
		req.DateOfBirthday, req.Gender, req.Role, time.Now(), req.Id)

	if err != nil {
		m.log.Error(fmt.Sprintf("Profile ni yangilashda xatolik: %v", err))
		return nil, err
	}

	return &pb.UpdateUserProfileResponse{
		Message: "Profile muvafiqiyatli yangilandi.",
	}, nil
}

func (m *MedicineUser) LogoutUser(ctx context.Context, token *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	_, err := m.db.ExecContext(ctx, `
	update 
		refresh_token 
	set 
		deleted_at=$1 
	where 
		refresh=$2`, time.Now(), token.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.LogoutUserResponse{
		Message: "Logout succsessfully.",
	}, nil
}

func (m *MedicineUser) GetByUserId(ctx context.Context, req *pb.UserId) (*pb.FLResponse, error) {
	query := `select 
				first_name,last_name
			from
				users
			where
				id=$1`

	var res pb.FLResponse
	err := m.db.QueryRowContext(ctx, query, req.Userid).Scan(&res.FirstName, &res.LastName)
	if err != nil {
		m.log.Error(fmt.Sprintf("first name va last name ni olishda xatolik: %v", err))
		return nil, err
	}
	return &pb.FLResponse{
		FirstName: res.FirstName,
		LastName:  res.LastName,
	}, nil
}

func (u *MedicineUser) IdCheck(req *pb.UserId) (*pb.Response, error) {
	query := `select id from users`

	rows, err := u.db.Query(query)
	if err != nil {
		return &pb.Response{B: false}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return &pb.Response{B: false}, err
		}
		if id == req.Userid {
			return &pb.Response{B: true}, nil
		}
	}

	if err = rows.Err(); err != nil {
		return &pb.Response{B: false}, err
	}

	return &pb.Response{B: false}, nil
}

func (u *MedicineUser) NotificationsAdd(ctx context.Context,req *pb.NotificationsAddRequest)(*pb.NotificationsAddResponse,error){
	query:=`insert into notification(
					id, user_id, message, created_at, updated_at
				)values(
					$1,$2,$3,$4,$5)`

	id:=uuid.NewString()
	vaqt:=time.Now().Format("2006/01/02")
	_,err:=u.db.ExecContext(ctx,query,id,req.UserId,req.Message,vaqt,vaqt)
	if err!=nil{
		return nil,err
	}

	return &pb.NotificationsAddResponse{
		Message: "Notification yuborildi.",
	},nil
}

func (u *MedicineUser) NotificationsGet(ctx context.Context,req *pb.NotificationsGetRequest)(*pb.NotificationsGetResponse,error){
	query:=`select 
				message,created_at
			from
				notification
			where
				user_id=$1`

	rows,err:=u.db.QueryContext(ctx,query,req.UserId)
	if err!=nil{
		return nil,err
	}
	defer rows.Close()

	var resp pb.NotificationsGetResponse

	for rows.Next() {
		var notification pb.Notification

		err := rows.Scan(&notification.Message, &notification.CreatedAt)
		if err != nil {
			return nil, err
		}

		resp.Notifications = append(resp.Notifications, &notification)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &resp,nil
}

func (u *MedicineUser) NotificationsPut(ctx context.Context,req *pb.NotificationsPutRequest)(*pb.NotificationsPutResponse,error){
	query:=`update 
				notification
			set
				message=$1
			where
				user_id=$2 and created_at=$3`

	_,err:=u.db.ExecContext(ctx,query,req.Message,req.UserId,req.CreatedAt)
	if err!=nil{
		return nil,err
	}
	return &pb.NotificationsPutResponse{
		Message: "Siz yuborgan notification uzgartirildi >>> "+req.Message,
	},nil
}