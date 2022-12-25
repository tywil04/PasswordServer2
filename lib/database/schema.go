package database

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
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
	ClientKeyDerivation string `bson:"clientkeyderivation"`
	ClientDigest        string `bson:"clientdigest"`
	ClientIterations    int    `bson:"clientiterations"`
	ClientKeyLength     int    `bson:"clientkeylength"`
	ServerKeyDerivation string `bson:"serverkeyderivation"`
	ServerDigest        string `bson:"serverdigest"`
	ServerIterations    int    `bson:"serveriterations"`
	ServerKeyLength     int    `bson:"serverkeylength"`
	MasterHash          []byte `bson:"masterhash"`
	Salt                []byte `bson:"salt"`
}

type ProtectedDatabaseKey struct {
	ClientEncryption     string `bson:"clientencryption"`
	ProtectedDatabaseKey []byte `bson:"protecteddatabasekey"`
	Iv                   []byte `bson:"iv"`
}

type User struct {
	Email                 string                 `bson:"email"`
	MasterHashes          []MasterHash           `bson:"masterhashes"`
	ProtectedDatabaseKeys []ProtectedDatabaseKey `bson:"protecteddatabasekeys"`
	Entries               []Entry                `bson:"entries"`
	SessionTokens         []SessionToken         `bson:"sessiontokens"`
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

func NewProtectedDatabaseKey() ProtectedDatabaseKey {
	return ProtectedDatabaseKey{}
}

func NewUser() User {
	return User{
		MasterHashes:          []MasterHash{},
		ProtectedDatabaseKeys: []ProtectedDatabaseKey{},
		Entries:               []Entry{},
		SessionTokens:         []SessionToken{},
	}
}

var Users *mongo.Collection

func LoadCollections() {
	Users = Database.Collection("users")
}
