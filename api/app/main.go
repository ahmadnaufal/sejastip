package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sejastip.id/api/storage"

	"sejastip.id/api/delivery"
	"sejastip.id/api/handler"
	"sejastip.id/api/repository"
	"sejastip.id/api/usecase"

	"github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/joeshaw/envdecode"
	"github.com/subosito/gotenv"
)

type Config struct {
	Database struct {
		User                   string `env:"DATABASE_USER,required"`
		Password               string `env:"DATABASE_PASSWORD,required"`
		Host                   string `env:"DATABASE_HOST,default=127.0.0.1"`
		Port                   string `env:"DATABASE_PORT,default=3306"`
		Name                   string `env:"DATABASE_NAME,default=sejastip"`
		Pool                   int    `env:"DATABASE_POOL,default=5"`
		Charset                string `env:"DATABASE_CHARSET,default=utf8"`
		CloudSQLEnabled        bool   `env:"DATABASE_CLOUD_SQL_ENABLED,default=false"`
		CloudSQLConnectionName string `env:"DATABASE_CLOUD_SQL_CONNECTION_NAME"`
	}

	GCS struct {
		Enabled  bool   `env:"GCS_ENABLED,default=false"`
		BucketID string `env:"GCS_BUCKET_ID,required"`
	}

	Port string `env:"PORT,required"`

	JWTPrivateKey string `env:"JWT_PRIVATE_KEY,required"`
}

var config Config

func init() {
	gotenv.Load()
	err := envdecode.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func NewMysqlConnection(c *Config) (*sql.DB, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.Charset,
	)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(c.Database.Pool)
	err = db.Ping()
	return db, err
}

func NewCloudSqlConnection(c *Config) (*sql.DB, error) {
	cfg := mysql.Cfg(c.Database.CloudSQLConnectionName, c.Database.User, c.Database.Password)
	cfg.DBName = c.Database.Name
	cfg.ParseTime = true
	db, err := mysql.DialCfg(cfg)
	if err != nil {
		log.Fatal("error connecting to cloudsql: ", err)
	}

	db.SetMaxOpenConns(c.Database.Pool)
	err = db.Ping()
	return db, err
}

func main() {
	var (
		db  *sql.DB
		err error
	)
	if config.Database.CloudSQLEnabled {
		db, err = NewCloudSqlConnection(&config)
	} else {
		db, err = NewMysqlConnection(&config)
	}

	if err != nil {
		log.Fatal("error connecting to mysql server: ", err)
	}

	userRepo := repository.NewMysqlUser(db)
	bankRepo := repository.NewMysqlBank(db)
	countryRepo := repository.NewMysqlCountry(db)
	productRepo := repository.NewMysqlProduct(db)
	addressRepo := repository.NewMysqlUserAddress(db)
	transactionRepo := repository.NewMysqlTransaction(db)
	deviceRepo := repository.NewMysqlDevice(db)

	appStorage := storage.NewLocalStorage()
	if config.GCS.Enabled {
		appStorage = storage.NewGCS(config.GCS.BucketID)
	}

	uuc := usecase.NewUserUsecase(&usecase.UserProvider{UserRepository: userRepo})
	uh := delivery.NewUserHandler(uuc)

	auc := usecase.NewAuthUsecase(&usecase.AuthProvider{
		UserRepository: userRepo,
		JWTPrivateKey:  config.JWTPrivateKey,
	})
	ah := delivery.NewAuthHandler(auc)

	buc := usecase.NewBankUsecase(&usecase.BankProvider{
		BankRepo: bankRepo,
		Storage:  appStorage,
	})
	bh := delivery.NewBankHandler(buc)

	cuc := usecase.NewCountryUsecase(&usecase.CountryProvider{
		CountryRepo: countryRepo,
		Storage:     appStorage,
	})
	ch := delivery.NewCountryHandler(cuc)

	puc := usecase.NewProductUsecase(&usecase.ProductProvider{
		ProductRepo: productRepo,
		UserRepo:    userRepo,
		CountryRepo: countryRepo,
		Storage:     appStorage,
	})
	ph := delivery.NewProductHandler(puc)

	uauc := usecase.NewUserAddressUsecase(&usecase.UserAddressProvider{
		UserAddressRepo: addressRepo,
	})
	uah := delivery.NewUserAddressHandler(uauc)

	tc := usecase.NewTransactionUsecase(&usecase.TransactionProvider{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
		ProductRepo:     productRepo,
		AddressRepo:     addressRepo,
		CountryRepo:     countryRepo,
	})
	th := delivery.NewTransactionHandler(tc)

	dc := usecase.NewDeviceUsecase(&usecase.DeviceProvider{
		DeviceRepo: deviceRepo,
	})
	dh := delivery.NewDeviceHandler(dc)

	h := handler.NewHandler(config.JWTPrivateKey, &uh, &ah, &bh, &ch, &ph, &uah, &th, &dh)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Port),
		Handler:      h,
		ReadTimeout:  300 * time.Second,
		WriteTimeout: 300 * time.Second,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		log.Printf("Sejastip API is available at %s\n", s.Addr)
		if serr := s.ListenAndServe(); serr != http.ErrServerClosed {
			log.Fatal(serr)
		}
	}(s)

	<-sigChan

	log.Println("Sejastip API stopped.")
}
