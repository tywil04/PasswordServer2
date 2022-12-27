package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdditionalField struct {
	Field string `bson:"field"`
	Value string `bson:"value"`
}

type Entry struct {
	Username         []byte            `bson:"username"`
	Password         []byte            `bson:"password"`
	AdditionalFields []AdditionalField `bson:"additionalfields"`
}

type SessionToken struct {
	N      []byte    `bson:"n"`
	E      int       `bson:"e"`
	Expiry time.Time `bson:"expiry"`
}

type MasterHash struct {
	Hash []byte `bson:"halt"`
	Salt []byte `bson:"salt"`
}

type ProtectedDatabaseKey struct {
	Key []byte `bson:"key"`
	Iv  []byte `bson:"iv"`
}

type ClientConfigMasterKeyMasterHash struct {
	KeyFunction string `bson:"keyfunction"`
	Digest      string `bson:"digest"`
	Iterations  int    `bson:"iterations"`
}

type ClientConfigDatabaseKey struct {
	EncryptionFunction string `bson:"encryptionfunction"`
	Size               int    `bson:"size"`
}

type ClientConfigHash struct {
	Digest string `bson:"digest"`
}

type ClientConfig struct {
	MasterKey   ClientConfigMasterKeyMasterHash `bson:"masterkey"`
	MasterHash  ClientConfigMasterKeyMasterHash `bson:"masterhash"`
	DatabaseKey ClientConfigDatabaseKey         `bson:"databasekey"`
	Hash        ClientConfigHash                `bson:"hash"`
}

type ConfigProfile struct {
	MasterHash           MasterHash           `bson:"masterhash"`
	ProtectedDatabaseKey ProtectedDatabaseKey `bson:"protecteddatabasekey"`
	Config               ClientConfig         `bson:"clientconfig"`
}

type User struct {
	Email          string          `bson:"email"`
	ConfigProfiles []ConfigProfile `bson:"configprofiles"`
	Entries        []Entry         `bson:"entries"`
	SessionTokens  []SessionToken  `bson:"sessiontokens"`
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

func NewConfigProfile() ConfigProfile {
	return ConfigProfile{}
}

func NewProtectedDatabaseKey() ProtectedDatabaseKey {
	return ProtectedDatabaseKey{}
}

func NewUser() User {
	return User{
		ConfigProfiles: []ConfigProfile{},
		Entries:        []Entry{},
		SessionTokens:  []SessionToken{},
	}
}

var Users *mongo.Collection

func SetupSchema() {
	Users = Database.Collection("users")

	Users.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
