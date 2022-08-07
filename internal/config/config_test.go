package config_test

import (
	"log"
	"os"
	"testing"

	"github.com/ernur-eskermes/product-store/internal/config"
	"github.com/stretchr/testify/require"
)

type env struct {
	mongoURI      string
	mongoUser     string
	mongoPassword string
	mongoDatabase string
	grpcPort      string
}

func setEnv(env env) {
	var err error
	if env.mongoURI != "" {
		if err = os.Setenv("MONGO_URI", env.mongoURI); err != nil {
			log.Fatal(err)
		}
	}
	if env.mongoUser != "" {
		if err = os.Setenv("MONGO_USER", env.mongoUser); err != nil {
			log.Fatal(err)
		}
	}
	if env.mongoPassword != "" {
		if err = os.Setenv("MONGO_PASSWORD", env.mongoPassword); err != nil {
			log.Fatal(err)
		}
	}
	if env.mongoDatabase != "" {
		if err = os.Setenv("MONGO_DATABASE", env.mongoDatabase); err != nil {
			log.Fatal(err)
		}
	}
	if env.grpcPort != "" {
		if err = os.Setenv("GRPC_PORT", env.grpcPort); err != nil {
			log.Fatal(err)
		}
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		name     string
		env      env
		want     *config.Config
		expErr   string
		unsetEnv func()
	}{
		{
			name: "test_config",
			env: env{
				grpcPort:      "9000",
				mongoDatabase: "test_database",
				mongoPassword: "test_password",
				mongoUser:     "test_user",
				mongoURI:      "test_uri",
			},
			want: &config.Config{
				GRPC: config.GRPCConfig{
					Port: 9000,
				},
				Mongo: config.MongoConfig{
					URI:      "test_uri",
					Database: "test_database",
					User:     "test_user",
					Password: "test_password",
				},
			},
		},
		{
			name: "invalid_grpc_port",
			env: env{
				grpcPort:      "asd",
				mongoDatabase: "test_database",
				mongoPassword: "test_password",
				mongoUser:     "test_user",
				mongoURI:      "test_uri",
			},
			expErr: "envconfig.Process: assigning GRPC_PORT to Port: converting 'asd' to type int. details: strconv.ParseInt: parsing \"asd\": invalid syntax",
		},
		{
			name: "empty_mongo_uri",
			env: env{
				grpcPort:      "9000",
				mongoDatabase: "test_database",
				mongoPassword: "test_password",
				mongoUser:     "test_user",
			},
			expErr: "required key MONGO_URI missing value",
		},
	}

	for _, s := range cases {
		t.Run(s.name, func(t *testing.T) {
			setEnv(s.env)

			cfg, err := config.New()
			if err != nil {
				require.EqualError(t, err, s.expErr)
			} else {
				require.Equal(t, s.want, cfg)
			}
		})
	}
}
