package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"trelloBE/common"
	"trelloBE/controller"
	"trelloBE/db"
	"trelloBE/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	util.ReadEnv()
	common.Init()

	var privKeyBytes []byte
	var pubKeyBytes []byte
	privKeyBytes, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		if !os.IsNotExist(err) {
			util.Log.Fatalln("Failed to read private.pem file, err:", err)
		}

		util.Log.Println("Creating private & public keys...")
		var privateKey *rsa.PrivateKey
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			util.Log.Fatalln("Failed to generate private key, err:", err)
		}

		var publicKey rsa.PublicKey = privateKey.PublicKey
		pubKeyBytes = x509.MarshalPKCS1PublicKey(&publicKey)

		privKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)

		var privateKeyBlock *pem.Block = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privKeyBytes,
		}
		privKeyFile, err := os.Create("private.pem")
		if err != nil {
			util.Log.Fatalln("Failed to create private pem file, err:", err)
		}

		if err = pem.Encode(privKeyFile, privateKeyBlock); err != nil {
			util.Log.Fatalln("Failed to encode private key block to pem file, err:", err)
		}
		util.Log.Println("Private key written to file...")

		var publicKeyBlock *pem.Block = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubKeyBytes,
		}
		pubKeyFile, err := os.Create("public.pem")
		if err != nil {
			util.Log.Fatalln("Failed to create public pem file, err:", err)
		}

		if err = pem.Encode(pubKeyFile, publicKeyBlock); err != nil {
			util.Log.Fatalln("Failed to encode public key block to pem file, err:", err)
		}
		util.Log.Println("Public key written to file...")
	}

	if privKeyPem, _ := pem.Decode(privKeyBytes); privKeyPem != nil {
		common.PrivateKey, err = x509.ParsePKCS1PrivateKey(privKeyPem.Bytes)
		if err != nil {
			util.Log.Fatalln("Failed to parse private key, err:", err)
		}
	} else {
		util.Log.Fatalln("Failed to decode private key pem")
	}

	common.PublicKey = &common.PrivateKey.PublicKey
	util.Log.Println("Private & public keys loaded...")

	if err := db.InitDb(); err != nil {
		util.Log.Fatalln("Failed to connect to db, err: ", err)
	}
	if errString, err := db.RunConfig(); err != nil {
		util.Log.Fatalln("Failed to run config on db, errString: ", errString, "err: ", err)
	}

	var r *chi.Mux = chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// TODO: Change when going to prod
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Use(middleware.Logger)
	r.Route("/trelloapi", func(r chi.Router) {
		r.Get("/ping", controller.PongCheck)

		r.Route("/"+os.Getenv("API_VERSION"), func(r chi.Router) {

			r.Route("/auth", func(r chi.Router) {
				r.Use(util.LimitBody1Mb)

				r.Group(func(r chi.Router) {
					r.Use(util.JsonRoute)

					r.Post("/register", controller.Register)
					r.Post("/login", controller.Login)
				})

				r.Get("/refresh", controller.RefreshAccessToken)
			})

			r.Group(func(r chi.Router) {
				// r.Use(util.JsonRoute)
				r.Use(util.AuthMiddleware)

				r.Post("/board/create", controller.CreateBoard)
			})
		})
	})

	util.Log.Println("Starting server on port:", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")), r); err != nil {
		util.Log.Fatalln("Failed to start server...", err)
	}
}
