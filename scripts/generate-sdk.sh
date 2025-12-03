#!/bin/bash
set -e

echo "=== Generating TypeScript SDK from OpenAPI ==="

# Check if openapi-generator is installed
if ! command -v openapi-generator &> /dev/null; then
    echo "Installing openapi-generator..."
    npm install -g @openapitools/openapi-generator-cli
fi

# Generate TypeScript SDK
openapi-generator-cli generate \
    -i api/openapi.yaml \
    -g typescript-axios \
    -o frontend/webapp/src/api/generated \
    --additional-properties=typescriptThreePlus=true,withNodeImports=true

echo "✅ TypeScript SDK generated successfully!"
echo "Location: frontend/webapp/src/api/generated"
