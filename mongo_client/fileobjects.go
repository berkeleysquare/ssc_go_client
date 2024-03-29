package mongo_client

import (
	"context"
	"encoding/csv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type SearchObject struct {
	Name     string `json:"name" bson:"_id"`
	Manifest string `json:"manifest"`
	Share    string `json:"share" bson:"storagelocation"`
}

var (
	timeout = 90 * time.Minute
)

func DisplaySearchObjects(w *csv.Writer, files []*SearchObject) error {
	lines := [][]string{}
	project := ""
	for fileIndex := range files {
		file := files[fileIndex]
		lines = append(lines, []string{file.Name, file.Manifest, project, file.Share})
	}
	return w.WriteAll(lines)
}

func queryPath(collection *mongo.Collection, name string, exts []string) ([]*SearchObject, error) {

	var extsBson bson.A
	if len(exts) == 0 || exts[0] == "*" {
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
							bson.D{{"componentpath", primitive.Regex{Pattern: strings.ToLower(name)}}},
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
					{"storagelocation", 1},
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
