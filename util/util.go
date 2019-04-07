package util

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"github.com/Masterminds/sprig"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func WrapErr(err error, msg string) error {
	return errors.WithStack(errors.Wrap(err, msg))
}

func PrintErr(err error, msg string) {
	log.Println(errors.WithStack(errors.Wrap(err, msg)))
}

func UUID() string {
	return uuid.New().String()
}
func JSON(v interface{}) []byte {
	output, _ := json.MarshalIndent(v, "", "  ")
	return output
}

func Proto(msg proto.Message) []byte {
	output, _ := proto.Marshal(msg)
	return output
}

func YAML(v interface{}) []byte {
	output, _ := yaml.Marshal(v)
	return output
}

func MustGetEnv(envKey, defaultValue string) string {
	val := os.Getenv(envKey)
	if val == "" {
		val = defaultValue
	}
	if val == "" {
		log.Fatalf("%q should be set", envKey)
	}
	return val
}

func XML(v interface{}) []byte {
	output, _ := xml.Marshal(v)
	return output
}

func Render(text string, data interface{}, w io.Writer) error {
	t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

var validate = validator.New()

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password must not be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return string(hash[:]), err
	}
	return string(hash[:]), nil
}

func ComparePasswordToHash(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func Validate(obj interface{}) error {
	return validate.Struct(obj)
}

func AsMap(obj interface{}) map[string]interface{} {
	struc := structs.New(obj)
	return struc.Map()
}
