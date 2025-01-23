package celcoin

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/contbank/grok"

	"github.com/google/uuid"
)

type requestIDKey string

func random(n float64) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return math.Round(r.Float64() * n)
}

func mod(dividendo float64, divisor float64) float64 {
	return math.Round(dividendo - (math.Floor(dividendo/divisor) * divisor))
}

// GeneratorCPF ...
func GeneratorCPF() string {
	cpfString := ""

	rand.Seed(time.Now().UTC().UnixNano())
	cpf := rand.Perm(9)
	cpf = append(cpf, verify(cpf, len(cpf)))
	cpf = append(cpf, verify(cpf, len(cpf)))

	for _, c := range cpf {
		cpfString += strconv.Itoa(c)
	}

	return cpfString
}

func verify(data []int, n int) int {
	var total int

	for i := 0; i < n; i++ {
		total += data[i] * (n + 1 - i)
	}

	total = total % 11
	if total < 2 {
		return 0
	}
	return 11 - total
}

// GeneratorCNPJ ...
func GeneratorCNPJ() string {
	var n float64
	var n9 float64
	var n10 float64
	var n11 float64
	var n12 float64

	n = 9
	n9 = 0
	n10 = 0
	n11 = 0
	n12 = 1

	var n1 = random(n)
	var n2 = random(n)
	var n3 = random(n)
	var n4 = random(n)
	var n5 = random(n)
	var n6 = random(n)
	var n7 = random(n)
	var n8 = random(n)

	var d1 = n12*2 + n11*3 + n10*4 + n9*5 + n8*6 + n7*7 + n6*8 + n5*9 + n4*2 + n3*3 + n2*4 + n1*5
	d1 = 11 - (mod(d1, 11))
	if d1 >= 10 {
		d1 = 0
	}
	var d2 = d1*2 + n12*3 + n11*4 + n10*5 + n9*6 + n8*7 + n7*8 + n6*9 + n5*2 + n4*3 + n3*4 + n2*5 + n1*6
	d2 = 11 - (mod(d2, 11))
	if d2 >= 10 {
		d2 = 0
	}

	resultado := fmt.Sprintf("%d%d.%d%d%d.%d%d%d/%d%d%d%d-%d%d", int(n1), int(n2), int(n3), int(n4), int(n5), int(n6), int(n7), int(n8), int(n9), int(n10), int(n11), int(n12), int(d1), int(d2))

	return resultado
}

// GeneratorCellphone ...
func GeneratorCellphone() string {
	phoneString := ""
	rand.Seed(time.Now().UTC().UnixNano())

	dddArray := [3]int{11, 21, 51}
	ddd := dddArray[rand.Intn(2)]

	phone := rand.Perm(8)

	for _, c := range phone {
		phoneString += strconv.Itoa(c)
	}

	phoneString = strconv.Itoa(ddd) + "9" + phoneString

	return phoneString
}

// OnlyDigits ...
func OnlyDigits(value string) string {

	var newValue string

	for _, c := range value {
		if unicode.IsDigit(c) {
			newValue += string(c)
		}
	}

	return newValue
}

// IsOnlyDigits ...
func IsOnlyDigits(value string) bool {
	for _, c := range value {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

// OnlyDate ...
func OnlyDate(datetime *time.Time) *time.Time {
	if datetime == nil {
		return nil
	}
	response := time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 00, 00, 00, 00, time.UTC)
	return &response
}

// RandStringBytes ...
func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// GetRequestID ...
func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value("Request-Id").(string)
	return requestID
}

// GenerateNewRequestID
func GenerateNewRequestID(ctx context.Context) context.Context {
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, requestIDKey("Request-Id"), requestID)
	return ctx
}

var GetClientID = func() string {
	return os.Getenv("CELCOIN_CLIENT_ID")
}

// GetEnvCelcoinClientID ...
func GetEnvCelcoinClientID() *string {
	clientID := GetClientID()
	return &clientID
}

var GetClientSecret = func() string {
	return os.Getenv("CELCOIN_CLIENT_SECRET")
}

// GetEnvCelcoinClientSecret ...
func GetEnvCelcoinClientSecret() *string {
	clientSecret := GetClientSecret()
	return &clientSecret
}

// NormalizeNameWithoutSpecialCharacters ...
func NormalizeNameWithoutSpecialCharacters(value *string) *string {
	if value == nil || (value != nil && *value == "") {
		return nil
	}
	newValue := strings.Replace(*value, "&", "e", -1)
	newValue = strings.Replace(newValue, ".", "", -1)
	newValue = strings.Replace(newValue, "-", "", -1)
	newValue = strings.Replace(newValue, "/", "", -1)
	return aws.String(grok.ToTitle(newValue))
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// setRequestHeader Set header with Bankly requirements
func setRequestHeader(request *http.Request, token string, apiVersion *string, headers *http.Header) *http.Request {
	if headers != nil {
		request.Header = *headers
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("Content-type", "application/json")
	if apiVersion != nil {
		request.Header.Add("api-version", *apiVersion)
	}
	return request
}

// NewContextRequestID ...
func NewContextRequestID(ctx context.Context) context.Context {
	requestID := uuid.New().String()
	ctx = context.WithValue(ctx, "Request-Id", requestID)
	return ctx
}

// NewRequestID ...
func NewRequestID() string {
	return uuid.New().String()
}

// ParseStringToCelcoinTime ... faz o parse de uma string para time.Time
func ParseStringToCelcoinTime(value string, layout string) (time.Time, error) {
	cleanStr := strings.TrimSuffix(value, "Z")
	t, err := time.Parse(layout, cleanStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("erro ao converter string para time.Time: %w", err)
	}
	return t, nil
}

/* PIX */

var requiredFieldsForPixCashOut = map[string][]string{
	"DYNAMIC_QRCODE": {"amount", "clientCode", "transactionIdentification", "endToEndId", "debitParty.account", "creditParty.bank", "creditParty.key", "creditParty.name"},
	"STATIC_QRCODE":  {"amount", "clientCode", "transactionIdentification", "endToEndId", "debitParty.account", "creditParty.bank", "creditParty.key", "creditParty.name", "creditParty.accountType"},
	"DICT":           {"amount", "clientCode", "endToEndId", "debitParty.account", "creditParty.bank", "creditParty.key"},
	"MANUAL":         {"amount", "clientCode", "initiationType", "paymentType", "debitParty.account", "creditParty.bank", "creditParty.account", "creditParty.branch", "creditParty.taxId", "creditParty.accountType"},
}

func validatePixCashOut(req PixCashOutRequest) error {
	// Validação de campos obrigatórios com base no requiredFieldsForPixCashOut
	fields, ok := requiredFieldsForPixCashOut[req.InitiationType]
	if !ok {
		return fmt.Errorf("unknown initiationType: %s", req.InitiationType)
	}

	for _, field := range fields {
		if !fieldExists(req, field) {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validação adicional por InitiationType
	switch req.InitiationType {
	case "MANUAL":
		if req.TransactionIdentification != "" || req.CreditParty.Key != "" || req.EndToEndId != "" || req.DebitParty.Account == "" {
			return fmt.Errorf("invalid fields for InitiationType MANUAL")
		}
	case "DICT":
		if req.TransactionIdentification != "" || req.CreditParty.Key == "" || req.EndToEndId == "" {
			return fmt.Errorf("invalid fields for InitiationType DICT")
		}
	case "STATIC_QRCODE":
		if len(req.TransactionIdentification) > 25 || req.CreditParty.Key == "" || req.EndToEndId == "" {
			return fmt.Errorf("invalid fields for InitiationType STATIC_QRCODE")
		}
	case "DYNAMIC_QRCODE":
		if len(req.TransactionIdentification) < 26 || len(req.TransactionIdentification) > 35 || req.CreditParty.Key == "" || req.EndToEndId == "" {
			return fmt.Errorf("invalid fields for InitiationType DYNAMIC_QRCODE")
		}
	case "PAYMENT_INITIATOR":
		if len(req.TransactionIdentification) > 25 || req.TaxIdPaymentInitiator == "" || req.EndToEndId == "" {
			return fmt.Errorf("invalid fields for InitiationType PAYMENT_INITIATOR")
		}
	default:
		return fmt.Errorf("unknown InitiationType: %s", req.InitiationType)
	}

	// Validação adicional por TransactionType
	switch req.TransactionType {
	case "TRANSFER":
		if req.VlcpAmount != 0 || req.VldnAmount != 0 || req.WithdrawalAgentMode != "" || req.WithdrawalServiceProvider != "" {
			return fmt.Errorf("invalid fields for TransactionType TRANSFER: vlcpAmount, vldnAmount, withdrawalAgentMode, and withdrawalServiceProvider must not be filled")
		}
	case "WITHDRAWAL":
		if req.VlcpAmount != 0 || req.VldnAmount == 0 || req.WithdrawalAgentMode == "" || req.WithdrawalServiceProvider == "" {
			return fmt.Errorf("invalid fields for TransactionType WITHDRAWAL: vlcpAmount must not be filled, and vldnAmount, withdrawalAgentMode, withdrawalServiceProvider must be filled")
		}
	case "CHANGE":
		if req.VlcpAmount == 0 || req.VldnAmount == 0 || req.WithdrawalAgentMode == "" || req.WithdrawalServiceProvider == "" {
			return fmt.Errorf("invalid fields for TransactionType CHANGE: all related fields must be filled")
		}
	default:
		return fmt.Errorf("unknown TransactionType: %s", req.TransactionType)
	}

	// Validação adicional por PaymentType
	switch req.PaymentType {
	case "IMMEDIATE":
		if req.Urgency != "HIGH" {
			return fmt.Errorf("invalid urgency for PaymentType IMMEDIATE: must be HIGH")
		}
	case "FRAUD":
		if req.Urgency != "NORMAL" {
			return fmt.Errorf("invalid urgency for PaymentType FRAUD: must be NORMAL")
		}
	case "SCHEDULED":
		if req.Urgency != "NORMAL" {
			return fmt.Errorf("invalid urgency for PaymentType SCHEDULED: must be NORMAL")
		}
	default:
		return fmt.Errorf("unknown PaymentType: %s", req.PaymentType)
	}

	return nil
}

// Função auxiliar para verificar se um campo existe no PixCashOutRequest
func fieldExists(req PixCashOutRequest, fieldPath string) bool {
	// Separar os campos aninhados por "."
	fields := strings.Split(fieldPath, ".")
	val := reflect.ValueOf(req)

	for _, field := range fields {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() != reflect.Struct {
			return false
		}

		val = val.FieldByName(strings.Title(field))
		if !val.IsValid() {
			return false
		}
	}

	// Verificar se o valor não é nulo ou zero
	return !isZeroValue(val)
}

// Função auxiliar para verificar se o valor é zero
func isZeroValue(val reflect.Value) bool {
	return (val.Kind() == reflect.Ptr && val.IsNil()) ||
		(val.Kind() == reflect.String && val.Len() == 0) ||
		(val.Kind() == reflect.Slice && val.Len() == 0) ||
		(val.IsValid() && reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface()))
}
