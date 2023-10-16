package godorjwt

import (
	"go.uber.org/goleak"
	"testing"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestNew(t *testing.T) {
	config := New("secret")
	if config.Secret != "secret" {
		t.Errorf("Expected secret to be secret, got %s", config.Secret)
	}
	if config.Expiry != 0 {
		t.Errorf("Expected expiry to be 0, got %d", config.Expiry)
	}
	if config.Algorithm != "" {
		t.Errorf("Expected algorithm to be empty, got %s", config.Algorithm)
	}
}

func TestEncode(t *testing.T) {
	config := New("secret")
	payload := map[string]any{"name": "John Doe"}
	token, _, _, err := Encode(payload, *config)
	if err != nil {
		t.Errorf("Expected error to be nil, got %s", err)
	}
	if token == "" {
		t.Errorf("Expected token to be not empty, got %s", token)
	}
}

func TestDecode(t *testing.T) {
	config := New("secret")
	payload := map[string]any{"name": "John Doe"}
	token, _, _, _ := Encode(payload, *config)
	decoded, err := Decode(token, *config)
	if err != nil {
		t.Errorf("Expected error to be nil, got %s", err)
	}
	if decoded["name"] != "John Doe" {
		t.Errorf("Expected name to be John Doe, got %s", decoded["name"])
	}
}

func TestDecodeWithEmptySecret(t *testing.T) {
	config := New("")
	payload := map[string]any{"name": "John Doe"}
	token, _, _, _ := Encode(payload, *config)
	_, err := Decode(token, *config)
	if err == nil {
		t.Errorf("Expected error to be not nil, got %s", err)
	}
}

func TestDecodeWithInvalidToken(t *testing.T) {
	config := New("secret")
	_, err := Decode("invalid", *config)
	if err == nil {
		t.Errorf("Expected error to be not nil, got %s", err)
	}
}

//
//func TestDecodeWithInvalidSigningMethod(t *testing.T) {
//	config := New("secret")
//	payload := map[string]any{"name": "John Doe"}
//	token, _, _, _ := Encode(payload, *config)
//	_, err := Decode(token, *config)
//	if err == nil {
//		t.Errorf("Expected error to be not nil, got %s", err)
//	}
//}

func TestDecodeWithEmptyToken(t *testing.T) {
	config := New("secret")
	_, err := Decode("", *config)
	if err == nil {
		t.Errorf("Expected error to be not nil, got %s", err)
	}
}

func TestDecodeWithEmptyConfig(t *testing.T) {
	_, err := Decode("", Config{})
	if err == nil {
		t.Errorf("Expected error to be not nil, got %s", err)
	}
}

//func TestConfig_Decoder(t *testing.T) {
//	config := New("secret")
//	decoder := config.Decoder()
//	if decoder == nil {
//		t.Errorf("Expected decoder to be not nil, got %s", decoder)
//	}
//}
//
//func TestConfig_DecoderWithEmptySecret(t *testing.T) {
//	config := New("")
//	decoder := config.Decoder()
//	if decoder != nil {
//		t.Errorf("Expected decoder to be nil, got %s", decoder)
//	}
//}
//
//func TestConfig_DecoderWithEmptyConfig(t *testing.T) {
//	config := Config{}
//	decoder := config.Decoder()
//	if decoder != nil {
//		t.Errorf("Expected decoder to be nil, got %s", decoder)
//	}
//}
//
//func TestConfig_DecoderWithInvalidSigningMethod(t *testing.T) {
//	config := New("secret")
//	config.Algorithm = "invalid"
//	decoder := config.Decoder()
//	if decoder != nil {
//		t.Errorf("Expected decoder to be nil, got %s", decoder)
//	}
//}
