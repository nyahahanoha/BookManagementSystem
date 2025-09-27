#!/bin/bash
curl -X POST -k https://books_api.nyahahanoha.net/put:$1 -H "Authorization: ${BOOKS_API_TOKEN}"
