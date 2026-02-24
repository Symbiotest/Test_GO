# test-golang-deep (chi + faux oapi-codegen)

Chaîne reproduite:

`main -> http server -> chi -> handler oapi (simulé) -> middleware oapi -> route handler -> controller`

Le controller exécute volontairement un `exec.Command("/bin/bash", "-c", "time sleep 0.01")` pour reproduire un appel indirect similaire au code du screen.

## Lancer les tests

```bash
cd test-golang-deep
go test ./...
```

## Lancer le serveur

```bash
cd test-golang-deep
go run .
```

Puis:

```bash
curl -X POST http://localhost:8080/customers \
  -H 'Content-Type: application/json' \
  -d '{"distributorId":"dist-42","distributorName":"ACME Distribution"}'
```
# Test_GO
