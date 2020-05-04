package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/meyskens/go-hcaptcha"
)

const itfWelcome = "687588438886842373"

var hc *hcaptcha.HCaptcha

func serve() {
	hc = hcaptcha.New(c.HCaptchaSiteSecret)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./www"))))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/index.html", handleHome)
	http.HandleFunc("/invite", handleInvite)
	if err := http.ListenAndServe(c.BindAddr, nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}

type homeTemplate struct {
	HCaptchaSiteKey string
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./www/index.html.tpl")
	if err != nil {
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, homeTemplate{
		HCaptchaSiteKey: c.HCaptchaSiteKey,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func handleInvite(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hcaptchaResponse, responseFound := r.Form["h-captcha-response"]
	if !responseFound || len(hcaptchaResponse) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !verifyCaptcha(ip, hcaptchaResponse[0]) {
		// todo: add error page
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i, err := dg.ChannelInviteCreate(itfWelcome, discordgo.Invite{
		MaxUses: 1,
		MaxAge:  60 * 60, // 1 hour
		Unique:  true,
	})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Invited user with code %q from IP %s", i.Code, ip)
	http.Redirect(w, r, "https://discord.gg/"+i.Code, http.StatusSeeOther)
}

func verifyCaptcha(ip, cResponse string) bool {
	resp, err := hc.Verify(cResponse, ip)
	if err != nil {
		log.Println(err)
		return false
	}
	return resp.Success
}
