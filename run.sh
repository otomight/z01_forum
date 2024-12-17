# create autosigned certificate
# and a pair of 4096 bits rsa keys without password (public and private)
openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
		-keyout server.key -out server.crt \
		-subj "/C=FR/ST=Normandie/L=Rouen/O=Zone01/OU=P12024/CN=localhost"
npx tsc
go build -o main
sudo ./main
