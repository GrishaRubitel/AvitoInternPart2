package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateCreationTime() (any, error) {
	time, err := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	return time, nil
}

func GenerateID() string {
	id := uuid.New()
	return id.String()
}

func TenderMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Tender{})
	if err != nil {
		return err
	}
	return nil
}

func CreateTenderFromMap(data map[string]string) (*Tender, error) {
	tender := &Tender{}

	if id, ok := data["id"]; ok {
		tender.Id = id
	} else {
		tender.Id = GenerateID()
	}

	if name, ok := data["name"]; ok {
		tender.Name = name
	} else {
		return nil, errors.New("eame is required")
	}

	if description, ok := data["description"]; ok {
		tender.Description = description
	} else {
		return nil, errors.New("eescription is required")
	}

	if serviceType, ok := data["serviceType"]; ok {
		switch TenderServiceType(serviceType) {
		case Construction, Delivery, Manufacture:
			tender.ServiceType = TenderServiceType(serviceType)
		default:
			return nil, errors.New("invalid ServiceType")
		}
	} else {
		return nil, errors.New("serviceType is required")
	}

	if status, ok := data["status"]; ok {
		switch TenderStatus(status) {
		case Created, Published, Closed:
			tender.Status = TenderStatus(status)
		default:
			return nil, errors.New("invalid Status")
		}
	} else {
		tender.Status = TenderStatus(Created)
	}

	if organizationID, ok := data["organizationId"]; ok {
		tender.OrganizationId = organizationID
	} else {
		return nil, errors.New("organizationID is required")
	}

	if versionStr, ok := data["version"]; ok {
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return nil, errors.New("invalid Version format")
		}
		tender.Version = version
	} else {
		tender.Version = 1
	}

	if createdAtStr, ok := data["createdAt"]; ok {
		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, errors.New("invalid CreatedAt format")
		}
		tender.CreatedAt = createdAt
	} else {
		createdAt, err := GenerateCreationTime()
		cA, ok := createdAt.(time.Time)
		if err != nil && ok {
			return nil, errors.New("createdAt is required")
		} else {
			tender.CreatedAt = cA
		}
	}

	return tender, nil
}

func SelectUserTenders(db *gorm.DB, limit int, offset int, username string) ([]Tender, error) {
	var tenders []Tender

	query := db.Table("tenders t").
		Joins("JOIN organization_responsible or2 ON t.organization_id = or2.organization_id").
		Joins("JOIN employee e ON or2.user_id = e.id").
		Offset(offset).
		Where("e.username = ?", username)

	if limit != 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&tenders)

	if result.Error != nil {
		return []Tender{}, errors.New("error while processing select")
	} else {
		return tenders, nil
	}
}

func FindTenderByID(tenders []Tender, tenderID string) (Tender, error) {
	for _, tender := range tenders {
		if tender.Id == tenderID {
			return tender, nil
		}
	}

	return Tender{}, errors.New("no tender found")
}

func IsValidTenderStatus(status string) bool {
	switch TenderStatus(status) {
	case Created, Published, Closed:
		return true
	default:
		return false
	}
}

func ProcessTendersByUserAndId(db *gorm.DB, data map[string]string) (int, Tender, error) {
	var tender []Tender

	username, ok := data["username"]
	if ok {
		username, err := FindUserInTable(db, username)
		if username.Username == "" {
			return http.StatusUnauthorized, Tender{}, err
		} else if err != nil {
			return http.StatusBadRequest, Tender{}, err
		}
	} else {
		return http.StatusUnauthorized, Tender{}, errors.New("username field is required")
	}

	tender, err := SelectUserTenders(db, 0, 0, username)
	if err != nil {
		return http.StatusBadRequest, Tender{}, err
	} else {
		tender, err := FindTenderByID(tender, data["tenderid"])
		if err != nil {
			return http.StatusBadRequest, Tender{}, err
		} else {
			return http.StatusOK, tender, nil
		}
	}
}

func JSONToMap(jsonStr string) (map[string]string, error) {
	var rawData map[string]interface{}
	result := make(map[string]string)
	jsonBytes := []byte(jsonStr)

	err := json.Unmarshal(jsonBytes, &rawData)
	if err != nil {
		return nil, errors.New("error while handling operation")
	} else {
		for key, value := range rawData {
			switch v := value.(type) {
			case string:
				result[key] = v
			case float64:
				result[key] = fmt.Sprintf("%v", v)
			default:
				return nil, errors.New("unsupported value type for key: " + key)
			}
		}
		return result, nil
	}
}
