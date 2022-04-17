package user

type Addresses struct {
	ValidAddress    []Address `json:"valid_address"`
	MaxAddressCount int       `json:"max_address_count"`
}

type Address struct {
	Id          string   `json:"id"`
	Gender      int      `json:"gender"`
	Mobile      string   `json:"mobile"`
	Location    location `json:"location"`
	UserName    string   `json:"user_name"`
	AddrDetail  string   `json:"addr_detail"`
	StationId   string   `json:"station_id"`
	StationName string   `json:"station_name"`
	IsDefault   bool     `json:"is_default"`
	CityNumber  string   `json:"city_number"`
}

type location struct {
	TypeCode string    `json:"typecode"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}
