package main

import (
	"fmt"

	"go-nextjs-dashboard/internal/hashing"
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

	// cfg, _ := config.Load()
	// logger, _ := sloglogger.New(cfg)
	// database, _ := db.Open(cfg, logger)
	hasher := hashing.NewArgon2IDHasher()
	hash := hashing.New(hasher)
	// gojwt := auth.NewGOJWT()
	// refreshStore := auth.NewGormRefreshStore(database, logger)
	// a := auth.New(hash, gojwt, refreshStore, logger)

	// uid := uuid.New()
	// token, _ := a.CreateRefreshToken(context.Background(), uid)

	// fmt.Println(uid, token)

	// uid := uuid.New()
	// s, _ := a.NewJWT(uid)
	// c, _ := a.ParseJWT(s)
	// fmt.Println(c, "uid", uid)

	p, _ := hash.Make("123456")
	ok, _ := hash.Check("123456", p)
	fmt.Println(p, ok)
}
