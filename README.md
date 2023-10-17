# godor-jwt

JWT helpers and middleware for godor (got the door :wink:)

## Installation

```bash
go get -u github.com/arcvats/godor-jwt
```

## Usage

godor-jwt provides helpers and middleware for Encoding and Decoding a JWT token with HMAC methods

### Config Struct

```go
type Config struct {
    Algorithm string // default is HS256
    Secret    string // required
    Expiry    int64 // default is 60 minutes
}
```

### JWT helpers

#### Encode

```go
package main
import (
    godorjwt "github.com/arcvats/godor-jwt"
)

func main() {
	config := godorjwt.Config{
		Algorithm: "HS256",
		Secret: "secret",
		Expiry: 60,
    }
	payload := map[string]any{
        "sub": "1234567890",
        "name": "John Doe",
        "admin": true,
    }
	token, jti, expiry, error := godorjwt.Encode(payload, config)
}
```

#### Decode

```go
package main
import (
    godorjwt "github.com/arcvats/godor-jwt"
)

func main() {
	config := godorjwt.Config{
		Secret: "secret",
    }
	tokenString := "some.jwt.token"
	decodedTokenClaims, error := godorjwt.Decode(tokenString, config)
}
```

### JWT middleware

```go
package main
import (
    godorjwt "github.com/arcvats/godor-jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	config := godorjwt.Config{
		Secret: "secret",
    }
	app.Use(config.Decoder())
	/* 
	    Token is acquired from the Authorization header or jwt Cookie
        Now you can access the decoded token claims in the context
        claims := c.Locals("decodedToken").(map[string]any)
    */
}
```

## License

MIT

