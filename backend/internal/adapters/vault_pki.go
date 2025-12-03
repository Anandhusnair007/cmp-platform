package adapters

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type VaultPKIAdapter struct {
	client     *vault.Client
	mountPath  string
	roleName   string
	commonName string
}

type VaultPKIConfig struct {
	Address   string `json:"address"`
	Token     string `json:"token"`
	MountPath string `json:"mount_path"`
	RoleName  string `json:"role_name"`
	CommonName string `json:"common_name"`
}

func NewVaultPKIAdapter(config map[string]interface{}) (*VaultPKIAdapter, error) {
	configJSON, _ := json.Marshal(config)
	var cfg VaultPKIConfig
	if err := json.Unmarshal(configJSON, &cfg); err != nil {
		return nil, fmt.Errorf("invalid vault config: %w", err)
	}

	// Create Vault client
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = cfg.Address
	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(cfg.Token)

	return &VaultPKIAdapter{
		client:     client,
		mountPath:  cfg.MountPath,
		roleName:   cfg.RoleName,
		commonName: cfg.CommonName,
	}, nil
}

func (a *VaultPKIAdapter) CreateCertificate(commonName string, sans []string, keyAlgo string, keySize int) (certPEM, chainPEM, keyPEM string, err error) {
	// Generate private key
	var privateKey interface{}
	switch keyAlgo {
	case "rsa":
		privateKey, err = rsa.GenerateKey(rand.Reader, keySize)
		if err != nil {
			return "", "", "", fmt.Errorf("failed to generate RSA key: %w", err)
		}
	case "ecdsa":
		// ECDSA implementation would go here
		return "", "", "", fmt.Errorf("ECDSA not yet implemented")
	default:
		return "", "", "", fmt.Errorf("unsupported key algorithm: %s", keyAlgo)
	}

	// Create CSR
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
		},
		DNSNames: sans,
	}

	if len(sans) == 0 {
		template.DNSNames = []string{commonName}
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create CSR: %w", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	// Prepare Vault request
	data := map[string]interface{}{
		"csr":           string(csrPEM),
		"common_name":   commonName,
		"alt_names":     fmt.Sprintf("%v", sans),
		"ttl":           "8760h", // 1 year
		"format":        "pem",
	}

	// Issue certificate from Vault
	path := fmt.Sprintf("%s/sign/%s", a.mountPath, a.roleName)
	secret, err := a.client.Logical().Write(path, data)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to sign certificate: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return "", "", "", fmt.Errorf("empty response from Vault")
	}

	certPEM, ok := secret.Data["certificate"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("invalid certificate in Vault response")
	}

	chainPEM, _ = secret.Data["ca_chain"].(string)
	if chainPEM == "" {
		chainPEM, _ = secret.Data["issuing_ca"].(string)
	}

	// Encode private key
	var keyBlock *pem.Block
	switch k := privateKey.(type) {
	case *rsa.PrivateKey:
		keyBytes := x509.MarshalPKCS1PrivateKey(k)
		keyBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		}
	}
	keyPEM = string(pem.EncodeToMemory(keyBlock))

	return certPEM, chainPEM, keyPEM, nil
}

func (a *VaultPKIAdapter) RevokeCertificate(serialNumber string) error {
	path := fmt.Sprintf("%s/revoke", a.mountPath)
	_, err := a.client.Logical().Write(path, map[string]interface{}{
		"serial_number": serialNumber,
	})
	return err
}

func (a *VaultPKIAdapter) GetStatus(certID string) (map[string]interface{}, error) {
	// Vault doesn't provide a direct status endpoint, return basic info
	return map[string]interface{}{
		"adapter": "vault-pki",
		"mount":   a.mountPath,
		"role":    a.roleName,
	}, nil
}
