package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	loginURL    = "https://sitnl.cfapps.eu12.hana.ondemand.com/login"
	imagesURL   = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/speaker/allPhotos"
	proposalURL = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/proposal"
	speakerURL  = "https://sitnl.cfapps.eu12.hana.ondemand.com/api/speaker"
	sponsorURL  = "https://presales-nl---build-apps-dev-sitnl-srv.cfapps.eu12.hana.ondemand.com/api/sponsors"

	imagesTargetDir = "../site/images/speakers"
	sponsorFile     = "../site/js/sponsors.js"
)

func main() {

	// Login
	cookies, err := performLogin(loginURL, os.Getenv("APP_USERNAME"), os.Getenv("APP_PASSWORD"))
	if err != nil {
		fmt.Println("Error during login:", err)
		return
	}

	// Download speaker images
	err = downloadAndUnpack(imagesURL, imagesTargetDir, cookies)
	if err != nil {
		fmt.Println("Error downloading and unpacking:", err)
	} else {
		fmt.Println("Successfully downloaded and unpacked speakers images to", imagesTargetDir)
	}

	// Fetch proposals
	proposals, err := fetchJson(proposalURL, cookies)
	if err != nil {
		fmt.Println("Error during retrieval of proposals:", err)
		return
	}

	// Fetch speakers
	speakers, err := fetchJson(speakerURL, cookies)
	if err != nil {
		fmt.Println("Error during retrieval of proposals:", err)
		return
	}

	// Clean up proposals
	cleanupProposals(proposals, speakers)
	writeAgenda(proposals)

	// Fetch and write agenda
	getSponsors()

}

func getSponsors() {

	content, err := fetchJson(sponsorURL, nil)
	if err != nil {
		fmt.Println("Error during retrieval of proposals:", err)
		return
	}

	// Marshal data back to JSON
	jsonContent, err := json.Marshal(content)
	if err != nil {
		fmt.Println("Error marshalling modified data:", err)
		return
	}

	// Write modified JSON to file
	err = ioutil.WriteFile(sponsorFile, []byte("var sponsorLineupJson = "+string(jsonContent)+";"), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Sponsors saved successfully!")
}

func downloadAndUnpack(url string, targetDir string, cookies []*http.Cookie) error {

	// Create target directory if it doesn't exist
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Create GET request for speaker images
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return err
	}

	// Add cookies from login response to GET request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Create HTTP client
	client := &http.Client{}

	// Perform GET request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful download
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d while downloading", resp.StatusCode)
	}

	// Create a new file to write the zip data
	outFile, err := os.Create(filepath.Join(targetDir, "allPhotos.zip"))
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Write the downloaded data to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write zip data to file: %w", err)
	}

	// Open the zip file
	reader, err := zip.OpenReader(filepath.Join(targetDir, "allPhotos.zip"))
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Extract all JPG files
	for _, f := range reader.File {
		if !f.FileInfo().IsDir() && filepath.Ext(f.Name) == ".jpg" {
			outFile, err := os.Create(filepath.Join(targetDir, f.Name))
			if err != nil {
				return fmt.Errorf("failed to create output file for %s: %w", f.Name, err)
			}
			defer outFile.Close()
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open file %s in zip: %w", f.Name, err)
			}
			defer rc.Close()
			_, err = io.Copy(outFile, rc)
			if err != nil {
				return fmt.Errorf("failed to copy file content %s: %w", f.Name, err)
			}
		}
	}

	// Remove the zip file
	os.Remove(filepath.Join(targetDir, "allPhotos.zip"))

	return nil
}

func writeAgenda(proposals []interface{}) {

	// Marshal modified data back to JSON
	agenda, err := json.Marshal(proposals)
	if err != nil {
		fmt.Println("Error marshalling modified data:", err)
		return
	}

	// Write modified JSON to file
	filePath := filepath.Join("..", "sitnl2024", "js", "agenda.js")
	err = ioutil.WriteFile(filePath, []byte("var proposalLineupJson = "+string(agenda)+";"), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Successfully wrote proposal data to agenda.js")

}

func performLogin(loginURL, username, password string) ([]*http.Cookie, error) {
	// Prepare login form data
	formData := url.Values{}
	formData.Add("username", username)
	formData.Add("password", password)

	// Create POST request body
	requestBody := bytes.NewBufferString(formData.Encode())

	// Create POST request
	req, err := http.NewRequest(http.MethodPost, loginURL, requestBody)
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		return nil, err
	}

	// Set content type for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create HTTP client
	client := &http.Client{}

	// Perform POST request (login)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing POST request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Login failed:", resp.Status)
		return nil, err
	}

	// Extract cookies from response
	return resp.Cookies(), nil
}

func fetchJson(url string, cookies []*http.Cookie) ([]interface{}, error) {

	// Create GET request for proposals
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating GET request:", err)
		return nil, err
	}

	// Add cookies from login response to GET request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Create HTTP client
	client := &http.Client{}

	// Perform GET request (proposals)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing GET request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Parse JSON response as an array
	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return nil, err
	}

	return data, nil

}

func cleanupProposals(proposals []interface{}, speakers []interface{}) {
	for i := 0; i < len(proposals); i++ {
		if proposal, ok := proposals[i].(map[string]interface{}); ok {
			cleanupProposal(proposal, speakers)
		}
	}
}

func cleanupProposal(proposals map[string]interface{}, speakers []interface{}) {
	delete(proposals, "ratings")
	delete(proposals, "acceptRejectErrors")
	delete(proposals, "acceptRejectMailSent")
	delete(proposals, "acceptRejectSuccesses")
	delete(proposals, "accepted")
	delete(proposals, "allSpeakersConfirmed")
	delete(proposals, "associatedSpeakers")
	delete(proposals, "confirmed")
	delete(proposals, "extraInfo")
	delete(proposals, "prerequisites")
	delete(proposals, "loginKey")

	// Access propSpeakers list and remove "loginKey" from each speaker
	if propSpeakers, ok := proposals["speakers"].([]interface{}); ok {
		for _, propSpeaker := range propSpeakers {
			if propSpeakerMap, ok := propSpeaker.(map[string]interface{}); ok {

				if hash, ok := propSpeakerMap["hash"].(string); ok {

					// Check if the speaker has an image
					imagePath := filepath.Join("..", "sitnl2024", "images", "speakers", hash+".jpg")
					_, err := os.Stat(imagePath)
					propSpeakerMap["hasPhoto"] = err == nil
					if err == nil {
						propSpeakerMap["photoUrl"] = hash
					}

					// Search for the speaker with the current hash
					for i := 0; i < len(speakers); i++ {
						if speakerMap, ok := speakers[i].(map[string]interface{}); ok {
							if speakerMap["hash"] == hash {
								copyWhenExists(speakerMap, propSpeakerMap, "twitterHandle")
								copyWhenExists(speakerMap, propSpeakerMap, "blueskyHandle")
								copyWhenExists(speakerMap, propSpeakerMap, "mastodonHandle")
								copyWhenExists(speakerMap, propSpeakerMap, "linkedInUrl")
								copyWhenExists(speakerMap, propSpeakerMap, "githubUrl")
							}
						}
					}
				}

				delete(propSpeakerMap, "loginKey")
				delete(propSpeakerMap, "physicallyPresent")
				delete(propSpeakerMap, "hash")
				delete(propSpeakerMap, "id")

			}
		}
	}
}

func copyWhenExists(speakerMap map[string]interface{}, propSpeakerMap map[string]interface{}, key string) {
	value, ok := speakerMap[key]
	if ok {
		propSpeakerMap[key] = value
	}
}
