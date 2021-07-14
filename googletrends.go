package main

import (
	"context"
	"reflect"

	"log"

	"github.com/groovili/gogtrends"
	"github.com/pkg/errors"
)

const (
	locUS  = "US"
	catAll = "all"
	langEn = "EN"
)



func (a *App) getGoogleTrends(keyword string) ([]*gogtrends.Timeline, error) {
	//Enable debug to see request-response
	//gogtrends.Debug(true)

	ctx := context.Background()
	// var sg = new(sync.WaitGroup)

	log.Println("Explore trends:")
	// get widgets for Golang keyword in programming category
	explore, err := gogtrends.Explore(ctx, &gogtrends.ExploreRequest{
		ComparisonItems: []*gogtrends.ComparisonItem{
			{
				Keyword: keyword,
				Geo:     locUS,
				Time:    "today 12-m",
			},
		},
		Property: "",
	}, langEn)
	handleError(err, "Failed to explore widgets")
	printItems(explore)

	log.Println("Interest over time:")
	overTime, err := gogtrends.InterestOverTime(ctx, explore[0], langEn)
	handleError(err, "Failed in call interest over time")
	printItems(overTime)
	return overTime, err


}

func handleError(err error, errMsg string) {
	if err != nil {
		log.Fatal(errors.Wrap(err, errMsg))
	}
}

func printItems(items interface{}) {
	ref := reflect.ValueOf(items)

	if ref.Kind() != reflect.Slice {
		log.Fatalf("Failed to print %s. It's not a slice type.", ref.Kind())
	}

	for i := 0; i < ref.Len(); i++ {
		log.Println(ref.Index(i).Interface())
	}
}
