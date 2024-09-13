package main

import (
	"errors"
	"fmt"
	models "grisha_rubitel/avito_p2/dbModels"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func CreateBid(db *gorm.DB, data map[string]string) (int, string, error) {
	code, bid, err := models.CreateBidFromMap(db, data)
	if err != nil {
		return code, "", err
	} else {
		result := db.Create(&bid)
		if result.Error != nil {
			return http.StatusBadRequest, "", errors.New("error while updating table")
		} else {
			resp, err := ToJson(bid)
			fmt.Println(resp)
			if err != nil {
				return http.StatusInternalServerError, "", err
			} else {
				return http.StatusOK, resp, nil
			}
		}
	}
}

func GetUserBids(db *gorm.DB, data map[string]string) (int, string, error) {
	var bids []models.Bid

	code, bids, err := models.SelectUserBid(db, data)
	if err != nil {
		return code, "", err
	} else {
		resp, err := ToJsonMulti(bids)
		if err != nil {
			return http.StatusInternalServerError, "", err
		} else {
			return http.StatusOK, resp, nil
		}
	}
}

func GetBidsForTender(db *gorm.DB, data map[string]string) (int, string, error) {
	var bids []models.Bid

	code, bids, err := models.SelectUserBid(db, data)
	if err != nil {
		return code, "", err
	} else {
		var respBid []models.Bid
		for _, bid := range bids {
			if bid.TenderId == data["tenderid"] {
				respBid = append(respBid, bid)
			}
		}
		resp, err := ToJsonMulti(respBid)
		if err != nil {
			return http.StatusBadRequest, "", err
		} else {
			return http.StatusOK, resp, nil
		}
	}
}

func GetBidStatus(db *gorm.DB, data map[string]string) (int, *models.Bid, error) {
	var bids []models.Bid

	code, bids, err := models.SelectUserBid(db, data)
	if err != nil {
		return code, &models.Bid{}, err
	} else {
		bid, err := models.FindBidByID(bids, data["bidid"])
		if err != nil {
			return http.StatusBadRequest, &models.Bid{}, err
		}
		return http.StatusOK, bid, nil
	}
}

func UpdateBidStatus(db *gorm.DB, data map[string]string) (int, string, error) {
	code, bid, err := GetBidStatus(db, data)
	if err != nil {
		return code, "", err
	} else {
		status := data["status"]
		if !models.IsValidBidStatus(status) {
			return http.StatusBadRequest, "", errors.New("invalid status type")
		} else {
			bid.Status = models.BidStatus(status)

			result := db.Save(bid).Error
			if result != nil {
				return http.StatusInternalServerError, "", errors.New("error while saving new data")
			} else {
				resp, err := ToJson(bid)
				if err != nil {
					return http.StatusInternalServerError, "", err
				} else {
					return http.StatusOK, resp, nil
				}
			}
		}
	}
}

func EditBid(db *gorm.DB, data map[string]string) (int, string, error) {
	var bids []models.Bid

	code, bids, err := models.SelectUserBid(db, data)
	if err != nil {
		return code, "", err
	} else {
		bid, err := models.FindBidByID(bids, data["bidid"])
		if err != nil {
			return http.StatusBadRequest, "", err
		} else {

			bidJSON, err := ToJson(bid)
			if err != nil {
				return http.StatusInternalServerError, "", err
			}
			bidMap, err := models.JSONToMap(bidJSON)
			if err != nil {
				return http.StatusInternalServerError, "", err
			}

			bidMap = UpdateMapFromAnother(bidMap, data)
			ver, err := strconv.Atoi(bidMap["version"])
			if err != nil {
				return http.StatusInternalServerError, "", errors.New("parsing error")
			}
			bidMap["version"] = strconv.Itoa(ver + 1)

			code, newBid, err := models.CreateBidFromMap(db, bidMap)
			if err != nil {
				return code, "", err
			} else {
				err = db.Save(&newBid).Error
				if err != nil {
					return http.StatusForbidden, "", errors.New("error while handling operation")
				} else {
					resp, err := ToJson(*newBid)
					if err != nil {
						return http.StatusInternalServerError, "", err
					} else {
						return http.StatusOK, resp, nil
					}
				}
			}
		}
	}
}

func SubmitBidDecision(db *gorm.DB, data map[string]string) (int, string, error) {
	code, bidsRevs, err := models.SelectUserBidReviews(db, data)
	if err != nil {
		return code, "", err
	} else {
		bidsRevs, err := models.FindReviewByID(bidsRevs, data["bidid"])
		if err != nil && bidsRevs == nil {
			code, bidsRevs, err = models.CreateBidReview(db, data)
			if err != nil {
				return http.StatusBadRequest, "", err
			}
		} else if err != nil && bidsRevs != nil {
			return http.StatusNotFound, "", err
		}

		if models.IsValidBidDecision(data["decision"]) {
			bidsRevs.Decision = models.BidDecision(data["decision"])
			resp, err := ToJson(bidsRevs)
			if err != nil {
				return http.StatusInternalServerError, "", err
			} else {
				return http.StatusOK, resp, nil
			}
		} else {
			return http.StatusBadRequest, "", errors.New("invalid decision")
		}
	}
}

// func SubmitBidFeedback(db *gorm.DB, data map[string]string) (int, string, error) {

// }

// func RollbackBid(db *gorm.DB, data map[string]string) (int, string, error) {

// }

// func GetBidReviews(db *gorm.DB, data map[string]string) (int, string, error) {

// }
