package main

import (
	"errors"
	"fmt"
	models "grisha_rubitel/avito_p2/dbModels"
	"net/http"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CheckServer(conn string) (int, string, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return http.StatusBadRequest, "DB connection not established", errors.New("DB connection error")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return http.StatusBadRequest, "Error while processing DB connection check", errors.New("DB connection error")
	}

	if err := sqlDB.Ping(); err != nil {
		return http.StatusBadRequest, "Error while pinging DB", errors.New("DB ping error")
	}

	return http.StatusOK, "Ok", nil
}

func GetTenders(db *gorm.DB, data map[string]string) (int, string, error) {
	var tenders []models.Tender

	limit, err := strconv.Atoi(data["limit"])
	if err != nil {
		return http.StatusBadRequest, "", errors.New("invalid type of limit")
	}

	offset, err := strconv.Atoi(data["offset"])
	if err != nil {
		return http.StatusBadRequest, "", errors.New("invalid type of offset")
	}

	query := db.Limit(limit).Offset(offset)

	serviceType, ok := data["service_type"]
	if ok && serviceType != "" {
		query.Where("service_type = ?", serviceType)
	}

	result := query.Find(&tenders)
	if result.Error != nil {
		return http.StatusForbidden, "", errors.New("error while fetching select's results")
	}
	tendersS, err := ToJsonMulti(tenders)
	if err != nil {
		return http.StatusBadRequest, "", errors.New("error while fetching select's results")
	} else {
		return http.StatusOK, tendersS, nil
	}

}

func CreateTender(db *gorm.DB, data map[string]string) (int, string, error) {
	newTender, err := models.CreateTenderFromMap(data)
	if err != nil {
		return http.StatusBadRequest, "", err
	} else {
		result := db.Create(&newTender)
		if result.Error != nil {
			return http.StatusBadRequest, "", errors.New("error while updating table")
		} else {
			resp, err := ToJson(newTender)
			fmt.Println(resp)
			if err != nil {
				return http.StatusInternalServerError, "", err
			} else {
				return http.StatusOK, resp, nil
			}
		}
	}
}

func GetUserTenders(db *gorm.DB, data map[string]string) (int, string, error) {
	var tenders []models.Tender

	limit, err := strconv.Atoi(data["limit"])
	if err != nil {
		return http.StatusBadRequest, "", errors.New("invalid type of limit")
	}

	offset, err := strconv.Atoi(data["offset"])
	if err != nil {
		return http.StatusBadRequest, "", errors.New("invalid type of offset")
	}

	username, err := models.FindUserInTable(db, data["username"])
	if username.Username == "" {
		return http.StatusUnauthorized, "", err
	} else if err != nil {
		return http.StatusBadRequest, "", err
	}

	tenders, err = models.SelectUserTenders(db, limit, offset, username.Username)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	resp, err := ToJsonMulti(tenders)
	if err != nil {
		return http.StatusInternalServerError, "", err
	} else {
		return http.StatusOK, resp, nil
	}
}

func GetTenderStatus(db *gorm.DB, data map[string]string) (int, string, error) {
	code, tender, err := models.ProcessTendersByUserAndId(db, data)
	if err != nil {
		return code, "", err
	} else {
		return code, string(tender.Status), nil
	}
}

func UpdateTenderStatus(db *gorm.DB, data map[string]string) (int, string, error) {
	status := data["status"]
	if !models.IsValidTenderStatus(status) {
		return http.StatusBadRequest, "", errors.New("invalid status")
	} else {
		code, tender, err := models.ProcessTendersByUserAndId(db, data)
		if err != nil {
			return code, "", err
		}

		if tender.Id == "" {
			return http.StatusNotFound, "", errors.New("tender not found")
		}

		tender.Status = models.TenderStatus(status)
		err = db.Save(&tender).Error
		if err != nil {
			return http.StatusForbidden, "", errors.New("aserror while handling operation")
		} else {
			resp, err := ToJson(tender)
			if err != nil {
				return http.StatusInternalServerError, "", err
			} else {
				return http.StatusOK, resp, nil
			}
		}
	}
}

func EditTender(db *gorm.DB, data map[string]string) (int, string, error) {
	code, tender, err := models.ProcessTendersByUserAndId(db, data)
	if err != nil {
		return code, "", err
	} else {
		tenderJSON, err := ToJson(tender)
		if err != nil {
			return http.StatusInternalServerError, "", err
		}
		tenderMap, err := models.JSONToMap(tenderJSON)
		if err != nil {
			return http.StatusInternalServerError, "", err
		}

		tenderMap = UpdateMapFromAnother(tenderMap, data)
		ver, err := strconv.Atoi(tenderMap["version"])
		if err != nil {
			return http.StatusInternalServerError, "", errors.New("parsing error")
		}
		tenderMap["version"] = strconv.Itoa(ver + 1)

		newTender, err := models.CreateTenderFromMap(tenderMap)
		if err != nil {
			return http.StatusBadRequest, "", err
		} else {
			err = db.Save(&newTender).Error
			if err != nil {
				return http.StatusForbidden, "", errors.New("error while handling operation")
			} else {
				resp, err := ToJson(*newTender)
				if err != nil {
					return http.StatusInternalServerError, "", err
				} else {
					return http.StatusOK, resp, nil
				}
			}
		}
	}
}

func RollbackTender(conn string) (int, string, error) {
	return 1, "", nil
}
