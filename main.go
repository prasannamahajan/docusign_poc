package main

import (
	"context"
	"fmt"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/legacy"
	"github.com/jfcote87/esign/v2/envelopes"
	"github.com/jfcote87/esign/v2/model"
	"github.com/sanity-io/litter"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

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
									AnchorString: "Please sign:",
									TabLabel:   "signature",
									AnchorXOffset:  "200",
									AnchorYOffset:  "-6",
									AnchorIgnoreIfNotPresent: false,
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
	//litter.Dump(rec)
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

func CreateEnvelope(w http.ResponseWriter, r *http.Request) {
	cfg := getCred()
	envelopeID := createEnvelope(cfg)
	rec := getRecipients(cfg, envelopeID)
	signer := rec.Signers[0]
	contributorReturnUrl := "http://localhost:8000/contributor_signed/" + envelopeID
	url := createViewUrl(cfg, envelopeID, signer,contributorReturnUrl)
	createdEnvelopeTemplate := `
	<html>
	<head>
	</head>
	<body>
	<h3>
		Hey {{.Name}}
	</h3>
	<h3>
		To contribute to project 
		 <a href="{{.URL}}">Sign</a> the contributor license agreement
	</h3>
	</body>
	</html>`
	t, err := template.New("test").Parse(createdEnvelopeTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, struct {
		Name string
		URL string
	}{
		Name : signer.Name,
		URL : url,
	})
}

func VoidEnvelope(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	envelopeID := args["envelopeID"]
	cfg := getCred()
	voidEnvelope(cfg, envelopeID)
}

func GenerateURLs(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	envelopeID := args["envelopeID"]
	contributorReturnUrl := "localhost:8000/contributor_signed/" + envelopeID
	clamanagerReturnUrl := "localhost:8000/cla_manager_signed/" + envelopeID
	cfg := getCred()
	rec := getRecipients(cfg, envelopeID)
	signer := rec.Signers[0]
	contribUrl := createViewUrl(cfg, envelopeID, signer,contributorReturnUrl)
	signer = rec.Signers[1]
	claManagerUrl := createViewUrl(cfg, envelopeID, signer,clamanagerReturnUrl)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "contributor signing url is %s\ncla manager signing url is %s\n", contribUrl, claManagerUrl)
}

func ContributorSigned(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	envelopeID := args["envelopeID"]
	event := r.FormValue("event")
	if event == "" {
		fmt.Fprintf(w, "got empty event. should not happen\n")
		return
	}
	if event != "signing_complete" {
		fmt.Fprintf(w, "received event : %s\n", event)
	}
	log.Printf("got response %s for envelopeID %s\n", event,args["envelopeID"])
	cfg := getCred()
	rec := getRecipients(cfg, envelopeID)
	signer := rec.Signers[1]
	clamanagerReturnUrl := "http://localhost:8000/cla_manager_signed/" + envelopeID
	claManagerUrl := createViewUrl(cfg, envelopeID, signer,clamanagerReturnUrl)
	contributorSignedTemplate := `
	<html>
	<head>
	</head>
	<body>
	<h3>
		Thank you for your submission. It has been sent to the {{.Name}} for countersignature. Once they have signed, you will be able to whitelist contributors to the project.
	</h3>
	<h3>
		cla manager signing url is <a href="{{.URL}}">Sign</a>
	</h3>
	</body>
	</html>`
	t, err := template.New("test").Parse(contributorSignedTemplate)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, struct {
		Name string
		URL string
	}{
		Name : signer.Name,
		URL : claManagerUrl,
	})
}

func ClaManagerSigned(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	event := r.FormValue("event")
	if event == "" {
		fmt.Fprintf(w, "got empty event. should not happen\n")
		return
	}
	if event != "signing_complete" {
		fmt.Fprintf(w, "received event : %s\n", event)
		return
	}
	log.Printf("recieved event %s", event)
	claManagerSignedTemplate := `
	<html>
	<head>
	</head>
	<body>
	<h3>
		Thank you for your submission. Signing process is complete.
	</h3>
	<h3>
		<a href="{{.URL}}">View signed document</a>
	</h3>
	</body>
	</html>`
	t, err := template.New("test").Parse(claManagerSignedTemplate)
	if err != nil {
		log.Fatal(err)
	}
	link := "http://localhost:8000/view/"+args["envelopeID"]
	t.Execute(w, struct {
		URL string
	} {
		URL : link,
	})
}

func ViewDocument(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	download, err :=	envelopes.New(getCred()).DocumentsGet("1",args["envelopeID"]).Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", download.ContentType)
	io.Copy(w, download.ReadCloser)
}

func test(w http.ResponseWriter, r *http.Request) {
	download, err :=	envelopes.New(getCred()).DocumentsGet("1","6aed5d3b-0818-4711-9a37-9251aaf83f99").Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", download.ContentType)
	io.Copy(w, download.ReadCloser)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/create",CreateEnvelope)
	router.HandleFunc("/test",test)
	router.HandleFunc("/view/{envelopeID}",ViewDocument)
	router.HandleFunc("/void/{envelopeID}", VoidEnvelope)
	router.HandleFunc("/genurl/{envelopeID}", GenerateURLs)
	router.HandleFunc("/contributor_signed/{envelopeID}", ContributorSigned)
	router.HandleFunc("/cla_manager_signed/{envelopeID}", ClaManagerSigned)
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("stating server on ",srv.Addr)
	log.Fatal(srv.ListenAndServe())
	/*
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
	 */
}

func createViewUrl(cfg esign.Credential, envelopeID string, signer model.Signer,returnURL string) string {
	viewUrl, err := envelopes.New(cfg).ViewsCreateRecipient(envelopeID,&model.RecipientViewRequest{
		AssertionID:               "",
		AuthenticationInstant:     "",
		AuthenticationMethod:      "None",
		ClientUserID:              signer.ClientUserID,
		Email:             		   signer.Email,
		PingFrequency:             "",
		PingURL:                   "",
		RecipientID:               "",
		ReturnURL:                 returnURL,
		SecurityDomain:            "",
		UserID:                    signer.UserID,
		UserName:                  signer.Name,
		XFrameOptions:             "",
		XFrameOptionsAllowFromURL: "",
	}).Do(context.TODO())
	if err != nil {
		return err.Error()
	}
	return viewUrl.URL
}

func getTabs(tabs []*Field) *model.Tabs {
	defaults := map[string]string{
		"date" : "2/24/2020",
		"signatory_name": "J F Kote",
		"signatory_email" : "jfkote@google.com",
		"corporation_name" : "Google",
		"cla_manager_name": "Prasanna Mahajan",
		"cla_manager_email": "prasannak@proximabiz.com",

	}
	var result model.Tabs
	for _,tab := range tabs {
		var tabValue model.TabValue
		if v,ok := defaults[tab.ID]; ok {
			tabValue = model.TabValue{Value:v}
		}
		var required model.TabRequired
		if tab.IsOptional {
			required = model.REQUIRED_FALSE
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
					AnchorXOffset: strconv.FormatInt(tab.OffsetX,10),
					AnchorYOffset:    strconv.FormatInt(tab.OffsetY, 10),
					AnchorIgnoreIfNotPresent:false,
					CustomTabID: tab.ID,
					TabLabel: tab.ID,
				},
				Name:              tab.Name,
			})
		case "text","text_unlocked","text_optional":
			textTab := model.Text{
				TabBase: model.TabBase{
					DocumentID:  "1",
					RecipientID: "1",
				},
				TabPosition: model.TabPosition{
					AnchorString:             tab.AnchorString,
					CustomTabID:              tab.ID,
					TabLabel:                 tab.ID,
					AnchorXOffset:            strconv.FormatInt(tab.OffsetX, 10),
					AnchorYOffset:            strconv.FormatInt(tab.OffsetY, 10),
					AnchorIgnoreIfNotPresent: false,
				},
				TabValue: tabValue,
				Height:   tab.Height,
				Width:    tab.Width,
				Locked:   model.DSBool(!tab.IsEditable),
				Required: required,
			}
			if tab.FieldType == "text_unlocked" {
				textTab.Locked = false
			}
			if tab.FieldType == "text_optional" {
				textTab.Required = model.REQUIRED_FALSE
			}
			result.TextTabs = append(result.TextTabs, textTab)
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
					AnchorXOffset: strconv.FormatInt(tab.OffsetX,10),
					AnchorYOffset:    strconv.FormatInt(tab.OffsetY, 10),
					AnchorIgnoreIfNotPresent:false,
				},
				Width:                        tab.Width,
				Locked:                       model.DSBool(!tab.IsEditable),
				Required:                     required,
				TabValue: tabValue,
			})
		}
	}
	return &result
}