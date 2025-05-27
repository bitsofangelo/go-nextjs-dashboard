package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"go-nextjs-dashboard/internal/auth"
	"go-nextjs-dashboard/internal/config"
	"go-nextjs-dashboard/internal/db"
	sloglogger "go-nextjs-dashboard/internal/logger/slog"
)

type T int

func main() {
	// var errs []error
	// errs = append(errs, errors.New("test"))
	// errs = append(errs, errors.New("test2"))
	// errs = append(errs, errors.New("test3"))
	//
	// fmt.Println(errors.Join(errs...))
	//

	cfg, _ := config.Load()
	logger, _ := sloglogger.New(cfg)
	database, _ := db.Open(cfg, logger)
	hasher := auth.NewArgonHasher()
	gojwt := auth.NewGOJWT()
	refreshStore := auth.NewGormRefreshStore(database, logger)
	a := auth.New(hasher, gojwt, refreshStore, logger)

	uid := uuid.New()
	token, _ := a.CreateRefreshToken(context.Background(), uid)

	fmt.Println(uid, token)

	// uid := uuid.New()
	// s, _ := a.NewJWT(uid)
	// c, _ := a.ParseJWT(s)
	// fmt.Println(c, "uid", uid)

	p, _ := a.HashPassword("123456")
	ok, _ := a.CheckPasswordHash("123456", p)
	fmt.Println(p, ok)
}
