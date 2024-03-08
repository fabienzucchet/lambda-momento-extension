build:
	cd src; GOOS=linux GOARCH=arm64 go build -o ../bin/extensions/lambda-cache-layer main.go

package: build
	cd bin; zip -r lambda-cache-layer.zip extensions

deploy: build package
	cd bin; aws lambda publish-layer-version --layer-name lambda-cache-layer --zip-file fileb://lambda-cache-layer.zip --compatible-runtimes go1.x python3.12 nodejs20.x --compatible-architectures arm64