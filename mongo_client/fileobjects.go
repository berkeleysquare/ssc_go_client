package mongo_client

import (
	"context"
	"encoding/csv"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type SearchObject struct {
	Name 		string 	`json:"name" bson:"_id"`
	Path 		string 	`json:"relativepath" bson:"relativepath"`
	Manifest  	string 	`json:"manifest" bson:"manifest"`
	Size	  	int 	`json:"filesize" bson:"filesize"`
	Checksum  	string 	`json:"hash" bson:"hash"`
}

var (
	timeout = 45 * time.Minute
)


func DisplaySearchObjects(w *csv.Writer, files []*SearchObject) error {
	lines := [][]string{}
	for fileIndex := range files {
		file := files[fileIndex]
		lines = append(lines, []string {file.Name, file.Manifest, strconv.Itoa(file.Size), file.Checksum})
	}
	return w.WriteAll(lines)
}


func queryPath(collection *mongo.Collection, name string, exts []string) ([]*SearchObject, error) {

	var extsBson bson.A
	if len(exts) == 0 || exts[0] == "*"{
		// match any -- will inquire about desired behavior
		extsBson = append(extsBson, bson.D{{"componentpath", primitive.Regex{Pattern: ".$", Options: "i"}}})
	} else {
		for _, ext := range exts {
			extsBson = append(extsBson, bson.D{{"componentpath", primitive.Regex{Pattern: "[.]" + ext + "$", Options: "i"}}})
		}
	}

	filter := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"$and",
						bson.A{
							bson.D{{"componentpath", primitive.Regex{Pattern: name, Options: "i"}}},
							bson.D{
								{"$or", extsBson},
							},
							bson.D{
								{"manifest",
									bson.D{
										{"$nin",
											bson.A{
												primitive.Null{},
												"",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"manifest", 1},
					{"filesize", 1},
					{"hash", 1},
				},
			},
		},
	}

	aggOptions := options.AggregateOptions{}
	aggOptions.MaxTime = &timeout

	var ret []*SearchObject
	cur, err := collection.Aggregate(context.TODO(), filter, &aggOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var result SearchObject
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &result)
	}

	return ret, nil
}

func queryProject(collection *mongo.Collection, project string) ([]*SearchObject, error) {

// 			{"$match", bson.D{{"manifest", primitive.Regex{Pattern: project + "-[0-9]*$"}}}}},

	filter := bson.A{
		bson.D{{"$match", bson.D{{"manifest", primitive.Regex{Pattern: "Sourcy-[0-9]*$"}}}}},
		bson.D{
			{"$addFields",
				bson.D{
					{"relativepath",
						bson.D{
							{"$substrCP",
								bson.A{
									"$_id",
									bson.D{
										{"$add",
											bson.A{
												bson.D{{"$strLenCP", "$locationpath"}},
												1,
											},
										},
									},
									bson.D{{"$strLenCP", "$_id"}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"relativepath", 1},
					{"manifest", 1},
					{"filesize", 1},
					{"hash", 1},
				},
			},
		},
	}

	aggOptions := options.AggregateOptions{}
	aggOptions.MaxTime = &timeout

	var ret []*SearchObject
	cur, err := collection.Aggregate(context.TODO(), filter, &aggOptions)
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &ret); err != nil {
		return nil, fmt.Errorf("could not marshall results in queryProject(%s) %v", project, err)
	}
	return ret, nil
}


