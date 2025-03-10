package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Auth struct {
	DcName       string
	RefreshToken string
	AccessToken  string
	StoredTime   int64
}

func GetAccessToken(dc string, password string) (string, error) {
	if dc == "" {
		return "", fmt.Errorf("Invalid dc have been given: %s", dc)
	}
	auth, err := FindByDC(dc)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", fmt.Errorf("Given dc_name %s is not exists, Use the login and store this dc authentication", dc)
	}
	if auth.AccessToken == "" || time.UnixMilli(auth.StoredTime).Add(45*time.Minute).Before(time.Now()) {
		accessToken, err := GenerateAccessToken(*auth, password)
		if err != nil {
			return "", err
		}
		if accessToken == "" {
			return "", errors.New("Unable to generate latest access_token")
		}
		auth.AccessToken, err = Encrypt(password, accessToken)
		if err != nil {
			return "", err
		}
		auth.StoredTime = time.Now().UnixMilli()
		UpdateAuth(*auth)
	}
	return Decrypt(password, auth.AccessToken)
}

func LoginToDC(dcName string, password string) (string, error) {
	auth, err := FindByDC(dcName)
	if err != nil {
		return "", err
	}
	if auth != nil {
		// Already logged in
		return "", nil
	}
	url := fmt.Sprintf("https://accounts.%s/oauth/v2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code&access_type=offline", dcName, viper.GetString(CLIENT_ID), viper.GetString(REDIRECT_URI), viper.GetString(SCOPE))
	fmt.Printf("Open this url in your browser => %s\n", url)

	loginCallbackWait := sync.WaitGroup{}
	loginCallbackWait.Add(1)

	startCallbackServer(dcName, password, &loginCallbackWait)

	loginCallbackWait.Wait()
	return url, nil
}

func GenerateAccessToken(auth Auth, password string) (string, error) {
	baseUrl := fmt.Sprintf("https://accounts.%s/oauth/v2/token", auth.DcName)

	refreshToken, err := Decrypt(password, auth.RefreshToken)
	params := url.Values{}
	params.Add(CLIENT_ID, viper.GetString(CLIENT_ID))
	params.Add(CLIENT_SECRET, viper.GetString(CLIENT_SECRET))
	params.Add(REDIRECT_URI, strings.TrimSpace(viper.GetString(REDIRECT_URI)))
	params.Add(GRANT_TYPE, "refresh_token")
	params.Add(REFRESH_TOKEN, refreshToken)

	url := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)

	return response["access_token"], nil
}

func generateAndStoreRefreshTokenWithCode(dcName string, password string, code string) error {
	accessTokenURL := fmt.Sprintf("https://accounts.%s/oauth/v2/token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&grant_type=authorization_code", dcName, viper.GetString(CLIENT_ID), viper.GetString(CLIENT_SECRET), code, viper.GetString(REDIRECT_URI))

	resp, err := http.Post(accessTokenURL, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)

	refreshToken := response["refresh_token"]
	if refreshToken == "" {
		return fmt.Errorf("Didn't get refresh token in token request, It may have already generated. Please delete the 'ZShell' application session in connected apps accounts.%s", dcName)
	}

	refreshToken, err = Encrypt(password, refreshToken)
	_, err = AddAuth(dcName, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func startCallbackServer(dcName string, password string, wg *sync.WaitGroup) {
	server := http.Server{
		Addr: fmt.Sprintf(":%s", viper.GetString(PORT)),
	}
	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			cobra.CheckErr(err)
		}
		code := u.Query().Get("code")
		err = generateAndStoreRefreshTokenWithCode(dcName, password, code)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = server.Close()
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Println("Done")
		wg.Done()
	})
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		cobra.CheckErr(err)
	}
}

func storeAuths(auths []Auth) error {
	authsJSON, err := json.Marshal(&auths)
	authFilePath, err := GetAuthDataFile()
	if err != nil {
		return err
	}
	err = os.WriteFile(authFilePath, authsJSON, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func FindByDC(dcName string) (*Auth, error) {
	auths, err := GetAllAuths()
	if err != nil {
		return nil, err
	}
	for _, auth := range auths {
		if auth.DcName == dcName {
			return &auth, nil
		}
	}
	return nil, nil
}

func UpdateAuth(auth Auth) error {
	auths, err := GetAllAuths()
	if err != nil {
		return err
	}

	for i := range auths {
		currentAuth := &auths[i]
		if currentAuth.DcName == auth.DcName {
			currentAuth.AccessToken = auth.AccessToken
			currentAuth.RefreshToken = auth.RefreshToken
			currentAuth.StoredTime = auth.StoredTime
			break
		}
	}
	storeAuths(auths)
	return nil
}

func AddAuth(dcName string, refreshToken string) (*Auth, error) {
	auth := Auth{
		DcName:       dcName,
		RefreshToken: refreshToken,
		StoredTime:   time.Now().UnixMilli(),
	}
	auths, err := GetAllAuths()
	if err != nil {
		return nil, err
	}
	auths = append(auths, auth)
	err = storeAuths(auths)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func GetAllAuths() ([]Auth, error) {
	authFilePath, err := GetAuthDataFile()
	if err != nil {
		return nil, nil
	}
	if !IsFileExists(authFilePath) {
		err = CreateDefaultAuthData()
		if err != nil {
			return nil, err
		}
	}
	data, err := os.ReadFile(authFilePath)
	if err != nil {
		return nil, err
	}
	var authData []Auth
	err = json.Unmarshal(data, &authData)
	if err != nil {
		return nil, err
	}
	return authData, nil
}

func CreateDefaultAuthData() error {
	authFilePath, err := GetAuthDataFile()
	if err != nil {
		return err
	}
	content := "[]"
	err = os.WriteFile(authFilePath, []byte(content), 0o755)
	if err != nil {
		return err
	}
	return nil
}
