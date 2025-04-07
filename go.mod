module github.com/muidea/magicToolkit

go 1.23.0

replace (
	github.com/muidea/magicCommon => ../magicCommon
	github.com/muidea/magicEngine => ../magicEngine
)

require (
	github.com/mark3labs/mcp-go v0.17.0
	github.com/muidea/magicCommon v1.3.122
	github.com/muidea/magicEngine v1.4.0
)

require github.com/yosida95/uritemplate/v3 v3.0.2 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/muidea/seelog v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.23.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
