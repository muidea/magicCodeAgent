module github.com/muidea/magicCodeAgent

go 1.24

replace (
	github.com/muidea/magicCommon => ../magicCommon
	github.com/muidea/magicEngine => ../magicEngine
)

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/mark3labs/mcp-go v0.17.0
	github.com/muidea/magicCommon v1.4.10
	github.com/muidea/magicEngine v1.4.0
	github.com/sashabaranov/go-openai v1.38.1
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/sync v0.12.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)

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
