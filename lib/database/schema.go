package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdditionalField struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

type Entry struct {
	Id               primitive.ObjectID `bson:"_id"`
	Username         []byte             `bson:"username"`
	Password         []byte             `bson:"password"`
	AdditionalFields []AdditionalField  `bson:"additionalfields"`
}

type SessionToken struct {
	Id     primitive.ObjectID `bson:"_id"`
	N      []byte             `bson:"n"`
	E      int                `bson:"e"`
	Expiry time.Time          `bson:"expiry"`
}

type MasterHash struct {
	Id   primitive.ObjectID `bson:"_id"`
	Hash []byte             `bson:"halt"`
	Salt []byte             `bson:"salt"`
}

type ProtectedDatabaseKey struct {
	Id  primitive.ObjectID `bson:"_id"`
	Key []byte             `bson:"key"`
	Iv  []byte             `bson:"iv"`
}

type ClientConfigMasterKeyMasterHash struct {
	Id          primitive.ObjectID `bson:"_id"`
	KeyFunction string             `bson:"keyfunction"`
	Digest      string             `bson:"digest"`
	Iterations  int                `bson:"iterations"`
}

type ClientConfigDatabaseKey struct {
	Id                 primitive.ObjectID `bson:"_id"`
	EncryptionFunction string             `bson:"encryptionfunction"`
	Size               int                `bson:"size"`
}

type ClientConfigHash struct {
	Id     primitive.ObjectID `bson:"_id"`
	Digest string             `bson:"digest"`
}

type ClientConfig struct {
	Id          primitive.ObjectID              `bson:"_id"`
	MasterKey   ClientConfigMasterKeyMasterHash `bson:"masterkey"`
	MasterHash  ClientConfigMasterKeyMasterHash `bson:"masterhash"`
	DatabaseKey ClientConfigDatabaseKey         `bson:"databasekey"`
	Hash        ClientConfigHash                `bson:"hash"`
}

type Credential struct {
	Id                   primitive.ObjectID   `bson:"_id"`
	MasterHash           MasterHash           `bson:"masterhash"`
	ProtectedDatabaseKey ProtectedDatabaseKey `bson:"protecteddatabasekey"`
	ClientConfigId       primitive.ObjectID   `bson:"clientconfigid"`
}

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	Email         string             `bson:"email"`
	Credentials   []Credential       `bson:"credentials"`
	Entries       []Entry            `bson:"entries"`
	SessionTokens []SessionToken     `bson:"sessiontokens"`
}

// functions that init a struct with default values
func NewAdditionalField() AdditionalField {
	return AdditionalField{}
}

func NewEntry() Entry {
	return Entry{AdditionalFields: []AdditionalField{}}
}

func NewSessionToken() SessionToken {
	return SessionToken{Expiry: time.Now()}
}

func NewMasterHash() MasterHash {
	return MasterHash{}
}

func NewClientConfigMasterKeyMasterHash() ClientConfigMasterKeyMasterHash {
	return ClientConfigMasterKeyMasterHash{}
}

func NewClientConfigDatabaseKey() ClientConfigDatabaseKey {
	return ClientConfigDatabaseKey{}
}

func NewClientConfigHash() ClientConfigHash {
	return ClientConfigHash{}
}

func NewClientConfig() ClientConfig {
	return ClientConfig{}
}

func NewProtectedDatabaseKey() ProtectedDatabaseKey {
	return ProtectedDatabaseKey{}
}

func NewCredential() Credential {
	return Credential{}
}

func NewUser() User {
	return User{
		Credentials:   []Credential{},
		Entries:       []Entry{},
		SessionTokens: []SessionToken{},
	}
}

var Users *mongo.Collection
var ClientConfigs *mongo.Collection

func SetupSchema() {
	Users = Database.Collection("users")
	ClientConfigs = Database.Collection("clientconfigs")

	Users.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
