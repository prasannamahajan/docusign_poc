package main

import (
	"context"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/legacy"
	"github.com/jfcote87/esign/v2/envelopes"
	"github.com/jfcote87/esign/v2/model"
	"github.com/sanity-io/litter"

	//"github.com/sanity-io/litter"
	"io/ioutil"
	"log"
	"os"
)

func createEnvelope(cred esign.Credential) (string){
	sv := envelopes.New(cred)
	f1, err := ioutil.ReadFile("letter.pdf")
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
					Tabs: &model.Tabs{
						SignHereTabs: []model.SignHere{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "signature",
									XPosition:  "100",
									YPosition:  "100",
								},
							},
						},
					},
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

