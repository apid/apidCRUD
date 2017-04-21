#! /bin/bash
comm -23 <(./all_funcs.sh) <(./tested_funcs.sh)
