package sdk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRecords(t *testing.T) {
	mux, client := setup(t)
	mux.HandleFunc("/api/zones/records/get", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"response": {
				"zone": {
					"name": "example.com",
					"type": "Primary",
					"internal": false,
					"dnssecStatus": "SignedWithNSEC3",
					"disabled": false
				},
				"records": [
					{
						"disabled": false,
						"name": "example.com",
						"type": "A",
						"ttl": 3600,
						"rData": {
							"ipAddress": "1.1.1.1"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "NS",
						"ttl": 3600,
						"rData": {
							"nameServer": "server1"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "SOA",
						"ttl": 900,
						"rData": {
							"primaryNameServer": "server1",
							"responsiblePerson": "hostadmin.example.com",
							"serial": 35,
							"refresh": 900,
							"retry": 300,
							"expire": 604800,
							"minimum": 900
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "NSEC3PARAM",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T11:45:31Z",
							"signatureInception": "2022-03-05T10:45:31Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "vJ/fXkGKsapdvWjDhcfHsBxpZhSzMRLZv3/bEGJ4N3/K7jiM92Ik336W680SI7g+NyPCQ3gqE7ta/JEL4bht4Q=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "SOA",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T12:53:39Z",
							"signatureInception": "2022-03-05T11:53:39Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "9PQHH3ZGCuFRYkn28SoilS8y8zszgeOpCfJpIOAaE5ao+iBPCXudHacr/EpgB2wLzXpRjR+WgiYjmJH17+6bKg=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "A",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T11:25:35Z",
							"signatureInception": "2022-03-05T10:25:35Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "dWjn5hTWuEq57ncwGdVq+kdbMuFtuxLuZhYCcQMdsTxYkM/64RrPY6eYwfYQ7+fY1+QBSX2WudAM4dzbmL/s2A=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "NS",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T11:25:35Z",
							"signatureInception": "2022-03-05T10:25:35Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "Yx+leBcYNFf0gUfN6rECWrUZwCDhJbAGk1BNOJN01nPakS5meSbDApUHJZeAzfSBcPzodK3ddmEuhho1MABaZw=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 86400,
						"rData": {
							"typeCovered": "DNSKEY",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 86400,
							"signatureExpiration": "2022-03-15T12:27:09Z",
							"signatureInception": "2022-03-05T11:27:09Z",
							"keyTag": 65078,
							"signersName": "example.com",
							"signature": "KWAK7o+FjJ2/6ZvX4C1wB41yRzlmec5pR2TTeNWlY/weg0MNKCLRs3uTopSjoTih+uq3IRR7Zx0iOcy7evOitA=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "RRSIG",
						"ttl": 86400,
						"rData": {
							"typeCovered": "DNSKEY",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 86400,
							"signatureExpiration": "2022-03-15T12:27:09Z",
							"signatureInception": "2022-03-05T11:27:09Z",
							"keyTag": 52896,
							"signersName": "example.com",
							"signature": "oHtt1gUmDXxI5GMfS+LJ6uxKUcuUu+5EELXdhLrbk5V/yganP6sMgA4hGkzokYM22LDowjSdO5qwzCW6IDgKxg=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "DNSKEY",
						"ttl": 86400,
						"rData": {
							"flags": "SecureEntryPoint, ZoneKey",
							"protocol": 3,
							"algorithm": "ECDSAP256SHA256",
							"publicKey": "dMRyc/Pji31mF3iHNrybPzbgvtb2NKtmXhjQq433BHI= ZveDa1z00VxDnugV1x7EDvpt+42TDh8OQwp1kOrpX0E=",
							"computedKeyTag": 65078,
							"dnsKeyState": "Ready",
							"computedDigests": [
								{
									"digestType": "SHA256",
									"digest": "BBE017B17E5CB5FFFF1EC2C7815367DF80D8E7EAEE4832D3ED192159D79B1EEB"
								},
								{
									"digestType": "SHA384",
									"digest": "0B0C9F1019BD3FE62C8B71F8C80E7A833BA468A7E303ABC819C0CB9BEDE8E26BB50CB1729547BFCCE2AE22390E44CDA3"
								}
							]
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "DNSKEY",
						"ttl": 86400,
						"rData": {
							"flags": "ZoneKey",
							"protocol": 3,
							"algorithm": "ECDSAP256SHA256",
							"publicKey": "IUvzTkf4JPg+7k57cQw7n7SR6/1dH7FaKxu9Cf+kcvo= UU+uoKRWnYAFHDNF0X3U8ZYetUyDF7fcNAwEaSQnIUM=",
							"computedKeyTag": 61009,
							"dnsKeyState": "Active"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "DNSKEY",
						"ttl": 3600,
						"rData": {
							"flags": "SecureEntryPoint, ZoneKey",
							"protocol": 3,
							"algorithm": "ECDSAP256SHA256",
							"publicKey": "KOJWFitKm58EgjO43GDnsFbnkGoqVKeLRkP8FGPAdhqA2F758Ta1mkxieEu0YN0EoX+u5bVuc5DEBFSv+U63CA==",
							"computedKeyTag": 15048,
							"dnsKeyState": "Published",
							"dnsKeyStateReadyBy": "2022-12-18T16:14:50.0328321Z",
							"computedDigests": [
								{
									"digestType": "SHA256",
									"digest": "8EAFAE3305DB57A27CA5A261525515461CB7232A34A44AD96441B88BCA9B9849"
								},
								{
									"digestType": "SHA384",
									"digest": "4A6DA59E91872B5B835FCEE5987B17151A6F10FE409B595BEEEDB28FE64315C9C268493B59A0BF72EA84BE0F20A33F96"
								}
							]
						},
						"dnssecStatus": "Unknown",
						"lastUsedOn": "0001-01-01T00:00:00"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "DNSKEY",
						"ttl": 86400,
						"rData": {
							"flags": "ZoneKey",
							"protocol": 3,
							"algorithm": "ECDSAP256SHA256",
							"publicKey": "337uQ11fdKbr6sKYq9mwwBC2xdnu0geuIkfHcIauKNI= rKk7pfVKlLfcGBOIn5hEVeod2aIRIyUiivdTPzrmpIo=",
							"computedKeyTag": 4811,
							"dnsKeyState": "Published"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "example.com",
						"type": "NSEC3PARAM",
						"ttl": 900,
						"rData": {
							"hashAlgorithm": "SHA1",
							"flags": "None",
							"iterations": 0,
							"salt": ""
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "*.example.com",
						"type": "A",
						"ttl": 3600,
						"rData": {
							"ipAddress": "7.7.7.7"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "*.example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "A",
							"algorithm": "ECDSAP256SHA256",
							"labels": 2,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T11:25:35Z",
							"signatureInception": "2022-03-05T10:25:35Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "ZoUNNEdb8XWqHHi5o4BcUe7deRVlJZLhQtc3sjRtuJ68DNPDmQ0GfCrNTigJcomspr7CYqWcXfoSOqu6f2AyyQ=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "4F3CNT8CU22TNGEC382JJ4GDE4RB47UB.example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "NSEC3",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T11:45:31Z",
							"signatureInception": "2022-03-05T10:45:31Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "piZeLYa6WpHyiJerPlXq2s+JKBjHznNALXHJCOfiQ4o/iTqWILoqYHfKB5AWrLwLmkxXcbKf63CnEMGlinRidg=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "4F3CNT8CU22TNGEC382JJ4GDE4RB47UB.example.com",
						"type": "NSEC3",
						"ttl": 900,
						"rData": {
							"hashAlgorithm": "SHA1",
							"flags": "None",
							"iterations": 0,
							"salt": "",
							"nextHashedOwnerName": "KG19N32806C832KIJDNGLQ8P9M2R5MDJ",
							"types": [
								"A"
							]
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "KG19N32806C832KIJDNGLQ8P9M2R5MDJ.example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "NSEC3",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T11:45:31Z",
							"signatureInception": "2022-03-05T10:45:31Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "i/PMxc1LFA9a8jLxju7SSpoY7y8aZYkAILcCRIxE3lTundPJmzFG0U9kve04kqT7+Klmzj3OzXnCvjTA54+DZA=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "KG19N32806C832KIJDNGLQ8P9M2R5MDJ.example.com",
						"type": "NSEC3",
						"ttl": 900,
						"rData": {
							"hashAlgorithm": "SHA1",
							"flags": "None",
							"iterations": 0,
							"salt": "",
							"nextHashedOwnerName": "MIFDNDT3NFF3OD53O7TLA1HRFF95JKUK",
							"types": [
								"NS",
								"DS"
							]
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "MIFDNDT3NFF3OD53O7TLA1HRFF95JKUK.example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "NSEC3",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T11:45:31Z",
							"signatureInception": "2022-03-05T10:45:31Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "mr37TDMmWJ3YLNtpYy++S9eAeHIXKajX6jB8zLscJyC1uI0OFnSTuesfhIlLDbj0SDgrzRQWsLmvMKzfq89TJA=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "MIFDNDT3NFF3OD53O7TLA1HRFF95JKUK.example.com",
						"type": "NSEC3",
						"ttl": 900,
						"rData": {
							"hashAlgorithm": "SHA1",
							"flags": "None",
							"iterations": 0,
							"salt": "",
							"nextHashedOwnerName": "ONIB9MGUB9H0RML3CDF5BGRJ59DKJHVK",
							"types": [
								"CNAME"
							]
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "ONIB9MGUB9H0RML3CDF5BGRJ59DKJHVK.example.com",
						"type": "RRSIG",
						"ttl": 900,
						"rData": {
							"typeCovered": "NSEC3",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 900,
							"signatureExpiration": "2022-03-15T11:45:31Z",
							"signatureInception": "2022-03-05T10:45:31Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "GGh/KkB6C2D55xRJa0zFbZ8As3DZK9btUamryZVmyo7FaLPyltkeRZor9OExgQ6HC1SLXNGJIfCO9cM4K6P8iw=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "ONIB9MGUB9H0RML3CDF5BGRJ59DKJHVK.example.com",
						"type": "NSEC3",
						"ttl": 900,
						"rData": {
							"hashAlgorithm": "SHA1",
							"flags": "None",
							"iterations": 0,
							"salt": "",
							"nextHashedOwnerName": "4F3CNT8CU22TNGEC382JJ4GDE4RB47UB",
							"types": [
								"A",
								"NS",
								"SOA",
								"DNSKEY",
								"NSEC3PARAM"
							]
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "sub.example.com",
						"type": "NS",
						"ttl": 3600,
						"rData": {
							"nameServer": "server1"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "sub.example.com",
						"type": "DS",
						"ttl": 3600,
						"rData": {
							"keyTag": 46125,
							"algorithm": "ECDSAP384SHA384",
							"digestType": "SHA1",
							"digest": "5590E425472785A16DC0F853000557DB5543C39E"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "sub.example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "NS",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T11:25:35Z",
							"signatureInception": "2022-03-05T10:25:35Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "hFzYTL9V0/0UQZlvZpRWCOvu/2udvhswKoxpe4+quNuC6K59W7uCJLuDm/z0aFK5nW8Of4oTk2YjSBZo0nBSlg=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "sub.example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "DS",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T12:53:39Z",
							"signatureInception": "2022-03-05T11:53:39Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "UYpUKV5Uq7DM3rltg3sPFOwYgRa2yBzT/j9U8xCh5oyXt27fIn3eemvqqe9qV4xeQaAN0QfQPkj9vmOZSAYafg=="
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "www.example.com",
						"type": "CNAME",
						"ttl": 3600,
						"rData": {
							"cname": "example.com"
						},
						"dnssecStatus": "Unknown"
					},
					{
						"disabled": false,
						"name": "www.example.com",
						"type": "RRSIG",
						"ttl": 3600,
						"rData": {
							"typeCovered": "CNAME",
							"algorithm": "ECDSAP256SHA256",
							"labels": 3,
							"originalTtl": 3600,
							"signatureExpiration": "2022-03-15T11:25:35Z",
							"signatureInception": "2022-03-05T10:25:35Z",
							"keyTag": 61009,
							"signersName": "example.com",
							"signature": "cAbYvDJhZGLS/uI5I4mSrh7S5gEUy6bmX2sY7zEd1XVFPqrUOZHbVZuwXPjA6r9/m0rCaww9RiG90JhNNDLEtA=="
						},
						"dnssecStatus": "Unknown"
					}
				]
			},
			"status": "ok"
}`)
	})

	records, _, err := client.RecordsAPI.ListRecords("example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(records) == 0 || records[0].Name != "example.com" {
		t.Errorf("unexpected records response: %+v", records)
	}
}

func TestCreateARecord(t *testing.T) {
	mux, client := setup(t)
	mux.HandleFunc("GET /api/zones/records/add", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if domain := q.Get("domain"); domain != "example.com" {
			t.Errorf("unexpected domain: wanted: %v, got: %v", "example.com", domain)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"response": {
				"zone": {
					"name": "example.com",
					"type": "Primary",
					"internal": false,
					"dnssecStatus": "SignedWithNSEC",
					"disabled": false
				},
				"addedRecord": {
					"disabled": false,
					"name": "example.com",
					"type": "A",
					"ttl": 3600,
					"rData": {
						"ipAddress": "3.3.3.3"
					},
					"dnssecStatus": "Unknown",
					"lastUsedOn": "0001-01-01T00:00:00"
				}
			},
			"status": "ok"
}`)
	})

	ipAddress := "3.3.3.3"
	record, _, err := client.RecordsAPI.CreateRecord(&RecordRequest{
		Token:     "test-token",
		Domain:    "example.com",
		Type:      "A",
		IPAddress: &ipAddress,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if record.Name != "example.com" {
		t.Errorf("unexpected record response: %+v", record)
	}
}

func TestCreateCNAMERecord(t *testing.T) {
	mux, client := setup(t)
	mux.HandleFunc("GET /api/zones/records/add", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"response": {
				"zone": {
					"name": "example.com",
					"type": "Primary",
					"internal": false,
					"dnssecStatus": "SignedWithNSEC",
					"disabled": false
				},
				"addedRecord": {
					"disabled": false,
					"name": "example.com",
					"type": "CNAME",
					"ttl": 3600,
					"rData": {
						"cname": "example.internal.com"
					},
					"dnssecStatus": "Unknown",
					"lastUsedOn": "0001-01-01T00:00:00"
				}
			},
			"status": "ok"
}`)
	})

	cname := "example.internal.com"
	record, _, err := client.RecordsAPI.CreateRecord(&RecordRequest{
		Token:  "test-token",
		Domain: "example.com",
		Type:   "CNAME",
		CNAME:  &cname,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if record.Name != "example.com" || *record.RData.CNAME != cname || record.Type != "CNAME" {
		t.Errorf("unexpected record response: %+v", record)
	}
}
func TestListZones(t *testing.T) {
	mux, client := setup(t)
	mux.HandleFunc("GET /api/zones/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"response": {
				"pageNumber": 1,
				"totalPages": 2,
				"totalZones": 12,
				"zones": [
					{
						"name": "",
						"type": "Secondary",
						"dnssecStatus": "SignedWithNSEC",
						"soaSerial": 1,
						"expiry": "2022-02-26T07:57:08.1842183Z",
						"isExpired": false,
						"syncFailed": false,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "0.in-addr.arpa",
						"type": "Primary",
						"internal": true,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa",
						"type": "Primary",
						"internal": true,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "127.in-addr.arpa",
						"type": "Primary",
						"internal": true,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "255.in-addr.arpa",
						"type": "Primary",
						"internal": true,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "example.com",
						"type": "Primary",
						"internal": false,
						"dnssecStatus": "SignedWithNSEC",
						"soaSerial": 1,
						"notifyFailed": false,
						"notifyFailedFor": [],
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "localhost",
						"type": "Primary",
						"internal": true,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "test0.com",
						"type": "Primary",
						"internal": false,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"notifyFailed": false,
						"notifyFailedFor": [],
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "test1.com",
						"type": "Primary",
						"internal": false,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"notifyFailed": false,
						"notifyFailedFor": [],
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					},
					{
						"name": "test2.com",
						"type": "Primary",
						"internal": false,
						"dnssecStatus": "Unsigned",
						"soaSerial": 1,
						"notifyFailed": false,
						"notifyFailedFor": [],
						"lastModified": "2022-02-26T07:57:08.1842183Z",
						"disabled": false
					}
				]
			},
			"status": "ok"
}`)
	})

	zones, _, err := client.ZonesAPI.ListZones()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if zones[0].Name != "" {
		t.Errorf("unexpected record response: %+v", zones[0])
	}
}

func TestLogin(t *testing.T) {
	_, client := setup(t)

	token, _, err := client.UsersAPI.Login("admin", "admin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token != "932b2a3495852c15af01598f62563ae534460388b6a370bfbbb8bb6094b698e9" {
		t.Errorf("unexpected record response: %+v", token)
	}
}

func setup(t *testing.T) (*http.ServeMux, *APIClient) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/user/login", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if user := q.Get("user"); user != "admin" {
			t.Errorf("wrong user, expected: %v, got: %v", "admin", user)
		}
		if pass := q.Get("pass"); pass != "admin" {
			t.Errorf("wrong pass, expected: %v, got: %v", "admin", pass)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"displayName": "Administrator",
			"username": "admin",
			"token": "932b2a3495852c15af01598f62563ae534460388b6a370bfbbb8bb6094b698e9",
			"status": "ok"
		}`)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	config := &Configuration{BaseURL: server.URL, User: "admin", Pass: "admin", Debug: true}
	client := NewAPIClient(config)

	return mux, client
}
