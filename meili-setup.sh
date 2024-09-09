#!/bin/bash

# Wait for MeiliSearch to start
until $(curl --output /dev/null --silent --head --fail http://localhost:7700); do
  printf '.'
  sleep 10
done

# Step 1: Check if the "manga" index exists
indexExists=$(curl -s 'http://localhost:7700/indexes/manga')

if [ -z "$indexExists" ]; then
  # Step 2: Create the index if it doesn't exist
  curl -X POST 'http://localhost:7700/indexes' \
  --header 'Content-Type: application/json' \
  --data '{"uid":"manga"}'
  echo "Index 'manga' created."
else
  echo "Index 'manga' already exists."
fi

# Step 3: Set filterable attributes
curl -X POST 'http://localhost:7700/indexes/manga/settings/filterable-attributes' \
--header 'Content-Type: application/json' \
--data '["genres.name", "type", "status"]'

echo "Filterable attributes set successfully!"
