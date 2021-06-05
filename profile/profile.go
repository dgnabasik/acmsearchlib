package profile

// profile.go manages users and encryption. See https://itnext.io/encrypt-data-with-a-password-in-go-b5366384e291
// Transform password to a suitable key using a key derivation function (KDF) which stretches the password to make it a suitable cryptographic key.
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
	"golang.org/x/crypto/scrypt" // KDF
)

/*************************************************************************************/

// Encrypt func uses AES symmetric-key.
func Encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := DeriveKey(key, nil)
	if err != nil {
		return nil, err
	}
	blockCipher, err := aes.NewCipher(key)
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
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

// Decrypt func.
func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]
	key, _, err := DeriveKey(key, salt)
	if err != nil {
		return nil, err
	}
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

// DeriveKey func. Each password has to be checked with the salt used to derive the key.
// The salt needs to be randomly generated. Tt doesn't need to be secret, it needs to be unique.
// Use 16384 (2^14) iterations for interactive logins. Use 1048576 (2^20) iterations for file encryption.
func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}
	const iterations = 16384
	key, err := scrypt.Key(password, salt, iterations, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}

// EncryptData func
func EncryptData(password, textdata string) string {
	pwd := []byte(password)
	data := []byte(textdata)

	ciphertext, err := Encrypt(pwd, data)
	if err != nil {
		log.Panic(err)
	}

	return hex.EncodeToString(ciphertext)
}

func DecryptData(password string, ciphertext []byte) (string, error) {
	pwd := []byte(password)

	plaintext, err := Decrypt(pwd, ciphertext)
	if err != nil {
		log.Panic(err)
	}

	return string(plaintext), err
}

/*************************************************************************************************/

// GetUserProfile func assumes unique case-insensitive userEmail.
func GetUserProfile(userName, pwdText string) (hd.UserProfile, error) {
	var user hd.UserProfile
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return user, err
	}
	defer db.Close()

	encryptedPwd := EncryptData(pwdText, pwdText)
	SELECT := "SELECT id, UserName, UserEmail, Password, AcmMemberId, DateUpdated FROM UserProfile WHERE LOWER(UserName)='" + strings.ToLower(userName) + "' AND Password='" + encryptedPwd + "'"
	err = db.QueryRow(context.Background(), SELECT).Scan(&user.ID, &user.UserName, &user.UserEmail, &user.Password, &user.AcmMemberId, &user.DateUpdated)
	dbx.CheckErr(err)
	if err != nil {
		return user, err
	}
	if dbx.NoRowsReturned(err) {
		return user, errors.New("username/password not found")
	}
	return user, nil
}

// InsertUserProfile func checks for unique username. Store encrypted password.
func InsertUserProfile(userName, userEmail, pwdText string, acmmemberid int) (hd.UserProfile, error) {
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
	encryptedPwd := EncryptData(pwdText, pwdText)
	INSERT := "INSERT INTO UserProfile (UserName, UserEmail, Password, AcmMemberId) VALUES ($1, $2, $3, $4) returning id"
	err = db.QueryRow(context.Background(), INSERT, userName, userEmail, encryptedPwd, acmmemberid).Scan(&id)
	dbx.CheckErr(err)

	existingUser = hd.UserProfile{ID: id, UserName: userName, UserEmail: userEmail, Password: encryptedPwd, AcmMemberId: acmmemberid, DateUpdated: time.Now().UTC()}
	return existingUser, nil
}
