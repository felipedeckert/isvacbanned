package service

type UrlServiceMock struct{}

func (u UrlServiceMock) getSteamID(param string) (string, error) {

	return "12345678901234567", nil
}
