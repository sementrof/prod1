package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gocql/gocql"
	Cassandra "github.com/sementrof/prod1/internal/cassandra"
	"github.com/sementrof/prod1/internal/kafka"
	"github.com/sementrof/prod1/internal/models"
	"github.com/sementrof/prod1/internal/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type ApiInterface interface {
	CreateUserPost(w http.ResponseWriter, r *http.Request)
	// GetUser(w http.ResponseWriter, r *http.Request)
	LoginUserPost(w http.ResponseWriter, r *http.Request)
	CreateOrganizationPost(w http.ResponseWriter, r *http.Request)
}

type ApiImplemented struct {
	log       *logrus.Logger
	repo      *repository.Repository
	producer  kafka.Producer
	cassandra Cassandra.Cassandra
}

func NewServer(logger *logrus.Logger, repo *repository.Repository, kafka *kafka.Producer, cass *Cassandra.Cassandra) *ApiImplemented {
	return &ApiImplemented{
		log:       logger,
		repo:      repo,
		producer:  *kafka,
		cassandra: *cass,
	}
}

func (im *ApiImplemented) logToKafka(event, message string, extra ...string) {
	logEntry := map[string]string{"event": event, "message": message}
	if len(extra) == 2 {
		logEntry[extra[0]] = extra[1]
	}

	logData, _ := json.Marshal(logEntry)
	_ = im.producer.Publish("logs", logData)
}

// CreateUserPost godoc
// @Summary Create a new user
// @Description Register a new user with name, surname, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserInput true "User input"
// @Success 201 {string} string "User created with ID"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /registration/user [post]
func (im *ApiImplemented) CreateUserPost(w http.ResponseWriter, r *http.Request) {

	var input models.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, models.ErrRequestPayload.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusUnprocessableEntity)
		im.log.Error(err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		im.log.Error(err)
		return
	}
	input.Password = string(hashedPassword)

	user := models.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Email:    input.Email,
		Password: input.Password,
	}

	userID, err := im.repo.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		im.logToKafka("error", "Organization created", "organization_id", userID)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": userID})
	im.logToKafka("success", "User created", "user_id", userID)
	uuid := gocql.TimeUUID()
	im.log.Infof("write in cassandra")
	im.cassandra.InsertLog(uuid, "susses", "create user")

}

// func (im *ApiImplemented) GetUser(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello World"))
// 	w.WriteHeader(http.StatusCreated)
// }

// CreateOrganizationPost godoc
// @Summary Create a new organization
// @Description Register a new organization
// @Tags organization
// @Accept json
// @Produce json
// @Param user body models.OrganizationInput true "Organization"
// // @Success 201 {string} string "Organization created with ID"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /registration/organization [post]
func (im *ApiImplemented) CreateOrganizationPost(w http.ResponseWriter, r *http.Request) {
	var input models.OrganizationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, models.ErrRequestPayload.Error(), http.StatusBadRequest)
		return
	}
	if err := validate.Struct(input); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusUnprocessableEntity)
		im.log.Error(err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		im.log.Error(err)
		return
	}
	input.Password = string(hashedPassword)
	organization := models.Organization{
		Name:                     input.Name,
		Address:                  input.Address,
		Email:                    input.Email,
		Password:                 input.Password,
		LisenceNumber:            input.LisenceNumber,
		IndividualTaxpayerNumber: input.IndividualTaxpayerNumber,
	}
	organizationID, err := im.repo.CreateOrganization(organization)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		im.logToKafka("error", "Organization created", "organization_id", organizationID)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": organizationID})
	im.logToKafka("success", "Organization created", "organization_id", organizationID)
}

// func (im *ApiImplemented) CreateCasePost(w http.ResponseWriter, r *http.Request) {
// 	var inputCase models.Case
// 	if err := json.NewDecoder(r.Body).Decode(&inputCase); err != nil {
// 		http.Error(w, models.ErrRequestPayload.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if r.Method == "Post" {
// 		err := r.ParseMultipartForm(32 << 20)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		file, filereader, err := r.FormFile("image")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}
// 		defer file.Close()
// 		fileBytes := new(bytes.Buffer)

// 		if _, err := io.Copy(fileBytes, file); err != nil {
// 			http.Error(w, "Ошибка при чтении файла: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		// Загружаем в MinIO
// 		objectID, err := im.minioClient.CreateOne(service.FileDataType{
// 			FileName: fileHeader.Filename,
// 			Data:     fileBytes.Bytes(),
// 		})
// 		if err != nil {
// 			http.Error(w, "Ошибка загрузки файла в MinIO: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		caseProblem := models.Case{
// 			Title:       inputCase.Title,
// 			Description: inputCase.Description,
// 			Photos:      inputCase.Photos,
// 			Coordinates: inputCase.Coordinates,
// 			Adress:      inputCase.Adress,
// 			Applicant:   inputCase.Applicant,
// 			Performer:   inputCase.Performer,
// 			Status:      inputCase.Status,
// 		}
// 		_, err = im.repo.CreateCase(caseProblem)
// 		if err != nil {
// 			http.Error(w, "Ошибка при сохранении кейса: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 		w.Write([]byte(fmt.Sprintf("Кейс успешно создан с ID изображения: %s", objectID)))

// 	}
// 	_ = service.FileDataType{} // временная заглушка, чтобы не удалялся импорт

// }
