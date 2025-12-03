package discovery

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cmp-platform/backend/internal/db"
)

// CertificateScanner discovers certificates in the network
type CertificateScanner struct {
	db           *db.DB
	scanInterval time.Duration
	workers      int
}

// NewCertificateScanner creates a new certificate scanner
func NewCertificateScanner(db *db.DB, scanInterval time.Duration, workers int) *CertificateScanner {
	return &CertificateScanner{
		db:           db,
		scanInterval: scanInterval,
		workers:      workers,
	}
}

// ScanTarget represents a target to scan for certificates
type ScanTarget struct {
	Host     string
	Port     int
	Protocol string // "tls" or "https"
}

// Start starts the certificate discovery scanner
func (cs *CertificateScanner) Start(ctx context.Context) {
	ticker := time.NewTicker(cs.scanInterval)
	defer ticker.Stop()

	// Run initial scan
	cs.performScan()

	for {
		select {
		case <-ctx.Done():
			log.Println("Certificate scanner stopped")
			return
		case <-ticker.C:
			cs.performScan()
		}
	}
}

// performScan performs certificate discovery scan
func (cs *CertificateScanner) performScan() {
	// Get scan targets from database or configuration
	targets := cs.getScanTargets()

	log.Printf("Scanning %d targets for certificates", len(targets))

	for _, target := range targets {
		certs, err := cs.scanTarget(target)
		if err != nil {
			log.Printf("Error scanning %s:%d: %v", target.Host, target.Port, err)
			continue
		}

		for _, cert := range certs {
			cs.storeCertificate(cert, target)
		}
	}
}

// scanTarget scans a target for certificates
func (cs *CertificateScanner) scanTarget(target ScanTarget) ([]*x509.Certificate, error) {
	address := fmt.Sprintf("%s:%d", target.Host, target.Port)

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true, // For discovery, we skip verification
	})
	defer tlsConn.Close()

	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}

	state := tlsConn.ConnectionState()
	return state.PeerCertificates, nil
}

// storeCertificate stores discovered certificate in database
func (cs *CertificateScanner) storeCertificate(cert *x509.Certificate, target ScanTarget) {
	fingerprint := calculateFingerprint(cert)

	// Check if certificate already exists
	var exists bool
	err := cs.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM certificates WHERE fingerprint = $1)",
		fingerprint,
	).Scan(&exists)

	if err != nil {
		log.Printf("Error checking certificate existence: %v", err)
		return
	}

	if exists {
		// Update last scanned timestamp
		_, _ = cs.db.Exec(
			"UPDATE certificates SET last_scanned_at = NOW(), source = $1 WHERE fingerprint = $2",
			fmt.Sprintf("%s:%d", target.Host, target.Port),
			fingerprint,
		)
		return
	}

	// Insert new certificate
	sansJSON, _ := json.Marshal(cert.DNSNames)

	query := `INSERT INTO certificates 
	          (id, fingerprint, subject, sans, issuer, not_before, not_after, 
	           key_algo, key_size, status, source, last_scanned_at, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW(), NOW())
	          ON CONFLICT (fingerprint) DO UPDATE SET 
	          last_scanned_at = NOW(), source = EXCLUDED.source`

	certID := generateCertID(cert)

	_, err = cs.db.Exec(query,
		certID,
		fingerprint,
		cert.Subject.CommonName,
		sansJSON,
		cert.Issuer.CommonName,
		cert.NotBefore,
		cert.NotAfter,
		cert.PublicKeyAlgorithm.String(),
		getKeySize(cert),
		"active",
		fmt.Sprintf("%s:%d", target.Host, target.Port),
	)

	if err != nil {
		log.Printf("Error storing certificate: %v", err)
		return
	}

	log.Printf("Discovered new certificate: %s from %s:%d", cert.Subject.CommonName, target.Host, target.Port)
}

// getScanTargets retrieves scan targets from database or config
func (cs *CertificateScanner) getScanTargets() []ScanTarget {
	// TODO: Load from database or configuration
	// For now, return empty list
	return []ScanTarget{}
}

func calculateFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}

func generateCertID(cert *x509.Certificate) string {
	fp := calculateFingerprint(cert)
	return "disc-" + fp[:16]
}

func getKeySize(cert *x509.Certificate) int {
	switch key := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		return key.N.BitLen()
	case *ecdsa.PublicKey:
		return key.Curve.Params().BitSize
	default:
		return 0
	}
}

