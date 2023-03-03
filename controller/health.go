package controller

import "net/http"

func PongCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<html>
		<body style="
			min-height: 100vh;
			display: flex;
			justify-content: center;
			align-teims: center;
		">
			<h1>Pong</h1>
			<h2>Hi!</h2>
		</body>
	</html>`))
}
