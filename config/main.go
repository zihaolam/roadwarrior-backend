package config

import (
	"os"
)

const BucketName = "roadwarrior-pothole-images-upload"
const Region = "ap-southeast-1"
const LocationPlaceIndex = "placeindexfbc375f3-dev"
const CloudfrontHostname = "https://dtt8282y8my28.cloudfront.net"

var ENV = os.Getenv("ENV")
