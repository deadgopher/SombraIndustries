package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type pilotData struct {
	id        string
	refresh   string `json:"refresh_token"`
	createdAt time.Time
	updatedAt time.Time
}

type Pilot struct {
	token ESIToken
	info  BasicInfo
	pic   Portrait
	data  *pilotData
}

func (x *Pilot) ID() string {
	return x.data.id
}

func (x Pilot) Create(i interface{}) (*Pilot, error) {

	if token, ok := i.([]byte); ok {
		if err := json.Unmarshal(token, &x.token); err != nil {
			return nil, err
		}
		charID, err := verifyPilot(x.token.AccessToken)
		if err != nil {
			return nil, err
		}
		x.data = &pilotData{
			id:      string(charID),
			refresh: x.token.Refresh,
		}
		if exist := x.alreadyInDB(); !exist {
			x.data.createdAt = time.Now()
			x.data.updatedAt = time.Now()
			if err := x.data.save(); err != nil {
				return nil, err
			}
		}
		if err := x.getESIData(); err != nil {
			return nil, err
		}

	}
	if id, ok := i.(string); ok {
		rows, err := eveHQ.Query(`
		SELECT *
		FROM pilots
		WHERE id = %v
		`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		if err := rows.Scan(&x.data); err != nil {
			return nil, err
		}
		if err := x.getESIData(); err != nil {
			return nil, err
		}

	}

	return &x, nil
}

func (x Pilot) Read() ([]*Pilot, error) {
	var pilots []*Pilot
	rows, err := eveHQ.Query(`
	SELECT *
	FROM pilots
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p *Pilot
		if err := rows.Scan(&x.data); err != nil {
			return nil, err
		}
		pilots = append(pilots, p)
	}
	return pilots, nil
}

func (x *Pilot) Destroy() error {
	_, err := eveHQ.Exec(`
	DELETE * 
	FROM pilots
	WHERE id = %v
	`, x.data.id)
	return err
}

func (x Pilot) Purge() error {
	_, err := eveHQ.Exec(`
	DELETE * 
	FROM pilots
	`)
	return err
}

func (x *Pilot) alreadyInDB() bool {

	var others int

	rows, _ := eveHQ.Query(`
	SELECT COUNT(*) 
	FROM pilots 
	WHERE pilot_id = %v
	`, x.data.id)
	defer rows.Close()
	rows.Scan(&others)
	return others != 0
}

func (x *pilotData) save() error {
	_, err := eveHQ.Query(`
	INSERT INTO pilots
	VALUES(%v,%v,%v,%v)
	`, x)
	return err
}

// ESI Requests ///
func (x *Pilot) esiGET(path string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+x.token.AccessToken)

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
	res, err := x.esiGET(fmt.Sprintf(ROOT+CHARACTERS+"/%v/", x.data.id))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res, &x.info); err != nil {
		return err
	}
	return nil
}
func (x *Pilot) getPortrait() error {
	res, err := x.esiGET(fmt.Sprintf(ROOT+CHARACTERS+"%v/%v/", x.data.id, "portrait"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, &x.pic)
	return err
}

// End ESI Requests ------^

// func (x *Pilot) getSkills() error {
// 	res, err := x.esiGET(fmt.Sprintf(ROOT+CHARACTERS+"%v/%v/", x.PilotID, "skills"))
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(res, &x.Skills)
// 	return err
// }

/// END ESI Requests ///

func (x *Pilot) getESIData() error {

	fmt.Println("getting basic info")
	if err := x.getBasicInfo(); err != nil {
		return err
	}

	fmt.Println("getting portrait")
	if err := x.getPortrait(); err != nil {
		return err
	}

	return nil
}

func verifyPilot(tok string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://login.eveonline.com/oauth/verify", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
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
