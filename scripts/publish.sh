#!/bin/bash

echo "Creating release $CI_COMMIT_TAG"

# create release
curl \
  -s \
  -w "%{http_code}" \
  -o response.json \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: token $CODEBERG_TOKEN" \
  --data "{\"tag_name\":\"$CI_COMMIT_TAG\",\"name\":\"$CI_COMMIT_TAG\",\"draft\":false}" \
  $CODEBERG_URL

ASSET_UPLOAD_URL=$(cat response.json | jq -r .upload_url)

# upload assets
echo "Uploading assets to $ASSET_UPLOAD_URL"

# find all archives in the build directory
shopt -s nullglob  # ensures empty array if no matches
ARTIFACTS=(build/*.tar.gz)

for ARTIFACT in "${ARTIFACTS[@]}"; do
  echo "Uploading: $ARTIFACT"

  curl \
    -s \
    -w "%{http_code}" \
    -o response.json \
    -X POST \
    -H "Content-Type: multipart/form-data" \
    -H "Authorization: token $CODEBERG_TOKEN" \
    -F "attachment=@./$ARTIFACT" \
    $ASSET_UPLOAD_URL
done
