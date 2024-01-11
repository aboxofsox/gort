#!/bin/bash

url='http://localhost:8080/create'

declare -a users=("John" "Jane" "Jack" "Jill" "Joe" "Jenny")

for user in "${users[@]}"
do
    id=$(uuidgen)
    curl -X POST -d "name=$user&id=$id" $url
done