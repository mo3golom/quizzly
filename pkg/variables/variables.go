package variables

import "fmt"

const (
	VariableTypeEnvironment VariableType = "environment"
)

var registry []Variable

var (
	// AppEnvironmentVariable окружение, в котором запущено приложение
	AppEnvironmentVariable = Environment[string]("ENV", "prod")

	AuthSecretKey       = Environment[string]("AUTH_SECRET_KEY", "")
	AuthCookieBlockKey  = Environment[string]("AUTH_COOKIE_BLOCK_KEY", "")
	AuthSenderFromEmail = Environment[string]("AUTH_SENDER_FROM_EMAIL", "")
	AuthSenderHost      = Environment[string]("AUTH_SENDER_HOST", "")
	AuthSenderPort      = Environment[string]("AUTH_SENDER_PORT", "")
	AuthSenderUser      = Environment[string]("AUTH_SENDER_USER", "")
	AuthSenderPassword  = Environment[string]("AUTH_SENDER_PASSWORD", "")

	S3Endpoint  = Environment[string]("S3_ENDPOINT", "")
	S3AccessKey = Environment[string]("S3_ACCESS_KEY", "")
	S3SecretKey = Environment[string]("S3_SECRET_KEY", "")
	S3Bucket    = Environment[string]("S3_BUCKET", "")
	S3UseSSL    = Environment[bool]("S3_USE_SSL", false)

	MetricsUser     = Environment[string]("METRICS_USER", "")
	MetricsPassword = Environment[string]("METRICS_PASSWORD", "")

	SupabaseReference = Environment[string]("SUPABASE_REFERENCE", "")
	SupabaseAnonKey   = Environment[string]("SUPABASE_ANON_KEY", "")
)

type (
	Variable interface {
		Name() string
		Type() VariableType
	}

	VariableType string

	DefaultVariable[T any] struct {
		name         string
		defaultValue T
		t            VariableType
	}

	StringVariable = DefaultVariable[string]
	BoolVariable   = DefaultVariable[bool]
)

func (v DefaultVariable[T]) Name() string {
	return v.name
}

func (v DefaultVariable[T]) Type() VariableType {
	return v.t
}

func (v DefaultVariable[T]) String() string {
	return fmt.Sprintf("variable '%s', default '%v'", v.name, v.defaultValue)
}

func Environment[T any](name string, defaultValue T) DefaultVariable[T] {
	v := DefaultVariable[T]{name: name, defaultValue: defaultValue, t: VariableTypeEnvironment}
	register(v)

	return v
}

func register(v Variable) {
	registry = append(registry, v)
}
