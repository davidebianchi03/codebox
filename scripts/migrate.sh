# !/bin/bash

atlas migrate apply \
  --url sqlite://codebox.db \
  --dir file://db/migrations