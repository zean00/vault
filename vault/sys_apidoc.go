package vault

import (
	"github.com/hashicorp/vault/apidoc/apidoc"
	"github.com/hashicorp/vault/logical/framework"
)

// Backend returns a new system backend for apidoc to parse
func Backend() *framework.Backend {
	return NewSystemBackend(&Core{}, nil).Backend
}

// ManualPaths returns a list of hand written paths. These are all paths
// that don't exist in the SystemBackend structure.
func ManualPaths() []apidoc.Path {
	return []apidoc.Path{
		sysGenerateRootAttempt(),
		sysGenerateRootUpdate(),
		sysLeader(),
		sysInit(),
		sysHealth(),
		sysRekeyInit(),
		sysRekeyUpdate(),
		//sysRekeyBackup(),
		sysRekeyRecoveryInit(),
		sysRekeyRecoveryUpdate(),
		sysRekeyRecoveryBackup(),
		//sysWrappingLookup(),
		//sysWrappingRewrap(),
		//sysWrappingUnwrap(),
		sealStatus(),
		seal(),
		stepDown(),
		unseal(),
	}
}

// The paths below were derived from API documentation because they're not present in help
// or in the System backend. Writing documentation this way, as separate structures, is not
// desireable long term. It would be better for this to be integrated framework structure.
// It would also be useful to extended the framework to accept more extensive documentation,
// such as example requests and reponses.

func sysGenerateRootAttempt() apidoc.Path {
	p := apidoc.NewPath("generate-root/attempt")

	// GET
	m := apidoc.NewMethod("GET", "Reads the configuration and process of the current root generation attempt.")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "progress": 1,
	  "required": 3,
	  "encoded_token": "",
	  "pgp_fingerprint": "816938b8a29146fbe245dd29e7cbaf8e011db793",
	  "complete": false
	}`)
	p.AddMethod(m)

	// PUT
	m = apidoc.NewMethod("PUT", "Initializes a new root generation attempt")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("otp", "string", "Specifies a base64-encoded 16-byte value."),
		apidoc.NewProperty("pgp_key", "string", "Specifies a base64-encoded PGP public key."),
	}
	m.AddResponse(200, `
	{
	    "started": true,
	    "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	    "progress": 1,
	    "required": 3,
	    "encoded_token": "",
	    "pgp_fingerprint": "",
	    "complete": false
	}`)
	p.AddMethod(m)

	// DELETE
	m = apidoc.NewMethod("DELETE", "Cancels any in-progress root generation attempt.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysGenerateRootUpdate() apidoc.Path {
	p := apidoc.NewPath("generate-root/update")

	// PUT
	m := apidoc.NewMethod("PUT", "Enter a single master key share to progress the root generation attempt.")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the attempt."),
	}
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "progress": 3,
	  "required": 3,
	  "pgp_fingerprint": "",
	  "complete": true,
	  "encoded_token": "FPzkNBvwNDeFh4SmGA8c+w=="
	}`)
	p.AddMethod(m)

	return p
}

func sysInit() apidoc.Path {
	p := apidoc.NewPath("init")

	m := apidoc.NewMethod("GET", sysHelp["init"][0])
	m.AddResponse(200, `{"initialized": true}`)
	p.AddMethod(m)

	m = apidoc.NewMethod("PUT", sysHelp["init"][0])
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("root_token_pgp_key", "string",
			"Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation."),
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the master key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("stored_shares", "number",
			"Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as secret_shares."),
		apidoc.NewProperty("recovery_pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as recovery_shares."),
	}
	m.AddResponse(200, `
		{
		  "keys": ["one", "two", "three"],
		  "keys_base64": ["cR9No5cBC", "F3VLrkOo", "zIDSZNGv"],
		  "root_token": "foo"
		}`)
	p.AddMethod(m)

	return p
}

func sysLeader() apidoc.Path {
	p := apidoc.NewPath("leader")
	m := apidoc.NewMethod("GET", "Check the high availability status and current leader of Vault")
	m.AddResponse(200, `
		{
            "ha_enabled": true,
            "is_self": false,
            "leader_address": "https://127.0.0.1:8200/",
            "leader_cluster_address": "https://127.0.0.1:8201/"
        }`)
	p.AddMethod(m)

	return p
}

func sealStatus() apidoc.Path {
	p := apidoc.NewPath("seal-status")
	m := apidoc.NewMethod("GET", sysHelp["seal-status"][0])
	m.AddResponse(200, `
		{
			  "type": "shamir",
			  "sealed": false,
			  "t": 3,
			  "n": 5,
			  "progress": 0,
			  "version": "0.9.0",
			  "cluster_name": "vault-cluster-d6ec3c7f",
			  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8",
			  "nonce": "ef05d55d-4d2c-c594-a5e8-55bc88604c24"
		}`)
	p.AddMethod(m)

	return p
}

func seal() apidoc.Path {
	p := apidoc.NewPath("seal")
	m := apidoc.NewMethod("GET", sysHelp["seal"][0])
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func stepDown() apidoc.Path {
	p := apidoc.NewPath("step-down")
	m := apidoc.NewMethod("PUT", "Causes the node to give up active status.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysHealth() apidoc.Path {
	p := apidoc.NewPath("health")
	m := apidoc.NewMethod("GET", "Returns the health status of Vault.")
	responses := []apidoc.Response{
		apidoc.NewResponse(200, "initialized, unsealed, and active", ""),
		apidoc.NewResponse(429, "unsealed and standby", ""),
		apidoc.NewResponse(472, "data recovery mode replication secondary and active", ""),
		apidoc.NewResponse(501, "not initialized", ""),
		apidoc.NewResponse(503, "sealed", ""),
	}
	m.Responses = responses
	p.AddMethod(m)

	m = apidoc.NewMethod("HEAD", "Returns the health status of Vault.")
	m.Responses = responses
	p.AddMethod(m)

	return p
}

func sysRekeyInit() apidoc.Path {
	p := apidoc.NewPath("rekey/init")

	// GET
	m := apidoc.NewMethod("GET", "Read the configuration and progress of the current rekey attempt.")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "t": 3,
	  "n": 5,
	  "progress": 1,
	  "required": 3,
	  "pgp_fingerprints": ["abcd1234"],
	  "backup": true
	}`)

	p.AddMethod(m)

	// PUT
	m = apidoc.NewMethod("PUT", "Initializes a new rekey attempt")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the master key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("backup", "boolean", "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys."),
	}
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	// DELETE
	m = apidoc.NewMethod("DELETE", "Cancels any in-progress rekey.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysRekeyUpdate() apidoc.Path {
	p := apidoc.NewPath("rekey/update")

	// PUT
	m := apidoc.NewMethod("PUT", "Enter a single master key share to progress the rekey of the Vault.")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the rekey attempt."),
	}
	m.AddResponse(200, `
	{
	  "complete": true,
	  "keys": ["one", "two", "three"],
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "pgp_fingerprints": ["abcd1234"],
	  "keys_base64": ["base64keyvalue"],
	  "backup": true
	}`)
	p.AddMethod(m)

	return p
}

func sysRekeyBackup() apidoc.Path {
	p := apidoc.NewPath("rekey/backup")

	// GET
	m := apidoc.NewMethod("PUT", "Return the backup copy of PGP-encrypted unseal keys.")
	m.AddResponse(200, `
	{
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "keys": {
		"abcd1234": "..."
	  }
	}`)

	p.AddMethod(m)

	// DELETE
	m = apidoc.NewMethod("DELETE", "Deletes the backup copy of PGP-encrypted unseal keys.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysRekeyRecoveryBackup() apidoc.Path {
	p := apidoc.NewPath("rekey-recovery-key/backup")

	// GET
	m := apidoc.NewMethod("GET", "Return the backup copy of PGP-encrypted recovery key shares.")
	m.AddResponse(200, `
	{
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "keys": {
		"abcd1234": "..."
	  }
	}`)

	p.AddMethod(m)

	// DELETE
	m = apidoc.NewMethod("DELETE", "Deletes the backup copy of PGP-encrypted recovery key shares.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysRekeyRecoveryInit() apidoc.Path {
	p := apidoc.NewPath("rekey-recovery-key/init")

	// GET
	m := apidoc.NewMethod("GET", "Read the configuration and progress of the current rekey attempt.")
	m.AddResponse(200, `
	{
	  "started": true,
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "t": 3,
	  "n": 5,
	  "progress": 1,
	  "required": 3,
	  "pgp_fingerprints": ["abcd1234"],
	  "backup": true
	}`)
	p.AddMethod(m)

	// PUT
	m = apidoc.NewMethod("PUT", "Initializes a new rekey attempt")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("secret_shares", "number",
			"Specifies the number of shares to split the recovery key into."),
		apidoc.NewProperty("secret_threshold", "number",
			"Specifies the number of shares required to reconstruct the recovery key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares."),
		apidoc.NewProperty("pgp_keys", "array/string",
			"Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares."),
		apidoc.NewProperty("backup", "boolean", "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys."),
	}
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	// DELETE
	m = apidoc.NewMethod("DELETE", "Cancels any in-progress rekey.")
	m.Responses = []apidoc.Response{apidoc.StdRespNoContent}
	p.AddMethod(m)

	return p
}

func sysRekeyRecoveryUpdate() apidoc.Path {
	p := apidoc.NewPath("rekey-recovery-key/update")

	// PUT
	m := apidoc.NewMethod("PUT", "Enter a single master key share to progress the rekey of the Vault.")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share."),
		apidoc.NewProperty("nonce", "string", "Specifies the nonce of the rekey attempt."),
	}
	m.AddResponse(200, `
	{
	  "complete": true,
	  "keys": ["one", "two", "three"],
	  "nonce": "2dbd10f1-8528-6246-09e7-82b25b8aba63",
	  "pgp_fingerprints": ["abcd1234"],
	  "keys_base64": ["base64keyvalue"],
	  "backup": true
	}`)
	p.AddMethod(m)

	return p
}

func sysWrappingLookup() apidoc.Path {
	p := apidoc.NewPath("wrapping/lookup")

	// POST
	m := apidoc.NewMethod("POST", "Look up wrapping properties for the given token.")
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("token", "string", "Specifies the wrapping token ID."),
	}
	m.AddResponse(200, `
	{
	  "request_id": "481320f5-fdf8-885d-8050-65fa767fd19b",
	  "lease_id": "",
	  "lease_duration": 0,
	  "renewable": false,
	  "data": {
		"creation_path": "sys/wrapping/wrap",
		"creation_time": "2016-09-28T14:16:13.07103516-04:00",
		"creation_ttl": 300
	  },
	  "wrap_info": null,
	  "warnings": null,
	  "auth": null
	}`)

	p.AddMethod(m)

	return p
}

func unseal() apidoc.Path {
	p := apidoc.NewPath("unseal")
	m := apidoc.NewMethod("PUT", sysHelp["unseal"][0])
	m.BodyFields = []apidoc.Property{
		apidoc.NewProperty("key", "string", "Specifies a single master key share. This is required unless reset is true."),
		apidoc.NewProperty("reset", "boolean", "Specifies if previously-provided unseal keys are discarded and the unseal process is reset."),
	}
	m.AddResponse(200, `
		{
		  "sealed": false,
		  "t": 3,
		  "n": 5,
		  "progress": 0,
		  "version": "0.6.2",
		  "cluster_name": "vault-cluster-d6ec3c7f",
		  "cluster_id": "3e8b3fec-3749-e056-ba41-b62a63b997e8"
		}`)

	p.AddMethod(m)

	return p
}
