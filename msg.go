package main

import (
    "fmt"
	"os"
	shell "github.com/ipfs/go-ipfs-api"
    "crypto/sha256"
    "crypto/cipher"
    "crypto/aes"
    "crypto/rand"
    "encoding/base64"
    "io"
    "net/http"
    "io/ioutil"
	"errors"
	"strings"
)

// Toma dos cadenas, cryptoText y keyString.
// cryptoText es el texto que se va a descifrar y keyString es la clave que se debe utilizar para el descifrado.
// La función dará salida a la cadena de texto sin formato resultante con una variable de error.
func decryptString(cryptoText string, keyString string) (plainTextString string, err error) {

   // Formatea el keyString para que sea de 32 bytes.
    newKeyString, err := hashTo32Bytes(keyString)

    // Encode la cryptoText a base64.
    cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
    block, err := aes.NewCipher([]byte(newKeyString))
    if err != nil {
        panic(err)
    }
    if len(cipherText) < aes.BlockSize {
        panic("cipherText too short")
    }

    iv := cipherText[:aes.BlockSize]
    cipherText = cipherText[aes.BlockSize:]
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(cipherText, cipherText)
    return string(cipherText), nil
}

// Toma dos cadenas, plainText y keyString.
// plainText es el texto que necesita ser cifrado por keyString.
// La función generará el texto criptográfico resultante y una variable de error.
func encryptString(plainText string, keyString string) (cipherTextString string, err error) {

    // Format the keyString so that it's 32 bytes.
    newKeyString, err := hashTo32Bytes(keyString)

    if err != nil {
        return "", err
    }

    key := []byte(newKeyString)
    value := []byte(plainText)

    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    cipherText := make([]byte, aes.BlockSize + len(value))

    iv := cipherText[:aes.BlockSize]
    if _, err = io.ReadFull(rand.Reader, iv); err != nil {
        return
    }

    cfb := cipher.NewCFBEncrypter(block, iv)
    cfb.XORKeyStream(cipherText[aes.BlockSize:], value)

    return base64.URLEncoding.EncodeToString(cipherText), nil
}

// Como no podemos usar una clave de longitud variable, debemos cortar la clave de los usuarios
// hasta 32 bytes o menos. Para ello la función toma un hash.
// de la clave y la corta a 32 bytes.
func hashTo32Bytes(input string) (output string, err error) {

    if len(input) == 0 {
        return "", errors.New("No input supplied")
    }

    hasher := sha256.New()
    hasher.Write([]byte(input))

    stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

    // Cut the length down to 32 bytes and return.
    return stringToSHA256[:32], nil
}

func main() {

    // Get the amount of arguments from the command line.
    argumentsCount := len(os.Args)

    // Expected usage:
    // encrypt.go -e|-d --key "key here" --value "value here"

    if argumentsCount != 6 {
		fmt.Printf("Author:\ngithub.com/aratan/.\n")
        fmt.Printf("Uso:\n-e encripta, -d desencripta.\n")
        fmt.Printf("--llave \"Contraseña\" que se cargará.\n")
        fmt.Printf("--txt \"Text a encriptar o desencriptar\".\n")
        return
    }

    // Set up some flags to check against arguments.
    encrypt := false
    decrypt := false
    key := false
    expectKeyString := 0
    keyString := false
    value := false
    expectValueString := 0
    valueString := false

    // Set the input variables up.
    encryptionFlag := ""
    stringToEncrypt := ""
    encryptionKey := ""

    // Get the arguments from the command line.
    // If any issues are detected, alert the user and exit.
    for index, element := range os.Args {

        if element == "-e" {
            // Ensure that decrypt has not also been set.
            if decrypt == true {
                fmt.Printf("Can't set both -e and -d.\nBye!\n")
                return
            }
            encrypt = true
            encryptionFlag = "-e"

        } else if element == "-d" {
            // Ensure that encrypt has not also been set.
            if encrypt == true {
                fmt.Printf("Can't set both -e and -d.\nBye!\n")
                return
            }
            decrypt = true
            encryptionFlag = "-d"

        } else if element == "--llave" {
            key = true
            expectKeyString++

        } else if element == "--txt" {
            value = true
            expectValueString++

        } else if expectKeyString == 1 {
            encryptionKey = os.Args[index]
            keyString = true
            expectKeyString = 0

        } else if expectValueString == 1 {
            stringToEncrypt = os.Args[index]
            valueString = true
            expectValueString = 0
        }

        if expectKeyString >= 2 {
            fmt.Printf("Something went wrong, too many keys entered.\bBye!\n")
            return

        } else if expectValueString >= 2 {
            fmt.Printf("Something went wrong, too many keys entered.\bBye!\n")
            return
        }
    }


    // On error, output some useful information.
    if !(encrypt == true || decrypt == true) || key == false || keyString == false || value == false || valueString == false {
        fmt.Printf("Incorrect usage!\n")
        fmt.Printf("---------\n")
        fmt.Printf("-e or -d -> %v\n", (encrypt == true || decrypt == true))
        fmt.Printf("--key -> %v\n", key)
        fmt.Printf("Key string? -> %v\n", keyString)
        fmt.Printf("--value -> %v\n", value)
        fmt.Printf("Value string? -> %v\n", valueString)
        fmt.Printf("---------")
        fmt.Printf("\nUsage:\n-e to encrypt, -d to decrypt.\n")
        fmt.Printf("--key \"I am a key\" to load the key.\n")
        fmt.Printf("--value \"I am a text to be encrypted or decrypted\".\n")
        return
    }


    // Check the encrpytion flag.
    if false == (encryptionFlag == "-e" || encryptionFlag == "-d") {
        fmt.Println("Disculpa el primer argumento debe ser -e o -d")
        fmt.Println("según quieras codificar o decodificar.")
        return
    }

    if encryptionFlag == "-e" {
        // Encrypt!

        //fmt.Printf("Encrypting '%s' with key '%s'\n", stringToEncrypt, encryptionKey)



		encryptedString, _ := encryptString(stringToEncrypt,encryptionKey)
		// añado subir ficheros cifrados a ipfs 
		sh := shell.NewShell("localhost:5001")
		cid, err := sh.Add(strings.NewReader(encryptedString))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
        fmt.Printf("'%s'\n", cid)
        //os.stdout("%s", cid)
		//fmt.Printf("https://gateway.pinata.cloud/ipns/QmUgddgEc71BH5movhDtLJ91tGy3pKs5iUEgsa69ewCNog \n")
        //os.stdout("'%s'\n", encryptedString)
        //fmt.Printf("'%s'\n", encryptedString)
        //fmt.Printf("Salida: '%s'\n", encryptedString)

    } else if encryptionFlag == "-d" {
        // Decrypt!
        //fmt.Printf("Decrypting '%s' with key '%s'\n", stringToEncrypt, encryptionKey)
////////
response, err := http.Get("https://gateway.pinata.cloud/ipns/QmUgddgEc71BH5movhDtLJ91tGy3pKs5iUEgsa69ewCNog")
if err != nil {
    fmt.Printf("%s", err)
    os.Exit(1)
} else {
    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Printf("%s", err)
        os.Exit(1)
    }
    leer := string(contents)
    //fmt.Println(strings.Index(leer, "Q"))
    //i := strings.Index(leer, "")
    //chars := leer[:i]
    //arefun := leer[i+1:]
    fmt.Println("esto es txt %s",leer) //chars)
    //fmt.Println(arefun)
    //fmt.Printf("%s\n", string(contents))
    stringToEncrypt = leer
}

////////
        decryptedString, _ := decryptString(stringToEncrypt,encryptionKey)
        fmt.Printf("'%s'\n", decryptedString)
        //fmt.Printf("Salida: '%s'\n", decryptedString)

    }
}