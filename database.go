package main

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// User model
type User struct {
    ID            primitive.ObjectID `bson:"_id,omitempty"`
    Username      string             `bson:"username"`
    LearningStyle string             `bson:"learning_style"`
    Metrics       struct {
        Engagement        int `bson:"engagement"`
        TimeSpent         int `bson:"time_spent"`
        ModulesCompleted  int `bson:"modules_completed"`
        TestsTaken        int `bson:"tests_taken"`
        Visual            struct {
            TimeSpent int `bson:"time_spent"`
            Switches  int `bson:"switches"`
        } `bson:"visual"`
        Text struct {
            Score    int `bson:"score"`
            Switches int `bson:"switches"`
        } `bson:"text"`
        Coding struct {
            Score              int `bson:"score"`
            TimeSpent          int `bson:"time_spent"`
            ExecutionFrequency int `bson:"execution_frequency"`
            NumberOfLines      int `bson:"number_of_lines"`
            NumberOfEdits      int `bson:"number_of_edits"`
            Switches           int `bson:"switches"`
        } `bson:"coding"`
        Tests struct {
            CompletionRate int `bson:"completion_rate"`
            AverageScores  int `bson:"average_scores"`
        } `bson:"tests"`
    } `bson:"metrics"`
    Courses []Course `bson:"courses"`
}

// Course model
type Course struct {
    Title    string    `bson:"title"`
    Modules  []Module `bson:"modules"`
}

// Module model
type Module struct {
    Title       string   `bson:"title"`
    Video       string   `bson:"video"`
    Animations  []string `bson:"animations"`
    Image       string   `bson:"image"`
    Text        string   `bson:"text"`
    Flowchart   string   `bson:"flowchart"`
    CodeSnippet string   `bson:"code_snippet"`
}

// MongoDB connection string
var connectionString = "mongodb://localhost:27017"

// MongoDB database name
var dbName = "your_db_name"

// Collection names
var userCollectionName = "users"
var courseCollectionName = "courses"

// MongoDB client
var client *mongo.Client

// Connect to MongoDB
func connectDB() {
    clientOptions := options.Client().ApplyURI(connectionString)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Connected to MongoDB")
}

// Close the MongoDB connection
func closeDB() {
    if err := client.Disconnect(context.Background()); err != nil {
        log.Fatal(err)
    }
    log.Println("Disconnected from MongoDB")
}

// User CRUD functions

// CreateUser creates a new user
func CreateUser(user User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := client.Database(dbName).Collection(userCollectionName)
    _, err := collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    return nil
}

// GetUser retrieves a user by ID
func GetUser(userID primitive.ObjectID) (*User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := client.Database(dbName).Collection(userCollectionName)
    var user User
    err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// UpdateUser updates a user
func UpdateUser(userID primitive.ObjectID, update bson.M) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := client.Database(dbName).Collection(userCollectionName)
    _, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, update)
    if err != nil {
        return err
    }
    return nil
}

// Course CRUD functions...

func main() {
    // Connect to MongoDB
    connectDB()
    defer closeDB()

    // Example usage
    user := User{
        Username:      "omer",
        LearningStyle: "null",
        Metrics: struct {
            Engagement        int
            TimeSpent         int
            ModulesCompleted  int
            TestsTaken        int
            Visual            struct {
                TimeSpent int
                Switches  int
            }
            Text struct {
                Score    int
                Switches int
            }
            Coding struct {
                Score              int
                TimeSpent          int
                ExecutionFrequency int
                NumberOfLines      int
                NumberOfEdits      int
                Switches           int
            }
            Tests struct {
                CompletionRate int
                AverageScores  int
            }
        }{},
        Courses: []Course{
            {
                Title: "algorithms",
                Modules: []Module{
                    {
                        Title:       "introduction to algorithms",
                        LearningStyle: "visual",
                        TestScore:      0.8,
                        CompletedQuestions: 0,
                        EngagementRate: 0,
                    },
                },
            },
        },
    }

    // Insert the user into the database
    err := CreateUser(user)
    if err != nil {
        log.Fatal(err)
    }

    // Get the user from the database
    retrievedUser, err := GetUser(user.ID)
    if err != nil {
        log.Fatal(err)
    }

    // Update the user's metrics
    update := bson.M{"$set": bson.M{"metrics.engagement": 10}}
    err = UpdateUser(retrievedUser.ID, update)
    if err != nil {
        log.Fatal(err)
    }
}
