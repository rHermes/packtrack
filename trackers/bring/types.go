package bring

import "time"

type APIResponseError struct {
	APIVersion     string `json:"apiVersion"`
	ConsignmentSet []struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"consignmentSet"`
}

type APIResponse struct {
	APIVersion     string           `json:"apiVersion"`
	ConsignmentSet []ConsignmentSet `json:"consignmentSet"`
}

type SenderAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	PostalCode   string `json:"postalCode"`
	City         string `json:"city"`
	CountryCode  string `json:"countryCode"`
	Country      string `json:"country"`
}

type RecipientAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	PostalCode   string `json:"postalCode"`
	City         string `json:"city"`
	CountryCode  string `json:"countryCode"`
	Country      string `json:"country"`
}

type RecipientHandlingAddress struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	PostalCode   string `json:"postalCode"`
	City         string `json:"city"`
	CountryCode  string `json:"countryCode"`
	Country      string `json:"country"`
}

type RecipientSignature struct {
	Name        string      `json:"name"`
	LinkToImage interface{} `json:"linkToImage"`
}

type EventSet struct {
	Description        string             `json:"description"`
	Status             string             `json:"status"`
	LmEventCode        interface{}        `json:"lmEventCode"`
	RecipientSignature RecipientSignature `json:"recipientSignature"`
	UnitID             string             `json:"unitId"`
	UnitInformationURL interface{}        `json:"unitInformationUrl"`
	UnitType           string             `json:"unitType"`
	PostalCode         string             `json:"postalCode"`
	City               string             `json:"city"`
	CountryCode        string             `json:"countryCode"`
	Country            string             `json:"country"`
	DateIso            time.Time          `json:"dateIso"`
	DisplayDate        string             `json:"displayDate"`
	DisplayTime        string             `json:"displayTime"`
	ConsignmentEvent   bool               `json:"consignmentEvent"`
	Insignificant      bool               `json:"insignificant"`
	GpsXCoordinate     string             `json:"gpsXCoordinate"`
	GpsYCoordinate     string             `json:"gpsYCoordinate"`
	GpsMapURL          string             `json:"gpsMapUrl"`
}

type AdditionalServiceSet struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	Amount          string `json:"amount"`
	CurrencyCode    string `json:"currencyCode"`
	LongDescription string `json:"longDescription"`
}

type PackageSet struct {
	StatusDescription        string                   `json:"statusDescription"`
	Descriptions             []interface{}            `json:"descriptions"`
	PackageNumber            string                   `json:"packageNumber"`
	PreviousPackageNumber    string                   `json:"previousPackageNumber"`
	ProductName              string                   `json:"productName"`
	ProductCode              string                   `json:"productCode"`
	ProductLink              string                   `json:"productLink"`
	Brand                    string                   `json:"brand"`
	LengthInCm               int                      `json:"lengthInCm"`
	WidthInCm                int                      `json:"widthInCm"`
	HeightInCm               int                      `json:"heightInCm"`
	VolumeInDm3              float64                  `json:"volumeInDm3"`
	WeightInKgs              float64                  `json:"weightInKgs"`
	ListPrice                interface{}              `json:"listPrice"`
	ContractPrice            interface{}              `json:"contractPrice"`
	CurrencyCode             interface{}              `json:"currencyCode"`
	PickupCode               interface{}              `json:"pickupCode"`
	ShelfNumber              interface{}              `json:"shelfNumber"`
	DateOfReturn             string                   `json:"dateOfReturn"`
	DateOfEstimatedDelivery  interface{}              `json:"dateOfEstimatedDelivery"`
	DateOfDelivery           interface{}              `json:"dateOfDelivery"`
	SenderName               string                   `json:"senderName"`
	SenderAddress            SenderAddress            `json:"senderAddress"`
	SenderHandlingAddress    interface{}              `json:"senderHandlingAddress"`
	RecipientName            interface{}              `json:"recipientName"`
	RecipientAddress         RecipientAddress         `json:"recipientAddress"`
	RecipientHandlingAddress RecipientHandlingAddress `json:"recipientHandlingAddress"`
	EventSet                 []EventSet               `json:"eventSet"`
	AdditionalServiceSet     []AdditionalServiceSet   `json:"additionalServiceSet"`
	RequestedPackage         interface{}              `json:"requestedPackage"`
}

type ConsignmentSet struct {
	ConsignmentID                 string                   `json:"consignmentId"`
	PreviousConsignmentID         string                   `json:"previousConsignmentId"`
	TotalWeightInKgs              float64                  `json:"totalWeightInKgs"`
	TotalVolumeInDm3              float64                  `json:"totalVolumeInDm3"`
	PackageSet                    []PackageSet             `json:"packageSet"`
	RecipientName                 interface{}              `json:"recipientName"`
	RecipientAddress              RecipientAddress         `json:"recipientAddress"`
	RecipientHandlingAddress      RecipientHandlingAddress `json:"recipientHandlingAddress"`
	SenderReference               string                   `json:"senderReference"`
	SenderCustomerNumber          string                   `json:"senderCustomerNumber"`
	SenderCustomerMasterNumber    string                   `json:"senderCustomerMasterNumber"`
	SenderName                    string                   `json:"senderName"`
	SenderAddress                 SenderAddress            `json:"senderAddress"`
	SenderHandlingAddress         interface{}              `json:"senderHandlingAddress"`
	SenderCustomerType            string                   `json:"senderCustomerType"`
	RecipientCustomerNumber       string                   `json:"recipientCustomerNumber"`
	RecipientCustomerMasterNumber string                   `json:"recipientCustomerMasterNumber"`
	RecipientCustomerType         string                   `json:"recipientCustomerType"`
	TotalListPrice                interface{}              `json:"totalListPrice"`
	TotalContractPrice            interface{}              `json:"totalContractPrice"`
	ListPricePackageCount         interface{}              `json:"listPricePackageCount"`
	ContractPricePackageCount     interface{}              `json:"contractPricePackageCount"`
	CurrencyCode                  interface{}              `json:"currencyCode"`
	IsPickupNoticeAvailable       bool                     `json:"isPickupNoticeAvailable"`
	ConsignmentActionSet          interface{}              `json:"consignmentActionSet"`
}
