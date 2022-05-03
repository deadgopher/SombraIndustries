package model

import "time"

const CHARACTERS = "characters/"

type MiningRecord struct {
	Date        time.Time `json:"data" bson:"data"`
	Quantity    int       `json:"quantity" bson:"quantity"`
	SolarSystem int       `json:"solar_system_id" bson:"solarSystemId"`
	Type        int       `json:"type_id" bson:"typeId"`
}

type ESIToken struct {
	AccessToken  string      `json:"access_token"`
	TokenExpires interface{} `json:"expires_in" bson:"expires_in"`
	Refresh      string      `json:"refresh_token" bson:"refresh_token"`
}
type Portrait struct {
	Small  string `json:"px128x128"`
	Medium string `json:"px256x256"`
	Large  string `json:"px512x512"`
	XSmall string `json:"px64x64"`
}
type BasicInfo struct {
	Alliance       int     `json:"alliance_id" bson:"alliance_id"`
	Birthday       string  `json:"birthday" bson:"birthday"`
	Bloodline      int     `json:"bloodline_id" bson:"bloodline_id"`
	CorpID         int     `json:"corporation_id" bson:"corporation_id"`
	Description    string  `json:"description" bson:"description"`
	Faction        int     `json:"faction_id" bson:"faction_id"`
	Gender         string  `json:"gender" bson:"gender"`
	Name           string  `json:"name" bson:"name"`
	Race           int     `json:"race_id" bson:"race_id"`
	SecurityStatus float32 `json:"security_status" bson:"security_status"`
	Title          string  `json:"title" bson:"title"`
}

type Standing struct {
	From     int     `json:"from_id" bson:"from_id"`
	FromType string  `json:"from_type" bson:"from_type"`
	Value    float32 `json:"standing" bson:"standing"`
}

type Fleet struct {
	Fleet int    `json:"fleet_id" bson:"fleet_id"`
	Role  string `json:"role" bson:"role"`
	Squad int    `json:"squad_id" bson:"squad_id"`
	Wing  int    `json:"wing_id" bson:"wing_id"`
}

type Transaction struct {
	ID         int     `json:"transaction_id" bson:"transaction_id"`
	Type       int     `json:"type_id" bson:"type_id"`
	Price      float64 `json:"unit_price" bson:"unit_price"`
	Client     int     `json:"client_id" bson:"client_id"`
	Data       string  `json:"data" bson:"data"`
	IsBuy      bool    `json:"is_buy" bson:"is_buy"`
	IsPersonal bool    `json:"is_personal" bson:"is_personal"`
	JournalRed int     `json:"journal_ref_id" bson:"journal_ref_id"`
	Location   int     `json:"location_id" bson:"location_id"`
	Quantity   int     `json:"quantity" bson:"quantity"`
}

type MarketOrder struct {
	ID           int     `json:"order_id" bson:"order_id"`
	Type         int     `json:"type_id" bson:"type_id"`
	Price        float64 `json:"price" bson:"price"`
	Duration     int     `json:"duration" bson:"duration"`
	Escrow       float64 `json:"escrow" bson:"escrow"`
	IsBuyOrder   bool    `json:"is_buy_order" bson:"is_buy_order"`
	IsCorp       bool    `json:"is_corporation" bson:"is_corporation"`
	Issused      string  `json:"issued" bson:"issued"`
	Location     int     `json:"location_id" bson:"location_id"`
	MinVolume    int     `json:"min_volume"`
	Range        string  `json:"range" bson:"range"`
	Region       int     `json:"region_id" bson:"region_id"`
	VolumeRemain int     `json:"volume_remain" bson:"volume_remain"`
	VolumeTotal  int     `json:"volume_total" bson:"volume_total"`
}

type Location struct {
	SolarSystem int `json:"solar_system_id" bson:"solar_system_id"`
	Station     int `json:"station_id" bson:"station_id"`
	Structure   int `json:"structure_id" bson:"structure_id"`
}
type Online struct {
	LastLogin  string `json:"last_login" bson:"last_login"`
	LastLogout string `json:"last_logout" bson:"last_logout"`
	Logins     int    `json:"logins" bson:"logins"`
	Status     bool   `json:"online" bson:"-"`
}
type Ship struct {
	Item int    `json:"ship_item_id" bson:"ship_item_id"`
	Name string `json:"ship_name" bson:"ship_name"`
	Type int    `json:"ship_type_id" bson:"ship_type_id"`
}

type Attributes struct {
	Charisma     int `json:"charisma" bson:"charisma"`
	Intelligence int `json:"intelligence" bson:"intelligence"`
	Memory       int `json:"memory" bson:"memory"`
	Perception   int `json:"perception" bson:"perception"`
	Willpower    int `json:"willpower" bson:"willpower"`
}

type Skill struct {
	ID           int `json:"skill_id" bson:"skill_id"`
	ActiveLevel  int `json:"active_skill_level" bson:"active_skill_level"`
	TrainedLevel int `json:"trained_skill_level" bson:"trained_skill_level"`
	Points       int `json:"skillpoints_in_skill" bson:"skillpoints_in_skill"`
}
type Skills struct {
	List              []*Skill `json:"list" bson:"list"`
	TotalPoints       int      `json:"total_sp" bson:"total_sp"`
	UnallocatedPoints int      `json:"unallocated_sp" bson:"unallocated_sp"`
}
