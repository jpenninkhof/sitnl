// Command agenda pulls the sitNL agenda, speakers, speaker photos and sponsors
// from the planning tool and the sponsor app, and writes them as static files:
//
//	../js/agenda.js        var proposalLineupJson = [...]
//	../js/sponsors.js      var sponsorLineupJson = [...]
//	../images/speakers/    <hash>.jpg
//
// Run it from the agenda/ folder:
//
//	go run .
//
// Credentials come from environment variables, or else from the first *.key
// file in this folder (git-ignored), with lines like:
//
//	Username: ...
//	Password: ...
//	OTP: 123456                      (a current code from the authenticator app), or
//	TOTP: <base32 secret>            (lets the script generate codes itself)
//
// Environment variables: APP_USERNAME, APP_PASSWORD, APP_OTP, APP_TOTP_SECRET.
//
// After a successful login the session is saved in .session (git-ignored) and
// reused on the next runs, so a one-time OTP is only needed again when the
// session expires (the planning tool issues sessions for 7 days).
//
// It exits non-zero on any failure, so a CI run never publishes partial or
// empty data. Output is deterministic (sorted), so an unchanged agenda gives
// an unchanged file and no git diff.
package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	loginURL     = "https://sitnl.cfapps.eu12.hana.ondemand.com/login"
	verify2FAURL = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/user/2fa/verify"
	statusURL    = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/user/2fa/status"
	baseURL      = "https://sitnl.cfapps.eu12.hana.ondemand.com/"
	sessionFile  = ".session"
	imagesURL    = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/speaker/allPhotos"
	proposalURL  = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/proposal"
	speakerURL   = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/speaker"
	sponsorURL   = "https://btp-ca-nl-dev-sitnl-srv.cfapps.eu10-004.hana.ondemand.com/api/sponsors"

	jsDir           = "../js"
	imagesTargetDir = "../images/speakers"
	agendaFile      = "../js/agenda.js"
	sponsorFile     = "../js/sponsors.js"
)

// Fields on a proposal that must never reach the public site.
var proposalPrivateFields = []string{
	"ratings",
	"acceptRejectErrors",
	"acceptRejectMailSent",
	"acceptRejectSuccesses",
	"accepted",
	"allSpeakersConfirmed",
	"associatedSpeakers",
	"confirmed",
	"extraInfo",
	"prerequisites",
	"loginKey",
}

// Fields on a speaker that must never reach the public site.
var speakerPrivateFields = []string{
	"loginKey",
	"physicallyPresent",
	"hash",
	"id",
	"email",
	"phone",
	"phoneNumber",
}

// Public profile fields copied from the speaker record onto the proposal speaker.
var speakerSocialFields = []string{
	"blueskyHandle",
	"mastodonHandle",
	"linkedInUrl",
	"githubUrl",
}

var client *http.Client

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}

func run() error {
	creds, err := loadCredentials()
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	client = &http.Client{
		Jar:     jar,
		Timeout: 2 * time.Minute,
		// Don't follow redirects: an unauthenticated API call redirects to the
		// login page, and that must show up as an error, not as HTML.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}

	if restoreSession() {
		fmt.Println("Reusing saved session")
	} else {
		if err := login(creds); err != nil {
			return fmt.Errorf("login: %w", err)
		}
		if err := saveSession(); err != nil {
			fmt.Fprintln(os.Stderr, "warning: could not save session:", err)
		}
		fmt.Println("Logged in, session saved to", sessionFile)
	}

	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		return err
	}

	if err := downloadPhotos(); err != nil {
		return fmt.Errorf("speaker photos: %w", err)
	}
	fmt.Println("Speaker photos written to", imagesTargetDir)

	proposals, err := fetchJSON(proposalURL)
	if err != nil {
		return fmt.Errorf("proposals: %w", err)
	}
	speakers, err := fetchJSON(speakerURL)
	if err != nil {
		return fmt.Errorf("speakers: %w", err)
	}

	agenda := cleanupProposals(proposals, speakers)
	sortAgenda(agenda)
	if err := writeJSFile(agendaFile, "proposalLineupJson", agenda); err != nil {
		return err
	}
	fmt.Printf("Wrote %d accepted sessions to %s\n", len(agenda), agendaFile)

	sponsors, err := fetchJSON(sponsorURL)
	if err != nil {
		return fmt.Errorf("sponsors: %w", err)
	}
	if err := writeJSFile(sponsorFile, "sponsorLineupJson", sponsors); err != nil {
		return err
	}
	fmt.Printf("Wrote %d sponsors to %s\n", len(sponsors), sponsorFile)

	return nil
}

type credentials struct {
	username, password, otp, totpSecret string
}

// loadCredentials reads APP_USERNAME, APP_PASSWORD and APP_TOTP_SECRET, and
// fills in anything missing from the first *.key file in the current folder.
func loadCredentials() (credentials, error) {
	c := credentials{
		username:   os.Getenv("APP_USERNAME"),
		password:   os.Getenv("APP_PASSWORD"),
		otp:        os.Getenv("APP_OTP"),
		totpSecret: os.Getenv("APP_TOTP_SECRET"),
	}

	if c.username == "" || c.password == "" || (c.otp == "" && c.totpSecret == "") {
		if files, _ := filepath.Glob("*.key"); len(files) > 0 {
			f, err := os.Open(files[0])
			if err != nil {
				return c, err
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				key, value, ok := strings.Cut(scanner.Text(), ":")
				if !ok {
					continue
				}
				value = strings.TrimSpace(value)
				switch strings.ToLower(strings.TrimSpace(key)) {
				case "username", "user":
					if c.username == "" {
						c.username = value
					}
				case "password":
					if c.password == "" {
						c.password = value
					}
				case "otp", "code":
					if c.otp == "" {
						c.otp = strings.ReplaceAll(value, " ", "")
					}
				case "totp", "totp secret", "secret", "2fa":
					if c.totpSecret == "" {
						c.totpSecret = value
					}
				}
			}
			if err := scanner.Err(); err != nil {
				return c, err
			}
		}
	}

	if c.username == "" || c.password == "" {
		return c, fmt.Errorf("no credentials: set APP_USERNAME/APP_PASSWORD or add a *.key file")
	}
	return c, nil
}

type loginResponse struct {
	Success       bool     `json:"success"`
	Requires2FA   bool     `json:"requires2FA"`
	RequiresSetup bool     `json:"requiresSetup"`
	Methods       []string `json:"methods"`
	Message       string   `json:"message"`
}

func login(c credentials) error {
	form := url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)

	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("wrong username or password")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}

	var lr loginResponse
	if err := json.Unmarshal(body, &lr); err != nil {
		// Older versions of the tool answered without JSON; treat as success
		// and let the first API call tell whether the session works.
		return nil
	}

	switch {
	case lr.RequiresSetup:
		return fmt.Errorf("2FA is not set up for this account yet; set it up in the planning tool first and save the secret")
	case lr.Requires2FA:
		switch {
		case c.totpSecret != "":
			return verify2FA(c.totpSecret)
		case c.otp != "":
			status, err := postCode(c.otp)
			if err != nil {
				return fmt.Errorf("2FA: %w", err)
			}
			if status != http.StatusOK {
				return fmt.Errorf("OTP rejected (status %d); codes are only valid for about 30 seconds, put a fresh one in the .key file and run again", status)
			}
			return nil
		default:
			return fmt.Errorf("2FA required: put a current code from your authenticator app in the .key file ('OTP: 123456') and run again")
		}
	}
	return nil
}

func postCode(code string) (int, error) {
	payload, _ := json.Marshal(map[string]string{"code": code})
	req, err := http.NewRequest(http.MethodPost, verify2FAURL, bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

type savedCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// restoreSession loads cookies from .session and checks with the planning
// tool that the session is still valid.
func restoreSession() bool {
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return false
	}
	var saved []savedCookie
	if json.Unmarshal(data, &saved) != nil || len(saved) == 0 {
		return false
	}
	u, _ := url.Parse(baseURL)
	cookies := make([]*http.Cookie, 0, len(saved))
	for _, c := range saved {
		cookies = append(cookies, &http.Cookie{Name: c.Name, Value: c.Value, Path: "/"})
	}
	client.Jar.SetCookies(u, cookies)

	resp, err := client.Get(statusURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}

func saveSession() error {
	u, _ := url.Parse(baseURL)
	var saved []savedCookie
	for _, c := range client.Jar.Cookies(u) {
		saved = append(saved, savedCookie{Name: c.Name, Value: c.Value})
	}
	data, err := json.Marshal(saved)
	if err != nil {
		return err
	}
	return os.WriteFile(sessionFile, data, 0o600)
}

func verify2FA(secret string) error {
	try := func(t time.Time) (int, error) {
		code, err := totp(secret, t)
		if err != nil {
			return 0, err
		}
		return postCode(code)
	}

	status, err := try(time.Now())
	if err != nil {
		return fmt.Errorf("2FA: %w", err)
	}
	if status == http.StatusOK {
		return nil
	}
	// A code generated at the very end of a 30-second window can expire on the
	// way. Wait for the next window and try once more.
	if status == http.StatusBadRequest || status == http.StatusUnauthorized {
		time.Sleep(time.Duration(30-time.Now().Unix()%30+1) * time.Second)
		if status, err = try(time.Now()); err == nil && status == http.StatusOK {
			return nil
		}
	}
	return fmt.Errorf("2FA code rejected (status %d); check the TOTP secret and the system clock", status)
}

// totp returns the RFC 6238 code (SHA-1, 6 digits, 30 s) for a base32 secret,
// the same code an authenticator app shows.
func totp(secret string, t time.Time) (string, error) {
	secret = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(secret), " ", ""))
	secret = strings.TrimRight(secret, "=")
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid TOTP secret: %w", err)
	}

	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(t.Unix()/30))
	mac := hmac.New(sha1.New, key)
	mac.Write(msg[:])
	sum := mac.Sum(nil)

	offset := sum[len(sum)-1] & 0x0f
	value := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	return fmt.Sprintf("%06d", value%1000000), nil
}

func get(u string) ([]byte, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("GET %s: not logged in (status %s)", u, resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: unexpected status %s", u, resp.Status)
	}
	return body, nil
}

func fetchJSON(u string) ([]interface{}, error) {
	body, err := get(u)
	if err != nil {
		return nil, err
	}
	var data []interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parsing JSON from %s: %w", u, err)
	}
	return data, nil
}

func downloadPhotos() error {
	if err := os.MkdirAll(imagesTargetDir, 0o755); err != nil {
		return err
	}

	body, err := get(imagesURL)
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("opening zip: %w", err)
	}

	for _, f := range reader.File {
		if f.FileInfo().IsDir() || filepath.Ext(f.Name) != ".jpg" {
			continue
		}
		// Only use the base name, so a path in the zip can never escape the target folder.
		target := filepath.Join(imagesTargetDir, filepath.Base(f.Name))
		if err := extract(f, target); err != nil {
			return err
		}
	}
	return nil
}

func extract(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening %s in zip: %w", f.Name, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("reading %s from zip: %w", f.Name, err)
	}

	// Leave identical files untouched, so git sees no change.
	if existing, err := os.ReadFile(target); err == nil && bytes.Equal(existing, data) {
		return nil
	}
	return os.WriteFile(target, data, 0o644)
}

func writeJSFile(path, varName string, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshalling %s: %w", path, err)
	}
	out := []byte("var " + varName + " = " + string(data) + ";\n")
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

func isAccepted(proposal map[string]interface{}) bool {
	switch v := proposal["accepted"].(type) {
	case bool:
		return v
	case string:
		return v == "accepted"
	}
	return false
}

func cleanupProposals(proposals, speakers []interface{}) []interface{} {
	speakersByHash := map[string]map[string]interface{}{}
	for _, s := range speakers {
		if m, ok := s.(map[string]interface{}); ok {
			if hash, ok := m["hash"].(string); ok {
				speakersByHash[hash] = m
			}
		}
	}

	cleaned := make([]interface{}, 0, len(proposals))
	for _, p := range proposals {
		proposal, ok := p.(map[string]interface{})
		if !ok || !isAccepted(proposal) {
			continue
		}
		cleanupProposal(proposal, speakersByHash)
		cleaned = append(cleaned, proposal)
	}
	return cleaned
}

func cleanupProposal(proposal map[string]interface{}, speakersByHash map[string]map[string]interface{}) {
	for _, field := range proposalPrivateFields {
		delete(proposal, field)
	}

	propSpeakers, ok := proposal["speakers"].([]interface{})
	if !ok {
		return
	}
	for _, ps := range propSpeakers {
		speaker, ok := ps.(map[string]interface{})
		if !ok {
			continue
		}

		if hash, ok := speaker["hash"].(string); ok {
			_, err := os.Stat(filepath.Join(imagesTargetDir, hash+".jpg"))
			speaker["hasPhoto"] = err == nil
			if err == nil {
				speaker["photoUrl"] = hash
			}

			if full, ok := speakersByHash[hash]; ok {
				for _, field := range speakerSocialFields {
					if v, ok := full[field]; ok {
						speaker[field] = v
					}
				}
			}
		}

		for _, field := range speakerPrivateFields {
			delete(speaker, field)
		}
	}
}

// sortAgenda orders sessions by start time, then location, then id, so the
// output does not depend on the order the API happens to return.
func sortAgenda(agenda []interface{}) {
	key := func(i int, field string) string {
		if m, ok := agenda[i].(map[string]interface{}); ok {
			if s, ok := m[field].(string); ok {
				return s
			}
		}
		return ""
	}
	sort.SliceStable(agenda, func(i, j int) bool {
		for _, field := range []string{"startTime", "location", "id"} {
			a, b := key(i, field), key(j, field)
			if a != b {
				return a < b
			}
		}
		return false
	})
}
