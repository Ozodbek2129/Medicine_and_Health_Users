package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"
	pb "user_medic/genproto/user"
	"user_medic/pkg/logger"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

func (s *MedicineUser) RefReshToken(ctx context.Context, req *pb.RefreshTokenRequest) error {
	err := godotenv.Load(".env")
	if err != nil {
		s.log.Error(fmt.Sprintf(".env faylini yuklashda xatolik yuz berdi: %v", err))
		return fmt.Errorf(".env faylini yuklashda xatolik yuz berdi")
	}

	const charset = "abcdefghijkQWERTYU*+_)(*&^%$#@lmnopqrstuvwxyz!@#$%^&" + "ABCD!@#$%^&*()_+EFGHIJKLMNOP(*&^%$QRSTUVW@#$%D@#$%^&*()_YZ0123456789"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	newSigningKey := stringWithCharset(10, charset, seededRand)

	err = os.Setenv("SIGNING_KEY", newSigningKey)
	if err != nil {
		s.log.Error(fmt.Sprintf("Yangi SIGNING_KEY sozlamasida xatolik yuz berdi: %v", err))
		return fmt.Errorf("yangi SIGNING_KEY sozlamasida xatolik yuz berdi")
	}

	err = godotenv.Write(map[string]string{
		"USER_SERVICE": os.Getenv("USER_SERVICE"),
		"USER_ROUTER":  os.Getenv("USER_ROUTER"),
		"DB_USER":      os.Getenv("DB_USER"),
		"DB_HOST":      os.Getenv("DB_HOST"),
		"DB_NAME":      os.Getenv("DB_NAME"),
		"DB_PASSWORD":  os.Getenv("DB_PASSWORD"),
		"DB_PORT":      os.Getenv("DB_PORT"),
		"SIGNING_KEY":  newSigningKey,
	}, ".env")

	if err != nil {
		s.log.Error(fmt.Sprintf(".env fayliga yozishda xatolik yuz berdi: %v", err))
		return fmt.Errorf(".env fayliga yozishda xatolik yuz berdi")
	}

	s.log.Info("SIGNING_KEY muvaffaqiyatli yangilandi")

	return nil
}

func stringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}


func (m *MedicineUser) GetUserProfile(ctx context.Context,req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse,error){
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

func (m *MedicineUser) UpdateUserProfile(ctx context.Context,req *pb.UpdateUserProfileRequest)(*pb.UpdateUserProfileResponse,error){
	query:=`UPDATE 
				users
			SET 
				email = $1, password_hash = $2, first_name = $3, 
				last_name = $4, date_of_birth = $5, gender = $6, role = $7, updated_at = $8 
			WHERE 
				id = $9`

	_,err:=m.db.ExecContext(ctx,query,req.Email,req.Password,req.FirstName,req.LastName,
							req.DateOfBirthday,req.Gender,req.Role,time.Now(),req.Id)

	if err!=nil{
		m.log.Error(fmt.Sprintf("Profile ni yangilashda xatolik: %v",err))
		return nil,err
	}

	return &pb.UpdateUserProfileResponse{
		Message: "Profile muvafiqiyatli yangilandi.",
	},nil
}

func (m *MedicineUser) LogoutUser(ctx context.Context, token *pb.LogoutUserRequest) (*pb.LogoutUserResponse,error) {
	_, err := m.db.ExecContext(ctx,`
	update 
		refresh_token 
	set 
		deleted_at=$1 
	where 
		refresh=$2`, time.Now(), token.RefreshToken)
	if err != nil {
		return nil,err
	}
	return &pb.LogoutUserResponse{
		Message: "Logout succsessfully.",
	},nil
}

