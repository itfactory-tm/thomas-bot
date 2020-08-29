package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"

	"github.com/meyskens/go-hcaptcha"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

// TODO: replace me
const itfWelcome = "687588438886842373"

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveHTTPCmdOptions struct {
	Token              string
	HCaptchaSiteKey    string `envconfig:"HCAPTCHA_SITE_KEY"`
	HCaptchaSiteSecret string `envconfig:"HCAPTCHA_SITE_SECRET"`
	BindAddr           string `default:":8080" envconfig:"BIND_ADDR"`

	hc *hcaptcha.HCaptcha
}

// NewServeCmd generates the `serve` command
func NewServeHTTPCmd() *cobra.Command {
	s := serveHTTPCmdOptions{}
	c := &cobra.Command{
		Use:     "serve-http",
		Short:   "Run only the HTTP server",
		Long:    `This is a separate instance got HTTP, included in the serve command`,
		RunE:    s.RunE,
		PreRunE: s.Validate,
	}

	// TODO: switch to viper
	err := envconfig.Process("thomasbot", &s)
	if err != nil {
		log.Fatalf("Error processing envvars: %q\n", err)
	}

	return c
}

func (s *serveHTTPCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if s.Token == "" {
		return errors.New("No token specified")
	}

	return nil
}

func (s *serveHTTPCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	s.hc = hcaptcha.New(s.HCaptchaSiteSecret)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./www"))))

	http.HandleFunc("/", s.handleHome)
	http.HandleFunc("/index.html", s.handleHome)
	http.HandleFunc("/invite", s.handleInvite)
	if err := http.ListenAndServe(s.BindAddr, nil); err != nil {
		return err
	}

	return nil
}

func (s *serveHTTPCmdOptions) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./www/index.html.tpl")
	if err != nil {
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, struct {
		HCaptchaSiteKey string
	}{
		HCaptchaSiteKey: s.HCaptchaSiteKey,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *serveHTTPCmdOptions) handleInvite(w http.ResponseWriter, r *http.Request) {
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

	if !s.verifyCaptcha(ip, hcaptchaResponse[0]) {
		// todo: add error page
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dg, err := discordgo.New("Bot " + s.Token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
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

func (s *serveHTTPCmdOptions) verifyCaptcha(ip, cResponse string) bool {
	resp, err := s.hc.Verify(cResponse, ip)
	if err != nil {
		log.Println(err)
		return false
	}
	return resp.Success
}
