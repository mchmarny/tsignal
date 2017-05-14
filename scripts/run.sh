#!/bin/bash

# get global vars
. scripts/config.sh

../tfeel-trader/tfeel-trader \
  --refresh=2 \
  --debug
