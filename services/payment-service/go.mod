module payment-service

go 1.23.0

toolchain go1.24.1

// Add replace directive for local proto package
replace payment-service/proto => ./proto

require (
	github.com/alicebob/miniredis/v2 v2.34.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/nats-io/nats.go v1.33.1
	github.com/redis/go-redis/v9 v9.8.0
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)

require (
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
