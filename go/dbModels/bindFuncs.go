package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func BindMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Tender{})
	if err != nil {
		return err
	}
	return nil
}

func ToJsonBid(tenders []Bid) (string, error) {
	jsonData, err := json.Marshal(tenders)
	if err != nil {
		return "", errors.New("error while formating query result to JSON")
	}

	return string(jsonData), nil
}

func CreateBidFromMap(db *gorm.DB, data map[string]string) (int, *Bid, error) {
	var bid Bid

	if id, ok := data["id"]; ok {
		bid.Id = id
	} else {
		bid.Id = GenerateID()
	}
	bid.Description = data["description"]

	if status, ok := data["status"]; ok {
		switch status {
		case string(BidStatusCreated), string(BidStatusPublished), string(BidStatusCanceled), string(BidStatusApproved), string(BidStatusRejected):
			bid.Status = BidStatus(status)
		default:
			return http.StatusBadRequest, nil, errors.New("invalid status value")
		}
	} else {
		bid.Status = BidStatusCreated
	}

	bid.TenderId = data["tenderId"]

	if authorType, ok := data["authorType"]; ok {
		switch authorType {
		case string(BidAuthorTypeOrganization), string(BidAuthorTypeUser):
			bid.AuthorType = BidAuthorType(authorType)
		default:
			return http.StatusBadRequest, nil, errors.New("invalid author type value")
		}
	} else {
		return http.StatusBadRequest, nil, errors.New("authorType field is required")
	}
	user := data["authorId"]

	_, err := FindUserInTable(db, user)
	if err != nil {
		return http.StatusUnauthorized, nil, err
	}

	bid.AuthorId = user

	if version, ok := data["version"]; ok {
		if v, err := strconv.Atoi(version); err == nil {
			bid.Version = v
		} else {
			return http.StatusBadRequest, nil, errors.New("invalid version value")
		}
	} else {
		bid.Version = 1
	}

	if createdAt, ok := data["createdAt"]; ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			bid.CreatedAt = t
		} else {
			createdAt, err := GenerateCreationTime()
			cA, ok := createdAt.(time.Time)
			if err != nil && ok {
				return http.StatusBadRequest, nil, errors.New("createdAt is required")
			} else {
				bid.CreatedAt = cA
			}
		}
	} else {
		createdAt, err := GenerateCreationTime()
		cA, ok := createdAt.(time.Time)
		if err != nil && ok {
			return http.StatusBadRequest, nil, errors.New("createdAt is required")
		} else {
			bid.CreatedAt = cA
		}
	}

	return http.StatusOK, &bid, nil
}

func CreateBidReviewFromMap(data map[string]string) (*BidReview, error) {
	var bidReview BidReview

	if id, ok := data["id"]; ok {
		bidReview.Id = id
	} else {
		bidReview.Id = GenerateID()
	}

	if description, ok := data["description"]; ok {
		bidReview.Description = description
	}

	if createdAt, ok := data["createdAt"]; ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			bidReview.CreatedAt = t
		} else {
			return &BidReview{}, errors.New("invalid createAdt value")
		}
	} else {
		createdAt, err := GenerateCreationTime()
		if err != nil {
			return &BidReview{}, err
		}
		cA, ok := createdAt.(time.Time)
		if !ok {
			return &BidReview{}, errors.New("invalid createAdt value")
		}
		bidReview.CreatedAt = cA
	}

	bidReview.BidId = data["bidid"]

	if decision, ok := data["decision"]; ok {
		if IsValidBidDecision(decision) {
			bidReview.Decision = BidDecision(decision)
		}
	}

	return &bidReview, nil
}

func SelectUserBid(db *gorm.DB, data map[string]string) (int, []Bid, error) {
	var bids []Bid

	username, err := FindUserInTable(db, data["username"])
	if username.Username == "" {
		return http.StatusUnauthorized, []Bid{}, err
	} else if err != nil {
		return http.StatusBadRequest, []Bid{}, err
	}

	query := db.Where("author_id = ?", username.ID)

	limit, err := strconv.Atoi(data["limit"])
	if err == nil && limit != 0 {
		query = query.Limit(limit)
	}

	offset, err := strconv.Atoi(data["offset"])
	if err == nil && offset != 0 {
		query = query.Offset(offset)
	}
	result := query.Find(&bids)

	if result.Error != nil {
		return http.StatusInternalServerError, []Bid{}, errors.New("error while processing select")
	} else {
		return http.StatusOK, bids, nil
	}
}

func SelectUserBidReviews(db *gorm.DB, data map[string]string) (int, []BidReview, error) {
	var bids []BidReview

	username, err := FindUserInTable(db, data["username"])
	if username.Username == "" {
		return http.StatusUnauthorized, []BidReview{}, err
	} else if err != nil {
		return http.StatusBadRequest, []BidReview{}, err
	}

	query := db.Table("bid_reviews br").Joins("join bids bs on br.bid_id = bs.id ").Where("bs.author_id = ?", username.ID)

	limit, err := strconv.Atoi(data["limit"])
	if err == nil && limit != 0 {
		query = query.Limit(limit)
	}

	offset, err := strconv.Atoi(data["offset"])
	if err == nil && offset != 0 {
		query = query.Offset(offset)
	}
	result := query.Find(&bids)

	if result.Error != nil {
		return http.StatusInternalServerError, []BidReview{}, errors.New("error while processing select")
	} else {
		return http.StatusOK, bids, nil
	}
}

func FindBidByID(bids []Bid, id string) (*Bid, error) {
	for _, bid := range bids {
		if strings.EqualFold(bid.Id, id) {
			return &bid, nil
		}
	}
	return nil, errors.New("bid not found")
}

func FindReviewByID(bids []BidReview, id string) (*BidReview, error) {
	for _, bid := range bids {
		if strings.EqualFold(bid.BidId, id) {
			return &bid, nil
		}
	}
	return nil, errors.New("review not found")
}

func IsValidBidStatus(status string) bool {
	switch BidStatus(status) {
	case BidStatusCreated, BidStatusPublished, BidStatusCanceled, BidStatusApproved, BidStatusRejected:
		return true
	default:
		return false
	}
}

func IsValidBidDecision(status string) bool {
	switch BidDecision(status) {
	case BidDecisionApproved, BidDecisionRejected:
		return true
	default:
		return false
	}
}

func CreateBidReview(db *gorm.DB, data map[string]string) (int, *BidReview, error) {
	bidReview, err := CreateBidReviewFromMap(data)
	if err != nil {
		return http.StatusInternalServerError, &BidReview{}, err
	}

	result := db.Create(&bidReview)
	if result.Error != nil {
		return http.StatusBadRequest, &BidReview{}, errors.New("error while creating bid review: " + result.Error.Error())
	}

	return http.StatusOK, bidReview, nil
}
