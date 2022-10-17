HOSTNAME=terraform.localhost.com
NAMESPACE=quelhasu
NAME=graphenedb
VERSION=1.1.66
OS_ARCH=amd64
OS=windows
BINARY=terraform-provider-${NAME}_v${VERSION}_${OS}_${OS_ARCH}

echo "Building project ..."
go build -o ./bin/${BINARY}
echo $APPDATA

echo "Creating local terraform directory ..."
mkdir -p "$APPDATA/.terraform.d/test"
mkdir -p "$APPDATA/terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}"
# echo "test"

echo "Copying terraform module ..."
cp ./bin/${BINARY} "$APPDATA/terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}"

# echo "DONE!"
