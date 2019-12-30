package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

const TR064ENDPOINT string = "/tr64desc.xml"

type UPnPClient struct {
	URL      string
	Username string
	Password string
}

func NewUPnPClient(url, username, password string) *UPnPClient {
	return &UPnPClient{
		URL:      fmt.Sprintf("https://%s/tr064", url),
		Username: username,
		Password: password,
	}
}

func (uc *UPnPClient) Execute() map[string]string {
	result := make(map[string]string)
	for _, service := range uc.parseServices() {
		for _, action := range uc.parseActions(service) {
			message := fmt.Sprintf(`
		<?xml version="1.0"?> 
        <s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" 
				s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"> 
            <s:Body><u:%s xmlns:u='%s'/></s:Body>
        </s:Envelope>`, action.Name, service.ServiceType)

			dr := newRequest(uc.Username, uc.Password, "POST", uc.URL+service.ControlURL, message)
			dr.Header.Add("Content-Type", "text/xml; charset='utf-8'")
			dr.Header.Add("SoapAction", fmt.Sprintf("%s#%s", service.ServiceType, action.Name))

			decoder := xml.NewDecoder(do(dr))
			for {
				t, _ := decoder.Token()
				if t == nil {
					break
				}
				switch se := t.(type) {
				case xml.StartElement:
					for _, argument := range action.Arguments {
						if se.Name.Local == argument.Name {
							t, _ = decoder.Token()
							switch element := t.(type) {
							case xml.CharData:
								result[argument.RelatedStateVariable] = string(element)
							}
						}
					}
				}
			}
		}
	}
	return result
}

func (uc *UPnPClient) parseServices() []Service {
	services := make([]Service, 0)

	dr := newRequest(uc.Username, uc.Password, "GET", uc.URL+TR064ENDPOINT, "")
	decoder := xml.NewDecoder(do(dr))
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "service" {
				var service Service
				if err := decoder.DecodeElement(&service, &se); err != nil {
					panic(err)
				}
				if IsServiceWhitelisted(service) {
					service.Actions = uc.parseActions(service)
					services = append(services, service)
				}
			}
		}
	}
	return services
}

func (uc *UPnPClient) parseActions(service Service) []Action {
	actions := make([]Action, 0)

	dr := newRequest(uc.Username, uc.Password, "GET", uc.URL+service.SCPDURL, "")
	decoder := xml.NewDecoder(do(dr))
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "action" {
				var action Action
				if err := decoder.DecodeElement(&action, &se); err != nil {
					panic(err)
				}
				if IsActionGetOnly(action) {
					actions = append(actions, action)
				}
			}
		}
	}
	return actions
}

func IsActionGetOnly(action Action) bool {
	match, _ := regexp.MatchString("^(Get)+[A-z]*", action.Name)
	if !match {
		return false
	}
	for _, a := range action.Arguments {
		if a.Direction == "in" {
			return false
		}
	}
	return len(action.Arguments) > 0
}

func IsServiceWhitelisted(service Service) bool {
	var whitelistedServices = []string{
		"urn:WANCIfConfig-com:serviceId:WANCommonInterfaceConfig1",
	}

	for _, entry := range whitelistedServices {
		if entry == service.ServiceId {
			return true
		}
	}
	return false
}
