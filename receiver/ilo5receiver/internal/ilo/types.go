package ilo

import "time"

// RedfishResponse is a generic wrapper for Redfish responses
type RedfishResponse[T any] struct {
	Members []T `json:"Members"`
	Count   int `json:"Members@odata.count"`
}

// System represents /redfish/v1/Systems/{id}
type System struct {
	ID           string `json:"Id"`
	HostName     string `json:"HostName"`
	SerialNumber string `json:"SerialNumber"`
	PowerState   string `json:"PowerState"` // On, Off
	Status       Status `json:"Status"`
}

// Status represents the health status common object
type Status struct {
	Health string `json:"Health"` // OK, Warning, Critical
	State  string `json:"State"`
}

// Chassis represents /redfish/v1/Chassis/{id}
type Chassis struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// Power represents /redfish/v1/Chassis/{id}/Power
type Power struct {
	PowerControl  []PowerControl `json:"PowerControl"`
	PowerSupplies []PowerSupply  `json:"PowerSupplies"`
	Voltages      []Voltage      `json:"Voltages"`
}

type PowerControl struct {
	PowerConsumedWatts float64 `json:"PowerConsumedWatts"`
	PowerCapacityWatts float64 `json:"PowerCapacityWatts"`
}

type PowerSupply struct {
	MemberID             string  `json:"MemberId"`
	Name                 string  `json:"Name"`
	LastPowerOutputWatts float64 `json:"LastPowerOutputWatts"`
	LineInputVoltage     float64 `json:"LineInputVoltage"`
	PowerCapacityWatts   float64 `json:"PowerCapacityWatts"`
	Status               Status  `json:"Status"`
}

type Voltage struct {
	Name            string  `json:"Name"`
	MemberID        string  `json:"MemberId"`
	ReadingVolts    float64 `json:"ReadingVolts"`
	PhysicalContext string  `json:"PhysicalContext"`
}

// Thermal represents /redfish/v1/Chassis/{id}/Thermal
type Thermal struct {
	Temperatures []Temperature `json:"Temperatures"`
	Fans         []Fan         `json:"Fans"`
}

type Temperature struct {
	Name            string  `json:"Name"`
	MemberID        string  `json:"MemberId"`
	ReadingCelsius  float64 `json:"ReadingCelsius"`
	PhysicalContext string  `json:"PhysicalContext"`
	SensorNumber    int     `json:"SensorNumber"`
	Oem             struct {
		Hpe struct {
			LocationXmm int `json:"LocationXmm"`
			LocationYmm int `json:"LocationYmm"`
		} `json:"Hpe"`
	} `json:"Oem"`
	UpperThresholdCritical *float64 `json:"UpperThresholdCritical"`
	UpperThresholdFatal    *float64 `json:"UpperThresholdFatal"`
	Status                 Status   `json:"Status"`
}

type Fan struct {
	Name         string  `json:"Name"`
	MemberID     string  `json:"MemberId"`
	Reading      float64 `json:"Reading"`      // Percentage usually, sometimes RPM depending on units
	ReadingUnits string  `json:"ReadingUnits"` // Percent or RPM
}

// StorageCollection for /redfish/v1/Systems/{id}/Storage
type StorageCollection struct {
	Members []struct {
		ID string `json:"@odata.id"`
	} `json:"Members"`
}

// Storage represents /redfish/v1/Systems/{id}/Storage/{storageId}
type Storage struct {
	ID     string `json:"Id"`
	Drives []struct {
		ID string `json:"@odata.id"`
	} `json:"Drives"`
}

// Drive represents /redfish/v1/Systems/{id}/Storage/{storageId}/Drives/{driveId} or linked Drive
type Drive struct {
	ID     string `json:"Id"`
	Status Status `json:"Status"`
}

// ClientConfig holds configuration for the client
type ClientConfig struct {
	Endpoint           string
	Username           string
	Password           string
	InsecureSkipVerify bool
	Timeout            time.Duration
}
