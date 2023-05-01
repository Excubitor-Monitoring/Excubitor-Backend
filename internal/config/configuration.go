package config

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flags "github.com/spf13/pflag"
	"os"
	"strings"
)

// k is the global configuration
var k = koanf.New(".")

var ErrInvalidConfigParameter = errors.New("invalid config parameter")

// InitConfig initializes the configuration.
func InitConfig() error {
	// Load default values
	err := k.Load(confmap.Provider(map[string]interface{}{
		"logging.log_level":                  "INFO",
		"logging.method":                     "CONSOLE",
		"http.port":                          "0.0.0.0",
		"http.cors.allowed_origins":          []string{"*"},
		"http.cors.allowed_methods":          []string{"GET", "POST"},
		"http.cors.allowed_headers":          []string{"Origin", "Content-Type", "Authorization"},
		"http.auth.jwt.access_token_secret":  "",
		"http.auth.jwt.refresh_token_secret": "",
	}, "."), nil)
	if err != nil {
		return err
	}

	// Configure and parse flagset
	f := flags.NewFlagSet("config", flags.ContinueOnError)
	f.Usage = func() {
		fmt.Println("Could not parse flags! For more information see 'excubitor --help'")
		os.Exit(1)
	}

	f.String("config", "config.yml", "Path to the config file.")
	f.String("host", "0.0.0.0", "Host the HTTP Server shall run on.")
	f.Int("port", 8080, "Port the HTTP Server shall run on.")
	if err := f.Parse(os.Args[1:]); err != nil {
		return err
	}

	configFile, err := f.GetString("config")
	if err != nil {
		return fmt.Errorf("could not init config file: %w", err)
	}

	// Load YAML Config file
	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		return fmt.Errorf("could not init config file: %w", err)
	}

	// Load environment variables
	err = k.Load(env.Provider("EXCUBITOR_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "EXCUBITOR_")), "_", ".", -1)
	}), nil)
	if err != nil {
		return fmt.Errorf("could not init environment variable configuration: %w", err)
	}

	// Load flagset into configuration
	err = k.Load(posflag.ProviderWithFlag(f, ".", nil, func(flag *flags.Flag) (string, interface{}) {
		switch flag.Name {
		case "host":
			return "http.host", posflag.FlagVal(f, flag)
		case "port":
			return "http.port", posflag.FlagVal(f, flag)
		default:
			return "", ""
		}
	}), nil)
	if err != nil {
		return fmt.Errorf("could not init configuration through flags: %w", err)
	}

	// Check config for errors
	err = checkConfig()
	if err != nil {
		return err
	}

	return nil
}

// GetConfig returns a reference to the global config object.
func GetConfig() *koanf.Koanf {
	return k
}

// checkConfig returns an error if certain configuration parameters are out of spec.
func checkConfig() error {
	if k.String("http.auth.jwt.access_token_secret") == "" {
		return fmt.Errorf("%w: %s", ErrInvalidConfigParameter, "access token secret is not set")
	}

	if k.String("http.auth.jwt.refresh_token_secret") == "" {
		return fmt.Errorf("%w: %s", ErrInvalidConfigParameter, "refresh token secret is not set")
	}

	if k.Int("http.port") < 1 || k.Int("http.port") > 65535 {
		return fmt.Errorf("%w: %s %d", ErrInvalidConfigParameter, "port needs to be at least 1 and lower than 65536. Is:", k.Int("http.port"))
	}

	return nil
}
