package client

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	// userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	// userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username 		string
	PersonalKey		[]byte
	DecryptionKey	userlib.PrivateKeyType
	SignatureKey	userlib.DSSignKey
	PersonalUUID 	uuid.UUID


	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
}

func (f FileMeta) newStructFile(start []byte, key []byte) (file File, err error) {
	file.Start, file.Key = start, key
	err = storeInDS(f.UUID, file, f.Key)
	return file, err
}

type Data struct {
	Encrypted		[]byte
	Authenticator	[]byte
}

// Returns true if user has been created
func userExists(username string) (exists bool) {
	strings.Compare("", "")
	_, e := userlib.KeystoreGet(username + "e")
	_, v := userlib.KeystoreGet(username + "v")
	return e && v
}

func compare(a, b []byte) bool {
	if len(a) == len(b) {
		for i := 0; i < len(a); i+=1 {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	} else {
		return false
	} 
}

func InitUser(username string, password string) (userdataptr *User, err error) {
	if len(username) == 0 {
		return nil, errors.New("Username can't be empty")
	}

	var userdata User
	userdata.Username = username

	var signatureKey, verificationKey, e = userlib.DSKeyGen()
	if e != nil {
		return nil, e
	}
	err = userlib.KeystoreSet(userdata.Username + "v", verificationKey)
	if err != nil {
		return nil, err
	}

	var encryptionKey, decryptionKey, e2 = userlib.PKEKeyGen()
	if e2 != nil {
		return nil, e2
	}
	err = userlib.KeystoreSet(userdata.Username + "e", encryptionKey)
	if err != nil {
		return nil, err
	}

	userB, passB :=  []byte(username), []byte(password)
	var orginalKey = userlib.Argon2Key(passB, userB, 64)

	userdata.PersonalUUID = uuid.New()
	userdata.SignatureKey = signatureKey
	userdata.PersonalKey = orginalKey[0:32]
	userdata.DecryptionKey = decryptionKey

	userUUID, e3 := uuid.FromBytes(userlib.Hash(userlib.Hash(userB))[:16])
		
	if e3 != nil {
		fmt.Println(e3)
		return nil, e3
	}	

	err = encryptStoreInDS(userdata.PersonalUUID, orginalKey, userdata.PersonalKey)
	if err != nil {
		return nil, err
	}

	err = storeInDS(userUUID, userdata, userdata.PersonalKey) 
	if err != nil {
		return nil, err
	}
	
	return &userdata, nil
}

func getKeyPair(key []byte) (sKey, mKey []byte) {
	return key[:16], key[16:32]
}

func decryptGetData(u uuid.UUID, key []byte) (data []byte, err error) {
	dKey, mKey := getKeyPair(key)

	// userlib.DebugMsg("DatastoreGetFailure")
	// userlib.DebugMsg("uuid in dgd: %s", u.String())

	bytes, ok := userlib.DatastoreGet(u)
	if !ok {
		return nil, errors.New("Data unavailable")
	}

	// // userlib.DebugMsg("GOOD TRIPLE ZERO ONE")

	var wrap Data
	err = json.Unmarshal(bytes, &wrap)
	if err != nil {
		return nil, err
	}

	// userlib.DebugMsg("GOOD TRIPLE ZERO TWO")

	m, err := userlib.HMACEval(mKey, wrap.Encrypted)
	if err != nil {
		return nil, err
	}

	// userlib.DebugMsg("GOOD TRIPLE ZERO THREE")

	if !userlib.HMACEqual(m, wrap.Authenticator) {
		return nil, errors.New("MACs do not match")
	}

	// userlib.DebugMsg("GOOD TRIPLE ZERO FOUR")

	data = userlib.SymDec(dKey, wrap.Encrypted)
	return data, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	if !userExists(username) {
		return nil, errors.New("No user with username " + username)
	}
	u, err := uuid.FromBytes(userlib.Hash(userlib.Hash([]byte(username)))[:16])
	if err != nil {
		return nil, err
	}

	seed := userlib.Argon2Key([]byte(password), []byte(username), 64)

	data, err := decryptGetData(u, seed[:32])
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	userdataptr = &user
	return userdataptr, nil

}

type File struct {
	Start	[]byte
	Key 	[]byte
}

/*func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	storageKey, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return err
	}
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	userlib.DatastoreSet(storageKey, contentBytes)
	return
}*/

func loadData(id []byte, key []byte) (data []byte, err error) {
	u, err := uuid.FromBytes(id[:16])
	if err != nil {
		return nil, err
	}
	return decryptGetData(u, key)
}

func (userdata *User) StoreFile(filename string, content []byte) error {
	storageKey, err := userdata.getFileMetaUUID(filename)
	// userlib.DebugMsg("StorageKey in StoreFile: %s", storageKey)

	if err != nil {
		return err
	}

	_, present := userlib.DatastoreGet(storageKey)
	var file File

	if (!present) {
		var f FileMeta
		f.Key, err = userdata.KeyGen()
		f.Successors = make(map[string]FileMeta)
		f.IsSuccessor = false
		f.UUID = uuid.New()
		if err != nil {
			return err
		}

		err = storeInDS(storageKey, f, userdata.PersonalKey)
		if err != nil {
			return err
		}

		key, err := userdata.KeyGen()
		if err != nil {
			return err
		}
		file, err = f.newStructFile(userlib.RandomBytes(64), key)
		if err != nil {
			return err
		}

		// userlib.DebugMsg("bp 1")

		err = storeData(file.Start, userlib.Hash(file.Start), key)
		if err != nil {
			return err
		}

		// userlib.DebugMsg("bp 2")
	
	} else {
		file, err = userdata.getFile(filename)
		if err != nil {
			return err
		}
	}
	
	file.deleteFile()

	// userlib.DebugMsg("bp 3")
	currId := userlib.Hash(file.Start)
	err = storeData(currId, content, file.Key)
	if err != nil {
		return err
	}

	// userlib.DebugMsg("bp 4")



	err = storeData(file.Start, userlib.Hash(currId), file.Key)

	// userlib.DebugMsg(storageKey.String())
	// userlib.DebugMsg(userdata.getFileMetaUUID(userlib.Hash(storageKey)).String())

	return err
}

func (user User) getFileMetaUUID(filename string) (u uuid.UUID, err error) {
	return uuid.FromBytes(userlib.Hash(append(userlib.Hash([]byte(user.Username)), userlib.Hash([]byte(filename))...))[:16])
}

func (user User) loadFileMeta(filename string) (ret FileMeta, err error) {
	u, err := user.getFileMetaUUID(filename)

	// userlib.DebugMsg("uuid in lfm: %s", u.String())

	if err != nil {
		return ret, err
	}

	// userlib.DebugMsg("GOOD ZERO ZERO ONE")

	bytes, err := decryptGetData(u, user.PersonalKey)
	if err != nil {
		return ret, err
	}

	// // userlib.DebugMsg("GOOD ZERO ZERO TWO")

	err = json.Unmarshal(bytes, &ret)
	return ret, err
}

func (user User) getFile(filename string) (ret File, err error) {
	// // userlib.DebugMsg("GOOD ZERO ZERO")

	fileInfo, err := user.loadFileMeta(filename)
	if err != nil {
		return ret, err
	}

	// userlib.DebugMsg("GOOD ZERO ONE")

	bytes, err := decryptGetData(fileInfo.UUID, fileInfo.Key)
	if err != nil {
		return ret, err
	}

	// userlib.DebugMsg("GOOD ZERO TWO")

	err = json.Unmarshal(bytes, &ret)
	return ret, err
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	file, err := userdata.getFile(filename)
	if err != nil {
		return err
	}

	// userlib.DebugMsg("GOOD ONE")
	// userlib.DebugMsg(string(file.Start))

	id, err := loadData(file.Start, file.Key)
	if err != nil {
		return err
	}

	// userlib.DebugMsg("GOOD TWO")

	err = storeData(id, content, file.Key)
	if err != nil {
		return err
	}

	// userlib.DebugMsg("GOOD THREE")

	err = storeData(file.Start, userlib.Hash(id), file.Key)
	return err
}

func storeData(bytes []byte, data []byte, key []byte) error {
	u, err := uuid.FromBytes(bytes[:16])
	// userlib.DebugMsg("uuid in sc: %s", u)
	if err != nil {
		return err
	}
	return encryptStoreInDS(u, data, key)
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	// userlib.DebugMsg("GOOD ZERO")

	file, err := userdata.getFile(filename)
	if err != nil {
		return nil, err
	}

	// userlib.DebugMsg("GOOD ONE")

	old, err := loadData(file.Start, file.Key)
	if err != nil {
		return nil, err
	}

	// userlib.DebugMsg("GOOD TWO")

	for id := userlib.Hash(file.Start); !compare(id, old); id = userlib.Hash(id) {
		new, err := loadData(id, file.Key)
		if err != nil {
			return nil, err
		}

		content = append(content, new...)
	}

	// userlib.DebugMsg("GOOD THREE")

	return content, nil
}

func (user *User) inviteStore(u uuid.UUID, invInfo InvitationMeta, rec string) error {
	eKey, _ := userlib.KeystoreGet(rec + "e")
	sKey := user.SignatureKey
	toEnc, err := json.Marshal(invInfo)
	if err != nil {
		return err
	}

	enc, err := userlib.PKEEnc(eKey, toEnc)
	if err != nil {
		return err
	}
	
	sign, err := userlib.DSSign(sKey, enc)
	if err != nil {
		return err
	}

	wrap := Data{enc, sign}
	bytes, err := json.Marshal(wrap)
	if err != nil {
		return err
	}

	userlib.DatastoreSet(u, bytes)
	return nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	if !userExists(recipientUsername) {
		return invitationPtr, errors.New("No user with username " + recipientUsername)
	}

	_, err = userdata.getFile(filename)
	if err != nil {
		return invitationPtr, err
	}

	fileInfo, err := userdata.loadFileMeta(filename)
	if err != nil {
		return invitationPtr, err
	}

	if fileInfo.IsSuccessor == true {
		invitationPtr = uuid.New()
		invInfo := InvitationMeta{fileInfo.UUID, fileInfo.Key}
		err = userdata.inviteStore(invitationPtr, invInfo, recipientUsername)
		return invitationPtr, err
	} else {
		k, err := userdata.KeyGen()
		if err != nil {
			return invitationPtr, err
		}

		childInfo := FileMeta { UUID: uuid.New(), Key: k }

		file, err := userdata.getFile(filename)
		if err != nil {
			return invitationPtr, err
		}

		_, err = childInfo.newStructFile(file.Start, file.Key)
		if err != nil {
			return invitationPtr, err
		}

		invitationPtr = uuid.New()

		invInfo := InvitationMeta{childInfo.UUID, childInfo.Key}
		err = userdata.inviteStore(invitationPtr, invInfo, recipientUsername)
		if err != nil {
			return invitationPtr, err
		}

		u, err := uuid.FromBytes(userlib.Hash(append(userlib.Hash([]byte(userdata.Username)), userlib.Hash([]byte(filename))...))[:16])

		parentInfo, err := userdata.loadFileMeta(filename)
		if err != nil {
			return invitationPtr, err
		}

		parentInfo.Successors[recipientUsername] = childInfo
		err = storeInDS(u, parentInfo, userdata.PersonalKey)
		return invitationPtr, err
	}
}

type InvitationMeta struct {
	UUID 	uuid.UUID
	Key 	[]byte
}

type FileMeta struct {
	UUID 			uuid.UUID
	IsSuccessor 	bool
	Key 			[]byte
	Successors		map[string] FileMeta 
}

func storeInDS(u uuid.UUID, object interface{}, key []byte) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	return encryptStoreInDS(u, bytes, key)
}

func encryptStoreInDS(u uuid.UUID, data []byte, key []byte) error {
	eKey, mKey := getKeyPair(key)
	enc := userlib.SymEnc(eKey, userlib.RandomBytes(16), data)
	m, err := userlib.HMACEval(mKey, enc)
	if err != nil {
		return err
	}

	wrap := Data{enc, m}
	bytes, err := json.Marshal(wrap)
	if err != nil {
		return err
	}

	userlib.DatastoreSet(u, bytes)
	return nil
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	if !userExists(senderUsername) {
		return errors.New("No user with username " + senderUsername)
	}

	u, err := userdata.getFileMetaUUID(filename)
	if err != nil {
		return err
	}

	_, ok := userlib.DatastoreGet(u)
	if ok {
		return errors.New("Cannot accept invitation for existing file")
	}

	// Load info
	var invInfo InvitationMeta

	dKey := userdata.DecryptionKey
	vKey, _ := userlib.KeystoreGet(senderUsername + "v")

	bytes, ok := userlib.DatastoreGet(invitationPtr)
	if !ok {
		return errors.New("Invitation Pointer doesn't point to invitation")
	}

	var wrap Data
	err = json.Unmarshal(bytes, &wrap)
	if err != nil {
		return err
	}

	err = userlib.DSVerify(vKey, wrap.Encrypted, wrap.Authenticator)
	if err != nil {
		return err
	}

	data, err := userlib.PKEDec(dKey, wrap.Encrypted)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &invInfo)
	if err != nil {
		return err
	}

	fileInfo := FileMeta {
		UUID: invInfo.UUID,
		IsSuccessor: true,
		Key: invInfo.Key }

	err = storeInDS(u, fileInfo, userdata.PersonalKey)
	if err != nil {
		return err
	}

	_, err = userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	userlib.DatastoreDelete(invitationPtr)
	return nil

}

func (user User) KeyGen() (key []byte, err error) {
	// userlib.DebugMsg("Begin KG")
	seed, err := decryptGetData(user.PersonalUUID, user.PersonalKey)
	if err != nil {
		return nil, err
	}

	seed, err = userlib.HashKDF(seed[48:64], []byte("seed"))
	if err != nil {
		return nil, err
	}

	err = encryptStoreInDS(user.PersonalUUID, seed, user.PersonalKey)
	// userlib.DebugMsg("End KG")
	return seed[:32], err
}

func (file File) deleteFile() error {
	// userlib.DebugMsg("begin delete")
	old, err := loadData(file.Start, file.Key)
	if err != nil {
		return err
	}

	for id := file.Start; !compare(id, old); id = userlib.Hash(id) {
		u, err := uuid.FromBytes(id[:16])
		if err != nil {
			return err
		}
		userlib.DatastoreDelete(u)
	}
	// userlib.DebugMsg("end delete")

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	if !userExists(recipientUsername) {
		return errors.New("No user with username " + recipientUsername)
	}

	file, err := userdata.getFile(filename)
	if err != nil {
		return err
	}

	fileInfo, err := userdata.loadFileMeta(filename)
	if err != nil {
		return err
	}

	content, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	file.Start = userlib.RandomBytes(64)
	file.Key, err = userdata.KeyGen()
	if err != nil {
		return err
	}

	err = storeInDS(fileInfo.UUID, file, fileInfo.Key)
	if err != nil {
		return err
	}

	err = userdata.StoreFile(filename, content)
	if err != nil {
		return err
	}

	userlib.DebugMsg("%d", len(fileInfo.Successors))
	flag := false
	for username, childInfo := range fileInfo.Successors {
		
		bytes, err := decryptGetData(childInfo.UUID, childInfo.Key)
		if err != nil {
			return err
		}

		var child File
		err = json.Unmarshal(bytes, &child)
		if err != nil {
			return err
		}
	
		if username == recipientUsername {
			child.deleteFile()
			userlib.DatastoreDelete(childInfo.UUID)
			flag = true
			delete(fileInfo.Successors, recipientUsername)
		} else {
			child.Start = file.Start
			child.Key = file.Key
			err = storeInDS(childInfo.UUID, child, childInfo.Key)
			if err != nil {
				return err
			}
		}
	}

	if flag == false {
		return errors.New("File not shared with " + recipientUsername)
	}

	return err
}
