package models

import (
	"time"
)

type User struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Icon        *string   `json:"icon"`
	GreenPoints int32     `json:"greenPoints"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Organization struct {
	Id                       string
	Name                     string
	Address                  string
	Email                    string
	Password                 string
	LisenceNumber            string
	IndividualTaxpayerNumber int
}

// type Case struct {
// 	Id          string
// 	Title       string
// 	Description string
// 	Photos      string
// 	Coordinates string
// 	Adress      string
// 	Applicant   string
// 	Performer   string
// 	Status      bool
// }

// not all fields are required

type UserInput struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Userlogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type OrganizationInput struct {
	Name                     string `json:"name" validate:"required"`
	Address                  string `json:"address" validate:"required"`
	Email                    string `json:"email" validate:"required,email"`
	Password                 string `json:"password" validate:"required"`
	LisenceNumber            string `json:"LisenceNumber" validate:"required,e164"`
	IndividualTaxpayerNumber int    `json:"IndividualTaxpayerNumber" validate:"required"`
}
