package services

import (
	"context"
	"errors"
	"log"
	"madaurus/dev/material/app/models"
	"madaurus/dev/material/app/shared"
	"madaurus/dev/material/app/utils"
	"time"

	"github.com/permitio/permit-golang/pkg/permit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// GetModulesByFilter Basic Usage  : GetModulesByFilter(ctx, collection, filterStruct, "public", nil) for public endpoints
// Advanced Usage: GetModulesByFilter(ctx, collection, filterStruct, "private", &teacherId) for private endpoints
func GetModulesByFilter(ctx context.Context, collection *mongo.Collection) ([]models.Module, error) {
	// variable must e pointer to slice
	modules := &[]models.Module{}
	opts := options.Find().SetProjection(bson.D{{"courses", 0}})
	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Errro While Getting Module", err.Error())

		return nil, err
	}
	cursorError := cursor.All(ctx, modules)
	if cursorError != nil {
		log.Printf("Error While Decoding the Modules %v ", cursorError)
		return nil, cursorError
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			log.Println("failed to close cursor")

		}
	}(cursor, ctx)

	return *modules, nil
}

func EditModuleVisibility(ctx context.Context, collection *mongo.Collection, moduleId string, visibility bool) error {
	filter := bson.D{{"_id", moduleId}}
	update := bson.D{{"$set", bson.D{{"isPublic", visibility}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateModule(ctx context.Context, collection *mongo.Collection, module models.Module) error {
	filter := bson.D{{"_id", module.ID}}
	update := bson.D{{"$set", bson.D{{"name", module.Name}, {"year", module.Year}, {"plan", module.Plan}, {"image", module.Image}}}}
	updatedAt := time.Now()
	module.UpdatedAt = updatedAt
	newModule, err := collection.UpdateOne(ctx, filter, update)

	if err != nil || newModule.ModifiedCount == 0 {
		log.Printf("Error in Mongo Module Update  : %v\n", err)
		return errors.New("unable to update the module")
	}
	return nil
}

func CreateModule(ctx context.Context, collection *mongo.Collection, module models.Module, permitApi *permit.Client, client *mongo.Client) error {
	module.ID = primitive.NewObjectID()
	module.Courses = []models.Course{}
	module.CreatedAt = time.Now()
	module.UpdatedAt = module.UpdatedAt
	session, errS := client.StartSession()
	transactionOption := options.Transaction().SetReadPreference(readpref.Primary())
	if errS != nil {
		log.Printf("Error While Trying to start session %v", errS.Error())
		return errors.New(shared.UNABLE_CREATE_MODULE)
	}
	errS = session.StartTransaction(transactionOption)
	if errS != nil {
		log.Printf("Error While Trying to start session %v", errS.Error())
		return errors.New(shared.UNABLE_CREATE_MODULE)
	}
	_, err := collection.InsertOne(ctx, module)
	if err != nil {
		session.AbortTransaction(ctx)
		log.Printf("error while trying to create the module")
		return errors.New(shared.UNABLE_CREATE_MODULE)
	}
	err = utils.CreateResourceInstance(permitApi, "modules", module.ID.Hex(), nil, nil, nil)
	if err != nil {
		session.AbortTransaction(ctx)
		return errors.New(shared.UNABLE_CREATE_MODULE)

	}
	err = session.CommitTransaction(ctx)
	if err != nil {
		return errors.New(shared.UNABLE_CREATE_MODULE)

	}
	defer func(ctx context.Context) {
		session.EndSession(ctx)
	}(ctx)
	return nil
}

func vGetModuleById(ctx context.Context, collection *mongo.Collection, moduleId primitive.ObjectID) (models.ExtendedModule, error) {
	// make aggregation to get the courses
	// then select sections from the courses
	// then select the lectures from the sections and videos from sections
	var module models.ExtendedModule

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"_id": moduleId},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "_id",
				"foreignField": "module_id",
				"as":           "courses",
			},
		},
		bson.M{
			"$unwind": "$courses",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "sections",
				"localField":   "courses._id",
				"foreignField": "course_id",
				"as":           "courses.sections",
			},
		},
		bson.M{
			"$unwind": "$courses.sections",
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "lectures",
				"localField":   "courses.sections._id",
				"foreignField": "section_id",
				"as":           "courses.sections.lectures",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "videos",
				"localField":   "courses.sections._id",
				"foreignField": "section_id",
				"as":           "courses.sections.videos",
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error While Getting module details")
		return models.ExtendedModule{}, err

	}
	errCursor := cursor.All(ctx, &module)
	if errCursor != nil {
		return models.ExtendedModule{}, errCursor

	}
	return module, nil
}

func GetModuleById(ctx context.Context, collection *mongo.Collection, moduleId primitive.ObjectID) (models.Module, error) {

	var module models.Module
	err := collection.FindOne(ctx, bson.D{{"_id", moduleId}}).Decode(&module)
	if err != nil {
		log.Printf("Error while retriving the single module:  %v", err.Error())
		return module, errors.New(shared.UNABLE_GET_MODULE)
	}
	return module, nil
}

func vDeleteModule(ctx context.Context, collection *mongo.Collection, moduleId primitive.ObjectID, teacherId *string) error {

	// Before Delete , get now the number of tha courses that this module have
	pipe := bson.A{
		bson.M{
			"$match": bson.M{"_id": bson.M{"$eq": moduleId}},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "_id",
				"foreignField": "module_id",
				"as":           "courses",
			},

			// Count
		},
		bson.M{
			"$project": bson.M{
				"_id":   0,
				"count": bson.M{"$size": "courses"},
			},
		},
		bson.M{
			"$group": bson.M{
				"total": bson.M{"$sum": "$count"},
			},
		},
	}
	result, err := collection.Aggregate(ctx, pipe)
	if err != nil {
		log.Printf("Error While Trying To Deeply delete Module %v", err.Error())
		return err

	}
	var count int64
	errR := result.Decode(&count)
	if errR != nil {
		log.Printf("Error While Trying To Deeply delete Module %v", errR.Error())
		return errR
	}
	if count > 0 {
		return errors.New("unable to delete this module")

	}
	var filter bson.D
	if (teacherId != nil) && (*teacherId != "") {

		filter = bson.D{{"_id", moduleId}, {"teacher_id", *teacherId}}
	} else {
		filter = bson.D{{"_id", moduleId}}
	}

	res, errD := collection.DeleteOne(ctx, filter)
	if errD != nil || res.DeletedCount == 0 {
		log.Printf("Error While Trying To delete Module %v", errD.Error())
		return errors.New("unable to delete this module")

	}
	return nil
}
func DeleteModule(ctx context.Context, collection *mongo.Collection, moduleId primitive.ObjectID) error {

	rs, err := collection.DeleteOne(ctx, bson.D{{"_id", moduleId}, {
		"courses", bson.D{{"$size", 0}},
	}})
	if err != nil {
		log.Printf("Error While Deleting the module %v", err.Error())
		return errors.New(shared.UNABLE_CREATE_MODULE)

	}
	if rs.DeletedCount < 1 {
		log.Printf("unable to delete the module ")
		return errors.New(shared.UNABLE_CREATE_MODULE)
	}

	return nil

}
func CreateManyModules(ctx context.Context, collection *mongo.Collection, modules []models.Module) error {
	var docs []interface{}
	for _, module := range modules {
		module.ID = primitive.NewObjectID()
		module.CreatedAt = time.Now()
		module.UpdatedAt = module.CreatedAt
		module.Courses = []models.Course{}
		docs = append(docs, module)
	}
	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("error while trying to create the module")
	}
	return err
}

func GetModulesByTeacher(ctx context.Context, collection *mongo.Collection, teacherId string) ([]models.Module, error) {
	var modules []models.Module
	filter := bson.D{{"teacher_id", teacherId}}
	log.Printf("teacher id %v", teacherId)
	opts := options.Find().SetProjection(bson.D{{"courses", 0}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		return nil, cursorError
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			log.Println("failed to close cursor")

		}
	}(cursor, ctx)
	return modules, nil
}
func ModuleSelector(ctx context.Context, collection *mongo.Collection, modulesId []string) ([]models.Module, error) {
	log.Printf("Debugging the Modules")
	var modules []models.Module
	modulesIds := make([]primitive.ObjectID, len(modulesId), len(modulesId))
	for _, module := range modulesId {
		modulesIdsObj, _ := primitive.ObjectIDFromHex(module)
		modulesIds = append(modulesIds, modulesIdsObj)
	}
	log.Printf("Len %v", modulesIds)
	opts := options.Find().SetProjection(bson.D{{"courses", 0}}).SetSort(bson.D{{"year", 1}})
	filter := bson.D{{"_id", bson.D{{"$in", modulesIds}}}}
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error While Getting the Modules %v", err.Error())
		return modules, errors.New(shared.UNABLE_GET_MODULES)
	}

	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		log.Printf("Error While Getting the Modules %v", cursorError.Error())
		return modules, cursorError
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			log.Println("failed to close cursor")

		}
	}(cursor, ctx)
	return modules, nil
}

func GetModuleByStudent(ctx context.Context, collection *mongo.Collection, year string) ([]models.Module, error) {
	var modules []models.Module
	filter := bson.D{{"year", year}}
	opts := options.Find().SetProjection(bson.D{{"courses", 0}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error While Getting the Modules %v", err.Error())
		return modules, errors.New(shared.UNABLE_GET_MODULES)
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		log.Printf("Error While Getting the Modules %v", cursorError.Error())
		return modules, errors.New(shared.UNABLE_GET_MODULES)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			log.Println("failed to close cursor")

		}
	}(cursor, ctx)
	return modules, nil
}

func GetModulesByAdmin(ctx context.Context, collection *mongo.Collection) ([]models.Module, error) {
	var modules []models.Module
	filter := bson.D{}
	opts := options.Find().SetProjection(bson.D{{"courses", 0}})
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error While Getting the Modules %v", err.Error())
		return modules, errors.New(shared.UNABLE_GET_MODULES)
	}
	cursorError := cursor.All(ctx, &modules)
	if cursorError != nil {
		log.Printf("Error While Getting the Modules %v", cursorError.Error())
		return modules, errors.New(shared.UNABLE_GET_MODULES)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

			log.Println("failed to close cursor")

		}

	}(cursor, ctx)
	return modules, nil
}

func GetModuleByIdOverview(ctx context.Context, collection *mongo.Collection, moduleId primitive.ObjectID) (models.Module, error) {
	var module models.Module
	opts := options.FindOne().SetProjection(bson.D{{"courses", 0}})
	err := collection.FindOne(ctx, bson.D{{"_id", moduleId}}, opts).Decode(&module)
	if err != nil {
		log.Printf("Error while retriving the single module:  %v", err.Error())
		return module, errors.New(shared.UNABLE_GET_MODULE)
	}
	return module, nil
}
