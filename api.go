package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func printBalance() {
	c := newClient()
	v := url.Values{}
	b := BalanceResp{}

	c.getRequest("getBalance", v, &b)

	fmt.Printf("%v\n", b.Balance.CurrentBalance)
}

func blockNumber(phone *string, note *string) {
	c := newClient()

	m := postParams{}
	m["callerid"] = *phone
	m["did"] = "all"
	m["routing"] = "sys:hangup"

	var filterNote = time.Now().Format("Jan.02/06")
	if *note != "" {
		filterNote += " - " + *note
	}
	m["note"] = filterNote

	s := SetCallerIDFilterResp{}

	c.postRequest("setCallerIDFiltering", m, &s)

	switch s.Status {
	case "used_filter":
		fmt.Println("Voip.ms reports a filter for this number already exists:", *phone)
		break
	case "success":
		fmt.Printf("Yay, %s blocked successfully with note [%s]: got filter id %s\n", *phone, filterNote, s.Filtering.String())
		break
	default:
		fmt.Println("Unknown result:")
		fmt.Printf("%+v\n", s)
	}
}

func getRecent() GetCallDataRecord {

	c := newClient()
	v := url.Values{}
	r := GetCallDataRecord{}

	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.Add(time.Hour * 24 * -1).Format("2006-01-02")
	_, offset := now.Zone()
	zoneDuration := time.Duration(offset) * time.Second
	v.Add("date_to", today)
	v.Add("date_from", yesterday)
	v.Add("timezone", fmt.Sprintf("%.2g", zoneDuration.Hours()))
	v.Add("answered", "1")
	v.Add("noanswer", "1")
	v.Add("busy", "1")
	v.Add("fail", "1")

	c.getRequest("getCDR", v, &r)

	return r
}

func printRecent(r GetCallDataRecord) {
	switch r.Status {
	case "success":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.AlignRight)
		header := []string{
			"Date",
			"CallerID",
			"Destination",
			"Description",
			"Disposition",
			"Duration",
		}
		fmt.Fprintln(w, strings.Join(header, "\t")+"\t")
		for _, cdr := range r.CallDataRecords {
			fmt.Fprintln(
				w,
				strings.Join([]string{
					cdr.Date,
					cdr.CallerID,
					string(cdr.Destination),
					cdr.Disposition,
					cdr.Description,
					cdr.Duration,
				}, "\t")+"\t")
		}
		w.Flush()
	case "no_cdr":
		fmt.Println("No recent calls found")
	default:
		fmt.Println("Non-success result getting CDR:", r.Status)
	}
}
