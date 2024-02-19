Esempio di utilizzo:

```go
func main() {
	// crea il tuo mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello server"))
	})
	// crea il Trust con i tuoi certificati self signed
	t, err := https.NewTrust("localhost.crt", "localhost.key", "localCA.crt", tls.RequireAndVerifyClientCert)
	if err != nil {
		log.Fatal(err)
	}
	// avvia il server HTTPS
	if e := t.StartServer(mux, 8080); e != nil {
		log.Fatal(e)
	}
}
```

La cartella `sample` contiene un esempio di esecuzione. È necessario avere i certificati self signed creati
con lo script `certs.sh`; utilizzando [Task](https://taskfile.dev/) viene fatto tutto in automatico come riportato
nell'esempio di output qui di seguito

```bash
$> task sample

task: [sample] go run main.go
2024/02/19 14:52:42 server started at port 8080
2024/02/19 14:52:43 ***** Trusted client:
2024/02/19 14:52:43 > hello server
2024/02/19 14:52:43 ***** Bad client error:
2024/02/19 14:52:43 http: TLS handshake error from [::1]:59274: remote error: tls: bad certificate
2024/02/19 14:52:43 Get "https://localhost:8080": tls: failed to verify certificate: x509: certificate signed by unknown authority
2024/02/19 14:52:43 Server closed, BYE :)
```