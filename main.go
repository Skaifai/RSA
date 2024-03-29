package main

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Choose primes:")

	p := new(big.Int)
	p.SetString(readAndCleanInput(reader), 10)
	for IsPrime(p) == false {
		fmt.Println(p.Text(10), "is not a prime number. Choose a prime.")
		p = new(big.Int).SetUint64(stringToInt(readAndCleanInput(reader)))
	}

	q := new(big.Int)
	for q.Cmp(p) == 0 {
		println("Choose another prime")
		q = new(big.Int).SetUint64(stringToInt(readAndCleanInput(reader)))
	}
	q.SetString(readAndCleanInput(reader), 10)
	for IsPrime(q) == false {
		fmt.Println(q, "is not a prime number. Choose a prime.")
		q = new(big.Int).SetUint64(stringToInt(readAndCleanInput(reader)))
	}

	var n big.Int
	n.Mul(p, q)
	var phi big.Int
	phi.Mul(p.Sub(p, big.NewInt(1)), q.Sub(q, big.NewInt(1)))
	fmt.Println("N:", n.String())
	fmt.Println("PHI:", phi.String())

	publicKeys := possiblePubKeys(&n, &phi)
	//fmt.Println(publicKeys)
	chosenPublicIndex := rand.Intn(len(publicKeys))
	publicKey := publicKeys[chosenPublicIndex]
	fmt.Println("PUBLIC KEY:", publicKey.String())

	privateKeys := possiblePrivateKeys(&publicKey, &phi, &n)
	//fmt.Println(privateKeys)
	chosenPrivateIndex := rand.Intn(len(privateKeys))
	privateKey := privateKeys[chosenPrivateIndex]
	fmt.Println("PRIVATE KEY:", privateKey.String())

	fmt.Print("Type the message you want to encrypt: ")
	M := new(big.Int).SetUint64(stringToInt(readAndCleanInput(reader)))
	fmt.Println("ORIGINAL MESSAGE:", M.String())
	encrypted := EncryptMessage(M, &publicKey, &n)
	fmt.Println("ENCRYPTED MESSAGE:", encrypted.String())
	decrypted := DecryptMessage(&encrypted, &privateKey, &n)
	fmt.Println("DECRYPTED MESSAGE:", decrypted.String())
}

func readAndCleanInput(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		println("Error during input read")
	}
	line = strings.Replace(line, "\n", "", 1)
	line = strings.Replace(line, "\r", "", 1)
	return line
}

func stringToInt(s string) uint64 {
	res, _ := strconv.ParseUint(s, 10, 64)
	return res
}

func possiblePubKeys(n *big.Int, phi *big.Int) []big.Int {
	var result []big.Int
	var one big.Int
	one.SetInt64(1)
	var i big.Int
	i.SetInt64(2)
	for i.Cmp(phi) < 0 {
		var firstGCD big.Int
		var secondGCD big.Int

		firstGCD.GCD(nil, nil, &i, n)
		secondGCD.GCD(nil, nil, &i, phi)

		if firstGCD.Cmp(&one) == 0 && secondGCD.Cmp(&one) == 0 {
			var toAppend big.Int
			toAppend.Set(&i)
			result = append(result, toAppend)
		}
		i.Add(&i, &one)
	}
	return result
}

func possiblePrivateKeys(publicKey *big.Int, phi *big.Int, n *big.Int) []big.Int {
	var result []big.Int
	var one big.Int
	one.SetInt64(1)
	var i big.Int
	i.Set(&one)
	var end big.Int
	end.Add(phi, big.NewInt(1000))
	for i.Cmp(&end) < 1 {
		var product big.Int
		product.Mul(&i, publicKey)
		var modulo big.Int
		modulo.Mod(&product, phi)
		if modulo.Cmp(&one) == 0 {
			var toAppend big.Int
			toAppend.Set(&i)
			result = append(result, toAppend)
		}
		i.Add(&i, &one)
	}
	return result
}

func GCDEuclidean(a, b uint64) uint64 {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}

func EncryptMessage(message *big.Int, publicKey *big.Int, n *big.Int) big.Int {
	//res := uint64(math.Pow(float64(message), float64(publicKey))) % n
	var res big.Int
	res.Exp(message, publicKey, n)
	return res
}

func DecryptMessage(message *big.Int, privateKey *big.Int, n *big.Int) big.Int {
	//res := uint64(math.Pow(float64(message), float64(privateKey))) % n
	if len(message.Bits()) == 0 {
		return *n
	}
	var res big.Int
	res.Exp(message, privateKey, n)
	return res
}

func IsPrime(number *big.Int) bool {
	var zero big.Int
	zero.SetInt64(0)
	var two big.Int
	two.SetInt64(2)
	var copiedNum big.Int
	copiedNum.Set(number)
	if copiedNum.Mod(number, &two).Cmp(&zero) == 0 {
		return false
	}
	copiedNum.Set(number)

	var half big.Int
	half.Sub(&copiedNum, &two)
	var divisor big.Int
	divisor.SetInt64(3)

	for divisor.Cmp(&half) < 0 {
		if copiedNum.Mod(&copiedNum, &divisor).Cmp(&zero) == 0 {
			return false
		}
		copiedNum.Set(number)
		divisor.Add(&divisor, &two)
	}

	return true
}
