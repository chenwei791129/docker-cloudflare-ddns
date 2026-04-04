package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

type config struct {
	Token         string
	ZoneID        string
	DomainName    string
	Proxied       bool
	TTL           int
	CheckURL      string
	CheckInterval time.Duration
}

type cloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type cloudflareResponse[T any] struct {
	Success bool              `json:"success"`
	Errors  []cloudflareError `json:"errors"`
	Result  T                 `json:"result"`
}

type dnsRecord struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type ddnsState struct {
	ip       string
	recordID string
}

type dnsUpdatePayload struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

func logf(format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05 -0700")
	fmt.Printf("[%s] %s\n", ts, fmt.Sprintf(format, args...))
}

func loadConfig() (config, error) {
	c := config{
		Token:      os.Getenv("CF_TOKEN"),
		ZoneID:     os.Getenv("CF_ZONE_ID"),
		DomainName: os.Getenv("CF_DOMAIN_NAME"),
		CheckURL:   os.Getenv("CHECK_URL"),
		TTL:        1,
		Proxied:    false,
	}

	var missing []string
	if c.Token == "" {
		missing = append(missing, "CF_TOKEN")
	}
	if c.ZoneID == "" {
		missing = append(missing, "CF_ZONE_ID")
	}
	if c.DomainName == "" {
		missing = append(missing, "CF_DOMAIN_NAME")
	}
	if len(missing) > 0 {
		return c, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	if c.CheckURL == "" {
		c.CheckURL = "http://whatismyip.akamai.com/"
	}

	if v := os.Getenv("TTL"); v != "" {
		ttl, err := strconv.Atoi(v)
		if err != nil {
			return c, fmt.Errorf("invalid TTL value: %s", v)
		}
		c.TTL = ttl
	}

	if v := os.Getenv("CF_PROXIED"); v != "" {
		proxied, err := strconv.ParseBool(v)
		if err != nil {
			return c, fmt.Errorf("invalid CF_PROXIED value: %s", v)
		}
		c.Proxied = proxied
	}

	intervalStr := os.Getenv("CHECK_INTERVAL")
	if intervalStr == "" {
		intervalStr = "5m"
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return c, fmt.Errorf("invalid CHECK_INTERVAL value: %s", intervalStr)
	}
	c.CheckInterval = interval

	return c, nil
}

func getPublicIP(checkURL string) (string, error) {
	resp, err := httpClient.Get(checkURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	ip := strings.TrimSpace(string(body))
	if net.ParseIP(ip) == nil || strings.Contains(ip, ":") {
		return "", fmt.Errorf("invalid IPv4 address: %s", ip)
	}

	return ip, nil
}

func formatCFErrors(errs []cloudflareError) string {
	if len(errs) == 0 {
		return "(no error details)"
	}
	msgs := make([]string, len(errs))
	for i, e := range errs {
		msgs[i] = fmt.Sprintf("code %d: %s", e.Code, e.Message)
	}
	return strings.Join(msgs, "; ")
}

func newCFRequest(token, method, reqURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func getRecordID(cfg config) (string, error) {
	reqURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s&type=A", cfg.ZoneID, url.QueryEscape(cfg.DomainName))

	req, err := newCFRequest(cfg.Token, "GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to query DNS records: %w", err)
	}
	defer resp.Body.Close()

	var cfResp cloudflareResponse[[]dnsRecord]
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return "", fmt.Errorf("failed to decode API response: %w", err)
	}

	if !cfResp.Success {
		return "", fmt.Errorf("API query failed, errors: %s", formatCFErrors(cfResp.Errors))
	}

	if len(cfResp.Result) == 0 {
		return "", fmt.Errorf("no A record found for %s", cfg.DomainName)
	}

	return cfResp.Result[0].ID, nil
}

func updateRecord(cfg config, recordID, ip string) error {
	reqURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", cfg.ZoneID, recordID)

	payload := dnsUpdatePayload{
		Type:    "A",
		Name:    cfg.DomainName,
		Content: ip,
		TTL:     cfg.TTL,
		Proxied: cfg.Proxied,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}

	req, err := newCFRequest(cfg.Token, "PUT", reqURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update DNS record: %w", err)
	}
	defer resp.Body.Close()

	var cfResp cloudflareResponse[dnsRecord]
	if err := json.NewDecoder(resp.Body).Decode(&cfResp); err != nil {
		return fmt.Errorf("failed to decode update response: %w", err)
	}

	if !cfResp.Success {
		return fmt.Errorf("update failed, errors: %s", formatCFErrors(cfResp.Errors))
	}

	return nil
}

func runUpdate(cfg config, state *ddnsState) {
	logf("Checking current public IP from %s...", cfg.CheckURL)

	ip, err := getPublicIP(cfg.CheckURL)
	if err != nil {
		logf("fetch public IP failed: %v", err)
		return
	}
	logf("public IP: %s", ip)

	if state.ip == ip {
		logf("Current public IP matches cached IP. No update required! IP: %s", ip)
		return
	}

	if state.recordID == "" {
		var err error
		state.recordID, err = getRecordID(cfg)
		if err != nil {
			logf("query record ID failed: %v", err)
			return
		}
	}

	logf("Updating new IP: %s", ip)
	if err := updateRecord(cfg, state.recordID, ip); err != nil {
		logf("update record failed: %v", err)
		return
	}

	state.ip = ip
	logf("Update record successfully, new IP: %s", ip)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		logf("configuration error: %v", err)
		os.Exit(1)
	}

	logf("Starting cloudflare-ddns (interval: %s)", cfg.CheckInterval)

	var state ddnsState

	// Run immediately on startup
	runUpdate(cfg, &state)

	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		runUpdate(cfg, &state)
	}
}
