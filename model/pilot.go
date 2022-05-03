package model

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPilot interface {
	Save() error
	Token() string
}

// ESI Requests ///
func (x *Pilot) esiGET(path string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	r := "https://https://esi.evetech.net/"
	req, err := http.NewRequest("GET", r+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+x.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (x *Pilot) getBasicInfo() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, &x.BasicInfo); err != nil {
		return err
	}
	return nil
}
func (x *Pilot) getPortrait() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/portrait/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Portrait)
	return err
}
func (x *Pilot) getMiningRecords() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/mining/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.MiningRecords)
	return err
}
func (x *Pilot) getFleet() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/fleet/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Fleet)
	return err
}
func (x *Pilot) getTransactions() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/wallet/transactions/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Transactions)
	return err
}
func (x *Pilot) getMarketOrders() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/orders/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.MarketOrders)
	return err
}
func (x *Pilot) getAttributes() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/attributes/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Attributes)
	return err
}
func (x *Pilot) getSkills() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/skills/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Skills)
	return err
}
func (x *Pilot) getStandings() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/standings/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Standings)
	return err
}
func (x *Pilot) getLocation() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/location/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Location)
	return err
}
func (x *Pilot) getOnlineStatus() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/online/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Online)
	return err
}
func (x *Pilot) getShip() error {
	res, err := x.esiGET(CHARACTERS + x.PilotID + "/ship/")
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.Ship)
	return err
}

/// END ESI Requests ///

func (x *Pilot) esiUpdate() error {
	if err := x.getBasicInfo(); err != nil {
		return err
	}
	if err := x.getPortrait(); err != nil {
		return err
	}
	if err := x.getMiningRecords(); err != nil {
		return err
	}
	if err := x.getFleet(); err != nil {
		return err
	}
	if err := x.getTransactions(); err != nil {
		return err
	}
	if err := x.getMarketOrders(); err != nil {
		return err
	}
	if err := x.getAttributes(); err != nil {
		return err
	}
	if err := x.getSkills(); err != nil {
		return err
	}
	if err := x.getStandings(); err != nil {
		return err
	}
	if err := x.getLocation(); err != nil {
		return err
	}
	if err := x.getOnlineStatus(); err != nil {
		return err
	}
	if err := x.getShip(); err != nil {
		return err
	}
	return nil
}

func (x *Pilot) Save() error {

	if err := x.esiUpdate(); err != nil {
		return err
	}

	filter := bson.M{
		"info": bson.M{
			"name": x.Name,
		}}

	var tmp *Pilot
	err := x.FindOne(x.ctx, filter).Decode(&tmp)
	if err == nil {
		x.ID = tmp.ID
		_, err := x.UpdateByID(x.ctx, x.ID, x)
		return err
	} else {
		if err == mongo.ErrNoDocuments {
			res, err := x.InsertOne(x.ctx, x)
			if err != nil {
				return err
			}
			x.ID = res.InsertedID.(primitive.ObjectID)
		} else {
			return err
		}
	}

	return nil
}

type Pilot struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

	PilotID   string `json:"CharacterID" bson:"id"`
	*ESIToken `json:"-" bson:"esiToken"`

	BasicInfo *BasicInfo `json:"info" bson:"info"`
	Portrait  *Portrait  `json:"portrait" bson:"portrait"`

	// How they get paid
	MiningRecords []*MiningRecord `json:"mining_records" bson:"mining_records"`
	Fleet         *Fleet          `json:"fleet" bson:"-"`
	Transactions  []*Transaction  `json:"transactions" bson:"transactions"`
	MarketOrders  []*MarketOrder  `json:"market_orders" bson:"market_orders"`
	// How we keep them honest

	Attributes *Attributes `json:"attributes" bson:"attributes"`
	Skills     *Skills     `json:"skills" bson:"skills"`
	Standings  []*Standing `json:"standings" bson:"standings"`
	Location   *Location   `json:"location" bson:"location"`
	Online     *Online     `json:"online" bson:"online"`
	Ship       *Ship       `json:"ship" bson:"ship"`

	*mongo.Collection `json:"-" bson:"-"`
	ctx               context.Context `json:"-" bson:"-"`

	CreatedAt time.Time `json:"_createdAt,omitempty" bson:"_createdAt,omitempty"`
	UpdatedAt time.Time `json:"_updatedAt,omitempty" bson:"_updatedAt,omitempty"`
}

func (x *Pilot) Delete() error {
	_, err := x.DeleteOne(x.ctx, bson.M{"_id": x.ID})
	return err
}
