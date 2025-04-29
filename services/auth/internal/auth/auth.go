package auth

type Storager interface {
	AddUser()
	GetUser()
	GetToken()
	AddToken()
}

type AuthProvider struct {
	s Storager
}

func NewAuth(storage Storager) *AuthProvider {
	return &AuthProvider{
		s: storage,
	}
}

func (a AuthProvider) Register() {}
func (a AuthProvider) Auth()     {}
func (a AuthProvider) Refresh()  {}
func (a AuthProvider) Validate() {}
