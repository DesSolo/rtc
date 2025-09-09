package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	clientv3 "go.etcd.io/etcd/client/v3"

	"rtc/internal/auth"
	"rtc/internal/config"
	"rtc/internal/provider"
	"rtc/internal/server"
	"rtc/internal/storage"
	"rtc/internal/storage/etcd"
	"rtc/internal/storage/postgres"
	"rtc/pkg/closer"
)

const (
	healthCheckTimeout = time.Second * 5
)

type container struct {
	config        *config.Config
	storage       storage.Storage
	valuesStorage storage.ValuesStorage
	provider      *provider.Provider
	jwtAuth       *auth.JWT
	server        *server.Server
}

func newContainer() *container {
	return &container{}
}

func (c *container) Config() *config.Config {
	if c.config == nil {
		configFilePath := os.Getenv("CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "config.yaml"
		}

		conf, err := config.NewConfigFromFile(configFilePath)
		if err != nil {
			fatal("failed to load config from file", err)
		}

		c.config = conf
	}

	return c.config
}

func (c *container) Storage() storage.Storage {
	if c.storage == nil {
		options := c.Config().Storage

		poolConfig, err := pgxpool.ParseConfig(options.DSN)
		if err != nil {
			fatal("failed to parse postgres config", err)
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			fatal("failed to connect to postgres", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			fatal("failed to ping postgres", err)
		}

		pgStorage := postgres.NewStorage(pool)

		closer.Add(pgStorage.Close)

		c.storage = pgStorage
	}

	return c.storage
}

func (c *container) ValuesStorage() storage.ValuesStorage {
	if c.valuesStorage == nil {
		options := c.Config().ValuesStorage

		client, err := clientv3.New(clientv3.Config{
			Endpoints:   options.Endpoints,
			DialTimeout: options.DialTimeout,
		})
		if err != nil {
			fatal("failed to connect to etcd", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
		defer cancel()

		if _, err := client.Status(ctx, options.Endpoints[0]); err != nil {
			fatal("failed to connect to etcd", err)
		}

		closer.Add(client.Close)

		c.valuesStorage = etcd.NewValuesStorage(client, etcd.WithPath(options.Path))
	}

	return c.valuesStorage
}

func (c *container) Provider() *provider.Provider {
	if c.provider == nil {
		c.provider = provider.NewProvider(c.Storage(), c.ValuesStorage())
	}

	return c.provider
}

func (c *container) JWTAuth() *auth.JWT {
	if c.jwtAuth == nil {
		options := c.Config().Server.Auth.JWT

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(options.PrivateKey))
		if err != nil {
			fatal("failed to parse private key", err)
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(options.PublicKey))
		if err != nil {
			fatal("failed to parse public key", err)
		}

		c.jwtAuth = auth.NewJWT(privateKey, publicKey, options.TTL)
	}

	return c.jwtAuth
}

func (c *container) Server() *server.Server {
	if c.server == nil {
		options := c.Config().Server

		c.server = server.NewServer(c.Provider(), c.JWTAuth(),
			server.WithAddress(options.Address),
			server.WithReadHeaderTimeout(options.ReadHeaderTimeout),
			loadServerAuth(c),
		)
	}

	return c.server
}

func fatal(message string, err error) {
	slog.Error(message, "err", err)
	os.Exit(1)
}
