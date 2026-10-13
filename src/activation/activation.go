package activation

import (
	"fmt"
	"os"
	"strconv"
	"time"

	turbojwt "github.com/MaminirinaEdwino/turbostackapi/src/TurboJwt"
	"github.com/MaminirinaEdwino/turbostackapi/src/utils"
)

func CheckActivationToken() map[string]interface{} {
	userDir, err := os.UserHomeDir()
	utils.ErrorChecker(err)

	tokenFile, err := os.ReadFile(userDir + "/.turbo_stack/token")

	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(userDir+"/.turbo_stack", os.ModePerm)
		file, _ := os.Create(userDir + "/.turbo_stack/token")
		defer file.Close()
		token, _ := turbojwt.Encode("secret", map[string]any{
			"subscription": "free",
		}, 10, 0)
		file.Write([]byte(token))
	}
	res, tokenErr := turbojwt.Verify("secret", string(tokenFile), 0)
	// nt, _ := turbojwt.Encode("secret", map[string]interface{}{"subscription": "pro"}, 0, 0)
	// SaveTokenToFile(nt)
	fval, _ := strconv.ParseFloat(fmt.Sprint(res["exp"]), 64)
	intVal := int64(fval)
	expTime := time.Unix(intVal, 0)
	iatfval, _ := strconv.ParseFloat(fmt.Sprint(res["iat"]), 64)
	iatintVal := int64(iatfval)
	iatTime := time.Unix(iatintVal, 0)
	// dateLayout := "Monday 01 January 2006 at 15:04:05"
	if tokenErr != nil {
		return map[string]interface{}{
			"subscription": res["subscription"],
			"exp":          expTime.Format(time.UnixDate),
			"iat":          iatTime.Format(time.UnixDate),
			"message":      "old subscription has expired",
			"expired":      "true",
		}
	}
	return map[string]interface{}{
		"subscription": res["subscription"],
		"exp":          expTime.Format(time.UnixDate),
		"iat":          iatTime.Format(time.UnixDate),
		"expired":      "false",
	}
}

func SaveTokenToFile(token string) {
	userDir, _ := os.UserHomeDir()
	tokenFile, _ := os.OpenFile(userDir+"/.turbo_stack/token", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	tokenFile.WriteString(token)
}
