package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/argon2"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	UnixTimestamp map[int64]int64
	Step          map[int64]int
	RequestId     map[int64]string
)

func Tracer() TracerModel {

	var model TracerModel
	pc, fileName, line, ok := runtime.Caller(1)
	if !ok {
		return model
	}

	model.FileName = fileName
	model.Line = line

	callerFunction := runtime.FuncForPC(pc)
	if callerFunction != nil {
		model.FunctionName = callerFunction.Name()
	}

	return model
}

func GetResponseIdAndLanguage(c *gin.Context) (int64, string) {

	language := getLanguage(c)
	responseId, exist := c.Get("response-id")
	if !exist {
		return time.Now().UnixNano(), language
	}

	int64Value, ok := responseId.(int64)
	if !ok {
		return time.Now().UnixNano(), language
	}

	return int64Value, language
}

func GetDuration(responseId int64) string {

	newUnixNano := time.Now().UnixNano()
	duration := UnixTimestamp[responseId]
	elapsed := newUnixNano - duration
	if UnixTimestamp == nil {
		UnixTimestamp = make(map[int64]int64)
	}
	UnixTimestamp[responseId] = newUnixNano
	ms := float64(elapsed) / float64(time.Millisecond)

	return fmt.Sprintf("%v", ms)
}

func GetStep(responseId int64) string {
	return strconv.Itoa(Step[responseId])
}

func GetNextStep(responseId int64) string {
	step := Step[responseId]
	if step == 0 {
		Step = make(map[int64]int)
		Step[responseId] = 1
		return strconv.Itoa(Step[responseId])
	}
	Step[responseId] = step + 1
	return strconv.Itoa(Step[responseId])
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func CreatePassword(password string, salt []byte) (string, error) {
	const (
		time    = 1         // Number of iterations
		memory  = 64 * 1024 // Memory usage in KB
		threads = 4         // Number of threads
		keyLen  = 32        // Length of the hash output
	)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s.%s", saltBase64, hashBase64), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := split(encodedHash, ".")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return sha256.Sum256(computedHash) == sha256.Sum256(expectedHash), nil
}

func IsValidEmail(email string) bool {
	// Regular expression for a valid email
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Match the email against the regex
	return re.MatchString(email)
}

func split(s, sep string) []string {
	var parts []string
	for len(s) > 0 {
		pos := strings.Index(s, sep)
		if pos == -1 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:pos])
		s = s[pos+len(sep):]
	}
	return parts
}

func getLanguage(c *gin.Context) string {

	language, exist := c.Get("Accept-Language")
	if !exist {
		return "EN"
	}

	strLanguage, ok := language.(string)
	if !ok {
		return "EN"
	}

	return strLanguage
}
