package datasource

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"testing"
)

var httpClient HTTPClient = &http.Client{}

func TestFetchCityByCEP(t *testing.T) {
	repository := NewCityRepository()
	result, err := repository.FetchCityByCEP("37503130")
	assert.Nil(t, err)
	assert.Equal(t, "Itajubá", result["localidade"])
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestMockTestFetchCityByCEP(t *testing.T) {

	mockHTTPClient := new(MockHTTPClient)
	mockHTTPClient.On("Get", "https://viacep.com.br/ws/08223110/json/").Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(`{"localidade": "São Paulo"}`)),
	}, nil)

	httpClient = mockHTTPClient

	repositoryMocked := NewCityRepositoryForTest(httpClient)
	result, err := repositoryMocked.FetchCityByCEP("08223110")

	mockHTTPClient.AssertExpectations(t)

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", result["localidade"])
}
