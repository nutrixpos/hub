// Package services contains the business logic of the core module of nutrix.
//
// The services in this package are used to interact with the persistence layer
// and perform operations on the data models of the core module.
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	pos_core_models "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SalesService contains the configuration and logger for the sales service.
type SalesService struct {
	// Logger is the logger for the sales service.
	Logger logger.ILogger
	// Config is the configuration for the sales service.
	Config config.Config
}

// format 2006-01-02
// GetSalesPerday returns a slice of models.SalesPerDay and the total count of records in the database.
// It takes two parameters, first and rows, which determine the offset and limit of the query.
// It returns an error if the query fails.
func (ss *SalesService) GetSalesPerday(page_number int, page_size int, tenant_id string) (salesPerDay []pos_core_models.SalesPerDay, totalRecords int, err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	salesPerDay = make([]pos_core_models.SalesPerDay, 0)

	deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		ss.Logger.Error(err.Error())
		return salesPerDay, totalRecords, err
	}

	collection := client.Database(ss.Config.Databases[0].Database).Collection(ss.Config.Databases[0].Tables["sales"])
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"date": -1})

	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"date": 1})
	findOptions.SetLimit(int64(page_size))

	cursor, err := collection.Find(ctx, bson.M{
		"tenant_id": tenant_id,
	}, findOptions)
	if err != nil {
		ss.Logger.Error(err.Error())
		return salesPerDay, totalRecords, err
	}
	defer cursor.Close(ctx)
	var tenant models.Tenant

	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&tenant); err != nil {
			return salesPerDay, totalRecords, err
		}

		for _, spd := range tenant.Sales {
			if spd.Refunds == nil {
				spd.Refunds = make([]pos_core_models.ItemRefund, 0)
			}
			salesPerDay = append(salesPerDay, spd)
		}

	}

	return salesPerDay, len(tenant.Sales), nil
}

func (ss *SalesService) InsertClientSalesOrders(tenant_id string, salesPerDayOrder []pos_core_models.SalesPerDayOrder) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	// connected to db

	collection := client.Database(ss.Config.Databases[0].Database).Collection(ss.Config.Databases[0].Tables["sales"])

	// check if document with tenant_id exists, otherwise create it with empty object value
	filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(ctx, bson.D{{Key: "tenant_id", Value: tenant_id}, {Key: "sales", Value: []bson.D{}}})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, sales_order := range salesPerDayOrder {

		filter := bson.M{
			"tenant_id": tenant_id,
			"sales": bson.M{
				"$elemMatch": bson.M{
					"date": sales_order.Order.SubmittedAt.Format("2006-01-02"),
				},
			},
		}

		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		if count == 0 {
			// Date doesn't exist → Add new sales entry
			update := bson.M{
				"$push": bson.M{
					"sales": pos_core_models.SalesPerDay{
						Id:           primitive.NewObjectID().Hex(),
						Date:         sales_order.Order.SubmittedAt.Format("2006-01-02"),
						Orders:       []pos_core_models.SalesPerDayOrder{sales_order},
						Refunds:      []pos_core_models.ItemRefund{},
						Costs:        sales_order.Order.Cost,
						TotalSales:   sales_order.Order.SalePrice,
						RefundsValue: 0,
					},
				},
			}
			_, err = collection.UpdateOne(context.TODO(), bson.M{"tenant_id": tenant_id}, update)
			if err != nil {
				return fmt.Errorf("failed to add order: %v", err)
			}

		} else {
			update := bson.M{
				"$push": bson.M{
					"sales.$[elem].orders": sales_order.Order,
				},
				"$inc": bson.M{
					"sales.$[elem].costs":       sales_order.Order.Cost,
					"sales.$[elem].total_sales": sales_order.Order.SalePrice,
				},
			}

			// Array filter: Ensure we only target the sales entry with the matching date
			arrayFilters := options.ArrayFilters{
				Filters: []interface{}{
					bson.M{"elem.date": sales_order.Order.SubmittedAt.Format("2006-01-02")},
				},
			}

			opts := options.Update().SetArrayFilters(arrayFilters)

			_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				return fmt.Errorf("failed to add order: %v", err)
			}
		}
	}

	return nil
}

func (ss *SalesService) InsertClientSalesRefunds(tenant_id string, salesPerDayRefunds []pos_core_models.LogOrderItemRefund) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ss.Config.Databases[0].Host, ss.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if ss.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	// connected to db

	collection := client.Database(ss.Config.Databases[0].Database).Collection(ss.Config.Databases[0].Tables["sales"])

	// check if document with tenant_id exists, otherwise create it with empty object value
	filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err = collection.InsertOne(ctx, bson.D{{Key: "tenant_id", Value: tenant_id}, {Key: "sales", Value: []bson.D{}}})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, refund := range salesPerDayRefunds {

		filter := bson.M{
			"tenant_id":  tenant_id,
			"sales.date": bson.M{"$exists": true, "$eq": refund.Date.Format("2006-01-02")},
		}

		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return err
		}

		if count == 0 {
			// Date doesn't exist → Add new sales entry
			update := bson.M{
				"$push": bson.M{
					"sales": bson.M{
						"date":   refund.Date.Format("2006-01-02"),
						"orders": []pos_core_models.SalesPerDayOrder{},
						"refunds": []pos_core_models.ItemRefund{
							{
								OrderId:         refund.OrderId,
								ItemId:          refund.ItemId,
								ProductId:       refund.ProductId,
								Amount:          refund.Amount,
								Reason:          refund.Reason,
								ItemCost:        refund.ItemCost,
								Destination:     refund.Destination,
								MaterialRerunds: refund.MaterialRerunds,
								ProductAdd:      refund.ProductAdd,
							},
						},
						"refunds_value": refund.Amount,
					},
				},
			}
			_, err = collection.UpdateOne(context.TODO(), bson.M{"tenant_id": tenant_id}, update)
			if err != nil {
				return fmt.Errorf("failed to add order: %v", err)
			}

		} else {
			update := bson.M{
				"$push": bson.M{
					"sales.$[elem].refunds": pos_core_models.ItemRefund{
						OrderId:         refund.OrderId,
						ItemId:          refund.ItemId,
						ProductId:       refund.ProductId,
						Amount:          refund.Amount,
						Reason:          refund.Reason,
						ItemCost:        refund.ItemCost,
						Destination:     refund.Destination,
						MaterialRerunds: refund.MaterialRerunds,
						ProductAdd:      refund.ProductAdd,
					},
				},
				"$inc": bson.M{"refunds_value": refund.Amount},
			}

			// Array filter: Ensure we only target the sales entry with the matching date
			arrayFilters := options.ArrayFilters{
				Filters: []interface{}{
					bson.M{"elem.date": refund.Date.Format("2006-01-02")},
				},
			}

			opts := options.Update().SetArrayFilters(arrayFilters)

			_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				return fmt.Errorf("failed to add order: %v", err)
			}
		}
	}

	return nil
}
