# End-to-end Encrypted File-Sharing System

Designed and developed an end-to-end encrypted file-sharing system that is secure against data store threats and revoked users.


## Details

Technologies: golang

The system includes Server APIs (Keystore, Datastore), a Client Application API, and cryptographic functions.

- Keystore: key-value store, trusted server where users can publish their public keys
- Datastore: key-value store, untrusted server that provides persistent storage


## Design

1) Data Structures
  - Record (each user): Username, PersonalKey, DecryptionKey, SignatureKey, and PersonalUUID
  - Data struct: Datastore content (Encrypted, Authenticator byte arrays)
  - File struct: basic file (starting ID, Key)
  - InvitationMeta struct: meta for a file invitation (UUID, Key)
  - File Meta struct: meta for a file (UUID, Successor status, Key, successor data)

2) User Authentication
- protocol: each user has a unique username and a password, username and the hash H(H(password || username))[1] (H(x) = hash of x,  x || y = the string concatenation of y to x) will be stored in the keystore. When a user logs in, the H(password input || username) will be compared to the hash stored in the datastore. If they match, the user will be authenticated. If not, access will be denied.
- Information stored in Datastore per user: File structs and files (i.e., the linked lists that comprise files), File namespaces, Key dictionaries, Login structs, Private encryption keys
- Information stored in Keystore per user: Public encryption keys
- Running multiple client instances simultaneously: Before any action, the client will “pull” to ensure that it has up-to-date information. After any action, the client will “push” to ensure that the system is updated and any subsequent action on any device will be pulling the updated version.

4) File Storage and Retrieval
- Storing and retrieving files from the server: Files will be stored as the union of two parts: the file data and the metadata. The metadata is the file struct, which will be stored in Datastore. The file data will be stored as a linked list of blocks, all of which will also be stored in Datastore. Files will be encrypted using a symmetric encryption scheme, whose key is stored in the key dictionary of any given user with access. File retrieval is performed by decrypting the ciphertext in the blocks of the linked list. Iterate through the linked list and stop when a block does not point to a next block.
- Supporting efficient file append: Files will be saved as a linked list of blocks of a fixed size. This way, whenever the file size increases, instead of having to find memory for the entire file and then relocate the file, we can just find an unused block of memory and add it to the linked list.

5) File Sharing and Revocation
- Sharing files with other users: Let User A be the owner of the file “FileA”. User A wants to share the file with User B. When User A shares “FileA” with User B, an invitation is generated and placed randomly in the datastore, the location of which is sent to User B. User B can then use the invitation to access the file.
- File revocation: Because anyone except the owner of a file revoking access is undefined behavior, we can simply check to make sure that the user attempting to revoke access is, in fact, the owner. If not, deny the revocation.
- ensuring revoked users can't take malicious actions on a file: When certain users no longer have access to the document, there are no malicious actions that can be taken. Even if they were to perform some malicious action, the action would be performed on the old location, which now consists of garbage. There is nothing there to read, edit, append to, or otherwise act upon.


## Organization
- implementation in `client/client.go`
- tests in `client_test/client_test.go`.


## Testing
run `go test -v` inside of the `client_test` directory
