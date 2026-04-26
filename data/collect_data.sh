#!/bin/bash

IMAGE="alpine"
MAX_CONTAINERS=1250
OUTPUT="data.csv"
NETWORK="none"

echo "[+] Starting "$MAX_CONTAINERS" containers..."

echo "n_containers,startup_ms" > "$OUTPUT"
docker pull "$IMAGE" > /dev/null 2>&1  # Pre-pull the image

launched=()
for i in $(seq 0 $MAX_CONTAINERS); do
	START=$(date +%s%N)
	CID=$(docker run -d --network "$NETWORK" "$IMAGE" sleep infinity)
	END=$(date +%s%N)
	MS=$(( (END - START) / 1000000 ))
	echo "$i,$MS" >> "$OUTPUT"
	launched+=("$CID")
done

# Kill everything
echo "[+] Cleaning up..."
docker rm -f "${launched[@]}" > /dev/null
