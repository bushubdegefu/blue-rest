package temps

import (
	"os"
	"text/template"
)

func HaproxyFrame() {
	// ####################################################
	//  rabbit template
	haproxy_tmpl, err := template.New("RenderData").Parse(haproxyCfgTemplate)
	if err != nil {
		panic(err)
	}

	haproxy_file, err := os.Create("haproxy.cfg")
	if err != nil {
		panic(err)
	}
	defer haproxy_file.Close()

	err = haproxy_tmpl.Execute(haproxy_file, RenderData)
	if err != nil {
		panic(err)
	}
}

func ServiceFrame() {
	// ####################################################
	//  rabbit template
	haproxy_tmpl, err := template.New("RenderData").Parse(systemdServiceTemplate)
	if err != nil {
		panic(err)
	}

	haproxy_file, err := os.Create("app.service")
	if err != nil {
		panic(err)
	}
	defer haproxy_file.Close()

	err = haproxy_tmpl.Execute(haproxy_file, RenderData)
	if err != nil {
		panic(err)
	}
}

var haproxyCfgTemplate = `
global
	log /dev/log	local0
	log /dev/log	local1 notice
	chroot /var/lib/haproxy
	stats socket /run/haproxy/admin.sock mode 660 level admin
	stats timeout 30s
	user haproxy
	group haproxy
	daemon

	# Default SSL material locations
	ca-base /etc/ssl/certs
	crt-base /etc/ssl/private

	# See: https://ssl-config.mozilla.org/#server=haproxy&server-version=2.0.3&config=intermediate
        ssl-default-bind-ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384
        ssl-default-bind-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256
        ssl-default-bind-options ssl-min-ver TLSv1.2 no-tls-tickets

defaults
	log	global
	mode	http
	option	httplog
	option	dontlognull
        timeout connect 5s
        timeout client  50s
        timeout server  50s


	# Logging Format
    	#log-format "%ci:%cp [%tr] %ft %b/%s %TR/%Tw/%Tc/%Tr/%Ta %ST %B %CC %CS %tsc %ac/%fc/%bc/%sc/%rc %sq/%bq %hr %hs %{+Q}r"

   	# Enable Stats Page
	stats enable
	stats uri /haproxy_stats
	stats realm HAProxy\ Statistics
	stats auth admin:admin

frontend app_frontend
    bind *:8500

    mode http
    option forwardfor
    option http-server-close
    http-request add-header X-Forwarded-Proto https

    # Websocket configuration
    acl is_websocket_ui hdr(Upgrade) -i WebSocket
    acl is_websocket_ui hdr_beg(Host) -i ws
    use_backend app_backend if is_websocket_ui

    default_backend app_backend

backend app_backend
    # Rate limiting configuration
    default-server maxconn 30000  check  observe layer7  error-limit 50  on-error mark-down  inter 1s  rise 30  slowstart 20s
    stick-table type ip size 100k expire 300s store http_req_rate(1s)
    tcp-request content track-sc0 src
    acl rate_limited sc0_http_req_rate gt 50
    http-request track-sc0 src
    http-request deny deny_status 429 if rate_limited
    server my_server localhost:7500

`

var systemdServiceTemplate = `
[Unit]
Description=My App Service
After=network.target

[Service]
ExecStart=/playground/app dev
WorkingDirectory=/playground/
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
`
