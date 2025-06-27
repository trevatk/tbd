package setup

const (
	defaultPort = "8080"
	defaultHost = "127.0.0.1"

	defaultKeyValueDir = "data"

	defaultLogLevel = "DEBUG"

	defaultSigningKey = "supersecret"

	defaultNameserver1 = "ns1.structx.io"
	defaultNameserver2 = "ns2.structx.io"
)

// Config service configuration
type Config struct {
	Auth       Auth
	Gateway    Gateway
	KeyValue   KeyValue
	Logger     Logger
	Nameserver Nameserver
}

// UnmarshalConfig read service config from env variables
func UnmarshalConfig() *Config {
	return &Config{
		Auth: Auth{
			SigningKey: envLookup("AUTH_SIGNING_KEY", defaultSigningKey),
		},
		Gateway: Gateway{
			Host: envLookup("GW_HOST", defaultHost),
			Port: envLookup("GW_PORT", defaultPort),
		},
		KeyValue: KeyValue{
			Dir: envLookup("KV_DIR", defaultKeyValueDir),
		},
		Logger: Logger{
			Level: envLookup("LOG_LEVEL", defaultLogLevel),
		},
		Nameserver: Nameserver{
			NS1: envLookup("NS_SERVER_1", defaultNameserver1),
			NS2: envLookup("NS_SERVER_2", defaultNameserver2),
		},
	}
}
