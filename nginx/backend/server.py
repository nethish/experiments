from http.server import BaseHTTPRequestHandler, HTTPServer
import os


class SimpleHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        name = os.getenv("app_name")
        self.send_response(200)
        self.end_headers()
        message = "Hello from Backend server: " + name
        print(message)
        self.wfile.write(bytes(message, encoding="utf-8"))


HTTPServer(("0.0.0.0", 5000), SimpleHandler).serve_forever()
