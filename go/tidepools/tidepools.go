package tidepools

// Define a struct for tide pool locations
type TidePool struct {
	Name       string
	City       string
	State      string
	Lat        float64
	Long       float64
	Station    string
	BackupLat  float64
	BackupLong float64
}

// Initialize a slice of TidePool with your locations data
var Locations = []TidePool{
	{"Point Loma Tide Pools", "San Diego", "CA", 32.6731, -117.2425, "9410170", 32.7157, -117.1611},
	{"Crystal Cove State Park", "Laguna Beach", "CA", 33.5665, -117.8090, "9410580", 33.5427, -117.7854},
	{"Leo Carrillo State Park", "Malibu", "CA", 34.0453, -118.9358, "9410230", 34.0259, -118.7798},
	{"Santa Rosa Island Tide Pools", "Channel Islands National Park", "CA", 33.9950, -120.0805, "9410840", 34.0147, -119.6982},
	{"Cape Perpetua Tide Pools", "Yachats", "OR", 44.2811, -124.1089, "9432780", 44.3118, -124.1037},
	{"Kalaloch Beach Tide Pools", "Forks", "WA", 47.6136, -124.3740, "9437540", 47.7109, -124.4154},
	{"Shi Shi Beach Tide Pools", "Neah Bay", "WA", 48.3687, -124.6252, "9443090", 48.3686, -124.6247},
	{"Ecola State Park Tide Pools", "Cannon Beach", "OR", 45.9273, -123.9788, "9435380", 45.8918, -123.9615},
	{"Cape Kiwanda Tide Pools", "Pacific City", "OR", 45.2100, -123.9680, "9435385", 45.2028, -123.9624},
	{"Second Beach Tide Pools", "La Push", "WA", 47.9023, -124.6356, "9444090", 47.9133, -124.6361},
	// Add additional locations here...
}

// Function to get the station numbers
// func GetAllStationNumbers() []string {
// 	var stations []string
// 	for _, location := range Locations {
// 		stations = append(stations, location.Station)
// 	}
// 	return stations
// }
