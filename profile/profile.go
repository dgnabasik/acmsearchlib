package profile

// profile.go manages users.

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
)

/*************************************************************************************/

// Version func
func Version() string {
	return "1.16.2"
}

// Encrypt func.  See https://itnext.io/encrypt-data-with-a-password-in-go-b5366384e291
func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key) // AES symmetric-key
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher) // Galois Counter Mode
	if err != nil {
		return nil, err
	}
	// the nonce doesn't have to be secret, it just has to be unique.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt func.
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// GenerateKey func. AES needs a key with length 32 bytes.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptData func
func EncryptData(textdata string) string {
	data := []byte(textdata)
	key, err := GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	ciphertext, err := Encrypt(key, data)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(ciphertext)
	/*fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))
	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		log.Fatal(err)
	}*/
}

/*************************************************************************************************/

// GetUserProfile func assumes unique case-insensitive userName.
func GetUserProfile(userName, pwdText string) (hd.UserProfile, error) {
	var user hd.UserProfile
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return user, err
	}
	defer db.Close()

	encryptedPwd := EncryptData(pwdText)
	SELECT := "SELECT id, UserName, Password, DateUpdated FROM UserProfile WHERE LOWER(UserName)='" + strings.ToLower(userName) + "' AND Password='" + encryptedPwd + "'"
	err = db.QueryRow(context.Background(), SELECT).Scan(&user.ID, &user.UserName, &user.Password, &user.DateUpdated)
	dbx.CheckErr(err)
	if err != nil {
		return user, err
	}
	if dbx.NoRowsReturned(err) {
		return user, errors.New("user/password not found")
	}
	return user, nil
}

// InsertUserProfile func checks for unique username. Store encrypted password.
func InsertUserProfile(userName, pwdText string) (hd.UserProfile, error) {
	existingUser, err := GetUserProfile(userName, pwdText)
	if err != nil && len(existingUser.UserName) > 0 {
		return hd.UserProfile{}, errors.New("duplicate user name")
	}
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.UserProfile{}, err
	}
	defer db.Close()

	var id int
	encryptedPwd := EncryptData(pwdText)
	INSERT := "INSERT INTO UserProfile (UserName, Password) VALUES ($1, $2) returning id"
	err = db.QueryRow(context.Background(), INSERT, userName, encryptedPwd).Scan(&id)
	dbx.CheckErr(err)

	existingUser = hd.UserProfile{ID: id, UserName: userName, Password: encryptedPwd, DateUpdated: time.Now().UTC()}
	return existingUser, nil
}
