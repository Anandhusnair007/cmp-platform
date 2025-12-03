package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	apiURL     = flag.String("api-url", "", "CMP API URL")
	apiToken   = flag.String("api-token", "", "CMP API token")
	agentID    = flag.String("agent-id", "", "Agent ID")
	hostname   = flag.String("hostname", "", "Agent hostname")
	port       = flag.Int("port", 8084, "Local agent port")
	certDir    = flag.String("cert-dir", "/var/lib/cmp-agent/certs", "Certificate storage directory")
)

type InstallJob struct {
	JobID    string `json:"job_id"`
	CertID   string `json:"cert_id"`
	Path     string `json:"path"`
	ReloadCmd *string `json:"reload_cmd"`
}

type CertData struct {
	CertPEM string `json:"cert_pem"`
	KeyPEM  string `json:"key_pem"`
	ChainPEM string `json:"chain_pem"`
}

func main() {
	flag.Parse()

	// Load from environment if not provided
	if *apiURL == "" {
		*apiURL = os.Getenv("CMP_API_URL")
	}
	if *apiToken == "" {
		*apiToken = os.Getenv("CMP_API_TOKEN")
	}
	if *agentID == "" {
		*agentID = os.Getenv("AGENT_ID")
		if *agentID == "" {
			*agentID = uuid.New().String()
		}
	}
	if *hostname == "" {
		*hostname = os.Getenv("AGENT_HOSTNAME")
		if *hostname == "" {
			var err error
			*hostname, err = os.Hostname()
			if err != nil {
				*hostname = "unknown"
			}
		}
	}

	log.Printf("Starting Linux Agent: ID=%s, Hostname=%s", *agentID, *hostname)

	// Ensure cert directory exists
	if err := os.MkdirAll(*certDir, 0755); err != nil {
		log.Fatalf("Failed to create cert directory: %v", err)
	}

	// Register with CMP
	if err := registerAgent(); err != nil {
		log.Printf("Warning: Failed to register agent: %v", err)
	}

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "healthy",
			"agent_id": *agentID,
			"hostname": *hostname,
		})
	})

	// Install endpoint (polled by CMP or called directly)
	router.POST("/install", handleInstall)

	// Check-in endpoint
	router.POST("/checkin", handleCheckin)

	// Start HTTP server
	log.Printf("Agent listening on port %d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
}

func registerAgent() error {
	// TODO: Implement agent registration with CMP API
	log.Println("Agent registration (placeholder)")
	return nil
}

func handleCheckin(c *gin.Context) {
	// Update last check-in time
	log.Println("Agent check-in received")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "agent_id": *agentID})
}

func handleInstall(c *gin.Context) {
	var job InstallJob
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received install job: %+v", job)

	// Fetch certificate from CMP API
	certData, err := fetchCertificate(job.CertID)
	if err != nil {
		log.Printf("Failed to fetch certificate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch certificate"})
		return
	}

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(job.Path), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
		return
	}

	// Write certificate file (combined cert + chain + key)
	certContent := certData.CertPEM
	if certData.ChainPEM != "" {
		certContent += "\n" + certData.ChainPEM
	}
	certContent += "\n" + certData.KeyPEM

	if err := os.WriteFile(job.Path, []byte(certContent), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write certificate file"})
		return
	}

	log.Printf("Certificate installed to %s", job.Path)

	// Execute reload command if provided
	if job.ReloadCmd != nil && *job.ReloadCmd != "" {
		log.Printf("Executing reload command: %s", *job.ReloadCmd)
		cmd := exec.Command("sh", "-c", *job.ReloadCmd)
		if err := cmd.Run(); err != nil {
			log.Printf("Warning: Reload command failed: %v", err)
			c.JSON(http.StatusOK, gin.H{
				"status":  "installed",
				"warning": "Reload command failed",
			})
			return
		}
		log.Printf("Reload command executed successfully")
	}

	c.JSON(http.StatusOK, gin.H{"status": "installed", "path": job.Path})
}

func fetchCertificate(certID string) (*CertData, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/certs/%s", *apiURL, certID), nil)
	if err != nil {
		return nil, err
	}

	if *apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+*apiToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch cert: %s", string(body))
	}

	var cert Certificate
	if err := json.NewDecoder(resp.Body).Decode(&cert); err != nil {
		return nil, err
	}

	return &CertData{
		CertPEM:  getStringValue(cert.CertPEM),
		KeyPEM:   "key-from-vault", // TODO: Retrieve key from secure storage
		ChainPEM: getStringValue(cert.ChainPEM),
	}, nil
}

// Helper to get string from pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Certificate model (should import from backend package in production)
type Certificate struct {
	ID       string  `json:"id"`
	CertPEM  *string `json:"cert_pem"`
	ChainPEM *string `json:"chain_pem"`
}

// GenerateCSR generates a certificate signing request
func GenerateCSR(commonName string, sans []string) (csrPEM []byte, privateKey *rsa.PrivateKey, err error) {
	// Generate private key
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Create CSR template
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
		},
		DNSNames: sans,
	}

	if len(sans) == 0 {
		template.DNSNames = []string{commonName}
	}

	// Create CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode to PEM
	csrPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	return csrPEM, privateKey, nil
}
