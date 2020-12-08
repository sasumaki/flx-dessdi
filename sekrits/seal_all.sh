#!/bin/bash
for filename in sekrits/unsealed/*.yaml; do
    [ -e "$filename" ] || continue
    kubeseal --format=yaml --cert=sekrits/pub-sealed-secrets.pem < ${filename} > ${filename//unsealed/staging}
done