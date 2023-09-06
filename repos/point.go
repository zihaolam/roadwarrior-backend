package repos

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/oklog/ulid/v2"
	"github.com/zihaolam/roadwarrior-backend/config"
	"github.com/zihaolam/roadwarrior-backend/database"
	"github.com/zihaolam/roadwarrior-backend/entities"
	"github.com/zihaolam/roadwarrior-backend/schemas"
	"github.com/zihaolam/roadwarrior-backend/tools"
)

type IPointRepo struct {
	TableName string
	TablePK   string
	table     dynamo.Table
}

var tableName = "RoadWarrior-Pothole-Points"
var TablePK = "point-PK"

func NewPointRepo() *IPointRepo {
	return &IPointRepo{TableName: "RoadWarrior-Pothole-Statistics", TablePK: TablePK, table: database.GetTable(tableName)}
}

func (repo *IPointRepo) Migrate() error {
	return database.CreateTable(tableName, entities.Point{})
}

func (repo *IPointRepo) Create(createPointBody *schemas.CreatePointSchema) (*entities.Point, error) {
	var point entities.Point
	result, err := tools.UploadBase64FileToS3(createPointBody.Base64PotholeImage)
	if err != nil {
		return nil, err
	}

	locationDetails, err := tools.ReverseGeocode(createPointBody.Coordinates)
	if err != nil {
		return nil, err
	}

	point.PK = TablePK
	point.SK = ulid.Make().String()
	point.Address = *locationDetails.Results[0].Place.Label
	point.Region = *locationDetails.Results[0].Place.Region
	point.SubRegion = *locationDetails.Results[0].Place.SubRegion
	point.ResourceUrl = fmt.Sprintf("%s/%s", config.CloudfrontHostname, result.Key)
	point.Coordinates = tools.FloatSliceToStringSlice(createPointBody.Coordinates)
	point.Count = createPointBody.Count

	err = repo.table.Put(point).Run()
	return &point, err
}

func (repo *IPointRepo) GetAll() ([]*entities.Point, error) {
	var results []*entities.Point
	pagingKey, err := repo.table.Get("PK", repo.TablePK).AllWithLastEvaluatedKey(&results)
	for pagingKey != nil {
		var newResults []*entities.Point
		pagingKey, err = repo.table.Get("PK", repo.TablePK).AllWithLastEvaluatedKey(&newResults)
		results = append(results, newResults...)
	}

	if results == nil {
		return []*entities.Point{}, err
	}

	return results, err
}

func (repo *IPointRepo) GetOne(keys *entities.PointKey) (*entities.Point, error) {
	var result entities.Point
	err := repo.table.Get("PK", keys.PK).Range("SK", dynamo.Equal, keys.SK).One(&result)
	return &result, err
}

func (repo *IPointRepo) UpdateOne(keys *entities.PointKey, updatePointBody *schemas.UpdatePointSchema) (*entities.Point, error) {
	var statToUpdate entities.Point
	err := repo.table.Get("PK", keys.PK).Range("SK", dynamo.Equal, keys.SK).One(&statToUpdate)
	if err != nil {
		return nil, err
	}

	locationDetails, err := tools.ReverseGeocode(updatePointBody.Coordinates)
	if err != nil {
		return nil, err
	}

	bucket, key, err := tools.GetBucketKeyFromS3ObjectUrl(statToUpdate.ResourceUrl)
	if err != nil {
		return nil, err
	}

	err = tools.DeleteObjectFromS3(bucket, key)
	if err != nil {
		return nil, err
	}

	result, err := tools.UploadBase64FileToS3(updatePointBody.Base64PotholeImage)
	if err != nil {
		return nil, err
	}

	statToUpdate.ResourceUrl = tools.GetS3ObjectUrl(&tools.S3Object{Bucket: config.BucketName, Key: result.Key}, config.Region)

	statToUpdate.Address = *locationDetails.Results[0].Place.Label
	statToUpdate.Region = *locationDetails.Results[0].Place.Region
	statToUpdate.SubRegion = *locationDetails.Results[0].Place.SubRegion
	statToUpdate.Coordinates = tools.FloatSliceToStringSlice(updatePointBody.Coordinates)
	statToUpdate.Count = updatePointBody.Count

	err = repo.table.Put(statToUpdate).Run()
	return &statToUpdate, err
}

func (repo *IPointRepo) DeleteOne(keys *entities.PointKey) error {
	return repo.table.Delete("PK", keys.PK).Range("SK", keys.SK).Run()
}

func (repo *IPointRepo) Replicate(points []*entities.Point) error {
	var AWSSession, _ = session.NewSessionWithOptions(session.Options{Profile: "localzone-project"})
	db := dynamo.New(AWSSession, &aws.Config{Region: aws.String("ap-southeast-1")})
	for _, point := range points {
		err := db.Table(tableName).Put(*point).Run()
		fmt.Println(point)
		if err != nil {
			return err
		}
	}
	return nil
}
