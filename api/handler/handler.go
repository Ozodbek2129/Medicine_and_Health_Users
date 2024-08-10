package handler

import (
	"log/slog"
	pb "user_medic/genproto/user"
)

type Handler struct{
	AuthUser pb.UserServiceClient
	Log *slog.Logger
}