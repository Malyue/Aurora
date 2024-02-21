package client

import (
	_message "Aurora/internal/apps/access-server/pkg/message"
	"Aurora/internal/apps/access-server/svc"
	"Aurora/pkg/hash"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type CredentialCrypto interface {
	EncryptCredentials(c *ClientAuthCredentials) ([]byte, error)

	DecryptCredentials(src []byte) (*ClientAuthCredentials, error)
}

type AesCBCCrypto struct {
	Key []byte
}

func NewAesCBCCrypto(key []byte) *AesCBCCrypto {
	keyLen := len(key)
	count := 0
	switch true {
	case keyLen <= 16:
		count = 16 - keyLen
	case keyLen <= 24:
		count = 24 - keyLen
	case keyLen <= 32:
		count = 32 - keyLen
	default:
		key = key[:32]
	}

	if count != 0 {
		key = append(key, bytes.Repeat([]byte{0}, count)...)
	}
	return &AesCBCCrypto{Key: key}
}

func (a *AesCBCCrypto) EncryptCredentials(c *ClientAuthCredentials) ([]byte, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	// generate random iv
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	encryptBody, err := a.Encrypt(jsonBytes, iv)
	if err != nil {
		return nil, err
	}

	var encrypt []byte
	encrypt = append(encrypt, iv...)
	encrypt = append(encrypt, encryptBody...)

	base64Bytes := make([]byte, base64.RawStdEncoding.EncodedLen(len(encrypt)))
	base64.RawStdEncoding.Encode(base64Bytes, encrypt)
	return base64Bytes, nil
}

func (a *AesCBCCrypto) DecryptCredentials(src []byte) (*ClientAuthCredentials, error) {
	encrypt := make([]byte, base64.RawStdEncoding.DecodedLen(len(src)))
	_, err := base64.RawStdEncoding.Decode(encrypt, src)
	if err != nil {
		return nil, err
	}

	// get iv
	var iv []byte
	iv = append(iv, encrypt[:aes.BlockSize]...)
	var encryptBody []byte
	encryptBody = append(encryptBody, encrypt[aes.BlockSize:]...)

	jsonBytes, err := a.Decrypt(encryptBody, iv)
	if err != nil {
		return nil, err
	}

	credentials := ClientAuthCredentials{}
	err = json.Unmarshal(jsonBytes, &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

// Encrypt implement AES CBC encrypt
// iv is a random num
func (a *AesCBCCrypto) Encrypt(src, iv []byte) ([]byte, error) {
	// create an aes block
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	// get the block size
	blockSize := block.BlockSize()

	// padding the src and make it can divide exactly blockSize
	padding := blockSize - len(src)%blockSize
	// PKCS#7
	if padding == 0 {
		padding = blockSize
	}
	src = append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)

	encryptData := make([]byte, len(src))

	if len(iv) != block.BlockSize() {
		iv = a.cbcIVPending(iv, blockSize)
	}

	CBCEncrypter := cipher.NewCBCEncrypter(block, iv)
	CBCEncrypter.CryptBlocks(encryptData, src)

	return encryptData, nil
}

func (a *AesCBCCrypto) Decrypt(src, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, len(src))
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		iv = a.cbcIVPending(iv, blockSize)
	}

	CBCDecrypter := cipher.NewCBCDecrypter(block, iv)
	CBCDecrypter.CryptBlocks(dst, src)

	length := len(dst)
	if length == 0 {
		return nil, errors.New("unpadding")
	}

	// get the padding num
	unpadding := int(dst[length-1])
	if length < unpadding {
		return nil, errors.New("unpadding")
	}
	res := dst[:(length - unpadding)]

	return res, nil
}

func (a *AesCBCCrypto) cbcIVPending(iv []byte, blockSize int) []byte {
	length := len(iv)
	if length < blockSize {
		return append(iv, bytes.Repeat([]byte{0}, blockSize-length)...)
	} else if length > blockSize {
		return iv[0:blockSize]
	}
	return iv
}

type Authenticator struct {
	credentialCrypto CredentialCrypto
	gateway          Gateway
	ctx              *svc.ServerCtx
}

func NewAuthenticator(gateway Gateway, key string, ctx *svc.ServerCtx) *Authenticator {
	k := sha512.New().Sum([]byte(key))
	return &Authenticator{
		credentialCrypto: NewAesCBCCrypto(k),
		gateway:          gateway,
		ctx:              ctx,
	}
}

func (a *Authenticator) MessageInterceptor(client Client, message *_message.Message) bool {
	if client.GetCredentials() == nil {
		return false
	}
	//switch message.Action {
	//
	//}

	if client.GetCredentials().Secrets == nil {
		_ = a.gateway.EnqueueMessage(client.GetInfo().ID, _message.NewMessage(message.GetSeq(), _message.ActionNotifyForbidden, "no credentials"))
		return true
	}

	secret := client.GetCredentials().Secrets.MessageDeliverSecret
	if secret == "" {
		_ = a.gateway.EnqueueMessage(client.GetInfo().ID, _message.NewMessage(message.GetSeq(), _message.ActionNotifyForbidden, "no message deliver secret"))
		return true
	}

	var ticket = message.Ticket
	if len(message.Ticket) != 40 {
		_ = a.gateway.EnqueueMessage(client.GetInfo().ID, _message.NewMessage(message.GetSeq(), _message.ActionNotifyForbidden, "invalid ticket"))
		return true
	}

	sum1 := hash.SHA1(secret + message.To)
	id := client.GetInfo().ID
	expectTicket := hash.SHA1(secret + id.UID() + sum1)

	if strings.ToUpper(ticket) != strings.ToUpper(expectTicket) {
		a.ctx.Logger.Errorf("invalid ticket, expected=%s, actually=%s, secret=%s, to=%s, from=%s", expectTicket, ticket, secret, message.To, id.UID())
		// invalid ticket
		_ = a.gateway.EnqueueMessage(client.GetInfo().ID, _message.NewMessage(message.GetSeq(), _message.ActionNotifyForbidden, "ticket expired"))
		return true
	}
	return false
}

func (a *Authenticator) ClientAuthMessageInterceptor(client Client, message *_message.Message) (intercept bool) {
	if message.Action != _message.ActionAuthenticate {
		return false
	}

	intercept = true

	var err error
	var errMsg string
	var newId ID
	var span int64
	var authCredentials *ClientAuthCredentials

	credential := EncryptedCredential{}
	err = message.Data.Deserialize(&credential)
	if err != nil {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	if len(credential.Credential) < 5 {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	authCredentials, err = a.credentialCrypto.DecryptCredentials([]byte(credential.Credential))
	if err != nil {
		errMsg = "invalid authenticate message"
		goto DONE
	}

	span = time.Now().UnixMilli() - authCredentials.Timestamp
	if span > 1500*1500 {
		errMsg = "credential expired"
		goto DONE
	}

	newId, err = a.updateClient(client, authCredentials)

DONE:
	ac, _ := json.Marshal(authCredentials)
	a.ctx.Logger.Debugf("credential: %s", string(ac))

	if err != nil || errMsg != "" {
		_ = a.gateway.EnqueueMessage(client.GetInfo().ID, _message.NewMessage(message.GetSeq(), _message.ActionNotifyError, errMsg))
		return
	}

	_ = a.gateway.EnqueueMessage(newId, _message.NewMessage(message.GetSeq(), _message.ActionNotifySuccess, nil))
	return
}

func (a *Authenticator) updateClient(client Client, authCredentials *ClientAuthCredentials) (ID, error) {
	client.SetCredentials(authCredentials)
	oldID := client.GetInfo().ID
	newID := NewID("", authCredentials.UserID, "")
	err := a.gateway.SetClientID(oldID, newID)
	//if err != nil && err.Error() == errClientAlreadyExist {
	//	// if userid and device is equals, return it directly
	//	if newID.Equals(oldID) {
	//		return newID, nil
	//	}
	//	err = a.gateway.SetClientID(newID, "")
	//	if err != nil {
	//		return "", err
	//	}
	//	//kickOut := _message.NewMessage(0, _message.ActionNotifyKickOut,&_message.Ki)
	//}
	return newID, err
}
