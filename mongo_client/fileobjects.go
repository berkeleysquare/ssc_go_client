package mongo_client

import (
	"context"
	"encoding/csv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SearchObject struct {
	Name 		string 	`json:"name" bson:"_id"`
	Manifest  	string 	`json:"manifest"`
}

var (
	timeout = 45 * time.Minute
)


func DisplaySearchObjects(w *csv.Writer, files []*SearchObject) error {
	lines := [][]string{}
	for fileIndex := range files {
		file := files[fileIndex]
		lines = append(lines, []string {file.Name, file.Manifest})
	}
	return w.WriteAll(lines)
}

func queryPath(collection *mongo.Collection, name string, ext string) ([]*SearchObject, error) {
	filter := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"$and",
						bson.A{
							bson.D{{"componentpath", primitive.Regex{Pattern: name, Options: "i"}}},
							bson.D{{"componentpath", primitive.Regex{Pattern: ext + "$", Options: "i"}}},
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

