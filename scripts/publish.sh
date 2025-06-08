#!/bin/bash

echo "Creating release $DRONE_TAG"

# create release
curl \
  -s \
  -w "%{http_code}" \
  -o response.json \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: token $CODEBERG_TOKEN" \
  --data "{\"tag_name\":\"$DRONE_TAG\",\"name\":\"$DRONE_TAG\",\"draft\":false}" \
  $CODEBERG_URL

ASSET_UPLOAD_URL=$(cat response.json | jq -r .upload_url)

# upload assets
echo "Uploading assets to $ASSET_UPLOAD_URL"

ARTIFACTS=("gowt-linux-amd64.tar.gz" "gowt-win-amd64.tar.gz")

for ARTIFACT in "${ARTIFACTS[@]}"; do
  echo "Uploading: $ARTIFACT"

  curl \
    -s \
    -w "%{http_code}" \
    -o response.json \
    -X POST \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: token $CODEBERG_TOKEN" \
    -F "attachment=@./build/$ARTIFACT" \
    $ASSET_UPLOAD_URL
done
