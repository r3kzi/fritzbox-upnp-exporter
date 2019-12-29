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

func newUPnPClient(url, username, password string) *UPnPClient {
	return &UPnPClient{
		URL:      fmt.Sprintf("https://%s/tr064", url),
		Username: username,
		Password: password,
	}
}

func (uc *UPnPClient) execute() map[string]string {
	dr := newRequest(uc.Username, uc.Password, "GET", uc.URL+TR064ENDPOINT, "")
	decoder := xml.NewDecoder(do(dr))

	actions := make([]Action, 0)
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
									action.Service = service
									actions = append(actions, action)
								}
							}
						}
					}
				}
			}
		}
	}

	result := make(map[string]string)
	for _, action := range actions {

		message := fmt.Sprintf(`
		<?xml version="1.0"?> 
        <s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" 
				s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"> 
            <s:Body><u:%s xmlns:u='%s'/></s:Body>
        </s:Envelope>`, action.Name, action.Service.ServiceType)

		dr := newRequest(uc.Username, uc.Password, "POST", uc.URL+action.Service.ControlURL, message)
		dr.Header.Add("Content-Type", "text/xml; charset='utf-8'")
		dr.Header.Add("SoapAction", fmt.Sprintf("%s#%s", action.Service.ServiceType, action.Name))

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
	return result
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
		//"urn:WANPPPConnection-com:serviceId:WANPPPConnection1",
	}

	for _, entry := range whitelistedServices {
		if entry == service.ServiceId {
			return true
		}
	}
	return false
}
