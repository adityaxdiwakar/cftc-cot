# cftc-cot
[![Build Status](https://travis-ci.com/adityaxdiwakar/cftc-cot.svg?branch=master)](https://travis-ci.com/adityaxdiwakar/cftc-cot)

Commitment of Traders Scraper and API open sourced by Loganov Data

## Open Sourced
This API was written, originally, specifically for Natural Gas commitments on NYMEX (Henry Hub). In the process, it made more sense to collect all the data. However, the data is stored inefficiently and the server used to use Python (Flask). This project rewrites the API in Golang using Gorilla's Mux backed by a PostgreSQL database.
