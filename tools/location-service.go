package tools

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/location"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/zihaolam/roadwarrior-backend/config"
)

var locationService *location.Client

func ReverseGeocode(coordinates []float64) (*location.SearchPlaceIndexForPositionOutput, error) {
	if locationService == nil {
		if config.ENV == "dev" {
			cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
				awsConfig.WithSharedConfigProfile("localzone-project"),
			)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			locationService = location.NewFromConfig(cfg)
		} else {
			locationService = location.New(location.Options{Region: "ap-southeast-1"})
		}
	}
	return locationService.SearchPlaceIndexForPosition(context.TODO(), &location.SearchPlaceIndexForPositionInput{IndexName: aws.String(config.LocationPlaceIndex), Position: coordinates, MaxResults: 1})
}
