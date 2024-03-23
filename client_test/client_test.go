package client_test

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	_ "encoding/hex"
	_ "errors"
	_ "strconv"
	_ "strings"
	"testing"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"
const contentFour = "no one likes 161"

const newPassword = "newpassword"
const password1 = "passwordone"
const password2 = "passwordtwo"
const password3 = "passwordthree"
const password4 = "passwordfour"
const password5 = "passwordfive"




// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"
	xFile := "xfile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Initiating single and multiple users", func() {
		Specify("Initiate a single, new user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Initiate multiple new users.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", newPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Charles.")
			charles, err = client.InitUser("charles", password1)
			Expect(err).To(BeNil())

			userlib.DebugMsg("%s%s", charles.Username, charlesFile)
		})

		Specify("Initiate multiple new users with the same password.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Initializing user Charles.")
			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Initiate new user with password of length zero.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", emptyString)
			Expect(err).To(BeNil())
		})

		Specify("Errors when username already in use.", func() {
			userlib.DebugMsg("Initializing user Alice twice.")
			client.InitUser("alice", defaultPassword)
			alice, err = client.InitUser("alice", newPassword)
			Expect(err).ToNot(BeNil(), "User Alice already initialized.")
		})

		Specify("Errors when username is of length equal to zero.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser(emptyString, defaultPassword)
			Expect(err).ToNot(BeNil(), "User Alice initialized with illegal username.")
		})

		Specify("Errors when user with specified username is not initialized.", func() {
			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).ToNot(BeNil(), "User Alice not initialized.")
		})

		Specify("Errors when user provides an invalid password.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", newPassword)
			Expect(err).ToNot(BeNil(), "Incorrect password.")
		})
	})

	Describe("Storing and Loading for single and multiple users", func() {
		BeforeEach(func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})


		Specify("Store content in file.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
		})

		Specify("Store and Load file content.", func() {
			userlib.DebugMsg("Hello world")
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("\n\n\nWe are now about to store the first file thing")
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("\n\n\nWe have now stored the file and are about to load")

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			userlib.DebugMsg("%s", err)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match original file.")
		})

		Specify("Overwrite file content.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			err = alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())
		})

		Specify("Load overwritten file content.", func() {
			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			alice.StoreFile(aliceFile, []byte(contentOne))
			err = alice.StoreFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil(), "Overwriting file content failed.")

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "File content not succesfully overwritten.")
		})

		Specify("Load multiple stored files' contents.", func() {
			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			alice.StoreFile(aliceFile, []byte(contentOne))
			err = alice.StoreFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading files...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match original file.")

			data, err = alice.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "Loaded file content doesn't match original file.")
		})

		Specify("Store and Load file with empty content.", func() {
			userlib.DebugMsg("Storing file data: %s", emptyString)
			err = alice.StoreFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(emptyString)), "Loaded file content doesn't match original file.")
		})

		Specify("Multiple users store and load files with same filename.", func() {
			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", newPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			alice.StoreFile(xFile, []byte(contentOne))
			err = bob.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading files...")
			data, err := alice.LoadFile(xFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match original file.")

			data, err = bob.LoadFile(xFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "Loaded file content doesn't match original file.")
		})

		Specify("Storing and Loading Files among multiple sessions", func() {
			userlib.DebugMsg("Getting 3 instances of user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alice storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop storing file %s with content: %s", xFile, contentTwo)
			err = aliceLaptop.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop loading file...")
			data, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "aliceLaptop's loaded file content does not match alice's stored file")

			userlib.DebugMsg("alicePhone loading file...")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			data, err = alicePhone.LoadFile(xFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("alicePhone overwrites content loaded from xFile")
			err = alicePhone.StoreFile(xFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop loading file...")
			data, err = aliceDesktop.LoadFile(xFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentThree)), "aliceDesktop's loaded file content does not reflect alicePhone's edits to xFile")
		})

		Specify("Errors when loading a non-existent file.", func() {
			userlib.DebugMsg("Loading file...")
			_, err := alice.LoadFile(aliceFile)
			Expect(err).ToNot(BeNil(), "File aliceFile does not exist.")
		})

		/*Specify("Manual edits to Datastore content do not impact integrity.", func() {
			userlib.DebugMsg("Initializing user Bob.")
			bob, err = client.InitUser("bob", newPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			alice.StoreFile(xFile, []byte(contentOne))
			err = bob.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Manually altering content in DataStore.")
			fUUID, err := uuid.FromBytes(userlib.Hash(append(userlib.Hash([]byte("alice")), userlib.Hash([]byte(xFile))...)))
			Expect(err).To(BeNil())
			userlib.DatastoreSet(fUUID, []byte(contentThree))

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(xFile)
			Expect(err).ToNot(BeNil(), "Datastore impacted by illegal edits.")
			Expect(data).ToNot(Equal([]byte(contentTwo)),
				"Alice's loaded file contents do not match originally stored content.")
		})*/
	})

	Describe("Sharing files among multiple users", func() {
		BeforeEach(func() {
			userlib.DebugMsg("Initializing users.")
			alice, err = client.InitUser("alice", password1)
			bob, err = client.InitUser("bob", password2)
			charles, err = client.InitUser("charles", password3)
			Expect(err).To(BeNil())
		})

		Specify("Share file between two users.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			
			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match orginal file.")
		})

		Specify("Sharing and editing files between mutiple users.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			
			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob overwrites file accepted from Alice.")
			err = bob.StoreFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loads orginally stored file, edited by Bob.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)), "Alice's file content is not updated with Bob's edits")

			userlib.DebugMsg("Bob creating invite for Charles.")
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepting invite from Bob.")
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles appends to accepted file.")
			err = charles.AppendToFile(charlesFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice loads orginally stored file, edited by Charles.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo + contentThree)), "Alice's file content is not updated with Charles' edits")

			userlib.DebugMsg("Bob loads stored file, edited by Charles.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo + contentThree)), "Bob's file content is not updated with Charles' edits")
		})

		Specify("Errors when sharing non-existent file", func() {
			userlib.DebugMsg("Alice creating invite for Bob.")
			_, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).ToNot(BeNil())
		})

		Specify("Errors when sharing file with non-existent user.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for non-existent user Frank.")
			_, err := alice.CreateInvitation(aliceFile, "frank")
			Expect(err).ToNot(BeNil())
		})

		Specify("Errors when accepting file with an existing filename.", func() {
			userlib.DebugMsg("Storing file data (%s) in 2 files", contentOne)
			err = alice.StoreFile(xFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = bob.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(xFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, xFile)
			Expect(err).ToNot(BeNil(), "Accepted fileanme is the same as filename from invitation.")
		})

		Specify("Errors when accepting file which user already possesses.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, aliceFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob creating invite for Alice.")
			invite, err = bob.CreateInvitation(aliceFile, "alice")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice accepting invite from Bob.")
			err = alice.AcceptInvitation("bob", invite, aliceFile)
			Expect(err).ToNot(BeNil())
		})

		Specify("Errors when invitation integrity is insecure.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DatastoreSet(invite, []byte(contentTwo))

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, aliceFile)
			Expect(err).ToNot(BeNil())
		})

	})


	Describe("Appending to files", func() {
		BeforeEach(func() {
			userlib.DebugMsg("Initializing user.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Append to a stored file.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo)), "aliceFile does not reflect edits from append")
		})

		Specify("Append empty byte to stored file.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", emptyString)
			err = alice.AppendToFile(aliceFile, []byte(emptyString))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "aliceFile improperly reflects edits from append")
		})

		Specify("Multiple appends to a stored file.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)), "aliceFile does not reflect edits from appends")
		})

		Specify("Overwrite file after appends.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentFour)
			err = alice.StoreFile(aliceFile, []byte(contentFour))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentFour)), "aliceFile does not reflect overwritten file content")
		})

		Specify("Errors when appending to non-existent file.", func() {
			userlib.DebugMsg("Appending file data: %s", contentOne)
			err = alice.AppendToFile(aliceFile, []byte(contentOne))
			Expect(err).ToNot(BeNil())
		})
	})




	Describe("Revoking access to files", func() {
		BeforeEach(func() {
			userlib.DebugMsg("Initializing users.")
			alice, err = client.InitUser("alice", password1)
			bob, err = client.InitUser("bob", password2)
			charles, err = client.InitUser("charles", password3)
			Expect(err).To(BeNil())
		})

		Specify("Errors when user attempts to revoke access to file to file user doesn't have stored/access to.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revoking Bob's access from %s.", xFile)
			err = alice.RevokeAccess(xFile, "bob")
			Expect(err).ToNot(BeNil(), "Alice revokes access to file that she does not have stored/access to.")
		})

		Specify("Errors when user attempts to revoke access from file thats has not been shared.", func() {
			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = alice.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revoking Bob's access from %s.", xFile)
			err = alice.RevokeAccess(xFile, "bob")
			Expect(err).ToNot(BeNil(), "Alice revokes access to file that she did not share with Bob.")
		})

		Specify("Errors when user attempts to revoke access from uninitialized user.", func() {
			userlib.DebugMsg("Storing files' data: %s and %s", contentOne, contentTwo)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			err = alice.StoreFile(xFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice revoking Charles's access from %s.", xFile)
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).ToNot(BeNil(), "Alice revokes access to file from unrecognized user.")
		})

		Specify("Errors when user attempts to accept revoked invitation to file access.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			// userlib.

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).ToNot(BeNil(), "Bob's access to file bobFile is revoked.")

			userlib.DebugMsg("Loading file...")
			data, err := bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil(), "File content cannot be loaded from file with unrevoked access.")
			Expect(data).ToNot(Equal([]byte(contentOne)), "File content cannot be loaded from file with unrevoked access.")
		})


		Specify("Errors when multiple users attempt to access revoked files.", func() {
			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice creating invite for Bob and Charles.")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			inviteC, err := alice.CreateInvitation(aliceFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob and Charles accepting invite from Alice.")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			err = charles.AcceptInvitation("alice", inviteC, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Both users load files...")
			data, err := bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match orginal file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)), "Loaded file content doesn't match orginal file.")

			userlib.DebugMsg("Alice revoking Bob's and Charles' access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Got through one")
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Got through both")

			userlib.DebugMsg("Checking that Bob/Charles cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil(), "Access to file bobFile revoked.")
			err = charles.AppendToFile(charlesFile, []byte(contentThree))
			Expect(err).ToNot(BeNil(), "Access to file charlesFile revoked.")

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil(), "Access to file bobFile revoked.")
			Expect(data).ToNot(Equal([]byte(contentOne + contentTwo)), "File should not reflect edits from a user with revoked access.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil(), "Access to file charlesFile revoked.")
			Expect(data).ToNot(Equal([]byte(contentOne + contentThree)), "File should not reflect edits from a user with revoked access.")

			userlib.DebugMsg("Checking that Bob/Charles cannot share or accept the file with revoked access.")
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).ToNot(BeNil(), "Access to file bobFile revoked.")
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).ToNot(BeNil(), "Access to file bobFile revoked.")
		
		})
	})


		









	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appended one but not two")

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appended both")

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})
})
