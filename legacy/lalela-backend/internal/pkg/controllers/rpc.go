package controllers

import (
	"lalela-backend/internal/pkg/middleware"
	"github.com/gorilla/rpc"
	"log"
)

func InitRPC(s *rpc.Server) {
	err := s.RegisterService(new(AuthCon), "")
	if err != nil {
		log.Print(middleware.NewError(err))
	}
	//s.RegisterService(new(DashsCon), "")
	s.RegisterService(new(DashCon), "")
	s.RegisterService(new(UserCon), "")
	s.RegisterService(new(PermissionCon), "")
	s.RegisterService(new(GroupsCon), "")
	s.RegisterService(new(KanbanCon), "")
}
