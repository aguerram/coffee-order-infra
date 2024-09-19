package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	appRoleName        = "inventory-service"
	k8sSecretName      = "inventory-approle-dev"
	k8sSecretNamespace = "ns-coffee-order-demo-dev"
	policyName         = "coffee-order-inventory-policy"
	secretIDTTL        = "0" // Set to 0 for no expiration
	tokenNumUses       = 0   // Set to 0 for unlimited uses
	tokenTTL           = "0" // Set to 0 for no expiration
	tokenMaxTTL        = "0" // Set to 0 for no expiration
	secretIDNumUses    = 0   // Set to 0 for unlimited uses
)

var secretPaths = []string{
	"secrets/data/coffee-demo/inventory-service/*",
}

func main() {
	vaultAddr := getEnvOrPanic("VAULT_ADDR")
	vaultToken := getEnvOrPanic("VAULT_TOKEN")
	vaultNamespace := getEnv("VAULT_NAMESPACE", "")

	client := &http.Client{}

	roleId, secretId := generateRoleIdAndSecretID(client, vaultAddr, vaultToken, vaultNamespace)
	fmt.Printf("RoleID: %s\n", roleId)
	fmt.Printf("SecretID: %s\n", secretId)

	//pin the secret id and role id in the env variables ROLE_ID and SECRET_ID
	fmt.Printf("----------------------------------\n")
	fmt.Printf(`
kubectl create secret generic %s \
--namespace=%s \
--from-literal=role_id=%s \
--from-literal=secret_id=%s
	`, k8sSecretName, k8sSecretNamespace, roleId, secretId)
}

func generateRoleIdAndSecretID(client *http.Client, vaultAddr, vaultToken, vaultNamespace string) (string, string) {
	// Create policy
	fmt.Println("---- Create policy ----")
	policyContent := createDynamicPolicyContent(secretPaths)
	createPolicy(client, vaultAddr, vaultToken, vaultNamespace, policyName, policyContent)

	// Enable AppRole auth method
	fmt.Println("---- Enable AppRole auth method ----")
	enableAppRole(client, vaultAddr, vaultToken, vaultNamespace)

	// Create AppRole
	fmt.Println("---- Create an AppRole ----")
	createAppRole(client, vaultAddr, vaultToken, vaultNamespace, appRoleName)

	// Get RoleID
	fmt.Println("---- Get the RoleID ----")
	roleID := getRoleID(client, vaultAddr, vaultToken, vaultNamespace, appRoleName)

	// Generate SecretID
	fmt.Println("---- Generate a SecretID ----")
	secretID := generateSecretID(client, vaultAddr, vaultToken, vaultNamespace, appRoleName)

	return roleID, secretID
}

func getEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func makeAPICall(client *http.Client, method, url, token, namespace string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Vault-Token", token)
	req.Header.Set("Content-Type", "application/json")
	if namespace != "" {
		req.Header.Set("X-Vault-Namespace", namespace)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("API call failed with status code %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func createPolicy(client *http.Client, vaultAddr, token, namespace, policyName, policyContent string) {
	url := fmt.Sprintf("%s/v1/sys/policies/acl/%s", vaultAddr, policyName)
	data := map[string]string{"policy": policyContent}
	jsonData, _ := json.Marshal(data)
	_, err := makeAPICall(client, "PUT", url, token, namespace, jsonData)
	if err != nil {
		panic(err)
	}
}

func createAppRole(client *http.Client, vaultAddr, token, namespace, appRoleName string) {
	url := fmt.Sprintf("%s/v1/auth/approle/role/%s", vaultAddr, appRoleName)
	data := map[string]interface{}{
		"secret_id_ttl":      secretIDTTL,
		"token_num_uses":     tokenNumUses,
		"token_ttl":          tokenTTL,
		"token_max_ttl":      tokenMaxTTL,
		"secret_id_num_uses": secretIDNumUses,
		"policies":           []string{policyName},
	}
	jsonData, _ := json.Marshal(data)
	_, err := makeAPICall(client, "POST", url, token, namespace, jsonData)
	if err != nil {
		panic(err)
	}
}

func getRoleID(client *http.Client, vaultAddr, token, namespace, appRoleName string) string {
	url := fmt.Sprintf("%s/v1/auth/approle/role/%s/role-id", vaultAddr, appRoleName)
	body, err := makeAPICall(client, "GET", url, token, namespace, nil)
	if err != nil {
		panic(err)
	}
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	return response["data"].(map[string]interface{})["role_id"].(string)
}

func generateSecretID(client *http.Client, vaultAddr, token, namespace, appRoleName string) string {
	url := fmt.Sprintf("%s/v1/auth/approle/role/%s/secret-id", vaultAddr, appRoleName)
	body, err := makeAPICall(client, "POST", url, token, namespace, nil)
	if err != nil {
		panic(err)
	}
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	return response["data"].(map[string]interface{})["secret_id"].(string)
}

func createDynamicPolicyContent(paths []string) string {
	var policyBuilder strings.Builder
	for _, path := range paths {
		policyBuilder.WriteString(fmt.Sprintf(`
path "%s" {
  capabilities = ["read"]
}
`, path))
	}
	return policyBuilder.String()
}

func enableAppRole(client *http.Client, vaultAddr, token, namespace string) {
	url := fmt.Sprintf("%s/v1/sys/auth/approle", vaultAddr)
	data := map[string]string{"type": "approle"}
	jsonData, _ := json.Marshal(data)
	_, err := makeAPICall(client, "POST", url, token, namespace, jsonData)
	if err != nil {
		if strings.Contains(err.Error(), "path is already in use at approle/") {
			fmt.Println("Warning: AppRole auth method is already enabled")
		} else {
			panic(err)
		}
	}
}
