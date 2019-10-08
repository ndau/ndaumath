#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import http.server
from http.server import HTTPServer, BaseHTTPRequestHandler
import socketserver

PORT = 8080

Handler = http.server.SimpleHTTPRequestHandler

Handler.extensions_map = {
    '.manifest': 'text/cache-manifest',
    '.html': 'text/html',
    '.png': 'image/png',
    '.jpg': 'image/jpg',
    '.svg':	'image/svg+xml',
    '.css':	'text/css',
    '.js':	'application/x-javascript',
    # This is important. Javascript will refuse to use anything unless it comes back with application/wasm.
    '.wasm': 'application/wasm',
    '': 'application/octet-stream',  # Default
}

httpd = socketserver.TCPServer(("", PORT), Handler)

print("serving at port", PORT)
httpd.serve_forever()
