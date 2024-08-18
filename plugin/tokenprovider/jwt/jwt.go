package jwt

import (
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"todolist/plugin/tokenprovider"
)

type jwtProvider struct {
	name   string
	secret string
}

func NewJWTProvider(name string) *jwtProvider {
	return &jwtProvider{name: name}
}

func (p *jwtProvider) GetPrefix() string {
	//TODO implement me
	return p.name
}

func (p *jwtProvider) Get() interface{} {
	//TODO implement me
	return p
}

func (p *jwtProvider) Name() string {
	//TODO implement me
	return p.name
}

func (p *jwtProvider) InitFlags() {
	//TODO implement me
	flag.StringVar(&p.secret, "jwt-secret", "200lab.io", "Secret key for generating JWT")
}

func (p *jwtProvider) Configure() error {
	//TODO implement me
	return nil
}

func (p *jwtProvider) Run() error {
	//TODO implement me
	return nil
}

func (p *jwtProvider) Stop() <-chan bool {
	//TODO implement me
	c := make(chan bool)
	go func() {
		c <- true
	}()
	return c
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (*tokenprovider.Token, error) {
	// generate JWT
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		data,
		jwt.StandardClaims{
			IssuedAt:  now.Local().Unix(),
			ExpiresAt: now.Local().Add(time.Second * time.Duration(expiry)).Unix(),
			Id:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})
	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &tokenprovider.Token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(myToken string) (*tokenprovider.TokenPayload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, tokenprovider.ErrInvalidToken
	}

	if !res.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	return &claims.Payload, nil
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

type myClaims struct {
	Payload tokenprovider.TokenPayload `json:"payload"`
	jwt.StandardClaims
}
