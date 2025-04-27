package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/sementrof/prod1/internal/models"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	conn *pgx.Conn
	log  *logrus.Logger
}

func NewRepository(conn *pgx.Conn, logger *logrus.Logger) *Repository {
	return &Repository{
		conn: conn,
		log:  logger,
	}
}

func (r *Repository) CreateUser(user models.User) (string, error) {
	query := `INSERT INTO users(name, surname, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var userID string
	err := r.conn.QueryRow(context.Background(), query, user.Name, user.Surname, user.Email, user.Password).Scan(&userID)
	if err != nil {
		r.log.Errorf("Error creating user: %v", err)
		return "", err
	}
	return userID, nil
}

func (r *Repository) GetUserByEmail(user models.User, email string) (models.User, error) {
	query := `SELECT id, name, surname, email, password, COALESCE(icon, '') AS icon, greenPoints, created_at, updated_at FROM users WHERE email = $1;`
	err := r.conn.QueryRow(context.Background(), query, email).Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.Password, &user.Icon, &user.GreenPoints, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.log.Errorf("Error getting user by email: %v", err)
		return user, err
	}
	return user, nil
}

func (r *Repository) CreateOrganization(organization models.Organization) (string, error) {
	var organizationID string
	query := `INSERT INTO organizations(name, address, email, password, lisenceNumber, individualTaxpayerNumber) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.conn.QueryRow(context.Background(), query, organization.Name, organization.Address, organization.Email, organization.Password, organization.LisenceNumber, organization.IndividualTaxpayerNumber).Scan(&organizationID)
	if err != nil {
		r.log.Errorf("Error creating organization: %v", err)
	}
	return organizationID, nil
}

// func (r *Repository) CreateCase(problem models.Case) (string, error) {
// 	var caseID string
// 	query := `INSERT INTO organizations(title, description, photos , coordinates, adsress, applicant, performer, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
// 	err := r.conn.QueryRow(context.Background(), query, problem.Title, problem.Description, problem.Photos, problem.Coordinates, problem.Adress, problem.Applicant, problem.Performer, problem.Status).Scan(&caseID)
// 	if err != nil {
// 		r.log.Errorf("Error creating organization: %v", err)
// 	}
// 	return caseID, nil
// }
