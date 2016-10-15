cat wwbc-aws-key-pair.pem.dec | openssl enc -aes-256-cbc > wwbc-aws-key-pair.pem.enc
cat            classified.dec | openssl enc -aes-256-cbc > classified.enc
rm wwbc-aws-key-pair.pem.dec
rm classified.dec
