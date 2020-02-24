package main

import (
	"context"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/legacy"
	"github.com/jfcote87/esign/v2/envelopes"
	"github.com/jfcote87/esign/v2/model"
	"github.com/sanity-io/litter"
	"strconv"

	//"github.com/sanity-io/litter"
	"io/ioutil"
	"log"
	"os"
)

type Field struct {
	ID string
	Name string
	AnchorString string
	FieldType string
	IsOptional bool
	IsEditable bool
	Width        int32
	Height       int32
	OffsetX      int64
	OffsetY int64
}

var CclaFields =[]*Field {
	{
		ID:           "sign",
		Name:         "Please Sign",
		AnchorString: "Please sign:",
		FieldType:    "sign",
		IsOptional:   false,
		IsEditable:   false,
		Width:        0,
		Height:       0,
		OffsetX:      100,
		OffsetY:      -6,
	},
	{
		ID:           "date",
		Name:         "Date",
		AnchorString: "Date:",
		FieldType:    "date",
		IsOptional:   false,
		IsEditable:   false,
		Width:        0,
		Height:       0,
		OffsetX:      40,
		OffsetY:      -7,
	},
	{
		ID:           "signatory_name",
		Name:         "Signatory Name",
		AnchorString: "Signatory Name:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   false,
		Width:        355,
		Height:       20,
		OffsetX:      120,
		OffsetY:      -5,
	},
	{
		ID:           "signatory_email",
		Name:         "Signatory E-mail",
		AnchorString: "Signatory E-mail:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   false,
		Width:        355,
		Height:       20,
		OffsetX:      120,
		OffsetY:      -5,
	},
	{
		ID:           "signatory_title",
		Name:         "Signatory Title",
		AnchorString: "Signatory Title:",
		FieldType:    "text",
		IsOptional:   true,
		IsEditable:   true,
		Width:        355,
		Height:       20,
		OffsetX:      120,
		OffsetY:      -6,
	},
	{
		ID:           "corporation_name",
		Name:         "Corporation Name",
		AnchorString: "Corporation Name:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   false,
		Width:        355,
		Height:       20,
		OffsetX:      130,
		OffsetY:      -5,
	},
	{
		ID:           "corporation_address1",
		Name:         "Corporation Address1",
		AnchorString: "Corporation Address:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   true,
		Width:        230,
		Height:       20,
		OffsetX:      135,
		OffsetY:      -8,
	},
	{
		ID:           "corporation_address2",
		Name:         "Corporation Address2",
		AnchorString: "Corporation Address:",
		FieldType:    "text_unlocked",
		IsOptional:   false,
		IsEditable:   true,
		Width:        350,
		Height:       20,
		OffsetX:      0,
		OffsetY:      27,
	},
	{
		ID:           "corporation_address3",
		Name:         "Corporation Address3",
		AnchorString: "Corporation Address:",
		FieldType:    "text_unlocked",
		IsOptional:   true,
		IsEditable:   true,
		Width:        350,
		Height:       20,
		OffsetX:      0,
		OffsetY:      65,
	},
	{
		ID:           "cla_manager_name",
		Name:         "Initial CLA Manager Name",
		AnchorString: "Initial CLA Manager Name:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   false,
		Width:        385,
		Height:       20,
		OffsetX:      190,
		OffsetY:      -7,
	},
	{
		ID:           "cla_manager_email",
		Name:         "Initial CLA Manager Email",
		AnchorString: "Initial CLA Manager E-Mail:",
		FieldType:    "text",
		IsOptional:   false,
		IsEditable:   false,
		Width:        385,
		Height:       20,
		OffsetX:      190,
		OffsetY:      -7,
	},
}


func createEnvelope(cred esign.Credential) (string){
	sv := envelopes.New(cred)
	f1, err := ioutil.ReadFile("ccla.pdf")
	if err != nil {
		log.Fatalf("file read error: %v", err)
	}

	env := &model.EnvelopeDefinition{
		EventNotification:&model.EventNotification{
			EnvelopeEvents:                    []model.EnvelopeEvent{
				// Draft, Sent, Delivered, Completed, Declined, or Voided.
				{
					EnvelopeEventStatusCode: "Completed",
				},
				{
					EnvelopeEventStatusCode: "Declined",
				},
				/*
				{
					EnvelopeEventStatusCode: "Sent",
				},
				{
					EnvelopeEventStatusCode: "Delivered",
				},
				{
					EnvelopeEventStatusCode: "Voided",
				},
				 */
			},
			RecipientEvents: []model.RecipientEvent{
				// Send a webhook notification for the following recipient statuses: Sent, Delivered, Completed, Declined, AuthenticationFailed, and AutoResponded.
				{
					RecipientEventStatusCode : "Completed",
				},
				/*
				{
					RecipientEventStatusCode : "Sent",
				},
				{
					RecipientEventStatusCode : "Delivered",
				},
				{
					RecipientEventStatusCode : "Declined",
				},
				{
					RecipientEventStatusCode : "AuthenticationFailed",
				},
				{
					RecipientEventStatusCode : "AutoResponded",
				},
				*/
			},
			URL:                               "https://webhook.site/57da7464-a4c1-4105-874d-09c5bb565823",
		},
		EmailSubject: "[Go eSignagure SDK] - Please sign this doc",
		EmailBlurb:   "Please sign this test document",
		Status:       "sent",
		Documents: []model.Document{
			{
				DocumentBase64: f1,
				Name:           "invite letter.pdf",
				DocumentID:     "1",
			},
		},
		Recipients: &model.Recipients{
			Signers: []model.Signer{
				{
					Email:             "contributor@gmail.com",
					EmailNotification: nil,
					Name:              "J F Cote",
					ClientUserID: "1",
					RecipientID:       "1",
					RoutingOrder:"1",
					Tabs: getTabs(CclaFields),
				},
				{
					Email:             "clamanager@gmail.com",
					EmailNotification: nil,
					Name:              "Prasanna Mahajan",
					ClientUserID: "2",
					RecipientID:       "2",
					RoutingOrder: "2",
					Tabs: &model.Tabs{
						SignHereTabs: []model.SignHere{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "2",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "signature",
									XPosition:  "200",
									YPosition:  "100",
								},
							},
						},
					},
				},
			},
		},
	}
	envSummary, err := sv.Create(env).Do(context.TODO())
	if err != nil {
		log.Fatalf("create envelope error: %v", err)
	}
	log.Printf("New envelope id: %s", envSummary.EnvelopeID)
	return envSummary.EnvelopeID
}

func voidEnvelope(cfg esign.Credential, envelopeID string) {
	summ, err := envelopes.New(cfg).Update(envelopeID, &model.Envelope{
		Status: "voided",
		VoidedReason: "my void reason",
	}).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	litter.Dump(summ)
}

func getCred() esign.Credential {
	cfg := legacy.Config{
		AccountID:     os.Getenv("DOCUSIGN_ACCOUNT_ID"),
		IntegratorKey: os.Getenv("DOCUSIGN_INTEGRATOR_KEY"),
		UserName:      os.Getenv("DOCUSIGN_USERNAME"),
		Password:      os.Getenv("DOCUSIGN_PASSWORD"),
		//		Host:          os.Getenv("DOCUSIGN_ROOT_URL"),
		Host: "demo.docusign.net",
	}
	return cfg
}

func getRecipients(cred esign.Credential, envelopeID string) *model.Recipients{
	rec, err := envelopes.New(cred).RecipientsList(envelopeID).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	litter.Dump(rec)
	return rec

}

var helptext string = `
	required argument:
	syntax:
	program <cmd> <envelopeID>
	cmds = create/void/genurl
	create doesnt need envelopeID
`

func getCommand(args []string) string {
	if len(args) == 0 {
		log.Fatal(helptext)
	}
	return args[0]
}

func getEnvelopeID(args []string) string {
	if len(args) < 2 {
		log.Fatal(helptext)
	}
	return args[1]
}

func main() {
	args := os.Args[1:]
	var envelopeID string
	if len(args) == 0 {
		log.Fatal(helptext)
	}
	cmd := getCommand(args)
	cfg := getCred()
	switch cmd {
	case "create":
		envelopeID = createEnvelope(cfg)
		rec := getRecipients(cfg, envelopeID)
		signer := rec.Signers[0]
		createViewUrl(cfg, envelopeID, signer)
		signer = rec.Signers[1]
		createViewUrl(cfg, envelopeID, signer)
	case "void":
		envelopeID = getEnvelopeID(args)
		voidEnvelope(cfg, envelopeID)
	case "genurl":
		envelopeID = getEnvelopeID(args)
		rec := getRecipients(cfg, envelopeID)
		signer := rec.Signers[0]
		createViewUrl(cfg, envelopeID, signer)
		signer = rec.Signers[1]
		createViewUrl(cfg, envelopeID, signer)
	default:
		log.Fatal(helptext)
	}
}

func createViewUrl(cfg esign.Credential, envelopeID string, signer model.Signer) {
	viewUrl, err := envelopes.New(cfg).ViewsCreateRecipient(envelopeID,&model.RecipientViewRequest{
		AssertionID:               "",
		AuthenticationInstant:     "",
		AuthenticationMethod:      "None",
		ClientUserID:              signer.ClientUserID,
		Email:             		   signer.Email,
		PingFrequency:             "",
		PingURL:                   "",
		RecipientID:               "",
		ReturnURL:                 "https://localhost:3000/",
		SecurityDomain:            "",
		UserID:                    signer.UserID,
		UserName:                  signer.Name,
		XFrameOptions:             "",
		XFrameOptionsAllowFromURL: "",
	}).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	litter.Dump(viewUrl)
}

func getTabs(tabs []*Field) *model.Tabs {
	pageNumber := 3
	defaults := map[string]string{
		"date" : "2/24/2020",
		"signatory_name": "Prasanna Mahajan",
		"signatory_email" : "prasannak@proximabiz.com",
		"corporation_name" : "Fstack",
		"cla_manager_name": "Prasanna Mahaja",
		"cla_manager_email": "prasannak@proximabiz.com",

	}
	var result model.Tabs
	for _,tab := range tabs {
		tab.OffsetX = tab.OffsetX + 200
		var tabValue model.TabValue
		if v,ok := defaults[tab.ID]; ok {
			tabValue = model.TabValue{Value:v}
		}
		var isOptional model.TabRequired
		if tab.IsOptional {
			isOptional = model.REQUIRED_FALSE
		}
		switch  tab.FieldType {
		case "sign":
			result.SignHereTabs = append(result.SignHereTabs, model.SignHere{
				TabBase:           model.TabBase{
					DocumentID: "1",
					RecipientID: "1",
				},
				TabPosition:       model.TabPosition{
					AnchorString: tab.AnchorString,
					CustomTabID: tab.ID,
					TabLabel: tab.ID,
					XPosition:    strconv.FormatInt(tab.OffsetX,10),
					YPosition:    strconv.FormatInt(tab.OffsetY, 10),
					PageNumber: strconv.Itoa(pageNumber),
				},
				Name:              tab.Name,
			})
		case "text","text_unlocked","text_optional":
			result.TextTabs = append(result.TextTabs, model.Text{
				TabBase:           model.TabBase{
					DocumentID: "1",
					RecipientID: "1",
				},
				TabPosition:                  model.TabPosition{
					AnchorString: tab.AnchorString,
					CustomTabID: tab.ID,
					TabLabel: tab.ID,
					XPosition:    strconv.FormatInt(tab.OffsetX,10),
					YPosition:    strconv.FormatInt(tab.OffsetY, 10),
					PageNumber: strconv.Itoa(pageNumber),
				},
				TabValue:                     tabValue,
				Height:                       tab.Height,
				Width:                        tab.Width,
				Locked:                       model.DSBool(!tab.IsEditable),
				Required:                     isOptional,
			})
		case "date":
			result.DateTabs = append(result.DateTabs, model.Date{
				TabBase:           model.TabBase{
					DocumentID: "1",
					RecipientID: "1",
				},
				TabPosition:                  model.TabPosition{
					AnchorString: tab.AnchorString,
					CustomTabID: tab.ID,
					TabLabel: tab.ID,
					XPosition:    strconv.FormatInt(tab.OffsetX,10),
					YPosition:    strconv.FormatInt(tab.OffsetY, 10),
					PageNumber: strconv.Itoa(pageNumber),
				},
				Width:                        tab.Width,
				Locked:                       model.DSBool(!tab.IsEditable),
				Required:                     isOptional,
				TabValue: tabValue,
			})
		}
	}
	return &result
}