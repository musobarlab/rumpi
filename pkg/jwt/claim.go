package jwt

// Alg represent jwt algorithm
type Alg string

const (
	// HS256Key key
	HS256Key = "1234key"

	// HS256 const
	HS256 Alg = "HS256"

	// RS256 const
	RS256 Alg = "RS256"
)

// Claim model
type Claim struct {
	Issuer    string
	Audience  string
	Subject   string
	ExpiredAt int64
	IssuedAt  int64
	User      struct {
		ID       string
		FullName string
		Email    string
	}
	Alg Alg
}
